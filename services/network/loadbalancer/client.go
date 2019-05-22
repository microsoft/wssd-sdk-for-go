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

package network

import (
	"context"
)

// Service interface
type Service interface {
	Get(context.Context, string, string) (loadbalancer.LoadBalancer, error)
	CreateOrUpdate(context.Context, string, string, loadbalancer.LoadBalancer) (loadbalancer.LoadBalancer, error)
	Delete(context.Context, string, string) (loadbalancer.LoadBalancer, error)
}

// LoadBalancerClient structure
type LoadBalancerClient struct {
	BaseClient
	group    string
	internal Service
}

// NewLoadBalancerClient method returns new client
func NewLoadBalancerClient(subID, group string) (*LoadBalancerClient, error) {
	c, err := newLoadBalancerClient(subID)
	if err != nil {
		return nil, err
	}

	return &LoadBalancerClient{group: group, internal: c}, nil
}

// Get methods invokes the client Get method
func (c *LoadBalancerClient) Get(ctx context.Context, name string) (*Spec, error) {
	id, err := c.internal.Get(ctx, c.group, name)
	if err != nil && errors.IsNotFound(err) {
		return &Spec{&loadbalancer.LoadBalancer{}}, nil
	} else if err != nil {
		return nil, err
	}

	return &Spec{&id}, nil
}

// Ensure methods invokes create or update on the client
func (c *LoadBalancerClient) Ensure(ctx context.Context, name string, spec *Spec) error {
	result, err := c.internal.CreateOrUpdate(ctx, c.group, name, *spec.internal)
	if err != nil {
		return err
	}
	spec.internal = &result
	return nil
}
