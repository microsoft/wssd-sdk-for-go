// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package virtualmachinescaleset

import (
	"context"
	"github.com/microsoft/wssdagent/pkg/errors"
	pb "github.com/microsoft/wssdagent/rpc/compute"
)

type VirtualMachineScaleSetProvider struct {
	client *Client
}

func NewVirtualMachineScaleSetProvider() *VirtualMachineScaleSetProvider {
	return &VirtualMachineScaleSetProvider{
		client: NewClient(),
	}
}

func (vmssProv *VirtualMachineScaleSetProvider) Get(ctx context.Context, vmsss []*pb.VirtualMachineScaleSet) ([]*pb.VirtualMachineScaleSet, error) {
	newvmsss := []*pb.VirtualMachineScaleSet{}
	if len(vmsss) == 0 {
		// Get Everything
		return vmssProv.client.Get(ctx, nil)
	}

	// Get only requested vmsss
	for _, vmss := range vmsss {
		newvmss, err := vmssProv.client.Get(ctx, vmss)
		if err != nil {
			return newvmsss, err
		}
		newvmsss = append(newvmsss, newvmss[0])
	}
	return newvmsss, nil
}

func (vmssProv *VirtualMachineScaleSetProvider) CreateOrUpdate(ctx context.Context, vmsss []*pb.VirtualMachineScaleSet) ([]*pb.VirtualMachineScaleSet, error) {
	newvmsss := []*pb.VirtualMachineScaleSet{}
	for _, vmss := range vmsss {
		newvmss, err := vmssProv.client.Create(ctx, vmss)
		if err != nil {
			if err != errors.AlreadyExists {
				vmssProv.client.Delete(ctx, vmss)
			}
			return newvmsss, err
		}
		newvmsss = append(newvmsss, newvmss)
	}

	return newvmsss, nil
}

func (vmssProv *VirtualMachineScaleSetProvider) Delete(ctx context.Context, vmsss []*pb.VirtualMachineScaleSet) error {
	for _, vmss := range vmsss {
		err := vmssProv.client.Delete(ctx, vmss)
		if err != nil {
			return err
		}
	}

	return nil
}
