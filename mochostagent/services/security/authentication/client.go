// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license

package authentication

import (
	"context"

	"github.com/microsoft/moc/pkg/auth"
	"github.com/microsoft/wssd-sdk-for-go/mochostagent/services/security"
	"github.com/microsoft/wssd-sdk-for-go/mochostagent/services/security/authentication/hostagent"
	"github.com/microsoft/wssd-sdk-for-go/mochostagent/services/security/authentication/nodeagent"
)

const (
	HostAgentSpec = "HostAgent"
	NodeAgentSpec = "NodeAgent"
)

// Service interface
type Service interface {
	Login(context.Context, string, *security.Identity) (*string, error)
}

// Client structure
type AuthenticationClient struct {
	security.BaseClient
	internal Service
}

// NewAuthenticationClient method returns new client used to connect to HostAgent
func NewAuthenticationClient(cloudFQDN string, authorizer auth.Authorizer) (*AuthenticationClient, error) {
	return NewAuthenticationClientToServer(cloudFQDN, authorizer, HostAgentSpec)
}

// NewAuthenticationClientToServer method returns new client used to connect to HostAgent or NodeAgent
func NewAuthenticationClientToServer(cloudFQDN string, authorizer auth.Authorizer, serverSpec string) (*AuthenticationClient, error) {
	var authClient Service
	var err error

	switch serverSpec {
	case NodeAgentSpec:
		authClient, err = nodeagent.NewAuthenticationClient(cloudFQDN, authorizer)
	case HostAgentSpec:
		authClient, err = hostagent.NewAuthenticationClient(cloudFQDN, authorizer)
	default:
		authClient, err = hostagent.NewAuthenticationClient(cloudFQDN, authorizer)
	}
	if err != nil {
		return nil, err
	}

	return &AuthenticationClient{internal: authClient}, nil
}

func (c *AuthenticationClient) Login(ctx context.Context, group string, identity *security.Identity) (*string, error) {
	return c.internal.Login(ctx, group, identity)
}
