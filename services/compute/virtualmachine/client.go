// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license

package virtualmachine

import (
	"context"
	"github.com/microsoft/wssd-sdk-for-go/pkg/auth"
	"github.com/microsoft/wssd-sdk-for-go/services/compute"
	"github.com/microsoft/wssd-sdk-for-go/services/compute/virtualmachine/internal"
)

type Service interface {
	Get(context.Context, string, string) (*[]compute.VirtualMachine, error)
	CreateOrUpdate(context.Context, string, string, *compute.VirtualMachine) (*compute.VirtualMachine, error)
	Delete(context.Context, string, string) error
}

type VirtualMachineClient struct {
	compute.BaseClient
	internal Service
}

func NewVirtualMachineClient(cloudFQDN string, authorizer auth.Authorizer) (*VirtualMachineClient, error) {
	c, err := internal.NewVirtualMachineClient(cloudFQDN, authorizer)
	if err != nil {
		return nil, err
	}

	return &VirtualMachineClient{internal: c}, nil
}

// Get methods invokes the client Get method
func (c *VirtualMachineClient) Get(ctx context.Context, group, name string) (*[]compute.VirtualMachine, error) {
	return c.internal.Get(ctx, group, name)
}

// CreateOrUpdate methods invokes create or update on the client
func (c *VirtualMachineClient) CreateOrUpdate(ctx context.Context, group, name string, compute *compute.VirtualMachine) (*compute.VirtualMachine, error) {
	return c.internal.CreateOrUpdate(ctx, group, name, compute)
}

// Delete methods invokes delete of the compute resource
func (c *VirtualMachineClient) Delete(ctx context.Context, group string, name string) error {
	return c.internal.Delete(ctx, group, name)
}
