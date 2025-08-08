// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license

package availabilityset

import (
	"context"

	"github.com/microsoft/moc/pkg/auth"
	"github.com/microsoft/wssd-sdk-for-go/services/compute"
	"github.com/microsoft/wssd-sdk-for-go/services/compute/availabilityset/internal"
)

type Service interface {
	Get(context.Context, string) (*[]compute.AvailabilitySet, error)
	CreateOrUpdate(context.Context, string, *compute.AvailabilitySet) (*compute.AvailabilitySet, error)
	Delete(context.Context, string) error
	AddVmToAvailabilitySet(context.Context, string, string) error
	RemoveVmFromAvailabilitySet(context.Context, string, string) error
}

type AvailabilitySetClient struct {
	compute.BaseClient
	internal Service
}

func NewAvailabilitySetClient(cloudFQDN string, authorizer auth.Authorizer) (*AvailabilitySetClient, error) {
	c, err := internal.NewAvailabilitySetWssdClient(cloudFQDN, authorizer)
	if err != nil {
		return nil, err
	}

	return &AvailabilitySetClient{internal: c}, nil
}

func (c *AvailabilitySetClient) Get(ctx context.Context, name string) (*[]compute.AvailabilitySet, error) {
	return c.internal.Get(ctx, name)
}

func (c *AvailabilitySetClient) CreateOrUpdate(ctx context.Context, name string, avset *compute.AvailabilitySet) (*compute.AvailabilitySet, error) {
	return c.internal.CreateOrUpdate(ctx, name, avset)
}

func (c *AvailabilitySetClient) Delete(ctx context.Context, name string) error {
	return c.internal.Delete(ctx, name)
}

func (c *AvailabilitySetClient) AddVmToAvailabilitySet(ctx context.Context, avset string, nodeagnetVMName string) error {
	return c.internal.AddVmToAvailabilitySet(ctx, avset, nodeagnetVMName)
}

func (c *AvailabilitySetClient) RemoveVmFromAvailabilitySet(ctx context.Context, avset string, nodeagnetVMName string) error {
	return c.internal.RemoveVmFromAvailabilitySet(ctx, avset, nodeagnetVMName)
}
