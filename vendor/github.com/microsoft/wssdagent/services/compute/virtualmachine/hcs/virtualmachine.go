// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package hcs

import (
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

func (vmp *VirtualMachineProvider) Get(vm []*pb.VirtualMachine) ([]*pb.VirtualMachine, error) {
	return vmp.client.Get(vm)
}

func (vmp *VirtualMachineProvider) CreateOrUpdate(vms []*pb.VirtualMachine) ([]*pb.VirtualMachine, error) {
	return vmp.client.Create(vms)
}

func (vmp *VirtualMachineProvider) Delete(vms []*pb.VirtualMachine) error {
	return vmp.client.Delete(vms)
}
