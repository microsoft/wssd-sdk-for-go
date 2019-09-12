// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package internal

import (
	"os"
	"path"

	pb "github.com/microsoft/wssdagent/rpc/compute"
)

type VirtualMachineScaleSetInternal struct {
	Vmss       *pb.VirtualMachineScaleSet
	Id         string
	ConfigPath string
}

func NewVirtualMachineScaleSetInternal(id, basepath string) *VirtualMachineScaleSetInternal {
	basevmpath := path.Join(basepath, id)
	os.MkdirAll(basevmpath, os.ModePerm)
	return &VirtualMachineScaleSetInternal{
		Id:         id,
		ConfigPath: basevmpath,
	}
}
