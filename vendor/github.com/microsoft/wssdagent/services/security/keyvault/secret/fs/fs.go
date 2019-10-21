// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package fs

import (
	pb "github.com/microsoft/wssdagent/rpc/security"
	fman "github.com/microsoft/wssdagent/services/security/keyvault/common/fs"
	"github.com/microsoft/wssdagent/services/security/keyvault/secret/internal"
)

type Client struct {
}

func NewClient() *Client {
	return &Client{}
}

// Create a Secret
func (c *Client) CreateSecret(secInternal *internal.SecretInternal) (err error) {
	sec := secInternal.Entity

	manager := fman.GetInstance()
	err = manager.AddSecretToVault(*sec)
	if err != nil {
		return
	}

	// Nil out the value so it is not placed in the store
	sec.Value = nil
	secInternal.Entity = sec
	
	return
}

// GetDecryptedSecret Get a Decrypted Secret
func (c *Client) GetDecryptedSecret(secInternal *internal.SecretInternal) (decSec *pb.Secret, err error) {
	secDef := secInternal.Entity
	manager := fman.GetInstance()
	decSec, err = manager.ShowSecretFromVault(*secDef)
	return
}

// CleanupSecret Delete a Virtual Secret
func (c *Client) CleanupSecret(secInternal *internal.SecretInternal) (err error) {
	secretToBeDeleted := secInternal.Entity
	manager := fman.GetInstance()
	err = manager.RemoveSecretFromVault(*secretToBeDeleted)
	if err != nil {
		return
	}
	return
}

// HasSecret
func (c *Client) HasSecret(secInternal *internal.SecretInternal) bool {
	manager := fman.GetInstance()
	return manager.IsValidSecret(secInternal.Entity)
}
