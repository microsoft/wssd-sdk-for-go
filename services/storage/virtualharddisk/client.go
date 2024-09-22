// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license

package virtualharddisk

import (
	"context"

	"github.com/microsoft/moc/pkg/auth"
	"github.com/microsoft/wssd-sdk-for-go/services/storage"
	"github.com/microsoft/wssd-sdk-for-go/services/storage/virtualharddisk/internal"
)

// Service interface
type Service interface {
	Get(context.Context, string, string) (*[]storage.VirtualHardDisk, error)
	CreateOrUpdate(context.Context, string, string, *storage.VirtualHardDisk, string) (*storage.VirtualHardDisk, error)
	Delete(context.Context, string, string) error
	Hydrate(context.Context, string, string, *storage.VirtualHardDisk) (*storage.VirtualHardDisk, error)
}

// Client structure
type VirtualHardDiskClient struct {
	storage.BaseClient
	internal Service
}

// NewClient method returns new client
func NewVirtualHardDiskClient(cloudFQDN string, authorizer auth.Authorizer) (*VirtualHardDiskClient, error) {
	c, err := internal.NewVirtualHardDiskClient(cloudFQDN, authorizer)
	if err != nil {
		return nil, err
	}

	return &VirtualHardDiskClient{internal: c}, nil
}

// Get methods invokes the client Get method
func (c *VirtualHardDiskClient) Get(ctx context.Context, container, name string) (*[]storage.VirtualHardDisk, error) {
	return c.internal.Get(ctx, container, name)
}

// CreateOrUpdate methods invokes create or update on the client
func (c *VirtualHardDiskClient) CreateOrUpdate(ctx context.Context, container, name string, storage *storage.VirtualHardDisk, path string) (*storage.VirtualHardDisk, error) {
	return c.internal.CreateOrUpdate(ctx, container, name, storage, path)
}

// Delete methods invokes delete of the storage resource
func (c *VirtualHardDiskClient) Delete(ctx context.Context, container, name string) error {
	return c.internal.Delete(ctx, container, name)
}

// The interface for the hydrate call takes the container name and the name of the disk file.
// Ultimately, we need the full path on disk to the disk file which we assemble from the path of the container plus the file name of the disk.
// (e.g. "C:\ClusterStorage\Userdata_1\abc123" for the container path and "my_disk.vhd" for the disk name)
func (c *VirtualHardDiskClient) Hydrate(ctx context.Context, container, name string, vhdDef *storage.VirtualHardDisk) (*storage.VirtualHardDisk, error) {
	return c.internal.Hydrate(ctx, container, name, vhdDef)
}
