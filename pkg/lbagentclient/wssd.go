// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license

package lbagentclient

import (
	"context"

	"github.com/microsoft/moc/pkg/auth"
	pbcom "github.com/microsoft/moc/rpc/common"
	pb "github.com/microsoft/moc/rpc/lbagent"
)

type client struct {
	pb.LoadBalancerAgentClient
}

// newClient - creates a client session with the backend wssdcloud agent
func newLoadBalancerAgentClient(subID string, authorizer auth.Authorizer) (*client, error) {
	c, err := GetLoadBalancerAgentClient(&subID, authorizer)
	if err != nil {
		return nil, err
	}
	return &client{c}, nil
}

func (c *client) Get(ctx context.Context, lbs []*pb.LoadBalancer) ([]*pb.LoadBalancer, error) {
	request := &pb.LoadBalancerRequest{
		OperationType: pbcom.Operation_POST,
		LoadBalancers: lbs,
	}

	response, err := c.LoadBalancerAgentClient.Get(ctx, request)
	if err != nil {
		return nil, err
	}
	return response.GetLoadBalancers(), nil
}

// Ensure methods invokes create or update on the client
func (c *client) CreateOrUpdate(ctx context.Context, lbs []*pb.LoadBalancer) ([]*pb.LoadBalancer, error) {
	request := &pb.LoadBalancerRequest{
		OperationType: pbcom.Operation_POST,
		LoadBalancers: lbs,
	}

	response, err := c.LoadBalancerAgentClient.Create(ctx, request)
	if err != nil {
		return nil, err
	}
	return response.GetLoadBalancers(), nil
}

// Delete methods invokes delete of the network resource
func (c *client) Delete(ctx context.Context, lbs []*pb.LoadBalancer) error {
	request := &pb.LoadBalancerRequest{
		OperationType: pbcom.Operation_DELETE,
		LoadBalancers: lbs,
	}

	_, err := c.LoadBalancerAgentClient.Delete(ctx, request)
	if err != nil {
		return err
	}

	return err
}

func (c *client) GetConfig(ctx context.Context, lbtype pb.LoadBalancerType) (string, error) {
	request := &pb.LoadBalancerConfigRequest{Loadbalancertype: lbtype}

	response, err := c.LoadBalancerAgentClient.GetConfig(ctx, request)
	if err != nil {
		return "", err
	}
	return response.GetConfig(), err
}

func (c *client) AddPeer(ctx context.Context, peers []string) ([]string, error) {
	request := &pb.LoadBalancerPeerRequest{
		Peers: peers,
	}

	response, err := c.LoadBalancerAgentClient.AddPeer(ctx, request)
	if err != nil {
		return nil, err
	}
	return response.GetPeers(), nil
}

func (c *client) RemovePeer(ctx context.Context, peers []string) error {
	request := &pb.LoadBalancerPeerRequest{
		Peers: peers,
	}

	_, err := c.LoadBalancerAgentClient.RemovePeer(ctx, request)
	if err != nil {
		return err
	}
	return nil
}

func (c *client) Resync(ctx context.Context, lbs []*pb.LoadBalancer, peers []string) ([]*pb.LoadBalancer, []string, error) {
	request := &pb.LoadBalancerResyncRequest{
		LoadBalancers: lbs,
		Peers:         peers,
	}

	response, err := c.LoadBalancerAgentClient.Resync(ctx, request)
	if err != nil {
		return nil, nil, err
	}
	return response.GetLoadBalancers(), response.GetPeers(), nil
}
