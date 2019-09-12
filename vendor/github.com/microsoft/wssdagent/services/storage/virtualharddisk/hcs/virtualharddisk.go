// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package hcs

import (
	//	"fmt"
	pb "github.com/microsoft/wssdagent/rpc/storage"
)

type VirtualHardDiskProvider struct {
	client *Client
}

func NewVirtualHardDiskProvider() *VirtualHardDiskProvider {
	return &VirtualHardDiskProvider{
		client: NewClient(),
	}
}

func (vhdp *VirtualHardDiskProvider) Get(vhdList []*pb.VirtualHardDisk) ([]*pb.VirtualHardDisk, error) {
	return vhdp.client.Get(vhdList)
}

func (vhdp *VirtualHardDiskProvider) CreateOrUpdate(vhdList []*pb.VirtualHardDisk) ([]*pb.VirtualHardDisk, error) {
	return vhdp.client.Create(vhdList)
}

func (vhdp *VirtualHardDiskProvider) Delete(vhdList []*pb.VirtualHardDisk) error {
	return vhdp.client.Delete(vhdList)
}
