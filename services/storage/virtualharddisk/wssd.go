// Copyright 2019 (c) Microsoft and contributors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package virtualharddisk

import (
	"fmt"
	"context"
	"github.com/microsoft/wssd-sdk-for-go/services/storage"

	wssdclient "github.com/microsoft/wssdagent/rpc/client"
	wssdstorage "github.com/microsoft/wssdagent/rpc/storage"
	log "k8s.io/klog"
)

type client struct {
	wssdstorage.VirtualHardDiskAgentClient
}

// newClient - creates a client session with the backend wssd agent
func newVirtualHardDiskClient(subID string) (*client, error) {
	c, err := wssdclient.GetVirtualHardDiskClient(&subID)
	if err != nil {
		return nil, err
	}
	return &client{c}, nil
}

// Get
func (c *client) Get(ctx context.Context, group, name string) (*[]storage.VirtualHardDisk, error) {
	request := getVirtualHardDiskRequest(wssdstorage.Operation_GET, name, nil)
	response, err := c.VirtualHardDiskAgentClient.Invoke(ctx, request)
	if err != nil {
		return nil, err
	}
	return getVirtualHardDisksFromResponse(response), nil
}

// CreateOrUpdate
func (c *client) CreateOrUpdate(ctx context.Context, group, name string, sg *storage.VirtualHardDisk) (*storage.VirtualHardDisk, error) {
	request := getVirtualHardDiskRequest(wssdstorage.Operation_POST, name, sg)
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
	request := getVirtualHardDiskRequest(wssdstorage.Operation_DELETE, name, nil)
	_, err := c.VirtualHardDiskAgentClient.Invoke(ctx, request)
	return err
}

func getVirtualHardDisksFromResponse(response *wssdstorage.VirtualHardDiskResponse) *[]storage.VirtualHardDisk {
	virtualHardDisks := []storage.VirtualHardDisk{}
	for _, vhd := range response.GetVirtualHardDiskSystems() {
		virtualHardDisks = append(virtualHardDisks, *(GetVirtualHardDisk(vhd)))
	}

	return &virtualHardDisks
}

func getVirtualHardDiskRequest(opType wssdstorage.Operation, name string, vhd *storage.VirtualHardDisk) *wssdstorage.VirtualHardDiskRequest {
	request := &wssdstorage.VirtualHardDiskRequest{
		OperationType:   opType,
		VirtualHardDiskSystems: []*wssdstorage.VirtualHardDisk{},
	}
	if vhd != nil {
		request.VirtualHardDiskSystems = append(request.VirtualHardDiskSystems, GetWssdVirtualHardDisk(vhd))
	} else if len(name) > 0 {
		request.VirtualHardDiskSystems = append(request.VirtualHardDiskSystems,
			&wssdstorage.VirtualHardDisk{
				Name: name,
			})
	}
	return request
}

func GetVirtualHardDisk(vhd *wssdstorage.VirtualHardDisk) *storage.VirtualHardDisk {

	return &storage.VirtualHardDisk{
		BaseProperties: storage.BaseProperties{
			ID : &vhd.Id,
			Name: &vhd.Name,
		},
		Source : &vhd.Source,
	}
}

func GetWssdVirtualHardDisk(vhd *storage.VirtualHardDisk) *wssdstorage.VirtualHardDisk {
	return &wssdstorage.VirtualHardDisk{
		//Id : *vhd.BaseProperties.ID,
		Name: *vhd.BaseProperties.Name,
		Source: *vhd.Source,
	}
}
