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
)

// Service interface
type Service interface {
	Get(context.Context, string) (storage.VirtualHardDisk, error)
	CreateOrUpdate(context.Context, string, string, storage.VirtualHardDisk) (storage.VirtualHardDisk, error)
	Delete(context.Context, string, string) error
}

// Client structure
type VirtualHardDiskClient struct {
	storage.BaseClient
	internal Service
}

// NewClient method returns new client
func NewVirtualHardDiskClient(cloudFQDN string) (*VirtualHardDiskClient, error) {
	c, err := newVirtualHardDiskClient(cloudFQDN)
	if err != nil {
		return nil, err
	}

	return &VirtualHardDiskClient{internal: c}, nil
}

// Get methods invokes the client Get method
func (c *VirtualHardDiskClient) Get(ctx context.Context, name string) (storage.VirtualHardDisk, error) {
	return c.internal.Get(ctx, name)
}

// CreateOrUpdate methods invokes create or update on the client
func (c *VirtualHardDiskClient) CreateOrUpdate(ctx context.Context, name string, id string, storage storage.VirtualHardDisk) (storage.VirtualHardDisk, error) {
	return c.internal.CreateOrUpdate(ctx, name, id, storage)
}

// Delete methods invokes delete of the storage resource
func (c *VirtualHardDiskClient) Delete(ctx context.Context, name string, id string) error {
	return c.internal.Delete(ctx, name, id)
}
