// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package vmms

import (
	"fmt"
	pb "github.com/microsoft/wssdagent/rpc/network"
)

type LoadBalancerProvider struct {
}

func NewLoadBalancerProvider() *LoadBalancerProvider {
	return &LoadBalancerProvider{}
}

func (*LoadBalancerProvider) Get([]*pb.LoadBalancer) ([]*pb.LoadBalancer, error) {
	return nil, fmt.Errorf("[LoadBalancerProvider] Get not implemented")
}

func (*LoadBalancerProvider) CreateOrUpdate([]*pb.LoadBalancer) ([]*pb.LoadBalancer, error) {
	return nil, fmt.Errorf("[LoadBalancerProvider] CreateOrUpdate not implemented")
}

func (*LoadBalancerProvider) Delete([]*pb.LoadBalancer) error {
	return fmt.Errorf("[LoadBalancerProvider] Delete not implemented")
}
