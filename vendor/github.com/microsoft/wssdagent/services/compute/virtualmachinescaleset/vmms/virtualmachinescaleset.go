// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package vmms

import (
	"fmt"
	pb "github.com/microsoft/wssdagent/rpc/compute"
)

type VirtualMachineScaleSetProvider struct {
}

// NewVirtualMachineScaleSetProvider creates a new vmms based provider
func NewVirtualMachineScaleSetProvider() *VirtualMachineScaleSetProvider {
	return &VirtualMachineScaleSetProvider{}
}

func (*VirtualMachineScaleSetProvider) Get([]*pb.VirtualMachineScaleSet) ([]*pb.VirtualMachineScaleSet, error) {
	return nil, fmt.Errorf("[VirtualHardDiskProvider] Get not implemented")
}

func (*VirtualMachineScaleSetProvider) CreateOrUpdate([]*pb.VirtualMachineScaleSet) ([]*pb.VirtualMachineScaleSet, error) {
	return nil, fmt.Errorf("[VirtualHardDiskProvider] CreateOrUpdate not implemented")
}

func (*VirtualMachineScaleSetProvider) Delete([]*pb.VirtualMachineScaleSet) error {
	return fmt.Errorf("[VirtualHardDiskProvider] Delete not implemented")
}
