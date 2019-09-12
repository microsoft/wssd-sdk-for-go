// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package hcs

import (
	"encoding/json"
	"fmt"
	log "k8s.io/klog"
	"reflect"

	"github.com/microsoft/wssdagent/common"
	"github.com/microsoft/wssdagent/pkg/wssdagent/apis/config"
	"github.com/microsoft/wssdagent/pkg/wssdagent/errors"
	"github.com/microsoft/wssdagent/pkg/wssdagent/store"
	pb "github.com/microsoft/wssdagent/rpc/compute"
	"github.com/microsoft/wssdagent/services/compute/virtualmachine"
	"github.com/microsoft/wssdagent/services/compute/virtualmachinescaleset/internal"
	"github.com/microsoft/wssdagent/services/network/virtualnetwork"
	"github.com/microsoft/wssdagent/services/network/virtualnetworkinterface"
	"github.com/microsoft/wssdagent/services/storage/virtualharddisk"
)

type client struct {
	vnetprovider virtualnetwork.VirtualNetworkProvider
	vnicprovider virtualnetworkinterface.VirtualNetworkInterfaceProvider
	vmprovider   virtualmachine.VirtualMachineProvider
	vhdprovider  virtualharddisk.VirtualHardDiskProvider
	store        *store.ConfigStore
	config       *config.ChildAgentConfiguration
}

func newClient() *client {
	cConfig := config.GetChildAgentConfiguration("VirtualMachineScaleSet")
	return &client{
		vmprovider:   (virtualmachine.GetVirtualMachineProvider(virtualmachine.HCSSpec)),
		vnicprovider: (virtualnetworkinterface.GetVirtualNetworkInterfaceProvider(virtualnetworkinterface.HCNSpec)),
		vnetprovider: (virtualnetwork.GetVirtualNetworkProvider(virtualnetwork.HCNSpec)),
		vhdprovider:  (virtualharddisk.GetVirtualHardDiskProvider(virtualharddisk.HCSSpec)),
		store:        store.NewConfigStore(cConfig.DataStorePath, reflect.TypeOf(internal.VirtualMachineScaleSetInternal{})),
		config:       cConfig,
	}
}

func (c *client) newVirtualMachineScaleSet(id string) *internal.VirtualMachineScaleSetInternal {
	return internal.NewVirtualMachineScaleSetInternal(id, c.config.DataStorePath)
}

// Create or Update a Virtual Machine Scale Set
func (c *client) Create(vmss *pb.VirtualMachineScaleSet) (*pb.VirtualMachineScaleSet, error) {
	vmssinternal, err := c.getVirtualMachineScaleSetInternal(vmss.Name)
	if err == nil {
		// The Vmss already exists. Look to update
		return c.update(vmssinternal.Vmss, vmss)
	}

	// 1.
	if len(vmss.Id) == 0 {
		vmss.Id = common.NewGuid()
	}
	vmssinternal = c.newVirtualMachineScaleSet(vmss.Id)

	// 2. Create replica Vms
	newvms, err := c.createReplicaVirtualMachines(1, vmss)
	if err != nil {
		c.cleanupReplicaVirtualMachines(1, vmss)
		return nil, err
	}

	// FIXME: Cleanup on failure

	// Assign the newly created virtual machines
	vmss.VirtualMachineSystems = newvms
	log.Infof("[VMSS][VMs][Create] [%v]", vmss)

	vmssinternal.Vmss = vmss
	// 3. Store the Vmss
	c.store.Add(vmss.Id, vmssinternal)
	return vmss, nil
}

// Update a Virtual Machine Scale Set
func (c *client) update(vmsscurrent, vmssnew *pb.VirtualMachineScaleSet) (*pb.VirtualMachineScaleSet, error) {
	// For now we only support replica count update
	if vmssnew.Sku.Capacity == vmsscurrent.Sku.Capacity {
		return vmsscurrent, nil
	}

	// 2. Create replica Vms
	newvms, err := c.createReplicaVirtualMachines(int(vmsscurrent.Sku.Capacity+1), vmssnew)
	if err != nil {
		c.cleanupReplicaVirtualMachines(int(vmsscurrent.Sku.Capacity+1), vmssnew)
		return nil, err
	}

	// FIXME: Cleanup on failure

	// Assign the newly created virtual machines
	vmsscurrent.VirtualMachineSystems = append(vmsscurrent.VirtualMachineSystems, newvms...)
	log.Infof("[VMSS][VMs][Update] [%v]", vmsscurrent)

	return vmsscurrent, nil
}

// Get a Virtual Machine scale set specified
func (c *client) Get(vmss *pb.VirtualMachineScaleSet) ([]*pb.VirtualMachineScaleSet, error) {
	vmssName := ""
	if vmss != nil {
		vmssName = vmss.Name
	}

	vmssval := []*pb.VirtualMachineScaleSet{}

	if len(vmssName) == 0 {
		vmssvals, err := c.store.List()
		if err != nil {
			return nil, err
		}
		// get everything
		if *vmssvals == nil || len(*vmssvals) == 0 {
			return nil, nil
		}

		for _, val := range *vmssvals {
			vmssint := val.(*internal.VirtualMachineScaleSetInternal)
			// Validate the Vms
			c.vmprovider.Get(vmssint.Vmss.GetVirtualMachineSystems())
			vmssval = append(vmssval, vmssint.Vmss)
		}

	} else {
		vmssget, err := c.getVirtualMachineScaleSetInternal(vmssName)
		if err != nil {
			return vmssval, err
		}
		vmssval = append(vmssval, vmssget.Vmss)
	}
	return vmssval, nil
}

// Delete a Virtual Machine scale set specified
func (c *client) Delete(vmss *pb.VirtualMachineScaleSet) error {
	vmssInt, err := c.getVirtualMachineScaleSetInternal(vmss.Name)
	if err != nil {
		return err
	}

	currentVmss := vmssInt.Vmss
	vnics := []string{}
	// 1. Get Endpoint info for each of the Vms
	for _, vm := range currentVmss.GetVirtualMachineSystems() {
		for _, vnic := range vm.Network.Interfaces {
			vnics = append(vnics, vnic.NetworkInterfaceId)
		}
	}

	// 2. Delete the Vms
	log.Infof("[VMSS][VMs][Delete] [%v]", currentVmss)
	err = c.vmprovider.Delete(currentVmss.GetVirtualMachineSystems())
	if err != nil && err != errors.NotFound {
		return err
	}

	// 3. Delete endpoints created for each of the Vms
	virtualnetworkinterface.DeleteVirtualNetworkInterface(c.vnicprovider, vnics)

	// 4. Delete Vhds for each replica vm

	return c.store.Delete(currentVmss.Id)
}

func (c *client) createReplicaVirtualMachines(startVmCount int, vmss *pb.VirtualMachineScaleSet) ([]*pb.VirtualMachine, error) {
	vms := []*pb.VirtualMachine{}
	vm := vmss.Virtualmachineprofile
	vmBytes, _ := json.Marshal(vm)
	ret := []*pb.VirtualMachine{}
	for vmcount := startVmCount; vmcount <= int(vmss.Sku.Capacity); vmcount++ {
		replicaVm := &pb.VirtualMachine{}
		json.Unmarshal(vmBytes, replicaVm)
		replicaVm.Name = fmt.Sprintf("vm_%s_%s_%d", vmss.Name, vmss.Virtualmachineprofile.Vmprefix, vmcount)
		replicaVm.Id = common.NewGuid()

		// 2.a Create endpoints for each replica Vm
		if err := c.addNetworkConfiguration(vmss, replicaVm, vmcount); err != nil {
			return ret, err
		}

		// 2.b Create Vhds for each replica Vm
		if err := c.addStorageConfiguration(vmss, replicaVm); err != nil {
			return ret, err
		}

		// 2.c Add OsConfiguration
		if err := c.addOSConfiguration(vmss, replicaVm); err != nil {
			return ret, err
		}

		vms = append(vms, replicaVm)
	}
	return c.vmprovider.CreateOrUpdate(vms)
}

func (c *client) cleanupReplicaVirtualMachines(startVmCount int, vmss *pb.VirtualMachineScaleSet) error {
	vms := []*pb.VirtualMachine{}
	vm := vmss.Virtualmachineprofile
	vmBytes, _ := json.Marshal(vm)
	for vmcount := startVmCount; vmcount <= int(vmss.Sku.Capacity); vmcount++ {
		replicaVm := &pb.VirtualMachine{}
		json.Unmarshal(vmBytes, replicaVm)
		replicaVm.Name = fmt.Sprintf("vm_%s_%s_%d", vmss.Name, vmss.Virtualmachineprofile.Vmprefix, vmcount)
		replicaVm.Id = common.NewGuid()

		// 2.a cleanup endpoints for each replica Vm
		if err := c.cleanupNetworkConfiguration(replicaVm); err != nil {
			// TODO: Log error during cleanup
		}

		// 2.b cleanup Vhds for each replica Vm
		// return c.cleanupStorageConfiguration(vmss, replicaVm);
		vms = append(vms, replicaVm)
	}
	return c.vmprovider.Delete(vms)
}

func (c *client) getVirtualMachineScaleSetInternal(name string) (*internal.VirtualMachineScaleSetInternal, error) {
	vmsint, err := c.store.List()
	if err != nil {
		return nil, err
	}
	if *vmsint == nil || len(*vmsint) == 0 {
		return nil, errors.NotFound
	}

	for _, val := range *vmsint {
		vmint := val.(*internal.VirtualMachineScaleSetInternal)
		if vmint.Vmss.Name == name {
			return vmint, nil
		}
	}
	return nil, errors.NotFound

}

func (c *client) addNetworkConfiguration(vmss *pb.VirtualMachineScaleSet, replicaVm *pb.VirtualMachine, vmcount int) error {
	if vmss.Virtualmachineprofile.GetNetwork() == nil {
		log.Infof("[VMSS][VMs][Create][addNetworkConfiguration] No network interface requested. VM will not have any connectivity")
		return nil
	}

	replicaVm.Network.Interfaces = []*pb.NetworkInterface{}
	vnicCount := 1
	for _, vnic := range vmss.Virtualmachineprofile.Network.Interfaces {
		vnicName := fmt.Sprintf("vnic_%s_%s_%d_%s_%d", vmss.Name, vmss.Virtualmachineprofile.Vmprefix, vmcount, vnic.Networkname, vnicCount)
		vnicCount++
		if err := virtualnetworkinterface.CreateVirtualNetworkInterface(c.vnicprovider, vnicName, vnic.Networkname); err != nil {
			return err
		}
		replicaVm.Network.Interfaces = append(replicaVm.Network.Interfaces, &pb.NetworkInterface{NetworkInterfaceId: vnicName})
	}
	return nil
}

func (c *client) cleanupNetworkConfiguration(replicaVm *pb.VirtualMachine) error {
	if replicaVm.Network.Interfaces == nil {
		return nil
	}

	vnics := []string{}
	for _, vnic := range replicaVm.Network.Interfaces {
		vnics = append(vnics, vnic.NetworkInterfaceId)
	}

	return virtualnetworkinterface.DeleteVirtualNetworkInterface(c.vnicprovider, vnics)
}

func (c *client) addStorageConfiguration(vmss *pb.VirtualMachineScaleSet, replicaVm *pb.VirtualMachine) error {
	if vmss.Virtualmachineprofile.Storage == nil {
		return errors.Wrap(errors.InvalidConfiguration, "Storage configuration missing")
	}

	// FIXME: Use create Vhds when the Apis are ready
	// For now hardcoding it to whatever is in the input schema
	// virtualharddisk.CreateVirtualHardDisk(c.vhdprovider, "", vmss.Virtualmachineprofile.Storage.Osdisk.Diskid, replicaVm.Storage.Osdisk.Diskid)
	// With this only 1 VM can be supported for Vmss

	replicaVm.Storage.Osdisk.Diskid = vmss.Virtualmachineprofile.Storage.Osdisk.Diskid
	return nil
}

func (c *client) addOSConfiguration(vmss *pb.VirtualMachineScaleSet, replicaVm *pb.VirtualMachine) error {
	if vmss.Virtualmachineprofile.Os == nil {
		return errors.Wrap(errors.InvalidConfiguration, "Operating System configuration missing")
	}
	replicaVm.Os.ComputerName = vmss.Virtualmachineprofile.Os.ComputerName
	return c.addUserConfiguration(vmss, replicaVm)
}

func (c *client) addUserConfiguration(vmss *pb.VirtualMachineScaleSet, replicaVm *pb.VirtualMachine) error {
	if vmss.Virtualmachineprofile.Os.Administrator == nil {
		return errors.Wrap(errors.InvalidConfiguration, "Administrator configuration missing in Virtual Machine Profile")
	}

	replicaVm.Os.Administrator = vmss.Virtualmachineprofile.Os.Administrator
	replicaVm.Os.Users = vmss.Virtualmachineprofile.Os.Users

	return nil
}
