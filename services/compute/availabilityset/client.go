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

// Get methods invokes the client Get method
func (c *AvailabilitySetClient) Get(ctx context.Context, name string) (*[]compute.AvailabilitySet, error) {
	return c.internal.Get(ctx, name)
}

// CreateOrUpdate methods invokes create or update on the client
func (c *AvailabilitySetClient) CreateOrUpdate(ctx context.Context, name string, avset *compute.AvailabilitySet) (*compute.AvailabilitySet, error) {
	return c.internal.CreateOrUpdate(ctx, name, avset)
}

// Delete methods invokes delete of the compute resource
func (c *AvailabilitySetClient) Delete(ctx context.Context, name string) error {
	return c.internal.Delete(ctx, name)
}
