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
)

type Service interface {
	Get(context.Context, string, string) (*[]compute.VirtualMachine, error)
	CreateOrUpdate(context.Context, string, string, *compute.VirtualMachine) (*compute.VirtualMachine, error)
	Delete(context.Context, string, string) error
	Hydrate(context.Context, string, string, *compute.VirtualMachine) (*compute.VirtualMachine, error)
	Start(context.Context, string, string) error
	Stop(context.Context, string, string) error
	StopGraceful(context.Context, string, string) error
	Pause(context.Context, string, string) error
	Save(context.Context, string, string) error
	RemoveIsoDisk(context.Context, string, string) error
	RepairGuestAgent(context.Context, string, string) error
	RunCommand(context.Context, string, string, *compute.VirtualMachineRunCommandRequest) (*compute.VirtualMachineRunCommandResponse, error)
	Validate(context.Context, string, string) error
	GetHyperVVmId(context.Context, string, string) (*compute.VirtualMachineHyperVVmId, error)
	HasHyperVVm(context.Context, string) (bool, error)
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
func (c *VirtualMachineClient) Get(ctx context.Context, group, name string) (*[]compute.VirtualMachine, error) {
	return c.internal.Get(ctx, group, name)
}

// CreateOrUpdate methods invokes create or update on the client
func (c *VirtualMachineClient) CreateOrUpdate(ctx context.Context, group, name string, compute *compute.VirtualMachine) (*compute.VirtualMachine, error) {
	return c.internal.CreateOrUpdate(ctx, group, name, compute)
}

// Delete methods invokes delete of the compute resource
func (c *VirtualMachineClient) Delete(ctx context.Context, group string, name string) error {
	return c.internal.Delete(ctx, group, name)
}

// Hydrate methods creates MOC representation of the VM resource
func (c *VirtualMachineClient) Hydrate(ctx context.Context, group string, name string, compute *compute.VirtualMachine) (*compute.VirtualMachine, error) {
	return c.internal.Hydrate(ctx, group, name, compute)
}

func (c *VirtualMachineClient) Start(ctx context.Context, group string, name string) (err error) {
	err = c.internal.Start(ctx, group, name)
	return
}

func (c *VirtualMachineClient) Stop(ctx context.Context, group string, name string) (err error) {
	err = c.internal.Stop(ctx, group, name)
	return
}

func (c *VirtualMachineClient) StopGraceful(ctx context.Context, group string, name string) (err error) {
	err = c.internal.StopGraceful(ctx, group, name)
	return
}

func (c *VirtualMachineClient) Restart(ctx context.Context, group string, name string) (err error) {
	err = c.internal.StopGraceful(ctx, group, name)
	if err != nil {
		return
	}
	err = c.internal.Start(ctx, group, name)
	return
}
func (c *VirtualMachineClient) Pause(ctx context.Context, group string, name string) (err error) {
	err = c.internal.Pause(ctx, group, name)
	return
}
func (c *VirtualMachineClient) Save(ctx context.Context, group string, name string) (err error) {
	err = c.internal.Save(ctx, group, name)
	return
}

func (c *VirtualMachineClient) RemoveIsoDisk(ctx context.Context, group string, name string) (err error) {
	err = c.internal.RemoveIsoDisk(ctx, group, name)
	return
}

func (c *VirtualMachineClient) RepairGuestAgent(ctx context.Context, group string, name string) (err error) {
	err = c.internal.RepairGuestAgent(ctx, group, name)
	return
}

// Validate methods invokes the validate Get method
func (c *VirtualMachineClient) Validate(ctx context.Context, group, name string) error {
	return c.internal.Validate(ctx, group, name)
}

func (c *VirtualMachineClient) Resize(ctx context.Context, group string, name string, newSize compute.VirtualMachineSizeTypes, newCustomSize *compute.VirtualMachineCustomSize) (err error) {
	return c.ResizeEx(ctx, group, name, newSize, newCustomSize, nil)
}

func (c *VirtualMachineClient) ResizeEx(ctx context.Context, group string, name string, newSize compute.VirtualMachineSizeTypes, newCustomSize *compute.VirtualMachineCustomSize, newVirtualMachineGPUs []*compute.VirtualMachineGPU) (err error) {
	vms, err := c.Get(ctx, group, name)
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

	_, err = c.CreateOrUpdate(ctx, group, name, &vm)
	return
}

func (c *VirtualMachineClient) DiskAttach(ctx context.Context, group string, vmName, diskName string) (err error) {
	vms, err := c.Get(ctx, group, vmName)
	if err != nil {
		return err
	}
	if vms == nil || len(*vms) == 0 {
		return errors.Wrapf(errors.NotFound, "Unable to find Virtual Machine [%s]", vmName)
	}

	vm := (*vms)[0]
	for _, disk := range *vm.StorageProfile.DataDisks {
		if *disk.VhdName == diskName {
			return errors.Wrapf(errors.AlreadyExists, "DataDisk [%s] is already attached to the VM [%s]", diskName, vmName)
		}
	}

	*vm.StorageProfile.DataDisks = append(*vm.StorageProfile.DataDisks, compute.DataDisk{VhdName: &diskName})

	_, err = c.CreateOrUpdate(ctx, group, vmName, &vm)
	if err != nil {
		return err
	}

	return

}

func (c *VirtualMachineClient) DiskDetach(ctx context.Context, group string, vmName, diskName string) (err error) {
	vms, err := c.Get(ctx, group, vmName)
	if err != nil {
		return err
	}
	if vms == nil || len(*vms) == 0 {
		return errors.Wrapf(errors.NotFound, "Unable to find Virtual Machine [%s]", vmName)
	}

	vm := (*vms)[0]
	for i, element := range *vm.StorageProfile.DataDisks {
		if *element.VhdName == diskName {
			*vm.StorageProfile.DataDisks = append((*vm.StorageProfile.DataDisks)[:i], (*vm.StorageProfile.DataDisks)[i+1:]...)
			break
		}
	}

	_, err = c.CreateOrUpdate(ctx, group, vmName, &vm)
	if err != nil {
		return err
	}
	return
}

func (c *VirtualMachineClient) NetworkInterfaceAdd(ctx context.Context, group string, vmName, nicName string) (err error) {
	vms, err := c.Get(ctx, group, vmName)
	if err != nil {
		return err
	}
	if vms == nil || len(*vms) == 0 {
		return errors.Wrapf(errors.NotFound, "Unable to find Virtual Machine [%s]", vmName)
	}

	vm := (*vms)[0]
	for _, nic := range *vm.NetworkProfile.NetworkInterfaces {
		if *nic.VirtualNetworkInterfaceReference == nicName {
			return errors.Wrapf(errors.AlreadyExists, "NetworkInterface [%s] is already attached to the VM [%s]", nicName, vmName)
		}
	}

	*vm.NetworkProfile.NetworkInterfaces = append(*vm.NetworkProfile.NetworkInterfaces, compute.NetworkInterfaceReference{VirtualNetworkInterfaceReference: &nicName})

	_, err = c.CreateOrUpdate(ctx, group, vmName, &vm)
	if err != nil {
		return err
	}
	return

}

func (c *VirtualMachineClient) NetworkInterfaceRemove(ctx context.Context, group string, vmName, nicName string) (err error) {
	vms, err := c.Get(ctx, group, vmName)
	if err != nil {
		return err
	}
	if vms == nil || len(*vms) == 0 {
		return errors.Wrapf(errors.NotFound, "Unable to find Virtual Machine [%s]", vmName)
	}

	vm := (*vms)[0]
	for i, element := range *vm.NetworkProfile.NetworkInterfaces {
		if *element.VirtualNetworkInterfaceReference == nicName {
			*vm.NetworkProfile.NetworkInterfaces = append((*vm.NetworkProfile.NetworkInterfaces)[:i], (*vm.NetworkProfile.NetworkInterfaces)[i+1:]...)
			break
		}
	}

	_, err = c.CreateOrUpdate(ctx, group, vmName, &vm)
	if err != nil {
		return err
	}
	return
}

func (c *VirtualMachineClient) NetworkInterfaceList(ctx context.Context, group string, vmName string) (err error) {
	vms, err := c.Get(ctx, group, vmName)
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

	return
}

func (c *VirtualMachineClient) NetworkInterfaceShow(ctx context.Context, group string, vmName, nicName string) (err error) {
	vms, err := c.Get(ctx, group, vmName)
	if err != nil {
		return err
	}
	if vms == nil || len(*vms) == 0 {
		return errors.Wrapf(errors.NotFound, "Unable to find Virtual Machine [%s]", vmName)
	}

	vm := (*vms)[0]
	for _, nic := range *vm.NetworkProfile.NetworkInterfaces {
		if *nic.VirtualNetworkInterfaceReference == nicName {
			// TODO - implement detailed show
			log.Printf("%+v\n", marshal.ToString(nic))
			break
		}
	}

	return
}

func (c *VirtualMachineClient) RunCommand(ctx context.Context, group, vmName string, request *compute.VirtualMachineRunCommandRequest) (response *compute.VirtualMachineRunCommandResponse, err error) {
	return c.internal.RunCommand(ctx, group, vmName, request)
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

func (c *VirtualMachineClient) GetHyperVVmId(ctx context.Context, group string, name string) (*compute.VirtualMachineHyperVVmId, error) {
	return c.internal.GetHyperVVmId(ctx, group, name)
}

func (c *VirtualMachineClient) HasHyperVVm(ctx context.Context, vmName string) (bool, error) {
	return c.internal.HasHyperVVm(ctx, vmName)
}
