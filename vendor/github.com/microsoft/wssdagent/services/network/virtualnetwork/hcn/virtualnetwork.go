// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package hcn

import (
	pb "github.com/microsoft/wssdagent/rpc/network"
)

type VirtualNetworkProvider struct {
	client *Client
}

func NewVirtualNetworkProvider() *VirtualNetworkProvider {
	return &VirtualNetworkProvider{
		client: NewClient(),
	}
}

func (vnetProv *VirtualNetworkProvider) Get(vnets []*pb.VirtualNetwork) ([]*pb.VirtualNetwork, error) {
	return vnetProv.client.Get(vnets)
}

func (vnetProv *VirtualNetworkProvider) CreateOrUpdate(vnets []*pb.VirtualNetwork) ([]*pb.VirtualNetwork, error) {
	return vnetProv.client.Create(vnets)
}

func (vnetProv *VirtualNetworkProvider) Delete(vnets []*pb.VirtualNetwork) error {
	return vnetProv.client.Delete(vnets)
}
