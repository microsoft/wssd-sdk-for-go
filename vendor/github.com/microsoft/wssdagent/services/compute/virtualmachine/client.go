// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package virtualmachine

import (
	"context"
	"github.com/microsoft/wssdagent/pkg/apis/config"
	"github.com/microsoft/wssdagent/pkg/cloudinit"
	"github.com/microsoft/wssdagent/pkg/errors"
	"github.com/microsoft/wssdagent/pkg/guid"
	"github.com/microsoft/wssdagent/pkg/marshal"
	"github.com/microsoft/wssdagent/pkg/store"
	"github.com/microsoft/wssdagent/pkg/trace"
	pb "github.com/microsoft/wssdagent/rpc/compute"
	networkpb "github.com/microsoft/wssdagent/rpc/network"
	storagepb "github.com/microsoft/wssdagent/rpc/storage"
	"github.com/microsoft/wssdagent/services/compute/virtualmachine/internal"
	"reflect"
	"strings"
	"sync"

	"github.com/microsoft/wssdagent/services/compute/virtualmachine/hcs"
	"github.com/microsoft/wssdagent/services/network/virtualnetworkinterface"
	"github.com/microsoft/wssdagent/services/storage/virtualharddisk"
)

const (
	HCSSpec  = "hcs"
	VMMSSpec = "vmms"
)

type Service interface {
	CreateVirtualMachine(*internal.VirtualMachineInternal, *networkpb.VirtualNetworkInterface, *storagepb.VirtualHardDisk) error
	CleanupVirtualMachine(*internal.VirtualMachineInternal) error
	HasVirtualMachine(*internal.VirtualMachineInternal) bool
}

type Client struct {
	internal     Service
	store        *store.ConfigStore
	config       *config.ChildAgentConfiguration
	vnicprovider *virtualnetworkinterface.VirtualNetworkInterfaceProvider
	vhdprovider  *virtualharddisk.VirtualHardDiskProvider
	mux          sync.Mutex
}

func NewClient() *Client {
	cConfig := config.GetChildAgentConfiguration("VirtualMachine")
	c := &Client{
		store:        store.NewConfigStore(cConfig.DataStorePath, reflect.TypeOf(internal.VirtualMachineInternal{})),
		config:       cConfig,
		vnicprovider: (virtualnetworkinterface.GetVirtualNetworkInterfaceProvider()),
		vhdprovider:  (virtualharddisk.GetVirtualHardDiskProvider()),
	}
	switch cConfig.ProviderSpec {
	case VMMSSpec:
	case HCSSpec:
	default:
		c.internal = hcs.NewClient()
	}
	return c
}

func (c *Client) newVirtualMachine(vm *pb.VirtualMachine) *internal.VirtualMachineInternal {
	return internal.NewVirtualMachineInternal(guid.NewGuid(), c.config.DataStorePath, vm)
}

// Create or Update the specified virtual compute(s)
func (c *Client) Create(ctx context.Context, vmDef *pb.VirtualMachine) (newvm *pb.VirtualMachine, err error) {
	ctx, span := trace.NewSpan(ctx, "VirtualMachine", "Create", marshal.ToString(vmDef))
	defer span.End(err)

	err = c.Validate(ctx, vmDef)
	if err != nil {
		return
	}
	vminternal := c.newVirtualMachine(vmDef)

	// 1. Generate bootstrap data
	switch vmDef.Os.GetOstype() {
	case pb.OperatingSystemType_WINDOWS:
		// TODO: 1.a Generate Answer File
		c.generateAnswerFile(vmDef, vminternal)
	case pb.OperatingSystemType_LINUX:
		// 1.b Generate Cloud Init Configuration
		c.generateCloudInitConfiguration(vmDef, vminternal)
	default:
		err = errors.Wrap(errors.NotSupported, "Unsupported Operating System")
		return
	}

	// 2. Get VHD
	vhd, err := c.vhdprovider.GetVirtualHardDisk(ctx, vmDef.Storage.Osdisk.Diskname)
	if err != nil {
		return
	}

	// 3. Get Vnic info
	var vmnic *networkpb.VirtualNetworkInterface
	if vmDef.Network != nil && len(vmDef.Network.Interfaces) > 0 {
		vnicName := vmDef.Network.Interfaces[0].NetworkInterfaceName
		vmnic, err = c.vnicprovider.GetVirtualNetworkInterfaceByName(ctx, vnicName)
		if err != nil {
			return
		}
	}

	err = c.internal.CreateVirtualMachine(vminternal, vmnic, vhd)
	if err != nil {
		return
	}
	newvm = vminternal.Entity

	err = c.store.Add(vminternal.Id, vminternal)
	if err != nil {
		return
	}

	return

}

// Get all/selected HCS virtual compute(s)
func (c *Client) Get(ctx context.Context, computeDef *pb.VirtualMachine) (vms []*pb.VirtualMachine, err error) {
	ctx, span := trace.NewSpan(ctx, "VirtualMachine", "Get", marshal.ToString(computeDef))
	defer span.End(err)

	c.mux.Lock()
	defer c.mux.Unlock()

	vmName := ""
	if computeDef != nil {
		vmName = computeDef.Name
	}

	vmsint, err := c.store.ListFilter("Name", vmName)
	if err != nil {
		return
	}

	for _, val := range *vmsint {
		vmint := val.(*internal.VirtualMachineInternal)
		vms = append(vms, vmint.Entity)
	}

	return
}

// Delete the specified virtual compute(s)
func (c *Client) Delete(ctx context.Context, computeDef *pb.VirtualMachine) (err error) {
	ctx, span := trace.NewSpan(ctx, "VirtualMachine", "Delete", marshal.ToString(computeDef))
	defer span.End(err)

	c.mux.Lock()
	defer c.mux.Unlock()

	vminternal, err := c.getVirtualMachineInternal(computeDef.Name)
	if err != nil {
		return
	}

	err = c.internal.CleanupVirtualMachine(vminternal)
	if err != nil {
		return
	}

	err = c.store.Delete(vminternal.Id)
	return
}

func (c *Client) Validate(ctx context.Context, computeDef *pb.VirtualMachine) (err error) {
	ctx, span := trace.NewSpan(ctx, "VirtualMachine", "Validate", marshal.ToString(computeDef))
	defer span.End(err)

	err = nil

	if computeDef == nil {
		err = errors.Wrapf(errors.InvalidInput, "Input Virtual Machine definition is nil")
		return
	}

	_, err = c.getVirtualMachineInternal(computeDef.Name)
	if err != nil && err == errors.NotFound {
		err = nil
	} else {
		err = errors.AlreadyExists
	}

	if computeDef.Storage == nil || computeDef.Storage.Osdisk == nil {
		err = errors.Wrapf(errors.InvalidInput, "Virtual Machine doesn't have storage definition")
		return
	}
	if computeDef.Os == nil {
		err = errors.Wrapf(errors.InvalidInput, "Virtual Machine doesn't have Os definition")
		return
	}
	if computeDef.Network == nil {
		//err = errors.Wrapf(errors.InvalidInput, "Virtual Machine doesn't have storage definition")
		//return
	}

	return
}

func (c *Client) getVirtualMachineInternal(name string) (*internal.VirtualMachineInternal, error) {
	vmsint, err := c.store.ListFilter("Name", name)
	if err != nil {
		return nil, err
	}
	if *vmsint == nil || len(*vmsint) == 0 {
		return nil, errors.NotFound
	}
	for _, val := range *vmsint {
		vmint := val.(*internal.VirtualMachineInternal)
		if vmint.Name == name {
			if c.internal.HasVirtualMachine(vmint) {
				return vmint, nil
			} else {
				// Store is out of sync with hcs
				c.store.Delete(vmint.Id)
				continue
			}
		}
	}

	return (*vmsint)[0].(*internal.VirtualMachineInternal), nil
}

func (c *Client) generateAnswerFile(vm *pb.VirtualMachine, vminternal *internal.VirtualMachineInternal) error {
	return nil
}

func (c *Client) generateCloudInitConfiguration(vm *pb.VirtualMachine, vminternal *internal.VirtualMachineInternal) error {
	// Initialize Cloud init user data
	hostname := vm.Os.ComputerName // TODO: Use the same logic as Azure VM hostnames
	if strings.Contains(hostname, "_") {
		return errors.Wrapf(errors.InvalidInput, "linux hostname cannot contain an underscore (_)")
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
		return err
	}
	return nil

}
