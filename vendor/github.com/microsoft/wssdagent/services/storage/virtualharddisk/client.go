// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package virtualharddisk

import (
	"context"
	"github.com/microsoft/wssdagent/pkg/apis/config"
	"github.com/microsoft/wssdagent/pkg/errors"
	"github.com/microsoft/wssdagent/pkg/guid"
	"github.com/microsoft/wssdagent/pkg/marshal"
	"github.com/microsoft/wssdagent/pkg/store"
	"github.com/microsoft/wssdagent/pkg/trace"
	pb "github.com/microsoft/wssdagent/rpc/storage"
	"os"
	"path/filepath"
	"reflect"
	"sync"

	"github.com/microsoft/wssdagent/services/storage/virtualharddisk/hcs"
	"github.com/microsoft/wssdagent/services/storage/virtualharddisk/internal"
)

const (
	HCSSpec  = "hcs"
	VMMSSpec = "vmms"
)

type Service interface {
	CreateVirtualHardDisk(*internal.VirtualHardDiskInternal) error
	CleanupVirtualHardDisk(*internal.VirtualHardDiskInternal) error
}

type Client struct {
	internal        Service
	store           *store.ConfigStore
	config          *config.ChildAgentConfiguration
	storageLocation string
	mux             sync.Mutex
}

func NewClient() *Client {
	cConfig := config.GetChildAgentConfiguration("VirtualHardDisk")
	c := &Client{
		store:  store.NewConfigStore(cConfig.DataStorePath, reflect.TypeOf(internal.VirtualHardDiskInternal{})),
		config: cConfig,
	}
	switch cConfig.ProviderSpec {
	case VMMSSpec:
	case HCSSpec:
	default:
		c.internal = hcs.NewClient()
	}

	agentConfig := config.GetAgentConfiguration()
	c.storageLocation = agentConfig.ImageStorePath
	return c
}

func (c *Client) newVirtualHardDisk(vhd *pb.VirtualHardDisk) *internal.VirtualHardDiskInternal {
	return internal.NewVirtualHardDiskInternal(guid.NewGuid(), c.config.DataStorePath, vhd)
}

// Create or Update the specified virtual virtualHardDisk(s)
func (c *Client) Create(ctx context.Context, vhdDef *pb.VirtualHardDisk) (newvhd *pb.VirtualHardDisk, err error) {
	ctx, span := trace.NewSpan(ctx, "VirtualHardDisk", "Create", marshal.ToString(vhdDef))
	defer span.End(err)

	err = c.Validate(ctx, vhdDef)
	if err != nil {
		return
	}

	vhdinternal := c.newVirtualHardDisk(vhdDef)

	if vhdDef.Path == "" {
		vhdDef.Path = c.generateFilePath(vhdDef.Name + ".vhdx") // TODO: append the correct extension
	}

	// We check in case the source is actually an Id, and swap the path if it is
	// we don't care about an error in this case
	vhdSource, _ := c.getVirtualHardDiskInternal(vhdDef.Source)
	if vhdSource != nil {
		vhdDef.Source = vhdSource.Entity.Path
	}

	err = c.internal.CreateVirtualHardDisk(vhdinternal)
	if err != nil {
		return
	}
	newvhd = vhdinternal.Entity

	err = c.store.Add(vhdinternal.Id, vhdinternal)

	return

}

// Get all/selected HCS virtual virtualHardDisk(s)
func (c *Client) Get(ctx context.Context, virtualHardDiskDef *pb.VirtualHardDisk) (vhds []*pb.VirtualHardDisk, err error) {
	ctx, span := trace.NewSpan(ctx, "VirtualHardDisk", "Get", marshal.ToString(virtualHardDiskDef))
	defer span.End(err)

	c.mux.Lock()
	defer c.mux.Unlock()

	vhdName := ""
	if virtualHardDiskDef != nil {
		vhdName = virtualHardDiskDef.Name
	}

	vhdsint, err := c.store.ListFilter("Name", vhdName)
	if err != nil {
		return
	}

	for _, val := range *vhdsint {
		vhdint := val.(*internal.VirtualHardDiskInternal)
		vhds = append(vhds, vhdint.Entity)
	}

	return
}

// Delete the specified virtual virtualHardDisk(s)
func (c *Client) Delete(ctx context.Context, virtualHardDiskDef *pb.VirtualHardDisk) (err error) {
	ctx, span := trace.NewSpan(ctx, "VirtualHardDisk", "Delete", marshal.ToString(virtualHardDiskDef))
	defer span.End(err)

	c.mux.Lock()
	defer c.mux.Unlock()

	vhdinternal, err := c.getVirtualHardDiskInternal(virtualHardDiskDef.Name)
	if err != nil {
		return
	}

	err = c.internal.CleanupVirtualHardDisk(vhdinternal)
	if err != nil {
		return
	}

	err = c.store.Delete(vhdinternal.Id)
	return
}

func (c *Client) getVirtualHardDiskInternal(name string) (*internal.VirtualHardDiskInternal, error) {
	vhdsint, err := c.store.ListFilter("Name", name)
	if err != nil {
		return nil, err
	}
	if *vhdsint == nil || len(*vhdsint) == 0 {
		return nil, errors.NotFound
	}

	return (*vhdsint)[0].(*internal.VirtualHardDiskInternal), nil
}

// Validate
func (c *Client) Validate(ctx context.Context, vhdDef *pb.VirtualHardDisk) (err error) {
	ctx, span := trace.NewSpan(ctx, "VirtualHardDisk", "Validate", marshal.ToString(vhdDef))
	defer span.End(err)

	err = nil

	if vhdDef == nil {
		err = errors.Wrapf(errors.InvalidInput, "Input vhd definition is nil")
		return
	}

	_, err = c.getVirtualHardDiskInternal(vhdDef.Name)
	if err != nil && err == errors.NotFound {
		err = nil
	} else {
		err = errors.AlreadyExists
	}
	if err != nil {
		return
	}

	if len(vhdDef.Source) == 0 {
		err = errors.Wrapf(errors.InvalidInput, "VhdDef Source not defined")
		return
	}
	return
}

func (c *Client) generateFilePath(fileName string) string {
	// Will not recreate or error if already exits
	os.MkdirAll(c.storageLocation, os.ModeDir)
	return filepath.Join(c.storageLocation, fileName)
}
