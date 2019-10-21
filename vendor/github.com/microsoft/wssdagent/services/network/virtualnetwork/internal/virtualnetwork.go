// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package internal

import (
	"os"
	"path"

	pb "github.com/microsoft/wssdagent/rpc/network"
)

type VirtualNetworkInternal struct {
	Entity     *pb.VirtualNetwork
	Name       string
	Id         string
	ConfigPath string
}

func NewVirtualNetworkInternal(id, basepath string, vnet *pb.VirtualNetwork) *VirtualNetworkInternal {
	basevnetpath := path.Join(basepath, id)
	os.MkdirAll(basevnetpath, os.ModePerm)
	vnet.Id = id
	return &VirtualNetworkInternal{
		Id:         id,
		ConfigPath: basevnetpath,
		Entity:     vnet,
		Name:       vnet.Name,
	}
}
