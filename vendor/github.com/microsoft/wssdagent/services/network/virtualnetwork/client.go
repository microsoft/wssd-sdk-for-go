// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package virtualnetwork

import (
	"context"
	"github.com/microsoft/wssdagent/pkg/apis/config"
	"github.com/microsoft/wssdagent/pkg/errors"
	"github.com/microsoft/wssdagent/pkg/guid"
	"github.com/microsoft/wssdagent/pkg/marshal"
	"github.com/microsoft/wssdagent/pkg/store"
	"github.com/microsoft/wssdagent/pkg/trace"
	pb "github.com/microsoft/wssdagent/rpc/network"
	"github.com/microsoft/wssdagent/services/network/virtualnetwork/internal"
	"reflect"
	"sync"

	hcn "github.com/microsoft/wssdagent/services/network/virtualnetwork/hcn"
)

const (
	HCNSpec  = "hcn"
	VMMSSpec = "vmms"
)

type Service interface {
	CreateVirtualNetwork(*internal.VirtualNetworkInternal) error
	CleanupVirtualNetwork(*internal.VirtualNetworkInternal) error
	HasVirtualNetwork(*pb.VirtualNetwork) bool
}

type Client struct {
	internal Service
	store    *store.ConfigStore
	config   *config.ChildAgentConfiguration
	mux      sync.Mutex
}

func NewClient() *Client {
	cConfig := config.GetChildAgentConfiguration("VirtualNetwork")
	c := &Client{
		store:  store.NewConfigStore(cConfig.DataStorePath, reflect.TypeOf(internal.VirtualNetworkInternal{})),
		config: cConfig,
	}
	switch cConfig.ProviderSpec {
	case VMMSSpec:
	case HCNSpec:
	default:
		c.internal = hcn.NewClient()
	}
	return c
}

func (c *Client) newVirtualNetwork(vnet *pb.VirtualNetwork) *internal.VirtualNetworkInternal {
	return internal.NewVirtualNetworkInternal(guid.NewGuid(), c.config.DataStorePath, vnet)
}

// Create or Update the specified virtual network(s)
func (c *Client) Create(ctx context.Context, vnetDef *pb.VirtualNetwork) (newvnet *pb.VirtualNetwork, err error) {
	ctx, span := trace.NewSpan(ctx, "VirtualNetwork", "Create", marshal.ToString(vnetDef))
	defer span.End(err)

	err = c.Validate(ctx, vnetDef)
	if err != nil {
		return
	}
	vnetinternal := c.newVirtualNetwork(vnetDef)

	err = c.internal.CreateVirtualNetwork(vnetinternal)
	if err != nil {
		return
	}
	newvnet = vnetinternal.Entity

	err = c.store.Add(vnetinternal.Id, vnetinternal)

	return

}

// Get all/selected HCS virtual network(s)
func (c *Client) Get(ctx context.Context, networkDef *pb.VirtualNetwork) (vnets []*pb.VirtualNetwork, err error) {
	ctx, span := trace.NewSpan(ctx, "VirtualNetwork", "Get", marshal.ToString(networkDef))
	defer span.End(err)

	c.mux.Lock()
	defer c.mux.Unlock()

	vnetName := ""
	if networkDef != nil {
		vnetName = networkDef.Name
	}

	vnetsint, err := c.store.ListFilter("Name", vnetName)
	if err != nil {
		return
	}

	for _, val := range *vnetsint {
		vnetint := val.(*internal.VirtualNetworkInternal)
		vnets = append(vnets, vnetint.Entity)
	}

	return
}

// Delete the specified virtual network(s)
func (c *Client) Delete(ctx context.Context, networkDef *pb.VirtualNetwork) (err error) {
	ctx, span := trace.NewSpan(ctx, "VirtualNetwork", "Delete", marshal.ToString(networkDef))
	defer span.End(err)

	c.mux.Lock()
	defer c.mux.Unlock()

	vnetinternal, err := c.getVirtualNetworkInternal(networkDef.Name)
	if err != nil {
		return
	}

	err = c.internal.CleanupVirtualNetwork(vnetinternal)
	if err != nil {
		// Log this error and continue
		// return
	}

	err = c.store.Delete(vnetinternal.Id)
	return
}

// Validate
func (c *Client) Validate(ctx context.Context, networkDef *pb.VirtualNetwork) (err error) {
	ctx, span := trace.NewSpan(ctx, "VirtualNetwork", "Validate", marshal.ToString(networkDef))
	defer span.End(err)

	err = nil

	if networkDef == nil {
		err = errors.Wrapf(errors.InvalidInput, "Input group definition is nil")
		return
	}

	_, err = c.getVirtualNetworkInternal(networkDef.Name)
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

func (c *Client) getVirtualNetworkInternal(name string) (*internal.VirtualNetworkInternal, error) {
	vnetsint, err := c.store.ListFilter("Name", name)
	if err != nil {
		return nil, err
	}
	if *vnetsint == nil || len(*vnetsint) == 0 {
		return nil, errors.NotFound
	}

	return (*vnetsint)[0].(*internal.VirtualNetworkInternal), nil
}

func (c *Client) pruneStore() {
	vnetsint, err := c.store.List()
	if err != nil {
		return
	}
	if *vnetsint == nil || len(*vnetsint) == 0 {
		return
	}

	for _, _ = range *vnetsint {

	}
}
