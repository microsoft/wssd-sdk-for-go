// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license

package key

import (
	"context"

	"github.com/microsoft/moc/pkg/auth"
	"github.com/microsoft/wssd-sdk-for-go/services/security"
	"github.com/microsoft/wssd-sdk-for-go/services/security/keyvault"
	"github.com/microsoft/wssd-sdk-for-go/services/security/keyvault/key/internal"
)

// Service interface
type Service interface {
	Get(context.Context, string, string) (*[]keyvault.Key, error)
	CreateOrUpdate(context.Context, *keyvault.Key) (*keyvault.Key, error)
	Delete(context.Context, *keyvault.Key) error
	RotateKey(context.Context, *keyvault.KeyOperationRequest) (*keyvault.KeyOperationResult, error)
	WrapKey(context.Context, *keyvault.KeyOperationRequest) (*keyvault.KeyOperationResult, error)
	UnwrapKey(context.Context, *keyvault.KeyOperationRequest) (*keyvault.KeyOperationResult, error)
}

// Client structure
type KeyClient struct {
	security.BaseClient
	internal Service
}

// NewClient method returns new client
func NewKeyClient(cloudFQDN string, authorizer auth.Authorizer) (*KeyClient, error) {
	c, err := internal.NewKeyClient(cloudFQDN, authorizer)
	if err != nil {
		return nil, err
	}

	return &KeyClient{internal: c}, nil
}

// Get methods invokes the client Get method
func (c *KeyClient) Get(ctx context.Context, name, vaultName string) (*[]keyvault.Key, error) {
	return c.internal.Get(ctx, name, vaultName)
}

// CreateOrUpdate methods invokes create or update on the client
func (c *KeyClient) CreateOrUpdate(ctx context.Context, key *keyvault.Key) (*keyvault.Key, error) {
	return c.internal.CreateOrUpdate(ctx, key)
}

// Delete methods invokes delete of the key resource
func (c *KeyClient) Delete(ctx context.Context, key *keyvault.Key) error {
	return c.internal.Delete(ctx, key)
}

// RotateKey methods invokes delete of the key resource
func (c *KeyClient) RotateKey(ctx context.Context, key *keyvault.Key) (*keyvault.KeyOperationResult, error) {
	keyOpRequest := &keyvault.KeyOperationRequest{Key: key}
	return c.internal.RotateKey(ctx, keyOpRequest)
}

// WrapKey methods invokes delete of the key resource
func (c *KeyClient) WrapKey(ctx context.Context, key *keyvault.Key, alg *keyvault.JSONWebKeyEncryptionAlgorithm) (*keyvault.KeyOperationResult, error) {
	keyOpRequest := &keyvault.KeyOperationRequest{Key: key, Algorithm: alg}
	return c.internal.WrapKey(ctx, keyOpRequest)
}

// UnwrapKey methods invokes delete of the key resource
func (c *KeyClient) UnwrapKey(ctx context.Context, key *keyvault.Key, alg *keyvault.JSONWebKeyEncryptionAlgorithm) (*keyvault.KeyOperationResult, error) {
	keyOpRequest := &keyvault.KeyOperationRequest{Key: key, Algorithm: alg}
	return c.internal.UnwrapKey(ctx, keyOpRequest)
}
