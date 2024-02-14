// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license

package logicalnetwork

import (
	"context"

	"github.com/microsoft/moc/pkg/auth"
	"github.com/microsoft/wssd-sdk-for-go/services/network"
	"github.com/microsoft/wssd-sdk-for-go/services/network/logicalnetwork/internal"
)

// Service interface
type Service interface {
	Get(context.Context, string) (*[]network.LogicalNetwork, error)
	CreateOrUpdate(context.Context, string, *network.LogicalNetwork) (*network.LogicalNetwork, error)
	Delete(context.Context, string) error
}

// Client structure
type LogicalNetworkClient struct {
	network.BaseClient
	internal Service
}

// NewClient method returns new client
func NewLogicalNetworkClient(cloudFQDN string, authorizer auth.Authorizer) (*LogicalNetworkClient, error) {
	c, err := internal.NewLogicalNetworkClient(cloudFQDN, authorizer)
	if err != nil {
		return nil, err
	}

	return &LogicalNetworkClient{internal: c}, nil
}

// Get methods invokes the client Get method
func (c *LogicalNetworkClient) Get(ctx context.Context, name string) (*[]network.LogicalNetwork, error) {
	return c.internal.Get(ctx, name)
}

// CreateOrUpdate methods invokes create or update on the client
func (c *LogicalNetworkClient) CreateOrUpdate(ctx context.Context, name string, network *network.LogicalNetwork) (*network.LogicalNetwork, error) {
	return c.internal.CreateOrUpdate(ctx, name, network)
}

// Delete methods invokes delete of the network resource
func (c *LogicalNetworkClient) Delete(ctx context.Context, name string) error {
	return c.internal.Delete(ctx, name)
}
