// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license
package internal

import (
	"context"

	"github.com/microsoft/moc/pkg/auth"
	"github.com/microsoft/moc/pkg/errors"
	wssdcompute "github.com/microsoft/moc/rpc/nodeagent/compute"
	"github.com/microsoft/wssd-sdk-for-go/services/compute"
)

type avsetStore struct {
	avsets map[string]*wssdcompute.AvailabilitySet
}

// this mock client operates on an in-memory store to allow for component testing
// at the calling layers
type mockClient struct {
	avsetStore
}

func NewAvailabilitySetMockClient(cloudFQDN string, authorizer auth.Authorizer) (*mockClient, error) {
	store := avsetStore{
		avsets: make(map[string]*wssdcompute.AvailabilitySet),
	}

	return &mockClient{store}, nil
}

func (c *mockClient) Get(ctx context.Context, name string) (*[]compute.AvailabilitySet, error) {
	// check if the name exists as a key in the store
	if _, ok := c.avsets[name]; ok {
		wssdavset := c.avsets[name]
		avset := getAvailabilitySet(wssdavset)
		// if it does, return the value
		return &[]compute.AvailabilitySet{*avset}, nil
	}

	return nil, errors.NotFound
}

func (c *mockClient) CreateOrUpdate(ctx context.Context, name string, avset *compute.AvailabilitySet) (*compute.AvailabilitySet, error) {
	wssdavset, err := getWssdAvailabilitySet(avset)
	if err != nil {
		return nil, err
	}

	// check if the name exists as a key in the store
	if _, ok := c.avsets[name]; ok {
		// if it does, check that the platform fault domain count is the same
		if c.avsets[name].PlatformFaultDomainCount != wssdavset.PlatformFaultDomainCount {
			return nil, errors.Wrapf(errors.InvalidInput, "PlatformFaultDomainCount cannot be changed")
		}

		// if it does, update the value
		c.avsets[name] = wssdavset
		return avset, nil
	}

	// if it doesn't, create it
	c.avsets[name] = wssdavset
	return avset, nil
}

func (c *mockClient) Delete(ctx context.Context, name string) error {
	// check if the name exists as a key in the store
	if _, ok := c.avsets[name]; ok {
		// if it does, check if it has any VM members
		if len(c.avsets[name].VirtualMachines) > 0 {
			return errors.Wrapf(errors.InUse, "AvailabilitySet %s has VM members, cannot delete an availability set with VM members", name)
		}

		// if it doesn't, delete it
		delete(c.avsets, name)
		return nil
	}

	return errors.NotFound
}
