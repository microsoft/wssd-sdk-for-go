// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package fs

import (
	fman "github.com/microsoft/wssdagent/services/security/keyvault/common/fs"
	"github.com/microsoft/wssdagent/services/security/keyvault/internal"
)

type Client struct {
}

func NewClient() *Client {
	return &Client{}
}

// Create a Key Vault
func (c *Client) CreateKeyVault(vaultInternal *internal.KeyVaultInternal) (err error) {
	vault := vaultInternal.Entity
	manager := fman.GetInstance()
	err = manager.AddVault(*vault)
	if err != nil {
		return
	}

	return
}

// Delete a KeyVault
func (c *Client) CleanupKeyVault(vaultInternal *internal.KeyVaultInternal) (err error) {
	vaultToBeDeleted := vaultInternal.Entity

	manager := fman.GetInstance()
	err = manager.RemoveVault(*vaultToBeDeleted)
	if err != nil {
		return
	}

	return
}
