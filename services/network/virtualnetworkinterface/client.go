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

package virtualnetworkinterface

import (
	"context"
	"github.com/microsoft/wssd-sdk-for-go/services/network"
)

// Service interface
type Service interface {
	Get(context.Context, string) (*[]network.VirtualNetworkInterface, error)
	CreateOrUpdate(context.Context, string, string, *network.VirtualNetworkInterface) (*network.VirtualNetworkInterface, error)
	Delete(context.Context, string, string) error
}

// VirtualNetworkInterfaceClient structure
type VirtualNetworkInterfaceClient struct {
	network.BaseClient
	internal Service
}

// NewVirtualNetworkInterfaceClient method returns new client
func NewVirtualNetworkInterfaceClient(cloudFQDN string) (*VirtualNetworkInterfaceClient, error) {
	c, err := newVirtualNetworkInterfaceClient(cloudFQDN)
	if err != nil {
		return nil, err
	}

	return &VirtualNetworkInterfaceClient{internal: c}, nil
}

// Get methods invokes the client Get method
func (c *VirtualNetworkInterfaceClient) Get(ctx context.Context, name string) (*[]network.VirtualNetworkInterface, error) {
	return c.internal.Get(ctx, name)
}

// CreateOrUpdate methods invokes create or update on the client
func (c *VirtualNetworkInterfaceClient) CreateOrUpdate(ctx context.Context, name string, id string, networkInterface *network.VirtualNetworkInterface) (*network.VirtualNetworkInterface, error) {
	return c.internal.CreateOrUpdate(ctx, name, id, networkInterface)
}

// Delete methods invokes delete of the network interface resource
func (c *VirtualNetworkInterfaceClient) Delete(ctx context.Context, name string, id string) error {
	return c.internal.Delete(ctx, name, id)
}
