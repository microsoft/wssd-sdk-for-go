// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package fs

import (
	"fmt"
	log "k8s.io/klog"
	"reflect"

	"github.com/microsoft/wssdagent/pkg/guid"
	"github.com/microsoft/wssdagent/pkg/apis/config"
	"github.com/microsoft/wssdagent/pkg/store"

	pb "github.com/microsoft/wssdagent/rpc/security"
	fman "github.com/microsoft/wssdagent/services/security/keyvault/common/fs"
	"github.com/microsoft/wssdagent/services/security/keyvault/secret/internal"
	"github.com/microsoft/wssdagent/services/security/keyvault"
)

type client struct {
	config *config.ChildAgentConfiguration
	store  *store.ConfigStore
}

func newClient() *client {
	cConfig := config.GetChildAgentConfiguration("Secret")
	return &client{
		store:  store.NewConfigStore(cConfig.DataStorePath, reflect.TypeOf(internal.SecretInternal{})),
		config: cConfig,
	}
}

func (c *client) newSecret(id string) *internal.SecretInternal {
	return internal.NewSecretInternal(id, c.config.DataStorePath)
}

// Create a Secret
func (c *client) Create(sec *pb.Secret) (*pb.Secret, error) {
	log.Infof("[Secret][Create] spec[%v]", sec)
	if len(sec.Id) == 0 {
		sec.Id = guid.NewGuid()

	}
	secinternal := c.newSecret(sec.Id)

	vaultId, err := getVaultIdFromName(&sec.VaultName)
	if err != nil {
		return nil, nil
	}
	sec.VaultId = *vaultId

	manager := fman.GetInstance()
	manager.AddSecretToVault(*sec)

	// 3. Save the config to the store
	c.store.Add(sec.Id, secinternal)

	return sec, nil
}

// Get a Secret specified by name
func (c *client) Get(secDef *pb.Secret) ([]*pb.Secret, error) {
	log.Infof("[Secret][Get] spec[%v]", secDef)

	secs := []*pb.Secret{}

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

		vaultId, err := getVaultIdFromName(&secDef.VaultName)
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

	secrets, err := c.Get(sec)
	if err != nil {
		return err
	}

	if len(secrets) == 0 {
		return fmt.Errorf("Secret [%s] was not found", sec.Name)
	}

	manager := fman.GetInstance()
	manager.RemoveSecretFromVault(*sec)

	return c.store.Delete(sec.Id)
}

func getVaultIdFromName(vaultName *string) (*string, error) {
	vaultManager := keyvault.GetKeyVaultProvider(keyvault.FSSpec)

	var queryVault []*pb.KeyVault
	queryVault = append(queryVault, &pb.KeyVault{
		Name: *vaultName,
	})
	vaults, err := vaultManager.Get(queryVault)
	if err != nil {
		return nil, err
	}

	if len(vaults) == 0 {
		return nil, fmt.Errorf("Associated Vault [%s] was not found", vaultName)
	}

	// Because of uniqueness only one will match the name
	return &vaults[0].Id, nil
}
