// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license

package internal

import (
	"context"

	"github.com/microsoft/moc/pkg/status"
	"github.com/microsoft/wssd-sdk-for-go/services/security/keyvault"

	"github.com/microsoft/moc/pkg/auth"
	"github.com/microsoft/moc/pkg/errors"
	wssdcommonproto "github.com/microsoft/moc/rpc/common"
	wssdsecurity "github.com/microsoft/moc/rpc/nodeagent/security"
	wssdclient "github.com/microsoft/wssd-sdk-for-go/pkg/client"
	log "k8s.io/klog"
)

type client struct {
	wssdsecurity.KeyAgentClient
}

// NewKeyClient - creates a client session with the backend wssd agent
func NewKeyClient(subID string, authorizer auth.Authorizer) (*client, error) {
	c, err := wssdclient.GetKeyClient(&subID, authorizer)
	if err != nil {
		return nil, err
	}
	return &client{c}, nil
}

// Get
func (c *client) Get(ctx context.Context, name string, vaultName string) (*[]keyvault.Key, error) {
	request := getKeyRequest(wssdcommonproto.Operation_GET, name, vaultName, nil, 0, nil)
	response, err := c.KeyAgentClient.Invoke(ctx, request)
	if err != nil {
		return nil, err
	}
	return getKeysFromResponse(response), nil
}

// CreateOrUpdate
func (c *client) CreateOrUpdate(ctx context.Context, keyIn *keyvault.Key) (*keyvault.Key, error) {
	err := c.validate(ctx, keyIn)
	if err != nil {
		return nil, err
	}

	request := getKeyRequest(wssdcommonproto.Operation_POST, "", "", nil, 0, keyIn)
	response, err := c.KeyAgentClient.Invoke(ctx, request)
	if err != nil {
		log.Errorf("[Key] Create failed with error %v", err)
		return nil, err
	}

	keys := getKeysFromResponse(response)

	if len(*keys) == 0 {
		return nil, errors.New("[Key][Create] Unexpected error: Creating a key returned no result")
	}

	return &((*keys)[0]), err
}

func (c *client) validate(ctx context.Context, key *keyvault.Key) (err error) {
	if key == nil || key.VaultName == nil || key.Name == nil || key.KeyType == nil || getMOCKeySize(key.KeySize) == wssdcommonproto.KeySize_K_UNKNOWN {
		return errors.Wrapf(errors.InvalidInput, "[Key][Create] Invalid Input")
	}

	return nil
}

// Delete methods invokes create or update on the client
func (c *client) Delete(ctx context.Context, key *keyvault.Key) error {
	keys, err := c.Get(ctx, *key.Name, *key.VaultName)
	if err != nil {
		return err
	}
	if len(*keys) == 0 {
		return errors.Wrapf(errors.NotFound, "Key [%s] not found", *key.Name)
	}

	request := getKeyRequest(wssdcommonproto.Operation_DELETE, "", "", nil, 0, &(*keys)[0])
	_, err = c.KeyAgentClient.Invoke(ctx, request)
	return err
}

// Rotates a key and returns the new key
func (c *client) RotateKey(ctx context.Context, keyReq *keyvault.KeyOperationRequest) (*keyvault.KeyOperationResult, error) {
	wssdReq := wssdsecurity.KeyOperationRequest{
		Key:           getWssdKey(keyReq.Key),
		OperationType: wssdcommonproto.ProviderAccessOperation_Key_Rotate}

	wssdRep, err := c.KeyAgentClient.Operate(ctx, &wssdReq)

	if err != nil {
		return nil, err
	}

	keyOpRes := keyvault.KeyOperationResult{
		Key:    getKey(wssdRep.GetKey()),
		Result: nil} // No result expected from rotate

	return &keyOpRes, nil
}

// Wraps a key and returns the result
func (c *client) WrapKey(ctx context.Context, keyReq *keyvault.KeyOperationRequest) (*keyvault.KeyOperationResult, error) {
	wssdReq := wssdsecurity.KeyOperationRequest{
		Key:           getWssdKey(keyReq.Key),
		Algorithm:     wssdcommonproto.Algorithm(wssdcommonproto.Algorithm_value[string(*keyReq.Algorithm)]),
		OperationType: wssdcommonproto.ProviderAccessOperation_Key_WrapKey,
		Data:          *keyReq.Data}

	wssdRep, err := c.KeyAgentClient.Operate(ctx, &wssdReq)

	if err != nil {
		return nil, err
	}

	if wssdRep.Data == "" {
		return nil, errors.New("[Key][Wrap] Unexpected error: Wrapping a key returned no result")
	}

	keyOpRes := keyvault.KeyOperationResult{
		Key:    nil, // No key changes expected
		Result: &wssdRep.Data}

	return &keyOpRes, nil
}

// Unwraps a key and returns the result
func (c *client) UnwrapKey(ctx context.Context, keyReq *keyvault.KeyOperationRequest) (*keyvault.KeyOperationResult, error) {
	wssdReq := wssdsecurity.KeyOperationRequest{
		Key:           getWssdKey(keyReq.Key),
		Algorithm:     wssdcommonproto.Algorithm(wssdcommonproto.Algorithm_value[string(*keyReq.Algorithm)]),
		OperationType: wssdcommonproto.ProviderAccessOperation_Key_UnwrapKey,
		Data:          *keyReq.Data}

	wssdRep, err := c.KeyAgentClient.Operate(ctx, &wssdReq)

	if err != nil {
		return nil, err
	}

	if wssdRep.Data == "" {
		return nil, errors.New("[Key][Wrap] Unexpected error: Unwrapping a key returned no result")
	}

	keyOpRes := keyvault.KeyOperationResult{
		Key:    nil, // No key changes expected
		Result: &wssdRep.Data}

	return &keyOpRes, nil
}

func getKeysFromResponse(response *wssdsecurity.KeyResponse) *[]keyvault.Key {
	Keys := []keyvault.Key{}
	for _, key := range response.GetKeys() {
		Keys = append(Keys, *(getKey(key)))
	}

	return &Keys
}

func getKeyRequest(opType wssdcommonproto.Operation, name, vaultName string, keyType *keyvault.JSONWebKeyType, keySize int32, key *keyvault.Key) *wssdsecurity.KeyRequest {
	request := &wssdsecurity.KeyRequest{
		OperationType: opType,
		Keys:          []*wssdsecurity.Key{},
	}

	if key != nil {
		request.Keys = append(request.Keys, getWssdKey(key))
	} else if len(name) > 0 {
		tempKey := &wssdsecurity.Key{
			Name:      name,
			VaultName: vaultName,
			Size:      getMOCKeySize(keySize)}

		if keyType != nil {
			tempKey.Type = wssdcommonproto.JsonWebKeyType(wssdcommonproto.JsonWebKeyType_value[string(*keyType)])
		}

		request.Keys = append(request.Keys, tempKey)
	}
	return request
}

func getKey(key *wssdsecurity.Key) *keyvault.Key {
	ct := key.CreationTime.AsTime()
	keyType := keyvault.JSONWebKeyType(wssdcommonproto.JsonWebKeyType_name[int32(key.Type)])
	return &keyvault.Key{
		ID:                &key.Id,
		Name:              &key.Name,
		VaultName:         &key.VaultName,
		CreationTime:      &ct,
		KeyVersion:        key.KeyVersion,
		KeyType:           &keyType,
		KeySize:           getKeySize(key.Size),
		ProvisioningState: status.GetProvisioningState(key.GetStatus().GetProvisioningStatus())}
}

func getWssdKey(key *keyvault.Key) *wssdsecurity.Key {
	keyOut := &wssdsecurity.Key{
		Name:       *key.Name,
		VaultName:  *key.VaultName,
		Type:       wssdcommonproto.JsonWebKeyType(wssdcommonproto.JsonWebKeyType_value[string(*key.KeyType)]),
		Size:       getMOCKeySize(key.KeySize),
		KeyVersion: key.KeyVersion}

	return keyOut
}

func getMOCKeySize(size int32) (ksize wssdcommonproto.KeySize) {
	switch size {
	case 256:
		ksize = wssdcommonproto.KeySize__256
	default:
		ksize = wssdcommonproto.KeySize_K_UNKNOWN
	}
	return
}

func getKeySize(ksize wssdcommonproto.KeySize) (size int32) {
	switch ksize {
	case wssdcommonproto.KeySize__256:
		size = 256
	default:
		size = -1
	}
	return
}
