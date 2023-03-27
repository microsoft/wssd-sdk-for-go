// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package internal

import (
	"context"

	hostcompute "github.com/microsoft/moc/rpc/mochostagent/compute"
	mochostagentclient "github.com/microsoft/wssd-sdk-for-go/mochostagent/pkg/client"

	"github.com/microsoft/moc/pkg/auth"
)

type client struct {
	hostcompute.VirtualMachineAgentClient
}

// newVirtualMachineClient - creates a client session with the backend host agent
func NewVirtualMachineClient(subID string, authorizer auth.Authorizer) (*client, error) {
	c, err := mochostagentclient.GetVirtualMachineClient(&subID, authorizer)
	if err != nil {
		return nil, err
	}
	return &client{c}, nil
}

func (c *client) RegisterVirtualMachine(ctx context.Context, request *hostcompute.RegisterVirtualMachineRequest) (response *hostcompute.RegisterVirtualMachineResponse, err error) {
	return c.VirtualMachineAgentClient.RegisterVirtualMachine(ctx, request)
}

func (c *client) DeregisterVirtualMachine(ctx context.Context, request *hostcompute.DeregisterVirtualMachineRequest) (response *hostcompute.DeregisterVirtualMachineResponse, err error) {
	return c.VirtualMachineAgentClient.DeregisterVirtualMachine(ctx, request)
}

func (c *client) RunCommand(ctx context.Context, request *hostcompute.VirtualMachineRunCommandRequest) (response *hostcompute.VirtualMachineRunCommandResponse, err error) {
	return c.VirtualMachineAgentClient.RunCommand(ctx, request)
}
