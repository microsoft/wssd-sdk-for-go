// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package internal

import (
	"os"
	"path"

	pb "github.com/microsoft/wssdagent/rpc/compute"
)

type VirtualMachineScaleSetInternal struct {
	Entity     *pb.VirtualMachineScaleSet
	Id         string
	ConfigPath string
	Name       string
}

func NewVirtualMachineScaleSetInternal(id, basepath string, vmss *pb.VirtualMachineScaleSet) *VirtualMachineScaleSetInternal {
	basevmpath := path.Join(basepath, id)
	os.MkdirAll(basevmpath, os.ModePerm)
	vmss.Id = id
	return &VirtualMachineScaleSetInternal{
		Id:         id,
		ConfigPath: basevmpath,
		Entity:     vmss,
		Name:       vmss.Name,
	}
}
