// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package virtualnetwork

import (
	pb "github.com/microsoft/wssdagent/rpc/network"
)

type VirtualNetworkProvider interface {
	CreateOrUpdate([]*pb.VirtualNetwork) ([]*pb.VirtualNetwork, error)
	Get([]*pb.VirtualNetwork) ([]*pb.VirtualNetwork, error)
	Delete([]*pb.VirtualNetwork) error
}
