// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package authentication

import (
	"context"
	"github.com/microsoft/wssdagent/pkg/apis/config"
	"github.com/microsoft/wssdagent/pkg/auth"
	pb "github.com/microsoft/wssdagent/rpc/security"
	//"github.com/microsoft/wssdagent/services/security/identity/internal"
	
	"github.com/microsoft/wssdagent/services/security/authentication/ikv"
)


const (
	IKVSpec = "ikv"
)

type Service interface {
	Login(context.Context, *pb.Identity, *auth.JwtAuthorizer) (*string, error)
}

type Client struct {
	internal Service
}

func NewClient() *Client {
	cConfig := config.GetChildAgentConfiguration("Authentication")
	c := &Client{}
	switch cConfig.ProviderSpec {
	case IKVSpec:
	default:
		c.internal = ikv.NewClient()
	}
	return c
}

func (c *Client) Login(ctx context.Context, identity *pb.Identity, authorizer *auth.JwtAuthorizer) (*string, error) {
	resultToken, err := c.internal.Login(ctx, identity, authorizer)
	if err != nil {
		return nil, err
	}
	return resultToken, nil
}