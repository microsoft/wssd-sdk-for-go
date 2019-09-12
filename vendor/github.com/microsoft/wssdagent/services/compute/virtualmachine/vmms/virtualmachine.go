// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package vmms

import (
	"fmt"

	pb "github.com/microsoft/wssdagent/rpc/compute"
)

type VirtualMachineProvider struct {
}

// NewVirtualMachineProvider creates a new vmms based provider
func NewVirtualMachineProvider() *VirtualMachineProvider {
	return &VirtualMachineProvider{}
}

func (*VirtualMachineProvider) Get(vms []*pb.VirtualMachine) ([]*pb.VirtualMachine, error) {

	return nil, nil

}

func (*VirtualMachineProvider) CreateOrUpdate(vms []*pb.VirtualMachine) ([]*pb.VirtualMachine, error) {
	return nil, fmt.Errorf("CreateOrUpdate not implemented")
}

func (*VirtualMachineProvider) Delete([]*pb.VirtualMachine) error {
	return fmt.Errorf("Delete not implemented")
}
