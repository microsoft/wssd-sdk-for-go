// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license

package authentication

import (
	"context"

	"github.com/microsoft/moc/pkg/auth"
	"github.com/microsoft/wssd-sdk-for-go/services/security"
	"github.com/microsoft/wssd-sdk-for-go/services/security/authentication/cloudagent"
	"github.com/microsoft/wssd-sdk-for-go/services/security/authentication/nodeagent"
)

const (
	CloudAgentSpec = "CloudAgent"
	NodeAgentSpec  = "NodeAgent"
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

// NewAuthenticationClient method returns new client used to connect to NodeAgent
func NewAuthenticationClient(cloudFQDN string, authorizer auth.Authorizer) (*AuthenticationClient, error) {
	return NewAuthenticationClientToServer(cloudFQDN, authorizer, NodeAgentSpec)
}

// NewAuthenticationClientToServer method returns new client used to connect to NodeAgent or CloudAgent
func NewAuthenticationClientToServer(cloudFQDN string, authorizer auth.Authorizer, serverSpec string) (*AuthenticationClient, error) {
	var authClient Service
	var err error

	switch serverSpec {
	case CloudAgentSpec:
		authClient, err = cloudagent.NewAuthenticationClient(cloudFQDN, authorizer)
	case NodeAgentSpec:
		authClient, err = nodeagent.NewAuthenticationClient(cloudFQDN, authorizer)
	default:
		authClient, err = nodeagent.NewAuthenticationClient(cloudFQDN, authorizer)
	}
	if err != nil {
		return nil, err
	}

	return &AuthenticationClient{internal: authClient}, nil
}

// Get methods invokes the client Get method
func (c *AuthenticationClient) Login(ctx context.Context, group string, identity *security.Identity) (*string, error) {
	return c.internal.Login(ctx, group, identity)
}
