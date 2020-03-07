// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license

package virtualharddisk

import (
	"context"
	"github.com/microsoft/wssd-sdk-for-go/pkg/auth"
	"github.com/microsoft/wssd-sdk-for-go/services/storage"
	"github.com/microsoft/wssd-sdk-for-go/services/storage/virtualharddisk/internal"
)

// Service interface
type Service interface {
	Get(context.Context, string, string) (*[]storage.VirtualHardDisk, error)
	CreateOrUpdate(context.Context, string, string, *storage.VirtualHardDisk) (*storage.VirtualHardDisk, error)
	Delete(context.Context, string, string) error
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
func (c *VirtualHardDiskClient) CreateOrUpdate(ctx context.Context, container, name string, storage *storage.VirtualHardDisk) (*storage.VirtualHardDisk, error) {
	return c.internal.CreateOrUpdate(ctx, container, name, storage)
}

// Delete methods invokes delete of the storage resource
func (c *VirtualHardDiskClient) Delete(ctx context.Context, container, name string) error {
	return c.internal.Delete(ctx, container, name)
}
