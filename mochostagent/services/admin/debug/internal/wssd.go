// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license

package internal

import (
	"context"

	"github.com/microsoft/moc/pkg/auth"
	wssdadmin "github.com/microsoft/moc/rpc/common/admin"
	wssdclient "github.com/microsoft/wssd-sdk-for-go/mochostagent/pkg/client"
)

type client struct {
	wssdadmin.DebugAgentClient
}

// NewDebugClient - creates a client session with the backend wssd hostagent
func NewDebugClient(subID string, authorizer auth.Authorizer) (*client, error) {
	c, err := wssdclient.GetDebugClient(&subID, authorizer)
	if err != nil {
		return nil, err
	}
	return &client{c}, nil
}

// Stacktrace
func (c *client) Stacktrace(ctx context.Context) (string, error) {
	request := getDebugRequest(wssdadmin.DebugOperation_STACKTRACE)
	response, err := c.DebugAgentClient.Invoke(ctx, request)
	if err != nil {
		return "", err
	}
	return response.Result, nil
}

func getDebugRequest(operation wssdadmin.DebugOperation) *wssdadmin.DebugRequest {
	return &wssdadmin.DebugRequest{OperationType: operation}
}
