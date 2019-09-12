// +build windows
// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package hcs

import (
	"encoding/json"
	"fmt"
	"github.com/Microsoft/hcsshim"
	log "k8s.io/klog"
	"os/exec"
	"reflect"
	"strings"
	"time"

	"github.com/microsoft/wssdagent/common"
	"github.com/microsoft/wssdagent/pkg/ssh"
	"github.com/microsoft/wssdagent/pkg/wssdagent/apis/config"
	"github.com/microsoft/wssdagent/pkg/wssdagent/errors"
	"github.com/microsoft/wssdagent/pkg/wssdagent/store"
	pb "github.com/microsoft/wssdagent/rpc/compute"
	"github.com/microsoft/wssdagent/services/compute/virtualmachine/cloudinit"
	schema "github.com/microsoft/wssdagent/services/compute/virtualmachine/hcs/internal"
	"github.com/microsoft/wssdagent/services/compute/virtualmachine/internal"

	"github.com/microsoft/wssdagent/services/network/virtualnetworkinterface"
	"github.com/microsoft/wssdagent/services/storage/virtualharddisk"
)

const (
	owner    string = "wssdagent"
	memoryMB int    = 8192
	cpuCount int    = 4
)

type client struct {
	vnicprovider virtualnetworkinterface.VirtualNetworkInterfaceProvider
	vhdprovider  virtualharddisk.VirtualHardDiskProvider
	config       *config.ChildAgentConfiguration
	store        *store.ConfigStore
}

func newClient() *client {
	cConfig := config.GetChildAgentConfiguration("VirtualMachine")
	client := &client{
		config:       cConfig,
		store:        store.NewConfigStore(cConfig.DataStorePath, reflect.TypeOf(internal.VirtualMachineInternal{})),
		vnicprovider: (virtualnetworkinterface.GetVirtualNetworkInterfaceProvider(virtualnetworkinterface.HCNSpec)),
		vhdprovider:  (virtualharddisk.GetVirtualHardDiskProvider(virtualharddisk.HCSSpec)),
	}

	return client
}

func (c *client) newVirtualMachine(id string) *internal.VirtualMachineInternal {
	return internal.NewVirtualMachineInternal(id, c.config.DataStorePath)
}

// Create a Virtual Machine
func (c *client) Create(vm *pb.VirtualMachine) (*pb.VirtualMachine, error) {
	log.Infof("[VirtualMachine][Create] [%v]", vm)
	if len(vm.Id) == 0 {
		vm.Id = common.NewGuid()
	}
	vminternal := c.newVirtualMachine(vm.Id)

	// 1. Generate bootstrap data
	switch vm.Os.GetOstype() {
	case pb.OperatingSystemType_WINDOWS:
		// TODO: 1.a Generate Answer File
	case pb.OperatingSystemType_LINUX:
		// 1.b Generate Cloud Init Configuration
		c.generateCloudInitConfiguration(vm, vminternal)
	default:
		return nil, errors.Wrap(errors.NotSupported, "Unsupported Operating System")
	}

	// 2. Get Vhd information
	vhdPath, err := virtualharddisk.GetVirtualHardDiskPath(c.vhdprovider, vm.Storage.Osdisk.Diskid)
	if err != nil {
		return nil, err
	}

	// 3. Get VNic Information
	vnicId := ""
	macAddress := ""
	if len(vm.Network.Interfaces) > 0 {
		vnicName := vm.Network.Interfaces[0].NetworkInterfaceId
		vnic, err := virtualnetworkinterface.GetVirtualNetworkInterfaceByName(c.vnicprovider, vnicName)
		if err != nil {
			return nil, err
		}
		macAddress = vnic.Macaddress
		vnicId = vnic.Id
	}

	// 4. Render VM Spec
	vmspec, err := hcsshim.CreateVirtualMachineSpec(vm.Name, vm.Id, vhdPath, vminternal.SeedIso, owner, memoryMB, cpuCount, vnicId, macAddress)
	if err != nil {
		return nil, err
	}

	log.Infof("[VirtualMachine][Create] hcs spec[%v]", vmspec)

	// 5. Create the VM
	if err = vmspec.Create(); err != nil {
		return nil, err
	}

	// 6. Start the VM
	if err = vmspec.Start(); err != nil {
		return nil, err
	}

	vminternal.Vm = vm

	// 7. Save the config to the store
	c.store.Add(vm.Id, vminternal)

	err = c.addSSHKeys(vm)
	if err != nil {
		// Ignore the error for now
		log.Warningf("Adding SSH Keys failed %v", err)

	}

	return vminternal.Vm, nil
}

// Get a Virtual Machine specified by id
func (c *client) Get(vm *pb.VirtualMachine) ([]*pb.VirtualMachine, error) {
	// Todo: can there be more than one match ?
	vms := []*pb.VirtualMachine{}
	log.Infof("[VirtualMachine][Get] spec[%v]", vm)
	vmname := ""
	if vm != nil {
		vmname = vm.Name
	}

	if len(vmname) == 0 {
		vmsint, err := c.store.List()
		if err != nil {
			return nil, err
		}
		// get everything
		if *vmsint == nil || len(*vmsint) == 0 {
			return nil, nil
		}

		for _, val := range *vmsint {
			vmint := val.(*internal.VirtualMachineInternal)
			if !hcsshim.HasVirtualMachine(vmint.Id) {
				// Store is out of sync with hcs
				c.store.Delete(vmint.Id)
				continue
			}

			vms = append(vms, vmint.Vm)
		}
	} else {
		vmint, err := c.getVirtualMachineInternal(vmname)
		if err != nil {
			return vms, err
		}
		vms = append(vms, vmint.Vm)
	}

	// FIXME: Validate if these VMs actually exists
	log.Infof("[VirtualMachine][Get] Found[%d], VMs[%v]", len(vms), vms)

	return vms, nil
}

// Delete a Virtual Machine
func (c *client) Delete(vm *pb.VirtualMachine) error {
	log.Infof("[VirtualMachine][Delete] spec[%v]", vm)
	// Check the internal store
	vmint, err := c.getVirtualMachineInternal(vm.Name)
	if err != nil {
		return err
	}

	// Check with hcs
	if hcsshim.HasVirtualMachine(vmint.Id) {
		hcsvm, err := getVirtualMachineSpec(vmint.Vm)
		if err != nil {
			return err
		}

		if err = hcsvm.Stop(); err != nil {
			log.Infof("Unable to stop the VM [%v]", err)
		}

		if err := hcsvm.Delete(); err != nil {
			return err
		}
	}
	return c.store.Delete(vm.Id)
}

////////////////////// Private Functions //////////////////////////////////

func (c *client) generateCloudInitConfiguration(vm *pb.VirtualMachine, vminternal *internal.VirtualMachineInternal) error {
	// Initialize Cloud init user data
	hostname := vm.Os.ComputerName // TODO: Use the same logic as Azure VM hostnames
	if strings.Contains(hostname, "_") {
		return fmt.Errorf("linux hostname cannot contain an underscore (_)")
	}
	vm.Os.ComputerName = hostname // Save it
	userdata := cloudinit.CreateUserdata(hostname)
	metadata := cloudinit.CreateMetadata(hostname)

	user := vm.Os.Administrator
	public_auth_keys := []string{}
	for _, keys := range vm.Os.Publickeys {
		public_auth_keys = append(public_auth_keys, keys.Keydata)
	}

	userdata.AddUser(user.Username, user.Username, user.Password, []string{"adm", "cdrom", "sudo", "lxd"}, public_auth_keys, nil)

	//for _, user = range vm.Os.Users {
	//		userdata.AddUser(user.Username, user.Username, user.Password, []string{"adm", "cdrom", "lxd"})
	//	}
	// TODO: Fixup the path

	if len(vm.Os.CustomData) != 0 {
		vendorData := cloudinit.CreateVendordata(vm.Os.CustomData)
		vendorData.RenderYAML(vminternal.VendorData)
		//script := "/opt/wssd/customdata.sh"
		//userdata.AddWriteFile("b64", vm.Os.StartupScript, "root:root", script, "0644")
		//userdata.AddRunCommand([]string{"sh", script})
	}

	err := userdata.RenderYAML(vminternal.UserData)
	if err != nil {
		return err
	}

	err = metadata.RenderJson(vminternal.MetaData)
	if err != nil {
		return err
	}

	err = userdata.GenerateSeedIso([]string{vminternal.UserData, vminternal.MetaData}, []string{vminternal.VendorData}, vminternal.SeedIso)
	if err != nil {
		log.Errorf("[VirtualMachine][Create] Seed iso creation failed [%v]", err)
		return err
	}
	log.Infof("[VirtualMachine][Create] Seed iso created [%v]", vminternal.SeedIso)
	return nil

}

// Conversion function
func getVirtualMachineSpec(vm *pb.VirtualMachine) (*hcsshim.VirtualMachineSpec, error) {
	vmspec, err := hcsshim.GetVirtualMachineSpec(vm.Id)
	if err != nil {
		return nil, err
	}
	vmspec.Name = vm.Name
	return vmspec, nil
}

// Conversion function
func getVirtualMachine(hcsvm *hcsshim.VirtualMachineSpec) (*pb.VirtualMachine, error) {
	vmspecString := hcsvm.String()
	log.Infof("[HCS][%s]", vmspecString)

	internalVmSchema := new(schema.ComputeSystem)
	if err := json.Unmarshal([]byte(vmspecString), internalVmSchema); err != nil {
		return nil, err
	}

	vm := &pb.VirtualMachine{
		Name: hcsvm.Name,
		Id:   hcsvm.ID,
		Storage: &pb.StorageConfiguration{
			Osdisk: &pb.Disk{
				Diskid: "", // internalVmSchema.VirtualMachine.Devices.Scsi["primary"].Attachments["0"].Path,
			},
		},
		Os:      &pb.OperatingSystemConfiguration{},
		Network: &pb.NetworkConfiguration{},
	}

	return vm, nil

}

func (c *client) getVirtualMachineInternal(name string) (*internal.VirtualMachineInternal, error) {
	vmsint, err := c.store.List()
	if err != nil {
		return nil, err
	}
	if *vmsint == nil || len(*vmsint) == 0 {
		return nil, errors.NotFound
	}

	for _, val := range *vmsint {
		vmint := val.(*internal.VirtualMachineInternal)
		if vmint.Vm.Name == name {
			if hcsshim.HasVirtualMachine(vmint.Id) {
				return vmint, nil
			} else {
				// Store is out of sync with hcs
				c.store.Delete(vmint.Id)
				continue
			}
		}
	}
	return nil, errors.NotFound

}

func (c *client) waitForSSH(vm *pb.VirtualMachine) (string, error) {
	if len(vm.Network.Interfaces) == 0 {
		return "", nil
	}
	vnicName := vm.Network.Interfaces[0].NetworkInterfaceId
	ip, err := virtualnetworkinterface.WaitForIPAddress(c.vnicprovider, vnicName)
	if err != nil {
		return "", err
	}
	if len(ip) == 0 {
		return "", fmt.Errorf("IPAddress not available for the Vm [%s][%s]", vm.Name, vnicName)
	}

	// Wait for SSH connectivity
	for i := 0; i < 10; i++ {
		err := exec.Command("scp", "-o", "StrictHostKeyChecking=no", fmt.Sprintf("%s@%s:~/.profile", vm.Os.Administrator.Username, ip), "tmp").Run()
		if err == nil {
			break
		}
		log.Infof("[VirtualMachine][waitForSSH] [%s][%v]", ip, err)
		time.Sleep(6 * time.Second)
	}

	return ip, nil
}

// TODO: Remove this workaround
func (c *client) addSSHKeys(vm *pb.VirtualMachine) error {
	ip, err := c.waitForSSH(vm)
	if err != nil {
		return err
	}

	// Upload keys to the Vms
	log.Infof("[VirtualMachine][addSSHKeys] Uploading id_rsa_test.pub to " + ip)
	err = exec.Command("scp", "-o", "StrictHostKeyChecking=no", "id_rsa_test.pub", fmt.Sprintf("%s@%s:~/.ssh/id_rsa.pub", vm.Os.Administrator.Username, ip)).Run()
	if err != nil {
		return err
	}

	log.Infof("[VirtualMachine][addSSHKeys] Uploading id_rsa_test to " + ip)
	err = exec.Command("scp", "-o", "StrictHostKeyChecking=no", "id_rsa_test", fmt.Sprintf("%s@%s:~/.ssh/id_rsa", vm.Os.Administrator.Username, ip)).Run()
	if err != nil {
		return err
	}
	log.Infof("[VirtualMachine][addSSHKeys] Fixing Persmission of private key on " + ip)
	err = ssh.ExecuteCommand(vm.Os.Administrator.Username, ip, "id_rsa_test", "chmod 600 ~/.ssh/id_rsa")
	if err != nil {
		return err
	}

	return nil
}
