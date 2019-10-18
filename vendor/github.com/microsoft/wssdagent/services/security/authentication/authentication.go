// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package authentication

import (
	"context"
	pb "github.com/microsoft/wssdagent/rpc/security"
	"github.com/microsoft/wssdagent/pkg/auth"
)


type AuthenticationProvider struct {
	client *Client
}

func NewAuthenticationProvider() *AuthenticationProvider {
	return &AuthenticationProvider{
		client: NewClient(),
	}
}

func (authenticationProv *AuthenticationProvider) Login(ctx context.Context, identity *pb.Identity, authorizer *auth.JwtAuthorizer) (*string, error) {
	resultToken, err := authenticationProv.client.Login(ctx, identity, authorizer)
	if err != nil {
		return nil, err
	}
	return resultToken, nil
}