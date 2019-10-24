// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license

package internal

import (
	"context"
	"github.com/microsoft/wssd-sdk-for-go/pkg/auth"
	"github.com/microsoft/wssd-sdk-for-go/services/network"

	wssdclient "github.com/microsoft/wssd-sdk-for-go/pkg/client"
	wssdnetwork "github.com/microsoft/wssdagent/rpc/network"
)

type client struct {
	wssdnetwork.LoadBalancerAgentClient
}

// NewLoadBalancerClient- creates a client session with the backend wssd agent
func NewLoadBalancerClient(subID string, authorizer auth.Authorizer) (*client, error) {
	c, err := wssdclient.GetLoadBalancerClient(&subID, authorizer)
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
