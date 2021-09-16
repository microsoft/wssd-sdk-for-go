// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license

package baremetalhostagentclient

import (
	"context"

	"github.com/microsoft/moc/pkg/auth"
	pb "github.com/microsoft/moc/rpc/baremetalhostagent"
)

type client struct {
	pb.BareMetalHostAgentClient
}

// newClient - creates a client session with the backend bare metal host agent
func newBareMetalHostAgentClient(subID string, authorizer auth.Authorizer) (*client, error) {
	c, err := GetBareMetalHostAgentClient(&subID, authorizer)
	if err != nil {
		return nil, err
	}
	return &client{c}, nil
}

// Method invokes update on the client
func (c *client) Update(ctx context.Context, bmh *pb.BareMetalHost) (*pb.BareMetalHost, error) {

	request := &pb.BareMetalHostRequest{
		BareMetalHost: bmh,
	}

	response, err := c.BareMetalHostAgentClient.CreateOrUpdate(ctx, request)
	if err != nil {
		return nil, err
	}
	return response.GetBareMetalHost(), nil
}
