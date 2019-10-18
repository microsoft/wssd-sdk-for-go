// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package virtualnetwork

import (
	"context"
	"github.com/microsoft/wssdagent/pkg/errors"
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

func (vnetProv *VirtualNetworkProvider) Get(ctx context.Context, vnets []*pb.VirtualNetwork) ([]*pb.VirtualNetwork, error) {
	newvnets := []*pb.VirtualNetwork{}
	if len(vnets) == 0 {
		// Get Everything
		return vnetProv.client.Get(ctx, nil)
	}

	// Get only requested vnets
	for _, vnet := range vnets {
		newvnet, err := vnetProv.client.Get(ctx, vnet)
		if err != nil {
			return newvnets, err
		}
		newvnets = append(newvnets, newvnet[0])
	}
	return newvnets, nil
}

func (vnetProv *VirtualNetworkProvider) CreateOrUpdate(ctx context.Context, vnets []*pb.VirtualNetwork) ([]*pb.VirtualNetwork, error) {
	newvnets := []*pb.VirtualNetwork{}
	for _, vnet := range vnets {
		newvnet, err := vnetProv.client.Create(ctx, vnet)
		if err != nil {
			if err != errors.AlreadyExists {
				vnetProv.client.Delete(ctx, vnet)
			}
			return newvnets, err
		}
		newvnets = append(newvnets, newvnet)
	}

	return newvnets, nil
}

func (vnetProv *VirtualNetworkProvider) Delete(ctx context.Context, vnets []*pb.VirtualNetwork) error {
	for _, vnet := range vnets {
		err := vnetProv.client.Delete(ctx, vnet)
		if err != nil {
			return err
		}
	}

	return nil
}
