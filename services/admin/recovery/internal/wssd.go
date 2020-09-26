// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license

package internal

import (
	"context"

	"github.com/microsoft/moc/pkg/auth"
	wssdadmin "github.com/microsoft/moc/rpc/common/admin"
	wssdclient "github.com/microsoft/wssd-sdk-for-go/pkg/client"
)

type client struct {
	wssdadmin.RecoveryAgentClient
}

// NewRecoveryClient - creates a client session with the backend wssd agent
func NewRecoveryClient(subID string, authorizer auth.Authorizer) (*client, error) {
	c, err := wssdclient.GetRecoveryClient(&subID, authorizer)
	if err != nil {
		return nil, err
	}
	return &client{c}, nil
}

// Backup
func (c *client) Backup(ctx context.Context, path string, configFilePath string, storeType string) error {
	request := getRecoveryRequest(wssdadmin.Operation_BACKUP, path, configFilePath, storeType)
	_, err := c.RecoveryAgentClient.Invoke(ctx, request)
	return err
}

// Restore
func (c *client) Restore(ctx context.Context, path string, configFilePath string, storeType string) error {
	request := getRecoveryRequest(wssdadmin.Operation_RESTORE, path, configFilePath, storeType)
	_, err := c.RecoveryAgentClient.Invoke(ctx, request)
	return err
}

func getRecoveryRequest(operation wssdadmin.Operation, path string, configFilePath string, storeType string) *wssdadmin.RecoveryRequest {
	return &wssdadmin.RecoveryRequest{OperationType: operation, Path: path, ConfigFilePath: configFilePath, StoreType: storeType}
}
