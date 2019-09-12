// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package internal

import (
	"os"
	"path"

	pb "github.com/microsoft/wssdagent/rpc/network"
)

type VirtualNetworkInterfaceInternal struct {
	VNic       *pb.VirtualNetworkInterface
	Id         string
	ConfigPath string
}

func NewVirtualNetworkInterfaceInternal(id, basepath string) *VirtualNetworkInterfaceInternal {
	basevnicpath := path.Join(basepath, id)
	os.MkdirAll(basevnicpath, os.ModePerm)
	return &VirtualNetworkInterfaceInternal{
		Id:         id,
		ConfigPath: basevnicpath,
	}
}
