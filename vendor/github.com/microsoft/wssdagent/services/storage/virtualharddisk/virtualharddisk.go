// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package virtualharddisk

import (
	"context"
	"github.com/microsoft/wssdagent/pkg/errors"
	pb "github.com/microsoft/wssdagent/rpc/storage"
)

type VirtualHardDiskProvider struct {
	client *Client
}

func NewVirtualHardDiskProvider() *VirtualHardDiskProvider {
	return &VirtualHardDiskProvider{
		client: NewClient(),
	}
}

func (vhdProv *VirtualHardDiskProvider) Get(ctx context.Context, vhds []*pb.VirtualHardDisk) ([]*pb.VirtualHardDisk, error) {
	newvhds := []*pb.VirtualHardDisk{}
	if len(vhds) == 0 {
		// Get Everything
		return vhdProv.client.Get(ctx, nil)
	}

	// Get only requested vhds
	for _, vhd := range vhds {
		newvhd, err := vhdProv.client.Get(ctx, vhd)
		if err != nil {
			return newvhds, err
		}
		newvhds = append(newvhds, newvhd[0])
	}
	return newvhds, nil
}

func (vhdProv *VirtualHardDiskProvider) CreateOrUpdate(ctx context.Context, vhds []*pb.VirtualHardDisk) ([]*pb.VirtualHardDisk, error) {
	newvhds := []*pb.VirtualHardDisk{}
	for _, vhd := range vhds {
		newvhd, err := vhdProv.client.Create(ctx, vhd)
		if err != nil {
			if err != errors.AlreadyExists {
				vhdProv.client.Delete(ctx, vhd)
			}
			return newvhds, err
		}
		newvhds = append(newvhds, newvhd)
	}

	return newvhds, nil
}

func (vhdProv *VirtualHardDiskProvider) Delete(ctx context.Context, vhds []*pb.VirtualHardDisk) error {
	for _, vhd := range vhds {
		err := vhdProv.client.Delete(ctx, vhd)
		if err != nil {
			return err
		}
	}

	return nil
}

// DeleteVirtualHardDisk helper to delete virtual hard disks
func (vhdProv *VirtualHardDiskProvider) DeleteVirtualHardDisk(ctx context.Context, vhdName string) error {
	return vhdProv.Delete(ctx, getVirtualHardDisk([]string{vhdName}))
}

func (vhdProv *VirtualHardDiskProvider) GetVirtualHardDisk(ctx context.Context, vhdName string) (*pb.VirtualHardDisk, error) {
	vhdsnew, err := vhdProv.Get(ctx, getVirtualHardDisk([]string{vhdName}))
	if err != nil {
		return nil, err
	}

	if len(vhdsnew) == 0 {
		return nil, errors.NotFound
	}

	return vhdsnew[0], nil
}

// CloneVirtualHardDisk would clone the vhd specified by VhdName to a newVhdName
func (vhdProv *VirtualHardDiskProvider) CloneVirtualHardDisk(ctx context.Context, vhdName, newVhdName string) (*pb.VirtualHardDisk, error) {
	vhd := &pb.VirtualHardDisk{Name: newVhdName, Source: vhdName}
	return vhdProv.client.Create(ctx, vhd)
}

func getVirtualHardDisk(vhds []string) []*pb.VirtualHardDisk {
	tmp := []*pb.VirtualHardDisk{}
	for _, vhd := range vhds {
		tmp = append(tmp, &pb.VirtualHardDisk{Name: vhd})
	}
	return tmp
}
