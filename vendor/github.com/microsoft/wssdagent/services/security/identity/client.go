// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package identity

import (
	"context"
	"github.com/microsoft/wssdagent/pkg/apis/config"
	"github.com/microsoft/wssdagent/pkg/errors"
	"github.com/microsoft/wssdagent/pkg/guid"
	"github.com/microsoft/wssdagent/pkg/marshal"
	"github.com/microsoft/wssdagent/pkg/store"
	"github.com/microsoft/wssdagent/pkg/trace"
	pb "github.com/microsoft/wssdagent/rpc/security"
	"github.com/microsoft/wssdagent/services/security/identity/internal"
	"reflect"
	"sync"

	"github.com/microsoft/wssdagent/services/security/identity/ikv"
)

const (
	IKVSpec = "ikv"
)

type Service interface {
	CreateIdentity(context.Context, *internal.IdentityInternal) error
	CleanupIdentity(context.Context, *internal.IdentityInternal) error
}

type Client struct {
	internal 		Service
	store    		*store.ConfigStore
	config   		*config.ChildAgentConfiguration
	mux      		sync.Mutex
}

func NewClient() *Client {
	cConfig := config.GetChildAgentConfiguration("Identity")
	c := &Client{
		store:   store.NewConfigStore(cConfig.DataStorePath, reflect.TypeOf(internal.IdentityInternal{})),
		config:  cConfig,
	}
	switch cConfig.ProviderSpec {
	case IKVSpec:
	default:
		c.internal = ikv.NewClient()
	}
	return c
}

func (c *Client) newIdentity(identity *pb.Identity) *internal.IdentityInternal {
	return internal.NewIdentityInternal(guid.NewGuid(), c.config.DataStorePath, identity)
}

// GetDataStorePath
func (c *Client) GetDataStorePath() string {
	return c.store.GetPath()
}

// Create or Update the specified identity(s)
func (c *Client) Create(ctx context.Context, identityDef *pb.Identity) (newidentity *pb.Identity, err error) {
	ctx, span := trace.NewSpan(ctx, "Identity", "Create", marshal.ToString(identityDef))
	defer span.End(err)

	err = c.Validate(ctx, identityDef)
	if err != nil {
		return
	}
	identityinternal := c.newIdentity(identityDef)

	err = c.internal.CreateIdentity(ctx, identityinternal)
	if err != nil {
		return
	}
	newidentity = identityinternal.Entity

	err = c.store.Add(identityinternal.Id, identityinternal)

	return

}

// Get all/selected identity(s)
func (c *Client) Get(ctx context.Context, securityDef *pb.Identity) (identities []*pb.Identity, err error) {
	ctx, span := trace.NewSpan(ctx, "Identity", "Get", marshal.ToString(securityDef))
	defer span.End(err)

	c.mux.Lock()
	defer c.mux.Unlock()

	identityName := ""
	if securityDef != nil {
		identityName = securityDef.Name
	}

	identitiesint, err := c.store.ListFilter("Name", identityName)
	if err != nil {
		return
	}

	for _, val := range *identitiesint {
		identityint := val.(*internal.IdentityInternal)
		identities = append(identities, identityint.Entity)
	}

	return
}

// Delete the specified virtual security(s)
func (c *Client) Delete(ctx context.Context, securityDef *pb.Identity) (err error) {
	ctx, span := trace.NewSpan(ctx, "Identity", "Delete", marshal.ToString(securityDef))
	defer span.End(err)

	c.mux.Lock()
	defer c.mux.Unlock()

	identityinternal, err := c.getIdentityInternal(securityDef.Name)
	if err != nil {
		return
	}

	err = c.internal.CleanupIdentity(ctx, identityinternal)
	if err != nil {
		// Log this error and continue
		// return
	}

	err = c.store.Delete(identityinternal.Id)
	return
}

// Validate
func (c *Client) Validate(ctx context.Context, securityDef *pb.Identity) (err error) {
	ctx, span := trace.NewSpan(ctx, "Identity", "Validate", marshal.ToString(securityDef))
	defer span.End(err)

	err = nil

	if securityDef == nil {
		err = errors.Wrapf(errors.InvalidInput, "Input group definition is nil")
		return
	}

	_, err = c.getIdentityInternal(securityDef.Name)
	if err != nil && err == errors.NotFound {
		err = nil
	} else {
		err = errors.AlreadyExists
	}

	if err != nil {
		return
	}

	return
}

func (c *Client) getIdentityInternal(name string) (*internal.IdentityInternal, error) {
	identitiesint, err := c.store.ListFilter("Name", name)
	if err != nil {
		return nil, err
	}
	if *identitiesint == nil || len(*identitiesint) == 0 {
		return nil, errors.NotFound
	}

	return (*identitiesint)[0].(*internal.IdentityInternal), nil
}

func (c *Client) pruneStore() {
	identitiesint, err := c.store.List()
	if err != nil {
		return
	}
	if *identitiesint == nil || len(*identitiesint) == 0 {
		return
	}

	for _, _ = range *identitiesint {

	}
}