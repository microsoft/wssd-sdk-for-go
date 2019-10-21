// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package virtualmachinescaleset

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/microsoft/wssdagent/pkg/apis/config"
	"github.com/microsoft/wssdagent/pkg/errors"
	"github.com/microsoft/wssdagent/pkg/guid"
	"github.com/microsoft/wssdagent/pkg/marshal"
	"github.com/microsoft/wssdagent/pkg/store"
	"github.com/microsoft/wssdagent/pkg/trace"
	pb "github.com/microsoft/wssdagent/rpc/compute"
	"github.com/microsoft/wssdagent/services/compute/virtualmachinescaleset/internal"
	"reflect"
	"sync"

	"github.com/microsoft/wssdagent/services/compute/virtualmachine"
	"github.com/microsoft/wssdagent/services/network/virtualnetwork"
	"github.com/microsoft/wssdagent/services/network/virtualnetworkinterface"
	"github.com/microsoft/wssdagent/services/storage/virtualharddisk"
)

type Client struct {
	store        *store.ConfigStore
	config       *config.ChildAgentConfiguration
	vnicprovider *virtualnetworkinterface.VirtualNetworkInterfaceProvider
	vhdprovider  *virtualharddisk.VirtualHardDiskProvider
	vnetprovider *virtualnetwork.VirtualNetworkProvider
	vmprovider   *virtualmachine.VirtualMachineProvider

	mux sync.Mutex
}

func NewClient() *Client {
	cConfig := config.GetChildAgentConfiguration("VirtualMachineScaleSet")
	return &Client{
		store:        store.NewConfigStore(cConfig.DataStorePath, reflect.TypeOf(internal.VirtualMachineScaleSetInternal{})),
		config:       cConfig,
		vnicprovider: (virtualnetworkinterface.GetVirtualNetworkInterfaceProvider()),
		vhdprovider:  (virtualharddisk.GetVirtualHardDiskProvider()),
		vmprovider:   (virtualmachine.GetVirtualMachineProvider()),
		vnetprovider: (virtualnetwork.GetVirtualNetworkProvider()),
	}
}

func (c *Client) newVirtualMachineScaleSet(vmss *pb.VirtualMachineScaleSet) *internal.VirtualMachineScaleSetInternal {
	return internal.NewVirtualMachineScaleSetInternal(guid.NewGuid(), c.config.DataStorePath, vmss)
}

// Create or Update the specified virtual compute(s)
func (c *Client) Create(ctx context.Context, vmssDef *pb.VirtualMachineScaleSet) (newvmss *pb.VirtualMachineScaleSet, err error) {
	ctx, span := trace.NewSpan(ctx, "VirtualMachineScaleSet", "Create", marshal.ToString(vmssDef))
	defer span.End(err)

	err = c.Validate(ctx, vmssDef)
	if err != nil && err != errors.AlreadyExists {
		return
	}

	vmssinternal, err := c.getVirtualMachineScaleSetInternal(vmssDef.Name)
	if err == nil {
		// The Vmss already exists. Look to update
		return c.update(ctx, vmssinternal.Entity, vmssDef)
	}

	vmssinternal = c.newVirtualMachineScaleSet(vmssDef)

	newvms, err := c.createReplicaVirtualMachines(ctx, 1, vmssDef)
	if err != nil {
		c.cleanupReplicaVirtualMachines(ctx, 1, vmssDef)
		return nil, err
	}

	newvmss = vmssinternal.Entity
	newvmss.VirtualMachineSystems = append(newvmss.VirtualMachineSystems, newvms...)

	err = c.store.Add(vmssinternal.Id, vmssinternal)

	return

}

// Get all/selected HCS virtual compute(s)
func (c *Client) Get(ctx context.Context, computeDef *pb.VirtualMachineScaleSet) (vmsss []*pb.VirtualMachineScaleSet, err error) {
	ctx, span := trace.NewSpan(ctx, "VirtualMachineScaleSet", "Get", marshal.ToString(computeDef))
	defer span.End(err)

	c.mux.Lock()
	defer c.mux.Unlock()

	vmssName := ""
	if computeDef != nil {
		vmssName = computeDef.Name
	}

	vmsssint, err := c.store.ListFilter("Name", vmssName)
	if err != nil {
		return
	}

	for _, val := range *vmsssint {
		vmssint := val.(*internal.VirtualMachineScaleSetInternal)
		vmsss = append(vmsss, vmssint.Entity)
	}

	return
}

// Delete the specified virtual compute(s)
func (c *Client) Delete(ctx context.Context, computeDef *pb.VirtualMachineScaleSet) (err error) {
	ctx, span := trace.NewSpan(ctx, "VirtualMachineScaleSet", "Delete", marshal.ToString(computeDef))
	defer span.End(err)

	c.mux.Lock()
	defer c.mux.Unlock()

	vmssinternal, err := c.getVirtualMachineScaleSetInternal(computeDef.Name)
	if err != nil {
		return
	}

	currentVmss := vmssinternal.Entity

	err = c.cleanupReplicaVirtualMachines(ctx, 1, currentVmss)
	if err != nil {
	}

	err = c.store.Delete(vmssinternal.Id)
	return
}

func (c *Client) Validate(ctx context.Context, vmssDef *pb.VirtualMachineScaleSet) (err error) {
	ctx, span := trace.NewSpan(ctx, "VirtualMachineScaleSet", "Validate", marshal.ToString(vmssDef))
	defer span.End(err)

	err = nil

	if vmssDef == nil {
		err = errors.Wrapf(errors.InvalidInput, "Input Virtual Machine definition is nil")
		return
	}

	_, err = c.getVirtualMachineScaleSetInternal(vmssDef.Name)
	if err != nil && err == errors.NotFound {
		err = nil
	} else {
		err = errors.AlreadyExists
	}

	return
}

func (c *Client) getVirtualMachineScaleSetInternal(name string) (*internal.VirtualMachineScaleSetInternal, error) {
	vmsssint, err := c.store.ListFilter("Name", name)
	if err != nil {
		return nil, err
	}
	if *vmsssint == nil || len(*vmsssint) == 0 {
		return nil, errors.NotFound
	}

	return (*vmsssint)[0].(*internal.VirtualMachineScaleSetInternal), nil
}

// Update a Virtual Machine Scale Set
func (c *Client) update(ctx context.Context, vmsscurrent, vmssnew *pb.VirtualMachineScaleSet) (*pb.VirtualMachineScaleSet, error) {
	// For now we only support replica count update
	if vmssnew.Sku.Capacity == vmsscurrent.Sku.Capacity {
		return vmsscurrent, nil
	}

	// 2. Create replica Vms
	newvms, err := c.createReplicaVirtualMachines(ctx, int(vmsscurrent.Sku.Capacity+1), vmssnew)
	if err != nil {
		c.cleanupReplicaVirtualMachines(ctx, int(vmsscurrent.Sku.Capacity+1), vmssnew)
		return nil, err
	}

	// FIXME: Cleanup on failure

	// Assign the newly created virtual machines
	vmsscurrent.VirtualMachineSystems = append(vmsscurrent.VirtualMachineSystems, newvms...)

	return vmsscurrent, nil
}

func (c *Client) createReplicaVirtualMachines(ctx context.Context, startVmCount int, vmss *pb.VirtualMachineScaleSet) ([]*pb.VirtualMachine, error) {
	vms := []*pb.VirtualMachine{}
	vm := vmss.Virtualmachineprofile
	vmBytes, _ := json.Marshal(vm)
	ret := []*pb.VirtualMachine{}
	for vmcount := startVmCount; vmcount <= int(vmss.Sku.Capacity); vmcount++ {
		replicaVm := &pb.VirtualMachine{}
		json.Unmarshal(vmBytes, replicaVm)
		replicaVm.Name = fmt.Sprintf("vm_%s_%s_%d", vmss.Name, vmss.Virtualmachineprofile.Vmprefix, vmcount)

		// 2.a Create endpoints for each replica Vm
		if err := c.addNetworkConfiguration(ctx, vmss, replicaVm, vmcount); err != nil {
			return ret, errors.Wrapf(err, "Adding Network Configuration Failed")
		}

		// 2.b Create Vhds for each replica Vm
		if err := c.addStorageConfiguration(ctx, vmss, replicaVm); err != nil {
			return ret, errors.Wrapf(err, "Adding Storag Configuration Failed")
		}

		// 2.c Add OsConfiguration
		if err := c.addOSConfiguration(vmss, replicaVm, vmcount); err != nil {
			return ret, err
		}

		vms = append(vms, replicaVm)
	}
	return c.vmprovider.CreateOrUpdate(ctx, vms)
}

func (c *Client) cleanupReplicaVirtualMachines(ctx context.Context, startVmCount int, vmss *pb.VirtualMachineScaleSet) error {
	vm := vmss.Virtualmachineprofile
	vmBytes, _ := json.Marshal(vm)
	for vmcount := startVmCount; vmcount <= int(vmss.Sku.Capacity); vmcount++ {
		replicaVm := &pb.VirtualMachine{}
		json.Unmarshal(vmBytes, replicaVm)
		replicaVm.Name = fmt.Sprintf("vm_%s_%s_%d", vmss.Name, vmss.Virtualmachineprofile.Vmprefix, vmcount)

		// 1. First Delete the Vm
		if err := c.vmprovider.Delete(ctx, []*pb.VirtualMachine{replicaVm}); err != nil {
			// If VM Deletion fails, the storage deletion would fail, so fail this call
			return err
		}

		// 2.a cleanup endpoints for each replica Vm
		if err := c.cleanupNetworkConfiguration(ctx, replicaVm); err != nil {
			// TODO: Log error during cleanup
		}

		// 2.b cleanup Vhds for each replica Vm
		if err := c.cleanupStorageConfiguration(ctx, vmss, replicaVm); err != nil {
		}
	}
	return nil

}

func (c *Client) addNetworkConfiguration(ctx context.Context, vmss *pb.VirtualMachineScaleSet, replicaVm *pb.VirtualMachine, vmcount int) error {
	if vmss.Virtualmachineprofile.GetNetwork() == nil {
		// log.Infof("[VMSS][VMs][Create][addNetworkConfiguration] No network interface requested. VM will not have any connectivity")
		return nil
	}

	replicaVm.Network.Interfaces = []*pb.NetworkInterface{}
	vnicCount := 1
	for _, vnic := range vmss.Virtualmachineprofile.Network.Interfaces {
		if len(vnic.Networkname) == 0 {
			return errors.Wrapf(errors.InvalidInput, "Network Name is missing in the vmss definition")
		}
		vnicName := fmt.Sprintf("vnic_%s_%s_%d_%s_%d", vmss.Name, vmss.Virtualmachineprofile.Vmprefix, vmcount, vnic.Networkname, vnicCount)
		vnicCount++
		if err := c.vnicprovider.CreateVirtualNetworkInterface(ctx, vnicName, vnic.Networkname); err != nil {
			return err
		}
		replicaVm.Network.Interfaces = append(replicaVm.Network.Interfaces, &pb.NetworkInterface{NetworkInterfaceName: vnicName})
	}
	return nil
}

func (c *Client) cleanupNetworkConfiguration(ctx context.Context, replicaVm *pb.VirtualMachine) error {
	if replicaVm.Network.Interfaces == nil {
		return nil
	}

	vnics := []string{}
	for _, vnic := range replicaVm.Network.Interfaces {
		vnics = append(vnics, vnic.NetworkInterfaceName)
	}

	return c.vnicprovider.DeleteVirtualNetworkInterface(ctx, vnics)
}

func (c *Client) addStorageConfiguration(ctx context.Context, vmss *pb.VirtualMachineScaleSet, replicaVm *pb.VirtualMachine) error {
	if vmss.Virtualmachineprofile.Storage == nil {
		return errors.Wrap(errors.InvalidConfiguration, "Storage configuration missing")
	}

	vhdName := fmt.Sprintf("vhd_%s", replicaVm.Name)

	_, err := c.vhdprovider.CloneVirtualHardDisk(ctx, vmss.Virtualmachineprofile.Storage.Osdisk.Diskname, vhdName)
	if err != nil {
		return err
	}
	replicaVm.Storage.Osdisk.Diskname = vhdName
	return nil
}

func (c *Client) cleanupStorageConfiguration(ctx context.Context, vmss *pb.VirtualMachineScaleSet, replicaVm *pb.VirtualMachine) error {
	vhdName := fmt.Sprintf("vhd_%s", replicaVm.Name)
	err := c.vhdprovider.DeleteVirtualHardDisk(ctx, vhdName)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) addOSConfiguration(vmss *pb.VirtualMachineScaleSet, replicaVm *pb.VirtualMachine, vmcount int) error {
	if vmss.Virtualmachineprofile.Os == nil {
		return errors.Wrap(errors.InvalidConfiguration, "Operating System configuration missing")
	}
	replicaVm.Os.ComputerName = fmt.Sprintf("%s-%d", vmss.Virtualmachineprofile.Os.ComputerName, vmcount)
	return c.addUserConfiguration(vmss, replicaVm)
}

func (c *Client) addUserConfiguration(vmss *pb.VirtualMachineScaleSet, replicaVm *pb.VirtualMachine) error {
	if vmss.Virtualmachineprofile.Os.Administrator == nil {
		return errors.Wrap(errors.InvalidConfiguration, "Administrator configuration missing in Virtual Machine Profile")
	}

	replicaVm.Os.Administrator = vmss.Virtualmachineprofile.Os.Administrator
	replicaVm.Os.Users = vmss.Virtualmachineprofile.Os.Users

	return nil
}
