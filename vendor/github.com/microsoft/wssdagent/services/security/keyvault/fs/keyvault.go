// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package fs

import (
	pb "github.com/microsoft/wssdagent/rpc/security"
)

type KeyVaultProvider struct {
	client *Client
}

func NewKeyVaultProvider() *KeyVaultProvider {
	return &KeyVaultProvider{
		client: NewClient(),
	}
}

func (svp *KeyVaultProvider) Get(keyvaults []*pb.KeyVault) ([]*pb.KeyVault, error) {
	return svp.client.Get(keyvaults)
}

func (svp *KeyVaultProvider) CreateOrUpdate(keyvaults []*pb.KeyVault) ([]*pb.KeyVault, error) {
	return svp.client.Create(keyvaults)
}

func (svp *KeyVaultProvider) Delete(keyvaults []*pb.KeyVault) error {
	return svp.client.Delete(keyvaults)
}
