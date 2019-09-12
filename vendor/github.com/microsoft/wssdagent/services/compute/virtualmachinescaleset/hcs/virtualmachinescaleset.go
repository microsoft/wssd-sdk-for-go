// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package hcs

import (
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

func (vmssp *VirtualMachineScaleSetProvider) Get(vmss []*pb.VirtualMachineScaleSet) ([]*pb.VirtualMachineScaleSet, error) {
	return vmssp.client.Get(vmss)
}

func (vmssp *VirtualMachineScaleSetProvider) CreateOrUpdate(vmss []*pb.VirtualMachineScaleSet) ([]*pb.VirtualMachineScaleSet, error) {
	return vmssp.client.Create(vmss)
}

func (vmssp *VirtualMachineScaleSetProvider) Delete(vmss []*pb.VirtualMachineScaleSet) error {
	return vmssp.client.Delete(vmss)
}
