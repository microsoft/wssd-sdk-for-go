// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package internal

import (
	"os"
	"path"

	pb "github.com/microsoft/wssdagent/rpc/network"
)

type VirtualNetworkInternal struct {
	VNet       *pb.VirtualNetwork
	Id         string
	ConfigPath string
}

func NewVirtualNetworkInternal(id, basepath string) *VirtualNetworkInternal {
	basevnetpath := path.Join(basepath, id)
	os.MkdirAll(basevnetpath, os.ModePerm)
	return &VirtualNetworkInternal{
		Id:         id,
		ConfigPath: basevnetpath,
	}
}
