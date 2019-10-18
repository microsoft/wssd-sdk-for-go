// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package loadbalancer

import (
	pb "github.com/microsoft/wssdagent/rpc/network"
)

type LoadBalancerProvider interface {
	CreateOrUpdate([]*pb.LoadBalancer) ([]*pb.LoadBalancer, error)
	Get([]*pb.LoadBalancer) ([]*pb.LoadBalancer, error)
	Delete([]*pb.LoadBalancer) error
}
