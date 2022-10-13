// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license

package sharedfolder

import (
	"context"

	"github.com/microsoft/moc/pkg/auth"
	"github.com/microsoft/wssd-sdk-for-go/services/storage"
	"github.com/microsoft/wssd-sdk-for-go/services/storage/sharedfolder/internal"
)

// Service interface
type Service interface {
	Get(context.Context, string) (*[]storage.SharedFolder, error)
	CreateOrUpdate(context.Context, string, *storage.SharedFolder) (*storage.SharedFolder, error)
	Delete(context.Context, string) error
}

// Client structure
type SharedFolderClient struct {
	storage.BaseClient
	internal Service
}

// NewClient method returns new client
func NewSharedFolderClient(cloudFQDN string, authorizer auth.Authorizer) (*SharedFolderClient, error) {
	c, err := internal.NewSharedFolderClient(cloudFQDN, authorizer)
	if err != nil {
		return nil, err
	}

	return &SharedFolderClient{internal: c}, nil
}

// Get method invokes the client Get method
func (c *SharedFolderClient) Get(ctx context.Context, name string) (*[]storage.SharedFolder, error) {
	return c.internal.Get(ctx, name)
}

// CreateOrUpdate method invokes create or update on the client
func (c *SharedFolderClient) CreateOrUpdate(ctx context.Context, name string, storage *storage.SharedFolder) (*storage.SharedFolder, error) {
	return c.internal.CreateOrUpdate(ctx, name, storage)
}

// Delete method invokes delete on the client
func (c *SharedFolderClient) Delete(ctx context.Context, name string) error {
	return c.internal.Delete(ctx, name)
}
