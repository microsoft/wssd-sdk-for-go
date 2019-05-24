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

package virtualnetwork

import (
	"context"
	"github.com/microsoft/wssd-sdk-for-go/services/network"
)

// Service interface
type Service interface {
	Get(context.Context, string) (network.VirtualNetwork, error)
	CreateOrUpdate(context.Context, string, string, network.VirtualNetwork) (network.VirtualNetwork, error)
	Delete(context.Context, string, string) (network.VirtualNetwork, error)
}

// Client structure
type VirtualNetworkClient struct {
	BaseClient
	internal Service
}

// NewClient method returns new client
func NewVirtualNetworkClient(cloudFQDN) (*Client, error) {
	c, err := newClient(cloudFQDN)
	if err != nil {
		return nil, err
	}

	return &Client{internal: c}, nil
}

// Get methods invokes the client Get method
func (c *VirtualNetworkClient) Get(ctx context.Context, name string) (*network.VirtualNetwork, error) {
	id, err := c.internal.Get(ctx, name)
	if err != nil && errors.IsNotFound(err) {
		return &network.VirtualNetwork{}, nil
	} else if err != nil {
		return nil, err
	}

	return &id, nil
}

// CreateOrUpdate methods invokes create or update on the client
func (c *VirtualNetworkClient) CreateOrUpdate(ctx context.Context, name string, id string, network network.VirtualNetwork) (network.VirtualNetwork, error) {
	return c.internal.CreateOrUpdate(ctx, name, network)
}

// Delete methods invokes delete of the network resource
func (c *VirtualNetworkClient) Delete(ctx context.Context, name string, id string) error {
	return c.internal.Delete(ctx, name, id)
}
