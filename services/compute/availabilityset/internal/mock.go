// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license
package internal

import (
	"context"

	"github.com/microsoft/moc/pkg/auth"
	"github.com/microsoft/moc/pkg/errors"
	"github.com/microsoft/wssd-sdk-for-go/services/compute"
)

type avsetStore struct {
	avsets map[string]*compute.AvailabilitySet
}

// this mock client operates on an in-memory store to allow for component testing
// at the calling layers
type mockClient struct {
	avsetStore
}

func NewAvailabilitySetMockClient(cloudFQDN string, authorizer auth.Authorizer) (*mockClient, error) {
	store := avsetStore{
		avsets: make(map[string]*compute.AvailabilitySet),
	}

	return &mockClient{store}, nil
}

func (c *mockClient) Get(ctx context.Context, name string) (*[]compute.AvailabilitySet, error) {
	// check if the name exists as a key in the store
	if _, ok := c.avsets[name]; ok {
		// if it does, return the value
		return &[]compute.AvailabilitySet{*c.avsets[name]}, nil
	}

	return nil, errors.NotFound
}

func (c *mockClient) CreateOrUpdate(ctx context.Context, name string, avset *compute.AvailabilitySet) (*compute.AvailabilitySet, error) {
	if avset == nil {
		return nil, errors.Wrapf(errors.InvalidInput, "AvailabilitySet cannot be nil")
	}

	if len(name) == 0 || len(name) > 200 {
		return nil, errors.Wrapf(errors.InvalidInput, "Name cannot be empty or more than 200 characters")
	}

	// check if the name exists as a key in the store
	if _, ok := c.avsets[name]; ok {
		// if it does, update the value
		c.avsets[name] = avset
		return avset, nil
	}

	// if it doesn't, create it
	c.avsets[name] = avset
	return avset, nil
}

func (c *mockClient) Delete(ctx context.Context, name string) error {
	// check if the name exists as a key in the store
	if _, ok := c.avsets[name]; ok {
		// if it does, check if it has any VM members
		if len(c.avsets[name].Properties.VirtualMachines) > 0 {
			return errors.InUse
		}

		// if it doesn't, delete it
		delete(c.avsets, name)
		return nil
	}

	return errors.NotFound
}
