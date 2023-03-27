// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license

package virtualmachine

import (
	"context"

	"github.com/microsoft/moc/rpc/mochostagent/compute"
	sdkCompute "github.com/microsoft/wssd-sdk-for-go/mochostagent/services/compute"
	"github.com/microsoft/wssd-sdk-for-go/mochostagent/services/compute/virtualmachine/internal"

	"github.com/microsoft/moc/pkg/auth"
)

type Service interface {
	RegisterVirtualMachine(context.Context, *compute.RegisterVirtualMachineRequest) (*compute.RegisterVirtualMachineResponse, error)
	DeregisterVirtualMachine(context.Context, *compute.DeregisterVirtualMachineRequest) (*compute.DeregisterVirtualMachineResponse, error)
	RunCommand(context.Context, *compute.VirtualMachineRunCommandRequest) (*compute.VirtualMachineRunCommandResponse, error)
}

type VirtualMachineClient struct {
	sdkCompute.BaseClient
	internal Service
}

func NewVirtualMachineClient(cloudFQDN string, authorizer auth.Authorizer) (*VirtualMachineClient, error) {
	c, err := internal.NewVirtualMachineClient(cloudFQDN, authorizer)
	if err != nil {
		return nil, err
	}

	return &VirtualMachineClient{internal: c}, nil
}

func (c *VirtualMachineClient) RegisterVirtualMachine(ctx context.Context, request *compute.RegisterVirtualMachineRequest) (*compute.RegisterVirtualMachineResponse, error) {
	return c.internal.RegisterVirtualMachine(ctx, request)
}

func (c *VirtualMachineClient) DeregisterVirtualMachine(ctx context.Context, request *compute.DeregisterVirtualMachineRequest) (*compute.DeregisterVirtualMachineResponse, error) {
	return c.internal.DeregisterVirtualMachine(ctx, request)
}

func (c *VirtualMachineClient) RunCommand(ctx context.Context, request *compute.VirtualMachineRunCommandRequest) (*compute.VirtualMachineRunCommandResponse, error) {
	return c.internal.RunCommand(ctx, request)
}
