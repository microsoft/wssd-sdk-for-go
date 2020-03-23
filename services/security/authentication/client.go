// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license

package authentication

import (
	"context"
	"github.com/microsoft/moc/pkg/auth"
	"github.com/microsoft/wssd-sdk-for-go/services/security"
	"github.com/microsoft/wssd-sdk-for-go/services/security/authentication/internal"
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

// NewClient method returns new client
func NewAuthenticationClient(cloudFQDN string, authorizer auth.Authorizer) (*AuthenticationClient, error) {
	c, err := internal.NewAuthenticationClient(cloudFQDN, authorizer)
	if err != nil {
		return nil, err
	}

	return &AuthenticationClient{internal: c}, nil
}

// Get methods invokes the client Get method
func (c *AuthenticationClient) Login(ctx context.Context, group string, identity *security.Identity) (*string, error) {
	return c.internal.Login(ctx, group, identity)
}
