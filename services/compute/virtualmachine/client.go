// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license

package virtualmachine

import (
	"context"
	"github.com/microsoft/moc/pkg/auth"
	"github.com/microsoft/moc/pkg/errors"

	"github.com/microsoft/wssd-sdk-for-go/services/compute"
	"github.com/microsoft/wssd-sdk-for-go/services/compute/virtualmachine/internal"
)

type Service interface {
	Get(context.Context, string, string) (*[]compute.VirtualMachine, error)
	CreateOrUpdate(context.Context, string, string, *compute.VirtualMachine) (*compute.VirtualMachine, error)
	Delete(context.Context, string, string) error
	Start(context.Context, string, string) error
	Stop(context.Context, string, string) error
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

func (c *VirtualMachineClient) Start(ctx context.Context, group string, name string) (err error) {
	err = c.internal.Start(ctx, group, name)
	return
}
func (c *VirtualMachineClient) Stop(ctx context.Context, group string, name string) (err error) {
	err = c.internal.Stop(ctx, group, name)
	return
}
func (c *VirtualMachineClient) Restart(ctx context.Context, group string, name string) (err error) {
	err = c.internal.Stop(ctx, group, name)
	if err != nil {
		return
	}
	err = c.internal.Start(ctx, group, name)
	return
}
func (c *VirtualMachineClient) Resize(ctx context.Context, group string, name string, newSize compute.VirtualMachineSizeTypes) (err error) {
	vms, err := c.Get(ctx, group, name)
	if err != nil {
		return
	}

	if vms == nil || len(*vms) == 0 {
		err = errors.Wrapf(errors.NotFound, "%s", name)
		return
	}

	vm := (*vms)[0]
	if vm.HardwareProfile.VMSize == newSize {
		// Nothing to do
		return
	}

	vm.HardwareProfile.VMSize = newSize
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
