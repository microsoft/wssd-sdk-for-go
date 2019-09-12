// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package virtualmachine

import (
	"os"
	"path"

	pb "github.com/microsoft/wssdagent/rpc/compute"
)

type VirtualMachineProvider interface {
	CreateOrUpdate([]*pb.VirtualMachine) ([]*pb.VirtualMachine, error)
	Get([]*pb.VirtualMachine) ([]*pb.VirtualMachine, error)
	Delete([]*pb.VirtualMachine) error
}

type VirtualMachineInternal struct {
	vm         *pb.VirtualMachine
	Id         string
	ConfigPath string
	SeedIso    string
	UserData   string
	MetaData   string
}

func NewVirtualMachineInternal(id, basepath string) *VirtualMachineInternal {
	basevmpath := path.Join(basepath, id)
	os.MkdirAll(basevmpath, os.ModePerm)
	basecloudinitpath := path.Join(basevmpath, "data")
	os.MkdirAll(basecloudinitpath, os.ModePerm)
	return &VirtualMachineInternal{
		Id:         id,
		ConfigPath: basevmpath,
		SeedIso:    path.Join(basecloudinitpath, "seed.iso"),
		UserData:   path.Join(basecloudinitpath, "user-data"),
		MetaData:   path.Join(basecloudinitpath, "meta-data"),
	}
}

// CreateVirtualNetworkInterface
func HasVirtualMachine(provider VirtualMachineProvider, vmName string) bool {
	vm := &pb.VirtualMachine{Name: vmName}
	_, err := provider.Get([]*pb.VirtualMachine{vm})

	if err != nil {
		return false
	}
	return true
}
