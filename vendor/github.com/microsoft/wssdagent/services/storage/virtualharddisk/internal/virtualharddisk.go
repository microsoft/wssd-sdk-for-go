// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package internal

import (
	"os"
	"path"

	pb "github.com/microsoft/wssdagent/rpc/storage"
)

type VirtualHardDiskInternal struct {
	Entity     *pb.VirtualHardDisk
	Id         string
	ConfigPath string
	Name       string
}

func NewVirtualHardDiskInternal(id, basepath string, vhd *pb.VirtualHardDisk) *VirtualHardDiskInternal {
	basevhdpath := path.Join(basepath, id)
	os.MkdirAll(basevhdpath, os.ModePerm)
	vhd.Id = id
	return &VirtualHardDiskInternal{
		Id:         id,
		ConfigPath: basevhdpath,
		Entity:     vhd,
		Name:       vhd.Name,
	}
}
