// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package vmms

import (
	"fmt"
	pb "github.com/microsoft/wssdagent/rpc/network"
)

type VirtualNetworkInterfaceProvider struct {
}

func NewVirtualNetworkInterfaceProvider() *VirtualNetworkInterfaceProvider {
	return &VirtualNetworkInterfaceProvider{}
}

func (*VirtualNetworkInterfaceProvider) Get([]*pb.VirtualNetworkInterface) ([]*pb.VirtualNetworkInterface, error) {
	return nil, fmt.Errorf("[VirtualNetworkInterfaceProvider] Get not implemented")
}

func (*VirtualNetworkInterfaceProvider) CreateOrUpdate([]*pb.VirtualNetworkInterface) ([]*pb.VirtualNetworkInterface, error) {
	return nil, fmt.Errorf("[VirtualNetworkInterfaceProvider] Get not implemented")
}

func (*VirtualNetworkInterfaceProvider) Delete([]*pb.VirtualNetworkInterface) error {
	return fmt.Errorf("[VirtualNetworkInterfaceProvider] Get not implemented")
}
