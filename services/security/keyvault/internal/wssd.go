// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license

package internal

import (
	"context"
	"fmt"
	"github.com/microsoft/wssd-sdk-for-go/pkg/auth"
	"github.com/microsoft/wssd-sdk-for-go/services/security"

	wssdclient "github.com/microsoft/wssd-sdk-for-go/pkg/client"
	wssdsecurity "github.com/microsoft/wssdagent/rpc/security"
	wssdcommonproto "github.com/microsoft/wssdagent/rpc/common"
	log "k8s.io/klog"
)

type client struct {
	wssdsecurity.KeyVaultAgentClient
}

// NewKeyVaultClientN- creates a client session with the backend wssd agent
func NewKeyVaultClient(subID string, authorizer auth.Authorizer) (*client, error) {
	c, err := wssdclient.GetKeyVaultClient(&subID, authorizer)
	if err != nil {
		return nil, err
	}
	return &client{c}, nil
}

// Get
func (c *client) Get(ctx context.Context, group, name string) (*[]security.KeyVault, error) {
	request := getKeyVaultRequest(wssdcommonproto.Operation_GET, name, nil)
	response, err := c.KeyVaultAgentClient.Invoke(ctx, request)
	if err != nil {
		return nil, err
	}
	return getKeyVaultsFromResponse(response), nil
}

// CreateOrUpdate
func (c *client) CreateOrUpdate(ctx context.Context, group, name string, sg *security.KeyVault) (*security.KeyVault, error) {
	request := getKeyVaultRequest(wssdcommonproto.Operation_POST, name, sg)
	response, err := c.KeyVaultAgentClient.Invoke(ctx, request)
	if err != nil {
		log.Errorf("[KeyVault] Create failed with error %v", err)
		return nil, err
	}

	vault := getKeyVaultsFromResponse(response)

	if len(*vault) == 0 {
		return nil, fmt.Errorf("[KeyVault][Create] Unexpected error: Creating a security returned no result")
	}

	return &((*vault)[0]), err
}

// Delete methods invokes create or update on the client
func (c *client) Delete(ctx context.Context, group, name string) error {
	vault, err := c.Get(ctx, group, name)
	if err != nil {
		return err
	}
	if len(*vault) == 0 {
		return fmt.Errorf("Keyvault [%s] not found", name)
	}

	request := getKeyVaultRequest(wssdcommonproto.Operation_DELETE, name, &(*vault)[0])
	_, err = c.KeyVaultAgentClient.Invoke(ctx, request)
	return err
}

func getKeyVaultsFromResponse(response *wssdsecurity.KeyVaultResponse) *[]security.KeyVault {
	vaults := []security.KeyVault{}
	for _, keyvaults := range response.GetKeyVaults() {
		vaults = append(vaults, *(getKeyVault(keyvaults)))
	}

	return &vaults
}

func getKeyVaultRequest(opType wssdcommonproto.Operation, name string, vault *security.KeyVault) *wssdsecurity.KeyVaultRequest {
	request := &wssdsecurity.KeyVaultRequest{
		OperationType: opType,
		KeyVaults:     []*wssdsecurity.KeyVault{},
	}
	if vault != nil {
		request.KeyVaults = append(request.KeyVaults, getWssdKeyVault(vault))
	} else if len(name) > 0 {
		request.KeyVaults = append(request.KeyVaults,
			&wssdsecurity.KeyVault{
				Name: name,
			})
	}
	return request
}

func getKeyVault(vault *wssdsecurity.KeyVault) *security.KeyVault {
	return &security.KeyVault{
		ID:   &vault.Id,
		Name: &vault.Name,
		//	Source : &vault.Source,
	}
}

func getWssdKeyVault(vault *security.KeyVault) *wssdsecurity.KeyVault {
	return &wssdsecurity.KeyVault{
		Name:    *vault.Name,
		Secrets: []*wssdsecurity.Secret{},
	}
}
