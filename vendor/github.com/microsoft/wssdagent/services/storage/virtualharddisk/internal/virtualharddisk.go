// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package internal

import (
	"os"
	"path"

	pb "github.com/microsoft/wssdagent/rpc/storage"
)

type VirtualHardDiskInternal struct {
	Vhd        *pb.VirtualHardDisk
	Id         string
	ConfigPath string
}

func NewVirtualHardDiskInternal(id, basepath string) *VirtualHardDiskInternal {
	basevhdpath := path.Join(basepath, id)
	os.MkdirAll(basevhdpath, os.ModePerm)
	return &VirtualHardDiskInternal{
		Id:         id,
		ConfigPath: basevhdpath,
	}
}
