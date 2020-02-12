// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license

package internal

import (
	"context"
	"fmt"

	"github.com/microsoft/wssd-sdk-for-go/pkg/auth"
	"github.com/microsoft/wssd-sdk-for-go/services/storage"

	wssdclient "github.com/microsoft/wssd-sdk-for-go/pkg/client"
	wssdcommonproto "github.com/microsoft/wssdagent/rpc/common"
	wssdstorage "github.com/microsoft/wssdagent/rpc/storage"
	log "k8s.io/klog"
)

type client struct {
	wssdstorage.VirtualHardDiskAgentClient
}

// NewVirtualHardDiskClient - creates a client session with the backend wssd agent
func NewVirtualHardDiskClient(subID string, authorizer auth.Authorizer) (*client, error) {
	c, err := wssdclient.GetVirtualHardDiskClient(&subID, authorizer)
	if err != nil {
		return nil, err
	}
	return &client{c}, nil
}

// Get
func (c *client) Get(ctx context.Context, group, name string) (*[]storage.VirtualHardDisk, error) {
	request := getVirtualHardDiskRequest(wssdcommonproto.Operation_GET, name, nil)
	response, err := c.VirtualHardDiskAgentClient.Invoke(ctx, request)
	if err != nil {
		return nil, err
	}
	return getVirtualHardDisksFromResponse(response), nil
}

// CreateOrUpdate
func (c *client) CreateOrUpdate(ctx context.Context, group, name string, sg *storage.VirtualHardDisk) (*storage.VirtualHardDisk, error) {
	request := getVirtualHardDiskRequest(wssdcommonproto.Operation_POST, name, sg)
	response, err := c.VirtualHardDiskAgentClient.Invoke(ctx, request)
	if err != nil {
		log.Errorf("[VirtualHardDisk] Create failed with error %v", err)
		return nil, err
	}

	vhd := getVirtualHardDisksFromResponse(response)

	if len(*vhd) == 0 {
		return nil, fmt.Errorf("[VirtualHardDisk][Create] Unexpected error: Creating a network interface returned no result")
	}

	return &((*vhd)[0]), err
}

// Delete methods invokes create or update on the client
func (c *client) Delete(ctx context.Context, group, name string) error {
	request := getVirtualHardDiskRequest(wssdcommonproto.Operation_DELETE, name, nil)
	_, err := c.VirtualHardDiskAgentClient.Invoke(ctx, request)
	return err
}

func getVirtualHardDisksFromResponse(response *wssdstorage.VirtualHardDiskResponse) *[]storage.VirtualHardDisk {
	virtualHardDisks := []storage.VirtualHardDisk{}
	for _, vhd := range response.GetVirtualHardDiskSystems() {
		virtualHardDisks = append(virtualHardDisks, *(getVirtualHardDisk(vhd)))
	}

	return &virtualHardDisks
}

func getVirtualHardDiskRequest(opType wssdcommonproto.Operation, name string, vhd *storage.VirtualHardDisk) *wssdstorage.VirtualHardDiskRequest {
	request := &wssdstorage.VirtualHardDiskRequest{
		OperationType:          opType,
		VirtualHardDiskSystems: []*wssdstorage.VirtualHardDisk{},
	}
	if vhd != nil {
		request.VirtualHardDiskSystems = append(request.VirtualHardDiskSystems, getWssdVirtualHardDisk(vhd))
	} else if len(name) > 0 {
		request.VirtualHardDiskSystems = append(request.VirtualHardDiskSystems,
			&wssdstorage.VirtualHardDisk{
				Name: name,
			})
	}
	return request
}

func getVirtualHardDisk(vhd *wssdstorage.VirtualHardDisk) *storage.VirtualHardDisk {

	return &storage.VirtualHardDisk{
		ID:   &vhd.Id,
		Name: &vhd.Name,
		VirtualHardDiskProperties: &storage.VirtualHardDiskProperties{
			Source:              &vhd.Source,
			Path:                &vhd.Path,
			DiskSizeGB:          &vhd.Size,
			Dynamic:             &vhd.Dynamic,
			Blocksizebytes:      &vhd.Blocksizebytes,
			Logicalsectorbytes:  &vhd.Logicalsectorbytes,
			Physicalsectorbytes: &vhd.Physicalsectorbytes,
			Controllernumber:    &vhd.Controllernumber,
			Controllerlocation:  &vhd.Controllerlocation,
			Disknumber:          &vhd.Disknumber,
			Vmname:              &vhd.Vmname,
			Scsipath:            &vhd.Scsipath,
			Virtualharddisktype: vhd.Virtualharddisktype.String(),
			ProvisioningState:   getVirtualHardDiskProvisioningState(vhd.Status.ProvisioningStatus),
		},
	}
}

func getVirtualHardDiskProvisioningState(status *wssdcommonproto.ProvisionStatus) *string {
	provisionState := wssdcommonproto.ProvisionState_UNKNOWN
	if status != nil {
		provisionState = status.CurrentState
	}
	stateString := provisionState.String()
	return &stateString
}

func getWssdVirtualHardDisk(vhd *storage.VirtualHardDisk) *wssdstorage.VirtualHardDisk {

	var disk wssdstorage.VirtualHardDisk

	if vhd.Name != nil {
		disk.Name = *vhd.Name
	}
	if vhd.Source != nil {
		disk.Source = *vhd.Source
	}
	if vhd.Path != nil {
		disk.Path = *vhd.Path
	}
	if vhd.DiskSizeGB != nil {
		disk.Size = *vhd.DiskSizeGB
	}
	if vhd.Dynamic != nil {
		disk.Dynamic = *vhd.Dynamic
	}
	if vhd.Blocksizebytes != nil {
		disk.Blocksizebytes = *vhd.Blocksizebytes
	}
	if vhd.Logicalsectorbytes != nil {
		disk.Logicalsectorbytes = *vhd.Logicalsectorbytes
	}
	if vhd.Physicalsectorbytes != nil {
		disk.Physicalsectorbytes = *vhd.Physicalsectorbytes
	}
	if vhd.Controllerlocation != nil {
		disk.Controllerlocation = *vhd.Controllerlocation
	}
	if vhd.Controllernumber != nil {
		disk.Controllernumber = *vhd.Controllernumber
	}
	if vhd.Disknumber != nil {
		disk.Disknumber = *vhd.Disknumber
	}
	if vhd.Vmname != nil {
		disk.Vmname = *vhd.Vmname
	}

	if vhd.Scsipath != nil {
		disk.Scsipath = *vhd.Scsipath
	}

	if vhd.Virtualharddisktype != "" {
		disk.Virtualharddisktype = getVirtualharddisktype(vhd.Virtualharddisktype)
	}

	return &disk
	// 	return &wssdstorage.VirtualHardDisk{
	// 		Name:                *vhd.Name,
	// 		Source:              *vhd.Source,
	// 		Path:                *vhd.Path,
	// 		Size:                *vhd.DiskSizeGB,
	// 		Dynamic:             *vhd.Dynamic,
	// 		Blocksizebytes:      *vhd.Blocksizebytes,
	// 		Logicalsectorbytes:  *vhd.Logicalsectorbytes,
	// 		Physicalsectorbytes: *vhd.Physicalsectorbytes,
	// 		Controllerlocation:  *vhd.Controllerlocation,
	// 		Controllernumber:    *vhd.Controllernumber,
	// 		Disknumber:          *vhd.Disknumber,
	// 		Vmname:              *vhd.Vmname,
	// 		Scsipath:            *vhd.Scsipath,
	// 		Virtualharddisktype: *vhd.Virtualharddisktype,
	// 	}
}

func getVirtualharddisktype(enum string) wssdstorage.VirtualHardDiskType {
	typevalue := wssdstorage.VirtualHardDiskType(0)
	typevTmp, ok := wssdstorage.VirtualHardDiskType_value[enum]
	if ok {
		typevalue = wssdstorage.VirtualHardDiskType(typevTmp)
	}
	return typevalue
}
