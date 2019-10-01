// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license

package secret

import (
	"context"
	"github.com/microsoft/wssd-sdk-for-go/services/security"
	"github.com/microsoft/wssd-sdk-for-go/services/security/keyvault"
	"github.com/microsoft/wssd-sdk-for-go/services/security/keyvault/secret/internal"
)

// Service interface
type Service interface {
	Get(context.Context, string, string, string) (*[]keyvault.Secret, error)
	CreateOrUpdate(context.Context, string, string, *keyvault.Secret) (*keyvault.Secret, error)
	Delete(context.Context, string, string, string) error
}

// Client structure
type SecretClient struct {
	security.BaseClient
	internal Service
}

// NewClient method returns new client
func NewSecretClient(cloudFQDN string) (*SecretClient, error) {
	c, err := internal.NewSecretClient(cloudFQDN)
	if err != nil {
		return nil, err
	}

	return &SecretClient{internal: c}, nil
}

// Get methods invokes the client Get method
func (c *SecretClient) Get(ctx context.Context, group, name, vaultName string) (*[]keyvault.Secret, error) {
	return c.internal.Get(ctx, group, name, vaultName)
}

// CreateOrUpdate methods invokes create or update on the client
func (c *SecretClient) CreateOrUpdate(ctx context.Context, group, name string, sec *keyvault.Secret) (*keyvault.Secret, error) {
	return c.internal.CreateOrUpdate(ctx, group, name, sec)
}

// Delete methods invokes delete of the keyvault resource
func (c *SecretClient) Delete(ctx context.Context, group, name, vaultName string) error {
	return c.internal.Delete(ctx, group, name, vaultName)
}
