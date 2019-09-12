// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package hcn

import (
	pb "github.com/microsoft/wssdagent/rpc/network"
)

type VirtualNetworkInterfaceProvider struct {
	client *Client
}

func NewVirtualNetworkInterfaceProvider() *VirtualNetworkInterfaceProvider {
	return &VirtualNetworkInterfaceProvider{
		client: NewClient(),
	}
}

func (vnetInterfaceProv *VirtualNetworkInterfaceProvider) Get(vnetInterfaces []*pb.VirtualNetworkInterface) ([]*pb.VirtualNetworkInterface, error) {
	return vnetInterfaceProv.client.Get(vnetInterfaces)
}

func (vnetInterfaceProv *VirtualNetworkInterfaceProvider) CreateOrUpdate(vnetInterfaces []*pb.VirtualNetworkInterface) ([]*pb.VirtualNetworkInterface, error) {
	return vnetInterfaceProv.client.Create(vnetInterfaces)
}

func (vnetInterfaceProv *VirtualNetworkInterfaceProvider) Delete(vnetInterfaces []*pb.VirtualNetworkInterface) error {
	return vnetInterfaceProv.client.Delete(vnetInterfaces)
}
