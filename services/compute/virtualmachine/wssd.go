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

package virtualmachine

import (
	"context"
	"github.com/microsoft/wssd-sdk-for-go/services/compute"

	wssdclient "github.com/microsoft/wssdagent/rpc/client"
	wssdcompute "github.com/microsoft/wssdagent/rpc/compute"
)

type client struct {
	wssdcompute.VirtualMachineAgentClient
}

// newClient - creates a client session with the backend wssd agent
func newClient(subID string) (*client, error) {
	c, err := wssdclient.GetVirtualMachineClient(subID)
	if err != nil {
		return nil, err
	}
	return &client{c}, nil
}

// Get
func (c *client) Get(ctx context.Context, group, name string) (compute.VirtualMachine, error) {
	request := &wssdnetwork.VirtualMachineRequest{Operation: wssdnetwork.Operation_GET}
	response, err := c.VirtualMachineAgentClient.Invoke(ctx, request, nil)
	return nil, err
}

// CreateOrUpdate
func (c *client) CreateOrUpdate(ctx context.Context, name string, id string, sg compute.VirtualMachine) (compute.VirtualMachine, error) {
	request := &wssdnetwork.VirtualMachineRequest{Operation: wssdnetwork.Operation_POST}
	response, err := c.VirtualMachineAgentClient.Invoke(ctx, request, nil)
	return nil, err
}

// Delete methods invokes create or update on the client
func (c *client) Delete(ctx context.Context, name string, id string) error {
	request := &wssdnetwork.VirtualMachineRequest{Operation: wssdnetwork.Operation_DELETE}
	response, err := c.VirtualMachineAgentClient.Invoke(ctx, request, nil)
	return nil, err
}
