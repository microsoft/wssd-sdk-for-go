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

package virtualharddisk

import (
	"context"
	"github.com/microsoft/wssd-sdk-for-go/services/storage"

	wssdclient "github.com/microsoft/wssdagent/rpc/client"
	wssdstorage "github.com/microsoft/wssdagent/rpc/storage"
)

type client struct {
	wssdstorage.VirtualHardDiskAgentClient
}

// newClient - creates a client session with the backend wssd agent
func newVirtualHardDiskClient(subID string) (*client, error) {
	c, err := wssdclient.GetVirtualHardDiskClient(&subID)
	if err != nil {
		return nil, err
	}
	return &client{c}, nil
}

// Get
func (c *client) Get(ctx context.Context, name string) (storage.VirtualHardDisk, error) {
	request := &wssdstorage.VirtualHardDiskRequest{OperationType: wssdstorage.Operation_GET}
	_, err := c.VirtualHardDiskAgentClient.Invoke(ctx, request, nil)
	return storage.VirtualHardDisk{}, err
}

// CreateOrUpdate
func (c *client) CreateOrUpdate(ctx context.Context, name string, id string, sg storage.VirtualHardDisk) (storage.VirtualHardDisk, error) {
	request := &wssdstorage.VirtualHardDiskRequest{OperationType: wssdstorage.Operation_POST}
	_, err := c.VirtualHardDiskAgentClient.Invoke(ctx, request, nil)
	return storage.VirtualHardDisk{}, err
}

// Delete methods invokes create or update on the client
func (c *client) Delete(ctx context.Context, name string, id string) error {
	request := &wssdstorage.VirtualHardDiskRequest{OperationType: wssdstorage.Operation_DELETE}
	_, err := c.VirtualHardDiskAgentClient.Invoke(ctx, request, nil)
	return err
}
