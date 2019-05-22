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

	wssdclient "github.com/microsoft/wssdagent/rpc/client"
	wssdnetwork "github.com/microsoft/wssdagent/rpc/network"
)

type client struct {
	wssdnetwork.VirtualNetworkInterfaceAgentClient
}

// newClient - creates a client session with the backend wssd agent
func newClient(subID string) (*client, error) {
	c, err := wssdclient.GetVirtualNetworkInterfaceClient(subID)
	if err != nil {
		return nil, err
	}
	return &client{c}, nil
}

// Get
func (c *client) Get(ctx context.Context, group, name string) (VirtualNetworkInterface, error) {
	return c.VirtualNetworkInterfaceAgentClient.Invoke(ctx, group, name, "")
}

// CreateOrUpdate
func (c *client) CreateOrUpdate(ctx context.Context, name string, id string, sg VirtualNetworkInterface) (VirtualNetworkInterface, error) {
	f, err := c.VirtualNetworkInterfaceAgentClient.Invoke(ctx, group, name, sg)
	if err != nil {
		return VirtualNetworkInterface{}, err
	}

	err = f.WaitForCompletionRef(ctx, c.Client)
	if err != nil {
		return VirtualNetworkInterface{}, err
	}

	return f.Result(c.VirtualNetworkInterfacesClient)
}

// Delete methods invokes create or update on the client
func (c *client) Delete(ctx context.Context, name string, id string) error {
	return c.VirtualNetworkInterfaceAgentClient.Invoke(ctx, group, name, "")
}
