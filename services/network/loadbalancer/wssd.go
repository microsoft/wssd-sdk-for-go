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

	wssdclient "github.com/microsoft/wssdagent/rpc/client"
	wssdnetwork "github.com/microsoft/wssdagent/rpc/network"
)

type client struct {
	wssdnetwork.LoadBalancerAgentClient
}

// newClient - creates a client session with the backend wssd agent
func newLoadBalancerClient(subID string) (*client, error) {
	c, err := wssdclient.GetLoadBalancerClient(&subID)
	if err != nil {
		return nil, err
	}
	return &client{c}, nil
}

// Get
func (c *client) Get(ctx context.Context, group, name string) (network.LoadBalancer, error) {
	lbrequest := &wssdnetwork.LoadBalancerRequest{OperationType: wssdnetwork.Operation_GET}
	_, err := c.LoadBalancerAgentClient.Invoke(ctx, lbrequest, nil)
	return network.LoadBalancer{}, err
}

// CreateOrUpdate
func (c *client) CreateOrUpdate(ctx context.Context, group, name string, sg network.LoadBalancer) (network.LoadBalancer, error) {
	lbrequest := &wssdnetwork.LoadBalancerRequest{OperationType: wssdnetwork.Operation_POST}
	_, err := c.LoadBalancerAgentClient.Invoke(ctx, lbrequest, nil)
	if err != nil {
		return network.LoadBalancer{}, err
	}

	// Convert _ output to LoadBalancer
	return network.LoadBalancer{}, err
}

// Delete methods invokes create or update on the client
func (c *client) Delete(ctx context.Context, group, name string) error {
	lbrequest := &wssdnetwork.LoadBalancerRequest{OperationType: wssdnetwork.Operation_DELETE}
	_, err := c.LoadBalancerAgentClient.Invoke(ctx, lbrequest, nil)

	// Convert _ output to LoadBalancer
	return err
}
