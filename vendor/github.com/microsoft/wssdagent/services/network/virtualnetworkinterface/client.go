// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package virtualnetworkinterface

import (
	"context"
	"github.com/microsoft/wssdagent/pkg/apis/config"
	"github.com/microsoft/wssdagent/pkg/errors"
	"github.com/microsoft/wssdagent/pkg/guid"
	"github.com/microsoft/wssdagent/pkg/marshal"
	"github.com/microsoft/wssdagent/pkg/store"
	"github.com/microsoft/wssdagent/pkg/trace"
	pb "github.com/microsoft/wssdagent/rpc/network"
	"github.com/microsoft/wssdagent/services/network/virtualnetworkinterface/internal"
	"reflect"
	"sync"

	hcn "github.com/microsoft/wssdagent/services/network/virtualnetworkinterface/hcn"
)

const (
	HCNSpec  = "hcn"
	VMMSSpec = "vmms"
)

type Service interface {
	CreateVirtualNetworkInterface(*internal.VirtualNetworkInterfaceInternal) error
	CleanupVirtualNetworkInterface(*internal.VirtualNetworkInterfaceInternal) error
}

type Client struct {
	internal Service
	store    *store.ConfigStore
	config   *config.ChildAgentConfiguration
	mux      sync.Mutex
}

func NewClient() *Client {
	cConfig := config.GetChildAgentConfiguration("VirtualNetworkInterface")
	c := &Client{
		store:  store.NewConfigStore(cConfig.DataStorePath, reflect.TypeOf(internal.VirtualNetworkInterfaceInternal{})),
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

func (c *Client) newVirtualNetworkInterface(vnic *pb.VirtualNetworkInterface) *internal.VirtualNetworkInterfaceInternal {
	return internal.NewVirtualNetworkInterfaceInternal(guid.NewGuid(), c.config.DataStorePath, vnic)
}

// Create or Update the specified virtual network(s)
func (c *Client) Create(ctx context.Context, vnicDef *pb.VirtualNetworkInterface) (newvnic *pb.VirtualNetworkInterface, err error) {
	ctx, span := trace.NewSpan(ctx, "VirtualNetworkInterface", "Create", marshal.ToString(vnicDef))
	defer span.End(err)

	err = c.Validate(ctx, vnicDef)
	if err != nil {
		return
	}
	vnicinternal := c.newVirtualNetworkInterface(vnicDef)

	err = c.internal.CreateVirtualNetworkInterface(vnicinternal)
	if err != nil {
		return
	}
	newvnic = vnicinternal.Entity

	err = c.store.Add(vnicinternal.Id, vnicinternal)

	return

}

// Get all/selected HCS virtual network(s)
func (c *Client) Get(ctx context.Context, vnicDef *pb.VirtualNetworkInterface) (vnics []*pb.VirtualNetworkInterface, err error) {
	ctx, span := trace.NewSpan(ctx, "VirtualNetworkInterface", "Get", marshal.ToString(vnicDef))
	defer span.End(err)

	c.mux.Lock()
	defer c.mux.Unlock()
	vnicName := ""
	if vnicDef != nil {
		vnicName = vnicDef.Name
	}

	vnicsint, err := c.store.ListFilter("Name", vnicName)
	if err != nil {
		return
	}

	for _, val := range *vnicsint {
		vnicint := val.(*internal.VirtualNetworkInterfaceInternal)
		vnics = append(vnics, vnicint.Entity)
	}

	return
}

// Delete the specified virtual network(s)
func (c *Client) Delete(ctx context.Context, vnicDef *pb.VirtualNetworkInterface) (err error) {
	ctx, span := trace.NewSpan(ctx, "VirtualNetworkInterface", "Delete", marshal.ToString(vnicDef))
	defer span.End(err)

	c.mux.Lock()
	defer c.mux.Unlock()
	vnicinternal, err := c.getVirtualNetworkInterfaceInternal(vnicDef.Name)
	if err != nil {
		return
	}

	err = c.internal.CleanupVirtualNetworkInterface(vnicinternal)
	if err != nil {
		return
	}

	err = c.store.Delete(vnicinternal.Id)
	return
}

func (c *Client) getVirtualNetworkInterfaceInternal(name string) (*internal.VirtualNetworkInterfaceInternal, error) {
	vnicsint, err := c.store.ListFilter("Name", name)
	if err != nil {
		return nil, err
	}
	if *vnicsint == nil || len(*vnicsint) == 0 {
		return nil, errors.NotFound
	}

	return (*vnicsint)[0].(*internal.VirtualNetworkInterfaceInternal), nil
}

// Validate
func (c *Client) Validate(ctx context.Context, vnicDef *pb.VirtualNetworkInterface) (err error) {
	ctx, span := trace.NewSpan(ctx, "VirtualNetworkInterface", "Validate", marshal.ToString(vnicDef))
	defer span.End(err)

	err = nil

	if vnicDef == nil {
		err = errors.Wrapf(errors.InvalidInput, "Input VNic definition is nil")
		return
	}

	_, err = c.getVirtualNetworkInterfaceInternal(vnicDef.Name)
	if err != nil && err == errors.NotFound {
		err = nil
	} else {
		err = errors.AlreadyExists
	}
	return
}
