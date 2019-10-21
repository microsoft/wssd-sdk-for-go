// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package internal

import (
	"os"
	"path"

	pb "github.com/microsoft/wssdagent/rpc/compute"
)

type VirtualMachineInternal struct {
	Entity     *pb.VirtualMachine
	Id         string
	Name       string
	ConfigPath string
	SeedIso    string
	UserData   string
	MetaData   string
	VendorData string
}

func NewVirtualMachineInternal(id, basepath string, vm *pb.VirtualMachine) *VirtualMachineInternal {
	basevmpath := path.Join(basepath, id)
	os.MkdirAll(basevmpath, os.ModePerm)
	basecloudinitpath := path.Join(basevmpath, "data")
	os.MkdirAll(basecloudinitpath, os.ModePerm)
	vm.Id = id
	return &VirtualMachineInternal{
		Id:         id,
		ConfigPath: basevmpath,
		Entity:     vm,
		Name:       vm.Name,
		SeedIso:    path.Join(basecloudinitpath, "seed.iso"),
		UserData:   path.Join(basecloudinitpath, "user-data"),
		MetaData:   path.Join(basecloudinitpath, "meta-data"),
		VendorData: path.Join(basecloudinitpath, "vendor-data"),
	}
}
