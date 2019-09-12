// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package vmms

import (
	"fmt"
	pb "github.com/microsoft/wssdagent/rpc/storage"
)

type VirtualHardDiskProvider struct {
}

func NewVirtualHardDiskProvider() *VirtualHardDiskProvider {
	return &VirtualHardDiskProvider{}
}

func (*VirtualHardDiskProvider) Get([]*pb.VirtualHardDisk) ([]*pb.VirtualHardDisk, error) {
	return nil, fmt.Errorf("[VirtualHardDiskProvider] Get not implemented")
}

func (*VirtualHardDiskProvider) CreateOrUpdate([]*pb.VirtualHardDisk) ([]*pb.VirtualHardDisk, error) {
	return nil, fmt.Errorf("[VirtualHardDiskProvider] CreateOrUpdate not implemented")
}

func (*VirtualHardDiskProvider) Delete([]*pb.VirtualHardDisk) error {
	return fmt.Errorf("[VirtualHardDiskProvider] Delete not implemented")
}
