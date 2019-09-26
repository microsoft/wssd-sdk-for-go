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

package keyvault

import (
	"context"
	"github.com/microsoft/wssd-sdk-for-go/services/security"
)

// Service interface
type Service interface {
	Get(context.Context, string, string) (*[]security.KeyVault, error)
	CreateOrUpdate(context.Context, string, string, *security.KeyVault) (*security.KeyVault, error)
	Delete(context.Context, string, string) error
}

// Client structure
type KeyVaultClient struct {
	security.BaseClient
	internal Service
}

// NewClient method returns new client
func NewKeyVaultClient(cloudFQDN string) (*KeyVaultClient, error) {
	c, err := newKeyVaultClient(cloudFQDN)
	if err != nil {
		return nil, err
	}

	return &KeyVaultClient{internal: c}, nil
}

// Get methods invokes the client Get method
func (c *KeyVaultClient) Get(ctx context.Context, group, name string) (*[]security.KeyVault, error) {
	return c.internal.Get(ctx, group, name)
}

// CreateOrUpdate methods invokes create or update on the client
func (c *KeyVaultClient) CreateOrUpdate(ctx context.Context, group, name string, keyvault *security.KeyVault) (*security.KeyVault, error) {
	return c.internal.CreateOrUpdate(ctx, group, name, keyvault)
}

// Delete methods invokes delete of the keyvault resource
func (c *KeyVaultClient) Delete(ctx context.Context, group, name string) error {
	return c.internal.Delete(ctx, group, name)
}
