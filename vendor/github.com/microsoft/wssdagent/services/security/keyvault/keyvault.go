// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package keyvault

import (
	"context"
	ch "github.com/microsoft/wssdagent/pkg/channel/vault"
	"github.com/microsoft/wssdagent/pkg/errors"
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

func (vaultProv *KeyVaultProvider) Get(ctx context.Context, vaults []*pb.KeyVault) ([]*pb.KeyVault, error) {
	newvaults := []*pb.KeyVault{}
	if len(vaults) == 0 {
		// Get Everything
		return vaultProv.client.Get(ctx, nil)
	}

	// Get only requested vaults
	for _, vault := range vaults {
		newvault, err := vaultProv.client.Get(ctx, vault)
		if err != nil {
			return newvaults, err
		}
		newvaults = append(newvaults, newvault[0])
	}
	return newvaults, nil
}

func (vaultProv *KeyVaultProvider) CreateOrUpdate(ctx context.Context, vaults []*pb.KeyVault) ([]*pb.KeyVault, error) {
	newvaults := []*pb.KeyVault{}
	for _, vault := range vaults {
		newvault, err := vaultProv.client.Create(ctx, vault)
		if err != nil {
			if err != errors.AlreadyExists {
				vaultProv.client.Delete(ctx, vault)
			}
			return newvaults, err
		}
		newvaults = append(newvaults, newvault)
	}

	return newvaults, nil
}

func (vaultProv *KeyVaultProvider) Delete(ctx context.Context, vaults []*pb.KeyVault) error {
	for _, vault := range vaults {
		err := vaultProv.client.Delete(ctx, vault)
		if err != nil {
			return err
		}
	}

	return nil
}

// GetVaultByName
func (vaultProv *KeyVaultProvider) GetVaultByName(ctx context.Context, vaultName string) (*pb.KeyVault, error) {
	vaults, err := vaultProv.Get(ctx, []*pb.KeyVault{&pb.KeyVault{Name: vaultName}})
	if err != nil {
		return nil, err
	}
	if len(vaults) == 0 {
		return nil, errors.NotFound
	}

	return vaults[0], nil

}

// GetChannel
func (vaultProv *KeyVaultProvider) GetChannel() *ch.Channel {
	return vaultProv.client.GetChannel()
}

// GetChannel
func (vaultProv *KeyVaultProvider) GetDataStorePath() string {
	return vaultProv.client.GetDataStorePath()
}
