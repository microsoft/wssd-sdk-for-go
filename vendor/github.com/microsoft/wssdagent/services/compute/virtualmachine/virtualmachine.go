// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package virtualmachine

import (
	"context"
	"github.com/microsoft/wssdagent/pkg/errors"
	pb "github.com/microsoft/wssdagent/rpc/compute"
)

type VirtualMachineProvider struct {
	client *Client
}

func NewVirtualMachineProvider() *VirtualMachineProvider {
	return &VirtualMachineProvider{
		client: NewClient(),
	}
}

func (vmProv *VirtualMachineProvider) Get(ctx context.Context, vms []*pb.VirtualMachine) ([]*pb.VirtualMachine, error) {
	newvms := []*pb.VirtualMachine{}
	if len(vms) == 0 {
		// Get Everything
		return vmProv.client.Get(ctx, nil)
	}

	// Get only requested vms
	for _, vm := range vms {
		newvm, err := vmProv.client.Get(ctx, vm)
		if err != nil {
			return newvms, err
		}
		newvms = append(newvms, newvm[0])
	}
	return newvms, nil
}

func (vmProv *VirtualMachineProvider) CreateOrUpdate(ctx context.Context, vms []*pb.VirtualMachine) ([]*pb.VirtualMachine, error) {
	newvms := []*pb.VirtualMachine{}
	for _, vm := range vms {
		newvm, err := vmProv.client.Create(ctx, vm)
		if err != nil {
			if err != errors.AlreadyExists {
				vmProv.client.Delete(ctx, vm)
			}
			return newvms, err
		}
		newvms = append(newvms, newvm)
	}

	return newvms, nil
}

func (vmProv *VirtualMachineProvider) Delete(ctx context.Context, vms []*pb.VirtualMachine) error {
	for _, vm := range vms {
		err := vmProv.client.Delete(ctx, vm)
		if err != nil {
			return err
		}
	}

	return nil
}

// CreateVirtualMachineInterface
func (vmProv *VirtualMachineProvider) HasVirtualMachie(ctx context.Context, vmName string) bool {
	vm := &pb.VirtualMachine{Name: vmName}
	_, err := vmProv.Get(ctx, []*pb.VirtualMachine{vm})

	if err != nil {
		return false
	}
	return true
}
