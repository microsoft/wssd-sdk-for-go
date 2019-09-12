// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package virtualmachinescaleset

import (
	pb "github.com/microsoft/wssdagent/rpc/compute"
)

type VirtualMachineScaleSetProvider interface {
	CreateOrUpdate([]*pb.VirtualMachineScaleSet) ([]*pb.VirtualMachineScaleSet, error)
	Get([]*pb.VirtualMachineScaleSet) ([]*pb.VirtualMachineScaleSet, error)
	Delete([]*pb.VirtualMachineScaleSet) error
}
