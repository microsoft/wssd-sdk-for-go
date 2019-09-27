// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package fs

import (
	"fmt"
	log "k8s.io/klog"
	"reflect"

	"github.com/microsoft/wssdagent/pkg/apis/config"
	"github.com/microsoft/wssdagent/pkg/errors"
	"github.com/microsoft/wssdagent/pkg/guid"
	"github.com/microsoft/wssdagent/pkg/store"

	pb "github.com/microsoft/wssdagent/rpc/security"
	"github.com/microsoft/wssdagent/services/security/keyvault"
	fman "github.com/microsoft/wssdagent/services/security/keyvault/common/fs"
	"github.com/microsoft/wssdagent/services/security/keyvault/secret/internal"
)

type client struct {
	config           *config.ChildAgentConfiguration
	keyvaultProvider keyvault.KeyVaultProvider
	store            *store.ConfigStore
}

func newClient() *client {
	cConfig := config.GetChildAgentConfiguration("Secret")
	return &client{
		store:            store.NewConfigStore(cConfig.DataStorePath, reflect.TypeOf(internal.SecretInternal{})),
		config:           cConfig,
		keyvaultProvider: keyvault.GetKeyVaultProvider(keyvault.FSSpec),
	}
}

func (c *client) newSecret(id string) *internal.SecretInternal {
	return internal.NewSecretInternal(id, c.config.DataStorePath)
}

func (c *client) validate(secDef *pb.Secret) (err error) {
	if secDef == nil || len(secDef.VaultName) == 0 {
		return errors.Wrapf(errors.InvalidInput, "Missing Vault Name")
	}
	return nil
}

// Create a Secret
func (c *client) Create(sec *pb.Secret) (*pb.Secret, error) {
	log.Infof("[Secret][Create] spec[%v]", sec)
	err := c.validate(sec)
	if err != nil {
		return nil, err
	}
	if len(sec.Id) == 0 {
		sec.Id = guid.NewGuid()
	}
	secinternal := c.newSecret(sec.Id)

	vaultId, err := c.getVaultIdFromName(&sec.VaultName)
	if err != nil {
		return nil, err
	}
	sec.VaultId = *vaultId

	manager := fman.GetInstance()
	err = manager.AddSecretToVault(*sec)
	if err != nil {
		return nil, err
	}

	// 3. Save the config to the store
	secinternal.Srt = sec
	err = c.store.Add(sec.Id, secinternal)
	if err != nil {
		return nil, err
	}

	return sec, nil
}

// Get a Secret specified by name
func (c *client) Get(secDef *pb.Secret) ([]*pb.Secret, error) {
	log.Infof("[Secret][Get] spec[%v]", secDef)
	secs := []*pb.Secret{}

	err := c.validate(secDef)
	if err != nil {
		return secs, err
	}

	if secDef == nil || len(secDef.Name) == 0 {
		secretsInternal, err := c.store.List()
		if err != nil {
			return nil, err
		}
		if *secretsInternal == nil || len(*secretsInternal) == 0 {
			return secs, nil
		}

		for _, srt := range *secretsInternal {
			srtInternal := srt.(*internal.SecretInternal)
			secs = append(secs, srtInternal.Srt)
		}
	} else {
		// Check the internal store
		_, err := c.getSecretInternalByName(secDef.Name)
		if err != nil {
			return nil, err
		}

		vaultId, err := c.getVaultIdFromName(&secDef.VaultName)
		if err != nil {
			return nil, err
		}
		secDef.VaultId = *vaultId

		log.Infof("[Secret][LOGFOR ID]", secDef)
		manager := fman.GetInstance()
		secretWithValue, err := manager.ShowSecretFromVault(*secDef)
		if err != nil {
			return nil, err
		}
		secs = append(secs, secretWithValue)
		/// QUERY VAULTS
	}

	return secs, nil
}

// Delete a Virtual Secret
func (c *client) Delete(sec *pb.Secret) error {
	log.Infof("[Secret][Delete] spec[%v]", sec)

	// Check the internal store
	secint, err := c.getSecretInternalByName(sec.Name)
	if err != nil {
		return err
	}

	secretToBeDeleted := secint.Srt
	manager := fman.GetInstance()
	err = manager.RemoveSecretFromVault(*secretToBeDeleted)
	if err != nil {
		return err
	}

	return c.store.Delete(secretToBeDeleted.Id)
}

func (c *client) getVaultIdFromName(vaultName *string) (*string, error) {
	if vaultName == nil || len(*vaultName) == 0 {
		return nil, errors.Wrapf(errors.InvalidInput, "Missing VaultName")
	}
	var queryVault []*pb.KeyVault
	queryVault = append(queryVault, &pb.KeyVault{
		Name: *vaultName,
	})
	vaults, err := c.keyvaultProvider.Get(queryVault)
	if err != nil {
		return nil, err
	}

	if len(vaults) == 0 {
		return nil, fmt.Errorf("Associated Vault [%s] was not found", vaultName)
	}

	// Because of uniqueness only one will match the name
	return &vaults[0].Id, nil
}

func (c *client) getSecretInternalByName(name string) (*internal.SecretInternal, error) {
	svList, err := c.store.List()
	if err != nil {
		return nil, err
	}

	if *svList == nil || len(*svList) == 0 {
		return nil, errors.NotFound
	}

	for _, sv := range *svList {
		svInt := sv.(*internal.SecretInternal)
		if svInt.Srt != nil && svInt.Srt.Name == name {
			return svInt, nil
		}
	}
	return nil, errors.NotFound
}
