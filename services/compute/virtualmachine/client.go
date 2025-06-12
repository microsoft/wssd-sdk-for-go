// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license

package virtualmachine

import (
	"context"
	"log"

	"github.com/microsoft/moc/pkg/auth"
	"github.com/microsoft/moc/pkg/errors"
	"github.com/microsoft/moc/pkg/marshal"

	"github.com/microsoft/wssd-sdk-for-go/services/compute"
	"github.com/microsoft/wssd-sdk-for-go/services/compute/virtualmachine/internal"
	"github.com/microsoft/wssd-sdk-for-go/services/network"
)

type Service interface {
	Get(context.Context, string, string, string) (*[]compute.VirtualMachine, error)
	CreateOrUpdate(context.Context, string, string, string, *compute.VirtualMachine) (*compute.VirtualMachine, error)
	Delete(context.Context, string, string, string) error
	Hydrate(context.Context, string, string, string, *compute.VirtualMachine) (*compute.VirtualMachine, error)
	Start(context.Context, string, string, string) error
	Stop(context.Context, string, string, string) error
	Pause(context.Context, string, string, string) error
	Save(context.Context, string, string, string) error
	RemoveIsoDisk(context.Context, string, string, string) error
	RepairGuestAgent(context.Context, string, string, string) error
	RunCommand(context.Context, string, string, string, *compute.VirtualMachineRunCommandRequest) (*compute.VirtualMachineRunCommandResponse, error)
	Validate(context.Context, string, string, string) error
}

type VirtualMachineClient struct {
	compute.BaseClient
	internal Service
}

func NewVirtualMachineClient(cloudFQDN string, authorizer auth.Authorizer) (*VirtualMachineClient, error) {
	c, err := internal.NewVirtualMachineClient(cloudFQDN, authorizer)
	if err != nil {
		return nil, err
	}

	return &VirtualMachineClient{internal: c}, nil
}

// Get methods invokes the client Get method
func (c *VirtualMachineClient) Get(ctx context.Context, group, name, id string) (*[]compute.VirtualMachine, error) {
	return c.internal.Get(ctx, group, name, id)
}

func (c *VirtualMachineClient) CreateOrUpdate(ctx context.Context, group, name, id string, compute *compute.VirtualMachine) (*compute.VirtualMachine, error) {
	return c.internal.CreateOrUpdate(ctx, group, name, id, compute)
}

func (c *VirtualMachineClient) Delete(ctx context.Context, group, name, id string) error {
	return c.internal.Delete(ctx, group, name, id)
}

func (c *VirtualMachineClient) Hydrate(ctx context.Context, group, name, id string, compute *compute.VirtualMachine) (*compute.VirtualMachine, error) {
	return c.internal.Hydrate(ctx, group, name, id, compute)
}

func (c *VirtualMachineClient) Start(ctx context.Context, group, name, id string) error {
	return c.internal.Start(ctx, group, name, id)
}

func (c *VirtualMachineClient) Stop(ctx context.Context, group, name, id string) error {
	return c.internal.Stop(ctx, group, name, id)
}

func (c *VirtualMachineClient) Restart(ctx context.Context, group, name, id string) error {
	err := c.internal.Stop(ctx, group, name, id)
	if err != nil {
		return err
	}
	return c.internal.Start(ctx, group, name, id)
}

func (c *VirtualMachineClient) Pause(ctx context.Context, group, name, id string) error {
	return c.internal.Pause(ctx, group, name, id)
}

func (c *VirtualMachineClient) Save(ctx context.Context, group, name, id string) error {
	return c.internal.Save(ctx, group, name, id)
}

func (c *VirtualMachineClient) RemoveIsoDisk(ctx context.Context, group, name, id string) error {
	return c.internal.RemoveIsoDisk(ctx, group, name, id)
}

func (c *VirtualMachineClient) RepairGuestAgent(ctx context.Context, group, name, id string) error {
	return c.internal.RepairGuestAgent(ctx, group, name, id)
}

// Validate methods invokes the validate Get method
func (c *VirtualMachineClient) Validate(ctx context.Context, group, name, id string) error {
	return c.internal.Validate(ctx, group, name, id)
}

func (c *VirtualMachineClient) Resize(ctx context.Context, group, name string, newSize compute.VirtualMachineSizeTypes, newCustomSize *compute.VirtualMachineCustomSize, id string) error {
	return c.ResizeEx(ctx, group, name, newSize, newCustomSize, nil, id)
}

func (c *VirtualMachineClient) ResizeEx(ctx context.Context, group string, name string, newSize compute.VirtualMachineSizeTypes, newCustomSize *compute.VirtualMachineCustomSize, newVirtualMachineGPUs []*compute.VirtualMachineGPU, id string) (err error) {
	vms, err := c.Get(ctx, group, name, id)
	if err != nil {
		return
	}

	if vms == nil || len(*vms) == 0 {
		err = errors.Wrapf(errors.NotFound, "%s", name)
		return
	}

	vm := (*vms)[0]

	if !isDifferentVmSize(vm.HardwareProfile.VMSize, newSize, vm.HardwareProfile.CustomSize, newCustomSize) && !isDifferentGpuList(vm.HardwareProfile.VirtualMachineGPUs, newVirtualMachineGPUs) {
		// Nothing to do
		return
	}

	vm.HardwareProfile.VMSize = newSize
	vm.HardwareProfile.CustomSize = newCustomSize
	vm.HardwareProfile.VirtualMachineGPUs = newVirtualMachineGPUs

	_, err = c.CreateOrUpdate(ctx, group, name, id, &vm)
	return
}

func (c *VirtualMachineClient) DiskAttach(ctx context.Context, group, vmName, diskName, id string) error {
	vms, err := c.Get(ctx, group, vmName, id)
	if err != nil {
		return err
	}
	if vms == nil || len(*vms) == 0 {
		return errors.Wrapf(errors.NotFound, "Unable to find Virtual Machine [%s]", vmName)
	}

	vm := (*vms)[0]
	for _, disk := range *vm.StorageProfile.DataDisks {
		if disk.Vhd.Name == diskName {
			return errors.Wrapf(errors.AlreadyExists, "DataDisk [%s] is already attached to the VM [%s]", diskName, vmName)
		}
	}

	*vm.StorageProfile.DataDisks = append(*vm.StorageProfile.DataDisks, compute.DataDisk{Vhd: &compute.VirtualHardDisk{Name: diskName}})

	_, err = c.CreateOrUpdate(ctx, group, vmName, id, &vm)
	if err != nil {
		return err
	}

	return nil
}

func (c *VirtualMachineClient) DiskDetach(ctx context.Context, group, vmName, diskName, id string) error {
	vms, err := c.Get(ctx, group, vmName, id)
	if err != nil {
		return err
	}
	if vms == nil || len(*vms) == 0 {
		return errors.Wrapf(errors.NotFound, "Unable to find Virtual Machine [%s]", vmName)
	}

	vm := (*vms)[0]
	for i, element := range *vm.StorageProfile.DataDisks {
		if element.Vhd.Name == diskName {
			*vm.StorageProfile.DataDisks = append((*vm.StorageProfile.DataDisks)[:i], (*vm.StorageProfile.DataDisks)[i+1:]...)
			break
		}
	}

	_, err = c.CreateOrUpdate(ctx, group, vmName, id, &vm)
	if err != nil {
		return err
	}
	return nil
}

func (c *VirtualMachineClient) NetworkInterfaceAdd(ctx context.Context, group, vmName, nicName, id string) error {
	vms, err := c.Get(ctx, group, vmName, id)
	if err != nil {
		return err
	}
	if vms == nil || len(*vms) == 0 {
		return errors.Wrapf(errors.NotFound, "Unable to find Virtual Machine [%s]", vmName)
	}

	vm := (*vms)[0]
	for _, nic := range *vm.NetworkProfile.NetworkInterfaces {
		if *nic.Vnic.Name == nicName {
			return errors.Wrapf(errors.AlreadyExists, "NetworkInterface [%s] is already attached to the VM [%s]", nicName, vmName)
		}
	}

	*vm.NetworkProfile.NetworkInterfaces = append(*vm.NetworkProfile.NetworkInterfaces, compute.NetworkInterfaceReference{Vnic: &network.VirtualNetworkInterface{Name: &nicName}})

	_, err = c.CreateOrUpdate(ctx, group, vmName, id, &vm)
	if err != nil {
		return err
	}
	return nil
}

func (c *VirtualMachineClient) NetworkInterfaceRemove(ctx context.Context, group, vmName, nicName, id string) error {
	vms, err := c.Get(ctx, group, vmName, id)
	if err != nil {
		return err
	}
	if vms == nil || len(*vms) == 0 {
		return errors.Wrapf(errors.NotFound, "Unable to find Virtual Machine [%s]", vmName)
	}

	vm := (*vms)[0]
	for i, element := range *vm.NetworkProfile.NetworkInterfaces {
		if *element.Vnic.Name == nicName {
			*vm.NetworkProfile.NetworkInterfaces = append((*vm.NetworkProfile.NetworkInterfaces)[:i], (*vm.NetworkProfile.NetworkInterfaces)[i+1:]...)
			break
		}
	}

	_, err = c.CreateOrUpdate(ctx, group, vmName, id, &vm)
	if err != nil {
		return err
	}
	return nil
}

func (c *VirtualMachineClient) NetworkInterfaceList(ctx context.Context, group, vmName, id string) error {
	vms, err := c.Get(ctx, group, vmName, id)
	if err != nil {
		return err
	}
	if vms == nil || len(*vms) == 0 {
		return errors.Wrapf(errors.NotFound, "Unable to find Virtual Machine [%s]", vmName)
	}

	vm := (*vms)[0]
	for _, element := range *vm.NetworkProfile.NetworkInterfaces {
		log.Printf("%+v\n", marshal.ToString(element))
	}

	return nil
}

func (c *VirtualMachineClient) NetworkInterfaceShow(ctx context.Context, group, vmName, nicName, id string) error {
	vms, err := c.Get(ctx, group, vmName, id)
	if err != nil {
		return err
	}
	if vms == nil || len(*vms) == 0 {
		return errors.Wrapf(errors.NotFound, "Unable to find Virtual Machine [%s]", vmName)
	}

	vm := (*vms)[0]
	for _, nic := range *vm.NetworkProfile.NetworkInterfaces {
		if *nic.Vnic.Name == nicName {
			// TODO - implement detailed show
			log.Printf("%+v\n", marshal.ToString(nic))
			break
		}
	}

	return nil
}

func (c *VirtualMachineClient) RunCommand(ctx context.Context, group, name, id string, request *compute.VirtualMachineRunCommandRequest) (*compute.VirtualMachineRunCommandResponse, error) {
	return c.internal.RunCommand(ctx, group, name, id, request)
}
func isDifferentVmSize(oldSizeType, newSizeType compute.VirtualMachineSizeTypes, oldCustomSize, newCustomSize *compute.VirtualMachineCustomSize) bool {
	if oldSizeType != newSizeType {
		return true
	}

	// same vm size type, check custom size
	// Note: fields in compute.VirtualMachineCustomSize are pointers, deference to compare the values
	switch newSizeType {
	case compute.VirtualMachineSizeTypesCustomNK:
		fallthrough
	case compute.VirtualMachineSizeTypesCustomGpupv:
		if *oldCustomSize.GpuCount != *newCustomSize.GpuCount {
			return true
		}
		fallthrough
	case compute.VirtualMachineSizeTypesCustom:
		if *oldCustomSize.CpuCount != *newCustomSize.CpuCount {
			return true
		}
		if *oldCustomSize.MemoryMB != *newCustomSize.MemoryMB {
			return true
		}
		return false
	default:
		return false
	}
}

func isDifferentGpuList(oldGpuList, newGpuList []*compute.VirtualMachineGPU) bool {
	// simultaneous addtion and removal of GPU is not supported
	return len(oldGpuList) != len(newGpuList)
}
