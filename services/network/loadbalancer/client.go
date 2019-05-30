// Copyright 2019 (c) Microsoft and contributors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package loadbalancer

import (
	"context"
	"github.com/microsoft/wssd-sdk-for-go/services/network"
)

// Service interface
type Service interface {
	Get(context.Context, string) (network.LoadBalancer, error)
	CreateOrUpdate(context.Context, string, string, network.LoadBalancer) (network.LoadBalancer, error)
	Delete(context.Context, string, string) error
}

// LoadBalancerClient structure
type LoadBalancerClient struct {
	network.BaseClient
	internal Service
}

// NewLoadBalancerClient method returns new client
func NewLoadBalancerClient(cloudFQDN string) (*LoadBalancerClient, error) {
	c, err := newLoadBalancerClient(cloudFQDN)
	if err != nil {
		return nil, err
	}

	return &LoadBalancerClient{internal: c}, nil
}

// Get methods invokes the client Get method
func (c *LoadBalancerClient) Get(ctx context.Context, name string) (network.LoadBalancer, error) {
	return c.internal.Get(ctx, name)
}

// Ensure methods invokes create or update on the client
func (c *LoadBalancerClient) CreateOrUpdate(ctx context.Context, name string, id string, lb network.LoadBalancer) (network.LoadBalancer, error) {
	return c.internal.CreateOrUpdate(ctx, name, id, lb)
}

// Delete methods invokes delete of the network resource
func (c *LoadBalancerClient) Delete(ctx context.Context, name string, id string) error {
	return c.internal.Delete(ctx, name, id)
}
