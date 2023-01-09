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
	wssdadmin.ValidationAgentClient
}

// NewValidationgingClient - creates a client session with the backend wssd agent
func NewValidationClient(subID string, authorizer auth.Authorizer) (*client, error) {
	c, err := wssdclient.GetValidationClient(&subID, authorizer)
	if err != nil {
		return nil, err
	}
	return &client{c}, nil
}

// Backup
func (c *client) Validate(ctx context.Context) error {
	request := getValidationRequest()
	_, err := c.ValidationAgentClient.Invoke(ctx, request)
	return err
}

func getValidationRequest() *wssdadmin.ValidationRequest {
	return &wssdadmin.ValidationRequest{}
}
