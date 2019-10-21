// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package secret

import (
	"context"
	"github.com/microsoft/wssdagent/pkg/apis/config"
	"github.com/microsoft/wssdagent/pkg/channel"
	"github.com/microsoft/wssdagent/pkg/errors"
	"github.com/microsoft/wssdagent/pkg/guid"
	"github.com/microsoft/wssdagent/pkg/marshal"
	"github.com/microsoft/wssdagent/pkg/store"
	"github.com/microsoft/wssdagent/pkg/trace"
	pb "github.com/microsoft/wssdagent/rpc/security"
	"github.com/microsoft/wssdagent/services/security/keyvault/secret/internal"
	"reflect"
	"sync"
	"time"

	"github.com/microsoft/wssdagent/services/security/keyvault"
	"github.com/microsoft/wssdagent/services/security/keyvault/secret/fs"
)

const (
	FSSpec = "fs"
)

type Service interface {
	CreateSecret(*internal.SecretInternal) error
	GetDecryptedSecret(*internal.SecretInternal) (*pb.Secret, error)
	HasSecret(*internal.SecretInternal) bool
	CleanupSecret(*internal.SecretInternal) error
}

type Client struct {
	internal         Service
	store            *store.ConfigStore
	config           *config.ChildAgentConfiguration
	keyvaultProvider *keyvault.KeyVaultProvider
	mux              sync.Mutex
}

func NewClient() *Client {
	cConfig := config.GetChildAgentConfiguration("Secret")
	keyvaultProvider := keyvault.GetKeyVaultProvider()
	c := &Client{
		// TODO: Move the secret store under the correct vaultId
		store:            store.NewConfigStore(cConfig.DataStorePath, reflect.TypeOf(internal.SecretInternal{})),
		config:           cConfig,
		keyvaultProvider: keyvaultProvider,
	}
	switch cConfig.ProviderSpec {
	case FSSpec:
	default:
		c.internal = fs.NewClient()
	}

	go c.monitorVaultNotifications()

	return c
}

func (c *Client) monitorVaultNotifications() {
	print("monitorVaultNotifications")
	vaultChannel := c.keyvaultProvider.GetChannel()

	for {
		result := 0

		notifyMessage, ok := <-vaultChannel.Notify
		if !ok {
			// Channel has been closed
			break
		}

		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		if notifyMessage.Operation == channel.DeleteOperation && len(notifyMessage.Name) > 0 {
			err := c.onVaultDelete(ctx, notifyMessage.Name)
			if err != nil {
				result = 1
			}
		}

		vaultChannel.Result <- result
	}
}

func (c *Client) newSecret(sec *pb.Secret) *internal.SecretInternal {
	return internal.NewSecretInternal(guid.NewGuid(), c.config.DataStorePath, sec)
}

// Create or Update the specified virtual security(s)
func (c *Client) Create(ctx context.Context, secDef *pb.Secret) (newsec *pb.Secret, err error) {
	ctx, span := trace.NewSpan(ctx, "Secret", "Create", marshal.ToString(secDef))
	defer span.End(err)

	err = c.Validate(ctx, secDef)
	if err != nil {
		if err == errors.AlreadyExists {
			newsec, err = c.update(ctx, secDef)
			return
		}
		return
	}
	secinternal := c.newSecret(secDef)

	err = c.internal.CreateSecret(secinternal)
	if err != nil {
		return
	}
	newsec = secinternal.Entity

	err = c.store.Add(secinternal.Id, secinternal)

	return

}

// update an existing secret with new value
func (c *Client) update(ctx context.Context, secDef *pb.Secret) (newsec *pb.Secret, err error) {
	ctx, span := trace.NewSpan(ctx, "Secret", "Update", marshal.ToString(secDef))
	defer span.End(err)

	c.mux.Lock()
	defer c.mux.Unlock()

	secinternal, err := c.getSecretInternal(secDef.Name)
	if err != nil {
		return
	}

	// Remove existing secret
	err = c.internal.CleanupSecret(secinternal)
	if err != nil {
		// Log this error and continue
		// return
	}

	// create a new secret
	secDef.Id = secinternal.Id
	secinternal.Entity = secDef
	err = c.internal.CreateSecret(secinternal)
	if err != nil {
		return
	}
	newsec = secinternal.Entity

	return

}

// Get all/selected HCS virtual security(s)
func (c *Client) Get(ctx context.Context, securityDef *pb.Secret) (secs []*pb.Secret, err error) {
	ctx, span := trace.NewSpan(ctx, "Secret", "Get", marshal.ToString(securityDef))
	defer span.End(err)

	c.mux.Lock()
	defer c.mux.Unlock()

	secName := ""
	vaultName := ""
	if securityDef != nil {
		secName = securityDef.Name
		vaultName = securityDef.VaultName
	}

	if len(vaultName) == 0 {
		err = errors.Wrapf(errors.InvalidInput, "vault-name is missing")
		return
	}
	_, err = c.keyvaultProvider.GetVaultByName(ctx, vaultName)
	if err != nil {
		return
	}

	secsint, err := c.store.ListFilter("Name", secName)
	if err != nil {
		return
	}

	for _, val := range *secsint {
		secint := val.(*internal.SecretInternal)

		if secint.VaultName != vaultName {
			continue
		}
		// Decrpt only for selective show of a secret
		if len(secName) > 0 {
			decryptedSecret, err1 := c.internal.GetDecryptedSecret(secint)
			if err1 != nil {
				err = err1
				return
			}
			secs = append(secs, decryptedSecret)
			break
		} else {
			secs = append(secs, secint.Entity)
		}
	}

	return
}

// Delete the specified virtual security(s)
func (c *Client) Delete(ctx context.Context, securityDef *pb.Secret) (err error) {
	ctx, span := trace.NewSpan(ctx, "Secret", "Delete", marshal.ToString(securityDef))
	defer span.End(err)

	c.mux.Lock()
	defer c.mux.Unlock()

	secinternal, err := c.getSecretInternal(securityDef.Name)
	if err != nil {
		return
	}

	err = c.internal.CleanupSecret(secinternal)
	if err != nil {
		// Log this error and continue
		// return
	}

	err = c.store.Delete(secinternal.Id)
	return
}

// Validate
func (c *Client) Validate(ctx context.Context, securityDef *pb.Secret) (err error) {
	ctx, span := trace.NewSpan(ctx, "Secret", "Validate", marshal.ToString(securityDef))
	defer span.End(err)

	err = nil

	if securityDef == nil {
		err = errors.Wrapf(errors.InvalidInput, "Input group definition is nil")
		return
	}

	if len(securityDef.VaultName) == 0 {
		err = errors.Wrapf(errors.InvalidInput, "Missing Vault Name")
		return
	}
	// Validate Vault name
	vault, err := c.keyvaultProvider.GetVaultByName(ctx, securityDef.VaultName)
	if err != nil {
		return
	}

	// Update vault Id
	securityDef.VaultId = vault.Id
	_, err = c.getSecretInternal(securityDef.Name)
	if err != nil && err == errors.NotFound {
		err = nil
	} else {
		err = errors.AlreadyExists
	}

	if err != nil {
		return
	}

	return
}

// Callback for vault deletion
func (c *Client) onVaultDelete(ctx context.Context, vaultName string) (err error) {
	ctx, span := trace.NewSpan(ctx, "Secret", "onVaultDelete", vaultName)
	defer span.End(err)

	secrets, err := c.getSecretInternalByVaultName(vaultName)
	if err != nil {
		return err
	}

	for _, secret := range secrets {
		err = c.Delete(ctx, secret.Entity)
		if err != nil {
			// Continue, just log the error
		}
	}
	return
}

func (c *Client) getSecretInternal(name string) (*internal.SecretInternal, error) {
	secsint, err := c.store.ListFilter("Name", name)
	if err != nil {
		return nil, err
	}
	if *secsint == nil || len(*secsint) == 0 {
		return nil, errors.NotFound
	}

	return (*secsint)[0].(*internal.SecretInternal), nil
}

func (c *Client) getSecretInternalByVaultName(vname string) ([]*internal.SecretInternal, error) {
	secrets := []*internal.SecretInternal{}
	secsint, err := c.store.ListFilterMany("VaultName", vname)
	if err != nil {
		return nil, err
	}
	if *secsint == nil || len(*secsint) == 0 {
		return nil, errors.NotFound
	}

	for _, val := range *secsint {
		secretInt := val.(*internal.SecretInternal)
		secrets = append(secrets, secretInt)
	}

	return secrets, nil
}

func (c *Client) pruneStore() {
	secsint, err := c.store.List()
	if err != nil {
		return
	}
	if *secsint == nil || len(*secsint) == 0 {
		return
	}

	for _, _ = range *secsint {

	}
}
