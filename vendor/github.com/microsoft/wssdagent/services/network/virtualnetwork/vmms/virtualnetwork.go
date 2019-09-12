// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package vmms

import (
	"fmt"
	pb "github.com/microsoft/wssdagent/rpc/network"
)

type VirtualNetworkProvider struct {
}

// NewVirtualNetworkProvider creates a new vmms based provider
func NewVirtualNetworkProvider() *VirtualNetworkProvider {
	return &VirtualNetworkProvider{}
}

func (*VirtualNetworkProvider) Get([]*pb.VirtualNetwork) ([]*pb.VirtualNetwork, error) {
	return nil, fmt.Errorf("[VirtualNetworkProvider] Get not implemented")
}

func (*VirtualNetworkProvider) CreateOrUpdate([]*pb.VirtualNetwork) ([]*pb.VirtualNetwork, error) {
	return nil, fmt.Errorf("[VirtualNetworkProvider] Get not implemented")
}

func (*VirtualNetworkProvider) Delete([]*pb.VirtualNetwork) error {
	return fmt.Errorf("[VirtualNetworkProvider] Get not implemented")
}
