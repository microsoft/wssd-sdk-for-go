// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license

package internal

import (
	"context"
	"fmt"
	"github.com/microsoft/wssd-sdk-for-go/pkg/auth"
	"github.com/microsoft/wssd-sdk-for-go/services/storage"

	wssdclient "github.com/microsoft/wssd-sdk-for-go/pkg/client"
	wssdstorage "github.com/microsoft/wssdagent/rpc/storage"
	wssdcommonproto "github.com/microsoft/wssdagent/rpc/common"
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
			Source: &vhd.Source,
			ProvisioningState: getVirtualHardDiskProvisioningState(vhd.ProvisionStatus),
		},
	}
}

func getVirtualHardDiskProvisioningState(status *wssdcommonproto.ProvisionStatus) (*string) {
	provisionState := wssdcommonproto.ProvisionState_UNKNOWN
	if status != nil {
		provisionState = status.CurrentState
	}
	stateString := provisionState.String()
	return &stateString
}


func getWssdVirtualHardDisk(vhd *storage.VirtualHardDisk) *wssdstorage.VirtualHardDisk {
	return &wssdstorage.VirtualHardDisk{
		Name:   *vhd.Name,
		Source: *vhd.Source,
	}
}
