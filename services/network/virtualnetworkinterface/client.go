// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license

package virtualnetworkinterface

import (
	"context"

	"github.com/microsoft/moc/pkg/auth"
	"github.com/microsoft/wssd-sdk-for-go/services/network"
	"github.com/microsoft/wssd-sdk-for-go/services/network/virtualnetworkinterface/internal"
)

// Service interface
type Service interface {
	Get(context.Context, string, string) (*[]network.VirtualNetworkInterface, error)
	CreateOrUpdate(context.Context, string, string, *network.VirtualNetworkInterface) (*network.VirtualNetworkInterface, error)
	Hydrate(context.Context, string, string, *network.VirtualNetworkInterface) (*network.VirtualNetworkInterface, error)
	// Hydrate(context.Context, string, string, string, string) (*network.VirtualNetworkInterface, error)
	Delete(context.Context, string, string) error
	Update(context.Context, string, string, *network.VirtualNetworkInterface) (*network.VirtualNetworkInterface, error)
}

// VirtualNetworkInterfaceClient structure
type VirtualNetworkInterfaceClient struct {
	network.BaseClient
	internal Service
}

// NewVirtualNetworkInterfaceClient method returns new client
func NewVirtualNetworkInterfaceClient(cloudFQDN string, authorizer auth.Authorizer) (*VirtualNetworkInterfaceClient, error) {
	c, err := internal.NewVirtualNetworkInterfaceClient(cloudFQDN, authorizer)
	if err != nil {
		return nil, err
	}

	return &VirtualNetworkInterfaceClient{internal: c}, nil
}

// Get methods invokes the client Get method
func (c *VirtualNetworkInterfaceClient) Get(ctx context.Context, group, name string) (*[]network.VirtualNetworkInterface, error) {
	return c.internal.Get(ctx, group, name)
}

// CreateOrUpdate methods invokes create or update on the client
func (c *VirtualNetworkInterfaceClient) CreateOrUpdate(ctx context.Context, group, name string, networkInterface *network.VirtualNetworkInterface) (*network.VirtualNetworkInterface, error) {
	return c.internal.CreateOrUpdate(ctx, group, name, networkInterface)
}

// Hydrate method invokes hydration of the networkinterface resource
func (c *VirtualNetworkInterfaceClient) Hydrate(ctx context.Context, group, name string, networkInterface *network.VirtualNetworkInterface) (*network.VirtualNetworkInterface, error) {
	return c.internal.Hydrate(ctx, group, name, networkInterface)
}

// Update supported settings on the networkinterface.
func (c *VirtualNetworkInterfaceClient) Update(ctx context.Context, group, name string, networkInterface *network.VirtualNetworkInterface) (*network.VirtualNetworkInterface, error) {
	return c.internal.Update(ctx, group, name, networkInterface)
}

// Delete methods invokes delete of the network interface resource
func (c *VirtualNetworkInterfaceClient) Delete(ctx context.Context, group, name string) error {
	return c.internal.Delete(ctx, group, name)
}
