// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license

package lbagentclient

import (
	"context"

	"github.com/microsoft/moc/pkg/auth"
	pb "github.com/microsoft/moc/rpc/lbagent"
)

// Service interface
type Service interface {
	Get(context.Context, []*pb.LoadBalancer) ([]*pb.LoadBalancer, error)
	CreateOrUpdate(context.Context, []*pb.LoadBalancer) ([]*pb.LoadBalancer, error)
	Delete(context.Context, []*pb.LoadBalancer) error
	GetConfig(context.Context, pb.LoadBalancerType) (string, error)
	AddPeer(context.Context, []string) ([]string, error)
	RemovePeer(context.Context, []string) error
	Resync(ctx context.Context, lbs []*pb.LoadBalancer, peers []string) ([]*pb.LoadBalancer, []string, error)
}

// LoadBalancerAgentClient structure
type LoadBalancerAgentClient struct {
	internal Service
}

// NewLoadBalancerAgentClient method returns new client
func NewLoadBalancerAgentClient(cloudFQDN string, authorizer auth.Authorizer) (*LoadBalancerAgentClient, error) {
	c, err := newLoadBalancerAgentClient(cloudFQDN, authorizer)
	if err != nil {
		return nil, err
	}

	return &LoadBalancerAgentClient{internal: c}, nil
}

// Get methods invokes the client Get method
func (c *LoadBalancerAgentClient) Get(ctx context.Context, lbs []*pb.LoadBalancer) ([]*pb.LoadBalancer, error) {
	return c.internal.Get(ctx, lbs)
}

// Ensure methods invokes create or update on the client
func (c *LoadBalancerAgentClient) CreateOrUpdate(ctx context.Context, lbs []*pb.LoadBalancer) ([]*pb.LoadBalancer, error) {
	return c.internal.CreateOrUpdate(ctx, lbs)
}

// Delete methods invokes delete of the network resource
func (c *LoadBalancerAgentClient) Delete(ctx context.Context, lbs []*pb.LoadBalancer) error {
	return c.internal.Delete(ctx, lbs)
}

// Delete methods invokes delete of the network resource
func (c *LoadBalancerAgentClient) GetConfig(ctx context.Context, lbtype pb.LoadBalancerType) (string, error) {
	return c.internal.GetConfig(ctx, lbtype)
}

func (c *LoadBalancerAgentClient) AddPeer(ctx context.Context, peers []string) ([]string, error) {
	return c.internal.AddPeer(ctx, peers)
}

func (c *LoadBalancerAgentClient) RemovePeer(ctx context.Context, peers []string) error {
	return c.internal.RemovePeer(ctx, peers)
}

func (c *LoadBalancerAgentClient) Resync(ctx context.Context, lbs []*pb.LoadBalancer, peers []string) ([]*pb.LoadBalancer, []string, error) {
	return c.internal.Resync(ctx, lbs, peers)
}
