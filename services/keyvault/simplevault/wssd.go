// Copyright 2019 (c) Microsoft and contributors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package simplevault

import (
	"fmt"
	"context"
	"github.com/microsoft/wssd-sdk-for-go/services/keyvault"

	wssdclient "github.com/microsoft/wssdagent/rpc/client"
	wssdkeyvault "github.com/microsoft/wssdagent/rpc/keyvault"
	log "k8s.io/klog"
)

type client struct {
	wssdkeyvault.SimpleVaultAgentClient
}

// newClient - creates a client session with the backend wssd agent
func newSimpleVaultClient(subID string) (*client, error) {
	c, err := wssdclient.GetSimpleVaultClient(&subID)
	if err != nil {
		return nil, err
	}
	return &client{c}, nil
}

// Get
func (c *client) Get(ctx context.Context, group, name string) (*[]keyvault.SimpleVault, error) {
	request := getSimpleVaultRequest(wssdkeyvault.Operation_GET, name, nil)
	response, err := c.SimpleVaultAgentClient.Invoke(ctx, request)
	if err != nil {
		return nil, err
	}
	return getSimpleVaultsFromResponse(response), nil
}

// CreateOrUpdate
func (c *client) CreateOrUpdate(ctx context.Context, group, name string, sg *keyvault.SimpleVault) (*keyvault.SimpleVault, error) {
	request := getSimpleVaultRequest(wssdkeyvault.Operation_POST, name, sg)
	response, err := c.SimpleVaultAgentClient.Invoke(ctx, request)
	if err != nil {
		log.Errorf("[SimpleVault] Create failed with error %v", err)
		return nil, err
	}

	vault := getSimpleVaultsFromResponse(response)
	
	if len(*vault) == 0 {
		return nil, fmt.Errorf("[SimpleVault][Create] Unexpected error: Creating a network interface returned no result")
	}
	
	return &((*vault)[0]), err
}

// Delete methods invokes create or update on the client
func (c *client) Delete(ctx context.Context, group, name string) error {
	request := getSimpleVaultRequest(wssdkeyvault.Operation_DELETE, name, nil)
	_, err := c.SimpleVaultAgentClient.Invoke(ctx, request)
	return err
}

func getSimpleVaultsFromResponse(response *wssdkeyvault.SimpleVaultResponse) *[]keyvault.SimpleVault {
	SimpleVaults := []keyvault.SimpleVault{}
	for _, keyvaults := range response.GetSimpleVaults() {
		SimpleVaults = append(SimpleVaults, *(getSimpleVault(keyvaults)))
	}

	return &SimpleVaults
}

func getSimpleVaultRequest(opType wssdkeyvault.Operation, name string, vault *keyvault.SimpleVault) *wssdkeyvault.SimpleVaultRequest {
	request := &wssdkeyvault.SimpleVaultRequest{
		OperationType:   opType,
		SimpleVaults: []*wssdkeyvault.SimpleVault{},
	}
	if vault != nil {
		request.SimpleVaults = append(request.SimpleVaults, getWssdSimpleVault(vault))
	} else if len(name) > 0 {
		request.SimpleVaults = append(request.SimpleVaults,
			&wssdkeyvault.SimpleVault{
				Name: name,
			})
	}
	return request
}

func getSimpleVault(vault *wssdkeyvault.SimpleVault) *keyvault.SimpleVault {

	return &keyvault.SimpleVault{
		BaseProperties: keyvault.BaseProperties{
			ID : &vault.Id,
			Name: &vault.Name,
		},
	//	Source : &vault.Source,
	}
}

func getWssdSimpleVault(vault *keyvault.SimpleVault) *wssdkeyvault.SimpleVault {
	return &wssdkeyvault.SimpleVault{
		//Id : *vault.BaseProperties.ID,
		Name: *vault.BaseProperties.Name,
	//	Source: *vault.Source,
		Secrets: []*wssdkeyvault.Secret{},
	}
}
