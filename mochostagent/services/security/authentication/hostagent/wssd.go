// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license

package hostagent

import (
	"context"

	"github.com/microsoft/moc/pkg/auth"
	wssdsecurity "github.com/microsoft/moc/rpc/mochostagent/security"
	wssdclient "github.com/microsoft/wssd-sdk-for-go/mochostagent/pkg/client"
	"github.com/microsoft/wssd-sdk-for-go/mochostagent/services/security"
)

type client struct {
	wssdsecurity.AuthenticationAgentClient
}

// NewAuthenticationClient creates a client session with the backend mochostagent
func NewAuthenticationClient(subID string, authorizer auth.Authorizer) (*client, error) {
	c, err := wssdclient.GetAuthenticationClient(&subID, authorizer)
	if err != nil {
		return nil, err
	}
	return &client{c}, nil
}

// Login
func (c *client) Login(ctx context.Context, group string, identity *security.Identity) (*string, error) {
	request := getAuthenticationRequest(identity)
	response, err := c.AuthenticationAgentClient.Login(ctx, request)
	if err != nil {
		return nil, err
	}
	return &response.Token, nil
}

func getAuthenticationRequest(identity *security.Identity) *wssdsecurity.AuthenticationRequest {
	request := &wssdsecurity.AuthenticationRequest{
		Identity: &wssdsecurity.Identity{
			Name:        *identity.Name,
			Certificate: *identity.Certificate,
		},
	}
	return request
}
