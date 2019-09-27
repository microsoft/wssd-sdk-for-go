// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package fs

import (
	log "k8s.io/klog"
	"reflect"

	"github.com/microsoft/wssdagent/pkg/apis/config"
	"github.com/microsoft/wssdagent/pkg/errors"
	"github.com/microsoft/wssdagent/pkg/guid"
	"github.com/microsoft/wssdagent/pkg/store"
	pb "github.com/microsoft/wssdagent/rpc/security"
	fman "github.com/microsoft/wssdagent/services/security/keyvault/common/fs"
	"github.com/microsoft/wssdagent/services/security/keyvault/internal"
)

type client struct {
	config *config.ChildAgentConfiguration
	store  *store.ConfigStore
}

func newClient() *client {
	cConfig := config.GetChildAgentConfiguration("KeyVault")
	return &client{
		store:  store.NewConfigStore(cConfig.DataStorePath, reflect.TypeOf(internal.KeyVaultInternal{})),
		config: cConfig,
	}
}

func (c *client) newKeyVault(id string) *internal.KeyVaultInternal {
	return internal.NewKeyVaultInternal(id, c.config.DataStorePath)
}

// Create a Key Vault
func (c *client) Create(vault *pb.KeyVault) (*pb.KeyVault, error) {
	log.Infof("[KeyVault][Create] spec[%v]", vault)
	if len(vault.Id) == 0 {
		vault.Id = guid.NewGuid()

	}
	vaultinternal := c.newKeyVault(vault.Id)

	manager := fman.GetInstance()
	err := manager.AddVault(*vault)
	if err != nil {
		return nil, err
	}

	vaultinternal.Skv = vault

	// 3. Save the config to the store
	c.store.Add(vault.Id, vaultinternal)

	return vault, nil
}

// Get a vault specified by name
func (c *client) Get(vaultDef *pb.KeyVault) ([]*pb.KeyVault, error) {
	log.Infof("[KeyVault][Get] spec[%v]", vaultDef)
	vaults := []*pb.KeyVault{}
	vaultName := ""
	if vaultDef != nil {
		vaultName = vaultDef.Name
	}
	if len(vaultName) == 0 {
		vaultsInternal, err := c.store.List()
		if err != nil {
			return nil, err
		}

		if *vaultsInternal == nil || len(*vaultsInternal) == 0 {
			return nil, nil
		}

		for _, kv := range *vaultsInternal {
			vaultInternal := kv.(*internal.KeyVaultInternal)
			vaults = append(vaults, vaultInternal.Skv)
		}
	} else {
		vaultInternal, err := c.getKeyVaultInternalByName(vaultName)
		if err != nil {
			return nil, err
		}
		vaults = append(vaults, vaultInternal.Skv)
	}

	return vaults, nil
}

// Delete a KeyVault
func (c *client) Delete(vault *pb.KeyVault) error {
	log.Infof("[KeyVault][Delete] spec[%v]", vault)

	// Check the internal store
	vaultint, err := c.getKeyVaultInternalByName(vault.Name)
	if err != nil {
		return err
	}
	// Because of uniqueness of name we know that there will only be one
	vaultToBeDeleted := vaultint.Skv

	manager := fman.GetInstance()
	err = manager.RemoveVault(*vaultToBeDeleted)
	if err != nil {
		return err
	}

	return c.store.Delete(vaultToBeDeleted.Id)
}

func (c *client) getKeyVaultInternalByName(name string) (*internal.KeyVaultInternal, error) {
	svList, err := c.store.List()
	if err != nil {
		return nil, err
	}

	if *svList == nil || len(*svList) == 0 {
		return nil, errors.NotFound
	}

	for _, sv := range *svList {
		svInt := sv.(*internal.KeyVaultInternal)
		if svInt.Skv.Name == name {
			return svInt, nil
		}
	}
	return nil, errors.NotFound
}
