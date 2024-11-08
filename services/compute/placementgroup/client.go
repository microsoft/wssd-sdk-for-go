// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license

package placementgroup

import (
	"context"

	"github.com/microsoft/moc/pkg/auth"
	"github.com/microsoft/wssd-sdk-for-go/services/compute"
	"github.com/microsoft/wssd-sdk-for-go/services/compute/placementgroup/internal"
)

type Service interface {
	Get(context.Context, string) (*[]compute.PlacementGroup, error)
	CreateOrUpdate(context.Context, string, *compute.PlacementGroup) (*compute.PlacementGroup, error)
	Delete(context.Context, string) error
}

type PlacementGroupClient struct {
	compute.BaseClient
	internal Service
}

func NewPlacementGroupClient(cloudFQDN string, authorizer auth.Authorizer) (*PlacementGroupClient, error) {
	c, err := internal.NewPlacementGroupWssdClient(cloudFQDN, authorizer)
	if err != nil {
		return nil, err
	}

	return &PlacementGroupClient{internal: c}, nil
}

func (c *PlacementGroupClient) Get(ctx context.Context, name string) (*[]compute.PlacementGroup, error) {
	return c.internal.Get(ctx, name)
}

func (c *PlacementGroupClient) CreateOrUpdate(ctx context.Context, name string, pgroup *compute.PlacementGroup) (*compute.PlacementGroup, error) {
	return c.internal.CreateOrUpdate(ctx, name, pgroup)
}

func (c *PlacementGroupClient) Delete(ctx context.Context, name string) error {
	return c.internal.Delete(ctx, name)
}
