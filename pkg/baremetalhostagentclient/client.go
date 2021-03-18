// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license

package baremetalhostagentclient

import (
	"context"

	"github.com/microsoft/moc/pkg/auth"
	pb "github.com/microsoft/moc/rpc/baremetalhostagent"
)

// Service interface
type Service interface {
	Update(context.Context, *pb.BareMetalHost) (*pb.BareMetalHost, error)
}

// BareMetalHostAgentClient structure
type BareMetalHostAgentClient struct {
	internal Service
}

// NewBareMetalHostAgentClient method returns new client
func NewBareMetalHostAgentClient(cloudFQDN string, authorizer auth.Authorizer) (*BareMetalHostAgentClient, error) {
	c, err := newBareMetalHostAgentClient(cloudFQDN, authorizer)
	if err != nil {
		return nil, err
	}

	return &BareMetalHostAgentClient{internal: c}, nil
}

// Update method invokes update on the client
func (c *BareMetalHostAgentClient) Update(ctx context.Context, bmh *pb.BareMetalHost) (*pb.BareMetalHost, error) {
	return c.internal.Update(ctx, bmh)
}
