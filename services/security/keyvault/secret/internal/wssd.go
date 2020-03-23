// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license

package internal

import (
	"context"
	"fmt"
	"github.com/microsoft/moc/pkg/auth"
	"github.com/microsoft/wssd-sdk-for-go/services/security/keyvault"

	wssdcommonproto "github.com/microsoft/moc/rpc/common"
	wssdsecurity "github.com/microsoft/moc/rpc/nodeagent/security"
	wssdclient "github.com/microsoft/wssd-sdk-for-go/pkg/client"
	log "k8s.io/klog"
)

type client struct {
	wssdsecurity.SecretAgentClient
}

// NewSecretClient - creates a client session with the backend wssd agent
func NewSecretClient(subID string, authorizer auth.Authorizer) (*client, error) {
	c, err := wssdclient.GetSecretClient(&subID, authorizer)
	if err != nil {
		return nil, err
	}
	return &client{c}, nil
}

// Get
func (c *client) Get(ctx context.Context, group, name, vaultName string) (*[]keyvault.Secret, error) {
	request := getSecretRequest(wssdcommonproto.Operation_GET, name, vaultName, nil)
	response, err := c.SecretAgentClient.Invoke(ctx, request)
	if err != nil {
		return nil, err
	}
	return getSecretsFromResponse(response), nil
}

// CreateOrUpdate
func (c *client) CreateOrUpdate(ctx context.Context, group, name string, sg *keyvault.Secret) (*keyvault.Secret, error) {
	err := c.validate(ctx, group, name, sg)
	if err != nil {
		return nil, err
	}
	request := getSecretRequest(wssdcommonproto.Operation_POST, name, *sg.VaultName, sg)
	response, err := c.SecretAgentClient.Invoke(ctx, request)
	if err != nil {
		log.Errorf("[Secret] Create failed with error %v", err)
		return nil, err
	}

	sec := getSecretsFromResponse(response)

	if len(*sec) == 0 {
		return nil, fmt.Errorf("[Secret][Create] Unexpected error: Creating a secret returned no result")
	}

	return &((*sec)[0]), err
}

func (c *client) validate(ctx context.Context, group, name string, sg *keyvault.Secret) (err error) {
	if sg == nil || sg.VaultName == nil || sg.Value == nil {
		return fmt.Errorf("[Secret][Create] Invalid Input")
	}

	if sg.Name == nil {
		sg.Name = &name
	}
	return nil
}

// Delete methods invokes create or update on the client
func (c *client) Delete(ctx context.Context, group, name, vaultName string) error {
	secret, err := c.Get(ctx, group, name, vaultName)
	if err != nil {
		return err
	}
	if len(*secret) == 0 {
		return fmt.Errorf("Keysecret [%s] not found", name)
	}

	request := getSecretRequest(wssdcommonproto.Operation_DELETE, name, vaultName, &(*secret)[0])
	_, err = c.SecretAgentClient.Invoke(ctx, request)
	return err
}

func getSecretsFromResponse(response *wssdsecurity.SecretResponse) *[]keyvault.Secret {
	Secrets := []keyvault.Secret{}
	for _, secrets := range response.GetSecrets() {
		Secrets = append(Secrets, *(getSecret(secrets)))
	}

	return &Secrets
}

func getSecretRequest(opType wssdcommonproto.Operation, name, vaultName string, sec *keyvault.Secret) *wssdsecurity.SecretRequest {
	request := &wssdsecurity.SecretRequest{
		OperationType: opType,
		Secrets:       []*wssdsecurity.Secret{},
	}
	if sec != nil {
		request.Secrets = append(request.Secrets, getWssdSecret(sec, opType))
	} else if len(name) > 0 {
		request.Secrets = append(request.Secrets,
			&wssdsecurity.Secret{
				Name:      name,
				VaultName: vaultName,
			})
	} else {
		request.Secrets = append(request.Secrets,
			&wssdsecurity.Secret{
				VaultName: vaultName,
			})

	}
	return request
}

func getSecret(sec *wssdsecurity.Secret) *keyvault.Secret {
	value := string(sec.Value)
	return &keyvault.Secret{
		ID:    &sec.Id,
		Name:  &sec.Name,
		Value: &value,
		SecretProperties: &keyvault.SecretProperties{
			FileName:  &sec.Filename,
			VaultName: &sec.VaultName,
		},
	}
}

func getWssdSecret(sec *keyvault.Secret, opType wssdcommonproto.Operation) *wssdsecurity.Secret {
	secret := &wssdsecurity.Secret{
		Name:      *sec.Name,
		VaultName: *sec.VaultName,
	}

	if opType == wssdcommonproto.Operation_POST {
		secret.Value = []byte(*sec.Value)
	}

	return secret
}
