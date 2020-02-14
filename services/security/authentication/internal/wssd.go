// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license

package internal

import (
	"context"
	"github.com/microsoft/wssd-sdk-for-go/pkg/auth"
	wssdclient "github.com/microsoft/wssd-sdk-for-go/pkg/client"
	"github.com/microsoft/wssd-sdk-for-go/services/security"
	wssdsecurity "github.com/microsoft/wssdagent/rpc/security"
	//log "k8s.io/klog"
)

type client struct {
	wssdsecurity.AuthenticationAgentClient
}

// NewAuthenticationClient creates a client session with the backend wssd agent
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
