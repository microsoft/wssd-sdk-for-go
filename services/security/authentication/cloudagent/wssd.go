// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license

package cloudagent

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/microsoft/moc/pkg/auth"
	wssdsecurity "github.com/microsoft/moc/rpc/cloudagent/security"
	wssdsdkauthport "github.com/microsoft/wssd-sdk-for-go/pkg/lbagentclient"
	"github.com/microsoft/wssd-sdk-for-go/services/security"
	"google.golang.org/grpc"
	//log "k8s.io/klog"
)

type client struct {
	wssdsecurity.AuthenticationAgentClient
}

func getDefaultAuthServerEndpoint(serverAddress *string) string {
	return fmt.Sprintf("%s:%d", *serverAddress, wssdsdkauthport.AuthPort)
}

func getAuthServerEndpoint(serverAddress *string) string {
	if !strings.Contains(*serverAddress, ":") {
		return getDefaultAuthServerEndpoint(serverAddress)
	}
	return *serverAddress
}

// getAuthenticationClient returns the secret client to communicate with the wssdagent
func getAuthenticationClient(serverAddress *string, authorizer auth.Authorizer) (wssdsecurity.AuthenticationAgentClient, error) {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(authorizer.WithTransportAuthorization()))
	opts = append(opts, grpc.WithPerRPCCredentials(authorizer.WithRPCAuthorization()))

	conn, err := grpc.Dial(getAuthServerEndpoint(serverAddress), opts...)
	if err != nil {
		log.Fatalf("Unable to get AuthenticationClient. Failed to dial: %v", err)
	}

	return wssdsecurity.NewAuthenticationAgentClient(conn), nil
}

// NewAuthenticationClient creates a client session with the backend wssd agent
func NewAuthenticationClient(subID string, authorizer auth.Authorizer) (*client, error) {
	c, err := getAuthenticationClient(&subID, authorizer)
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
