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
	Get(context.Context, string, string) (VirtualNetwork, error)
	CreateOrUpdate(context.Context, string, string, VirtualNetwork) (VirtualNetwork, error)
	Delete(context.Context, string, string) (VirtualNetwork, error)
}

// Client structure
type VirtualNetworkClient struct {
	BaseClient
	internal Service
}

// NewClient method returns new client
func NewVirtualNetworkClient(subID, group string) (*Client, error) {
	c, err := newClient(subID)
	if err != nil {
		return nil, err
	}

	return &Client{group: group, internal: c}, nil
}

// Get methods invokes the client Get method
func (c *VirtualNetworkClient) Get(ctx context.Context, name string) (*Spec, error) {
	id, err := c.internal.Get(ctx, c.group, name)
	if err != nil && errors.IsNotFound(err) {
		return &Spec{&VirtualNetwork{}}, nil
	} else if err != nil {
		return nil, err
	}

	return &Spec{&id}, nil
}

// CreateOrUpdate methods invokes create or update on the client
func (c *VirtualNetworkClient) CreateOrUpdate(ctx context.Context, name string, id string, network VirtualNetwork) (VirtualNetwork, error) {
	return c.internal.CreateOrUpdate(ctx, name, network)
}

// Delete methods invokes delete of the network resource
func (c *VirtualNetworkClient) Delete(ctx context.Context, name string, id string) error {
	return c.internal.Delete(ctx, name, id)
}
