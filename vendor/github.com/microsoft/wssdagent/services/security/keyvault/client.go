// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package keyvault

import (
	"context"
	"github.com/microsoft/wssdagent/pkg/apis/config"
	"github.com/microsoft/wssdagent/pkg/channel"
	ch "github.com/microsoft/wssdagent/pkg/channel/vault"
	"github.com/microsoft/wssdagent/pkg/errors"
	"github.com/microsoft/wssdagent/pkg/guid"
	"github.com/microsoft/wssdagent/pkg/marshal"
	"github.com/microsoft/wssdagent/pkg/store"
	"github.com/microsoft/wssdagent/pkg/trace"
	pb "github.com/microsoft/wssdagent/rpc/security"
	"github.com/microsoft/wssdagent/services/security/keyvault/internal"
	"reflect"
	"sync"

	"github.com/microsoft/wssdagent/services/security/keyvault/fs"
)

const (
	FSSpec = "fs"
)

type Service interface {
	CreateKeyVault(*internal.KeyVaultInternal) error
	CleanupKeyVault(*internal.KeyVaultInternal) error
}

type Client struct {
	internal Service
	store    *store.ConfigStore
	config   *config.ChildAgentConfiguration
	mux      sync.Mutex
	// channel for sending messages to secret provider
	channel ch.Channel
}

func NewClient() *Client {
	cConfig := config.GetChildAgentConfiguration("KeyVault")
	c := &Client{
		store:   store.NewConfigStore(cConfig.DataStorePath, reflect.TypeOf(internal.KeyVaultInternal{})),
		config:  cConfig,
		channel: ch.MakeChannel(),
	}
	switch cConfig.ProviderSpec {
	case FSSpec:
	default:
		c.internal = fs.NewClient()
	}
	return c
}

func (c *Client) newKeyVault(vault *pb.KeyVault) *internal.KeyVaultInternal {
	return internal.NewKeyVaultInternal(guid.NewGuid(), c.config.DataStorePath, vault)
}

// GetChannel
func (c *Client) GetChannel() *ch.Channel {
	return &c.channel
}

// GetChannel
func (c *Client) GetDataStorePath() string {
	return c.store.GetPath()
}

// Create or Update the specified virtual security(s)
func (c *Client) Create(ctx context.Context, vaultDef *pb.KeyVault) (newvault *pb.KeyVault, err error) {
	ctx, span := trace.NewSpan(ctx, "KeyVault", "Create", marshal.ToString(vaultDef))
	defer span.End(err)

	err = c.Validate(ctx, vaultDef)
	if err != nil {
		return
	}
	vaultinternal := c.newKeyVault(vaultDef)

	err = c.internal.CreateKeyVault(vaultinternal)
	if err != nil {
		return
	}
	newvault = vaultinternal.Entity

	err = c.store.Add(vaultinternal.Id, vaultinternal)

	return

}

// Get all/selected HCS virtual security(s)
func (c *Client) Get(ctx context.Context, securityDef *pb.KeyVault) (vaults []*pb.KeyVault, err error) {
	ctx, span := trace.NewSpan(ctx, "KeyVault", "Get", marshal.ToString(securityDef))
	defer span.End(err)

	c.mux.Lock()
	defer c.mux.Unlock()

	vaultName := ""
	if securityDef != nil {
		vaultName = securityDef.Name
	}

	vaultsint, err := c.store.ListFilter("Name", vaultName)
	if err != nil {
		return
	}

	for _, val := range *vaultsint {
		vaultint := val.(*internal.KeyVaultInternal)
		vaults = append(vaults, vaultint.Entity)
	}

	return
}

// Delete the specified virtual security(s)
func (c *Client) Delete(ctx context.Context, securityDef *pb.KeyVault) (err error) {
	ctx, span := trace.NewSpan(ctx, "KeyVault", "Delete", marshal.ToString(securityDef))
	defer span.End(err)

	c.mux.Lock()
	defer c.mux.Unlock()

	vaultinternal, err := c.getKeyVaultInternal(securityDef.Name)
	if err != nil {
		return
	}

	err = c.notifySecretProvider(ctx, vaultinternal.Name, channel.DeleteOperation)
	if err != nil {
		// Log this error and continue
		// return
	}

	err = c.internal.CleanupKeyVault(vaultinternal)
	if err != nil {
		// Log this error and continue
		// return
	}

	err = c.store.Delete(vaultinternal.Id)
	return
}

// Validate
func (c *Client) Validate(ctx context.Context, securityDef *pb.KeyVault) (err error) {
	ctx, span := trace.NewSpan(ctx, "KeyVault", "Validate", marshal.ToString(securityDef))
	defer span.End(err)

	err = nil

	if securityDef == nil {
		err = errors.Wrapf(errors.InvalidInput, "Input group definition is nil")
		return
	}

	_, err = c.getKeyVaultInternal(securityDef.Name)
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

func (c *Client) getKeyVaultInternal(name string) (*internal.KeyVaultInternal, error) {
	vaultsint, err := c.store.ListFilter("Name", name)
	if err != nil {
		return nil, err
	}
	if *vaultsint == nil || len(*vaultsint) == 0 {
		return nil, errors.NotFound
	}

	return (*vaultsint)[0].(*internal.KeyVaultInternal), nil
}

func (c *Client) pruneStore() {
	vaultsint, err := c.store.List()
	if err != nil {
		return
	}
	if *vaultsint == nil || len(*vaultsint) == 0 {
		return
	}

	for _, _ = range *vaultsint {

	}
}

func (c *Client) notifySecretProvider(ctx context.Context, vaultName string, oper channel.OperationType) (err error) {

	c.channel.Notify <- ch.MakeNotificationData(vaultName, oper)

	// Wait for response
	if <-c.channel.Result > 0 {
		err = errors.Failed
	}
	return
}
