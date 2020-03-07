// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license

package internal

import (
	"context"
	"fmt"

	"github.com/microsoft/wssd-sdk-for-go/pkg/auth"
	"github.com/microsoft/wssd-sdk-for-go/services/storage"

	wssdclient "github.com/microsoft/wssd-sdk-for-go/pkg/client"
	"github.com/microsoft/wssdagent/pkg/errors"
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
func (c *client) Get(ctx context.Context, containerName, name string) (*[]storage.VirtualHardDisk, error) {
	request, err := getVirtualHardDiskRequest(wssdcommonproto.Operation_GET, name, containerName, nil)
	if err != nil {
		return nil, err
	}
	response, err := c.VirtualHardDiskAgentClient.Invoke(ctx, request)
	if err != nil {
		return nil, err
	}
	return getVirtualHardDisksFromResponse(response), nil
}

// CreateOrUpdate
func (c *client) CreateOrUpdate(ctx context.Context, containerName, name string, sg *storage.VirtualHardDisk) (*storage.VirtualHardDisk, error) {
	request, err := getVirtualHardDiskRequest(wssdcommonproto.Operation_POST, name, containerName, sg)
	if err != nil {
		return nil, err
	}
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
func (c *client) Delete(ctx context.Context, containerName, name string) error {
	request, err := getVirtualHardDiskRequest(wssdcommonproto.Operation_DELETE, name, containerName, nil)
	if err != nil {
		return err
	}
	_, err = c.VirtualHardDiskAgentClient.Invoke(ctx, request)
	return err
}

func getVirtualHardDisksFromResponse(response *wssdstorage.VirtualHardDiskResponse) *[]storage.VirtualHardDisk {
	virtualHardDisks := []storage.VirtualHardDisk{}
	for _, vhd := range response.GetVirtualHardDiskSystems() {
		virtualHardDisks = append(virtualHardDisks, *(getVirtualHardDisk(vhd)))
	}

	return &virtualHardDisks
}

func getVirtualHardDiskRequest(opType wssdcommonproto.Operation, name, containerName string, vhd *storage.VirtualHardDisk) (*wssdstorage.VirtualHardDiskRequest, error) {
	request := &wssdstorage.VirtualHardDiskRequest{
		OperationType:          opType,
		VirtualHardDiskSystems: []*wssdstorage.VirtualHardDisk{},
	}

	wssdvhd := &wssdstorage.VirtualHardDisk{
		Name:          name,
		ContainerName: containerName,
	}
	var err error
	if vhd != nil {
		wssdvhd, err = getWssdVirtualHardDisk(containerName, vhd)
		if err != nil {
			return nil, err
		}
	}
	request.VirtualHardDiskSystems = append(request.VirtualHardDiskSystems, wssdvhd)
	return request, nil
}

func getVirtualHardDisk(vhd *wssdstorage.VirtualHardDisk) *storage.VirtualHardDisk {

	return &storage.VirtualHardDisk{
		ID:   &vhd.Id,
		Name: &vhd.Name,
		VirtualHardDiskProperties: &storage.VirtualHardDiskProperties{
			Source:              &vhd.Source,
			Path:                &vhd.Path,
			DiskSizeBytes:       &vhd.Size,
			Dynamic:             &vhd.Dynamic,
			Blocksizebytes:      &vhd.Blocksizebytes,
			Logicalsectorbytes:  &vhd.Logicalsectorbytes,
			Physicalsectorbytes: &vhd.Physicalsectorbytes,
			Controllernumber:    &vhd.Controllernumber,
			Controllerlocation:  &vhd.Controllerlocation,
			Disknumber:          &vhd.Disknumber,
			VirtualMachineName:  &vhd.VirtualmachineName,
			Scsipath:            &vhd.Scsipath,
			Virtualharddisktype: vhd.Virtualharddisktype.String(),
			ProvisioningState:   getVirtualHardDiskProvisioningState(vhd.Status.GetProvisioningStatus()),
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

func getWssdVirtualHardDisk(containerName string, vhd *storage.VirtualHardDisk) (*wssdstorage.VirtualHardDisk, error) {
	disk := wssdstorage.VirtualHardDisk{
		ContainerName: containerName,
	}

	if vhd.Name == nil {
		return nil, errors.Wrapf(errors.InvalidInput, "Missing Name")
	}

	disk.Name = *vhd.Name
	disk.Virtualharddisktype = getVirtualharddisktype(vhd.Virtualharddisktype)

	if disk.Virtualharddisktype == wssdstorage.VirtualHardDiskType_OS_VIRTUALHARDDISK {
		if vhd.Source == nil {
			return nil, errors.Wrapf(errors.InvalidInput, "Missing Source")
		}
		disk.Source = *vhd.Source
	} else {
		if vhd.DiskSizeBytes == nil {
			return nil, errors.Wrapf(errors.InvalidInput, "Missing DiskSize")
		}
		disk.Size = *vhd.DiskSizeBytes
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
		if vhd.VirtualMachineName != nil {
			disk.VirtualmachineName = *vhd.VirtualMachineName
		}
	}

	return &disk, nil
}

func getVirtualharddisktype(enum string) wssdstorage.VirtualHardDiskType {
	typevalue := wssdstorage.VirtualHardDiskType(0)
	typevTmp, ok := wssdstorage.VirtualHardDiskType_value[enum]
	if ok {
		typevalue = wssdstorage.VirtualHardDiskType(typevTmp)
	}
	return typevalue
}
