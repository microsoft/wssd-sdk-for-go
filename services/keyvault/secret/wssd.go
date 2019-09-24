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

package secret

import (
	"fmt"
	"context"
	"github.com/microsoft/wssd-sdk-for-go/services/keyvault"

	wssdclient "github.com/microsoft/wssdagent/rpc/client"
	wssdkeyvault "github.com/microsoft/wssdagent/rpc/keyvault"
	log "k8s.io/klog"
)

type client struct {
	wssdkeyvault.SecretAgentClient
}

// newClient - creates a client session with the backend wssd agent
func newSecretClient(subID string) (*client, error) {
	c, err := wssdclient.GetSecretClient(&subID)
	if err != nil {
		return nil, err
	}
	return &client{c}, nil
}

// Get
func (c *client) Get(ctx context.Context, group, name string, sg *keyvault.Secret) (*[]keyvault.Secret, error) {
	request := getSecretRequest(wssdkeyvault.Operation_GET, name, sg)
	response, err := c.SecretAgentClient.Invoke(ctx, request)
	if err != nil {
		return nil, err
	}
	return getSecretsFromResponse(response), nil
}

// CreateOrUpdate
func (c *client) CreateOrUpdate(ctx context.Context, group, name string, sg *keyvault.Secret) (*keyvault.Secret, error) {
	request := getSecretRequest(wssdkeyvault.Operation_POST, name, sg)
	response, err := c.SecretAgentClient.Invoke(ctx, request)
	if err != nil {
		log.Errorf("[Secret] Create failed with error %v", err)
		return nil, err
	}

	sec := getSecretsFromResponse(response)
	
	if len(*sec) == 0 {
		return nil, fmt.Errorf("[Secret][Create] Unexpected error: Creating a network interface returned no result")
	}
	
	return &((*sec)[0]), err
}

// Delete methods invokes create or update on the client
func (c *client) Delete(ctx context.Context, group, name string) error {
	request := getSecretRequest(wssdkeyvault.Operation_DELETE, name, nil)
	_, err := c.SecretAgentClient.Invoke(ctx, request)
	return err
}

func getSecretsFromResponse(response *wssdkeyvault.SecretResponse) *[]keyvault.Secret {
	Secrets := []keyvault.Secret{}
	for _, secrets := range response.GetSecrets() {
		Secrets = append(Secrets, *(getSecret(secrets)))
	}

	return &Secrets
}

func getSecretRequest(opType wssdkeyvault.Operation, name string, sec *keyvault.Secret) *wssdkeyvault.SecretRequest {
	request := &wssdkeyvault.SecretRequest{
		OperationType:   opType,
		Secrets: []*wssdkeyvault.Secret{},
	}
	if sec != nil {
		request.Secrets = append(request.Secrets, getWssdSecret(sec))
	} else if len(name) > 0 {
		request.Secrets = append(request.Secrets,
			&wssdkeyvault.Secret{
				Name: name,
			})
	}
	return request
}

func getSecret(sec *wssdkeyvault.Secret) *keyvault.Secret {

	return &keyvault.Secret{
		BaseProperties: keyvault.BaseProperties{
			ID : &sec.Id,
			Name: &sec.Name,
		},
		FileName : &sec.Filename,
		Value : &sec.Value,
		VaultName : &sec.VaultName,
	}
}

func getWssdSecret(sec *keyvault.Secret) *wssdkeyvault.Secret {
	return &wssdkeyvault.Secret{
		//Id : *vault.BaseProperties.ID,
		Name: *sec.BaseProperties.Name,
	//	Filename: *sec.FileName,
		Value: *sec.Value,
		VaultName: *sec.VaultName,
	}
}
