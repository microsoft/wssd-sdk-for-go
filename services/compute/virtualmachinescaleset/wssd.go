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

package virtualmachinescaleset

import (
	"context"

	"github.com/microsoft/wssd-sdk-for-go/services/compute"
	wssdclient "github.com/microsoft/wssdagent/rpc/client"
	wssdcompute "github.com/microsoft/wssdagent/rpc/compute"
)

type client struct {
	wssdcompute.VirtualMachineScaleSetAgentClient
}

// newClient - creates a client session with the backend wssd agent
func newVirtualMachineScaleSetClient(subID string) (*client, error) {
	c, err := wssdclient.GetVirtualMachineScaleSetClient(&subID)
	if err != nil {
		return nil, err
	}
	return &client{c}, nil
}

// Get
func (c *client) Get(ctx context.Context, name string) (compute.VirtualMachineScaleSet, error) {
	request := &wssdcompute.VirtualMachineScaleSetRequest{OperationType: wssdcompute.Operation_GET}
	_, err := c.VirtualMachineScaleSetAgentClient.Invoke(ctx, request, nil)
	return compute.VirtualMachineScaleSet{}, err
}

// CreateOrUpdate
func (c *client) CreateOrUpdate(ctx context.Context, name string, id string, sg compute.VirtualMachineScaleSet) (compute.VirtualMachineScaleSet, error) {
	request := &wssdcompute.VirtualMachineScaleSetRequest{OperationType: wssdcompute.Operation_POST}
	_, err := c.VirtualMachineScaleSetAgentClient.Invoke(ctx, request, nil)
	return compute.VirtualMachineScaleSet{}, err
}

// Delete methods invokes create or update on the client
func (c *client) Delete(ctx context.Context, name string, id string) error {
	request := &wssdcompute.VirtualMachineScaleSetRequest{OperationType: wssdcompute.Operation_DELETE}
	_, err := c.VirtualMachineScaleSetAgentClient.Invoke(ctx, request, nil)
	return err
}
