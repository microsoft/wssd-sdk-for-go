// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license

package internal

import (
	"context"
	"fmt"

	"github.com/microsoft/wssd-sdk-for-go/pkg/auth"
	"github.com/microsoft/wssd-sdk-for-go/services/network"

	"github.com/microsoft/moc/pkg/errors"
	wssdcommonproto "github.com/microsoft/moc/rpc/common"
	wssdnetwork "github.com/microsoft/moc/rpc/nodeagent/network"
	wssdclient "github.com/microsoft/wssd-sdk-for-go/pkg/client"
)

type client struct {
	wssdnetwork.LoadBalancerAgentClient
}

// NewLoadBalancerClient creates a client session with the backend wssd agent
func NewLoadBalancerClient(subID string, authorizer auth.Authorizer) (*client, error) {
	c, err := wssdclient.GetLoadBalancerClient(&subID, authorizer)
	if err != nil {
		return nil, err
	}
	return &client{c}, nil
}

// Get load balancers by name.  If name is nil, get all load balancers
func (c *client) Get(ctx context.Context, group, name string) (*[]network.LoadBalancer, error) {
	request, err := c.getLoadBalancerRequest(wssdcommonproto.Operation_GET, name, nil)
	if err != nil {
		return nil, err
	}
	response, err := c.LoadBalancerAgentClient.Invoke(ctx, request)
	if err != nil {
		return nil, err
	}
	outputLBs, err := c.getLoadBalancersFromResponse(group, response)
	if err != nil {
		return nil, err
	}
	return outputLBs, nil
}

// CreateOrUpdate creates a load balancer if it does not exist, or updates an existing load balancer
func (c *client) CreateOrUpdate(ctx context.Context, group, name string, inputLB *network.LoadBalancer) (*network.LoadBalancer, error) {
	request, err := c.getLoadBalancerRequest(wssdcommonproto.Operation_POST, name, inputLB)
	if err != nil {
		return nil, err
	}
	response, err := c.LoadBalancerAgentClient.Invoke(ctx, request)
	if err != nil {
		return nil, err
	}
	outputLBs, err := c.getLoadBalancersFromResponse(group, response)
	if err != nil {
		return nil, err
	}

	return &(*outputLBs)[0], nil
}

// Delete a load balancer
func (c *client) Delete(ctx context.Context, group, name string) error {
	networkLB, err := c.Get(ctx, group, name)
	if err != nil {
		return err
	}
	if len(*networkLB) == 0 {
		return fmt.Errorf("Load Balancer [%s] not found", name)
	}

	request, err := c.getLoadBalancerRequest(wssdcommonproto.Operation_DELETE, name, &(*networkLB)[0])
	if err != nil {
		return err
	}
	_, err = c.LoadBalancerAgentClient.Invoke(ctx, request)
	if err != nil {
		return err
	}
	return err
}

// getLoadBalancerRequest converts our internal representation of a load balancer (network.LoadBalancer) into a protobuf request (wssdnetwork.LoadBalancerRequest) that can be sent to wssdagent
func (c *client) getLoadBalancerRequest(opType wssdcommonproto.Operation, name string, networkLB *network.LoadBalancer) (*wssdnetwork.LoadBalancerRequest, error) {
	request := &wssdnetwork.LoadBalancerRequest{
		OperationType: opType,
		LoadBalancers: []*wssdnetwork.LoadBalancer{},
	}
	if networkLB != nil {
		wssdLB, err := c.getWssdLoadBalancer(networkLB)
		if err != nil {
			return nil, err
		}
		request.LoadBalancers = append(request.LoadBalancers, wssdLB)
	} else if len(name) > 0 {
		request.LoadBalancers = append(request.LoadBalancers,
			&wssdnetwork.LoadBalancer{
				Name: name,
			})
	}
	return request, nil
}

// getLoadBalancersFromResponse converts a protobuf response from wssdagent (wssdnetwork.LoadBalancerResponse) to out internal representation of a load balancer (network.LoadBalancer)
func (c *client) getLoadBalancersFromResponse(group string, response *wssdnetwork.LoadBalancerResponse) (*[]network.LoadBalancer, error) {
	networkLBs := []network.LoadBalancer{}

	for _, wssdLB := range response.GetLoadBalancers() {
		networkLB, err := c.getNetworkLoadBalancer(group, wssdLB)
		if err != nil {
			return nil, err
		}

		networkLBs = append(networkLBs, *networkLB)
	}

	return &networkLBs, nil
}

// getWssdLoadBalancer convert our internal representation of a loadbalancer (network.LoadBalancer) to the load balancer protobuf used by wssdagent (wssdnetwork.LoadBalancer)
func (c *client) getWssdLoadBalancer(networkLB *network.LoadBalancer) (wssdLB *wssdnetwork.LoadBalancer, err error) {
	/*
		if networkLB.LoadBalancerProperties == nil {
			return nil, errors.Wrapf(errors.InvalidInput, "Missing Load Balancer Properties")
		}
	*/
	wssdLB = &wssdnetwork.LoadBalancer{}
	if networkLB.Name == nil {
		return nil, errors.Wrapf(errors.InvalidInput, "Missing Load Balancer Name")
	}

	wssdLB.Name = *networkLB.Name

	return wssdLB, nil
}

// getNetworkLoadBalancer converts the load balancer protobuf returned from wssdagent (wssdnetwork.LoadBalancer) to our internal representation of a loadbalancer (network.LoadBalancer)
func (c *client) getNetworkLoadBalancer(group string, wssdLB *wssdnetwork.LoadBalancer) (networkLB *network.LoadBalancer, err error) {
	networkLB = &network.LoadBalancer{
		Name: &wssdLB.Name,
		ID:   &wssdLB.Id,
	}

	return networkLB, nil
}
