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

package virtualmachinescaleset

import (
	"context"
	log "k8s.io/klog"

	"github.com/microsoft/wssd-sdk-for-go/services/compute"
	"github.com/microsoft/wssd-sdk-for-go/services/compute/virtualmachine"
	wssdclient "github.com/microsoft/wssdagent/rpc/client"
	wssdcompute "github.com/microsoft/wssdagent/rpc/compute"
)

type client struct {
	wssdcompute.VirtualMachineScaleSetAgentClient
}

// newClient - creates a client session with the backend wssd agent
func newVirtualMachineScaleSetClient(subID string) (*client, error) {
	c, err := wssdclient.GetVirtualMachineScaleSetClient(&subID)
	if err != nil {
		return nil, err
	}
	return &client{c}, nil
}

// Get
func (c *client) Get(ctx context.Context, name string) (*[]compute.VirtualMachineScaleSet, error) {
	request := getVirtualMachineScaleSetRequest(wssdcompute.Operation_GET, name, nil)
	response, err := c.VirtualMachineScaleSetAgentClient.Invoke(ctx, request)
	if err != nil {
		return nil, err
	}
	log.Infof("[VirtualMachineScaleSet][Get] [%v]", response)
	return getVirtualMachineScaleSetFromResponse(response), nil
}

// CreateOrUpdate
func (c *client) CreateOrUpdate(ctx context.Context, name string, id string, sg *compute.VirtualMachineScaleSet) (*compute.VirtualMachineScaleSet, error) {
	request := getVirtualMachineScaleSetRequest(wssdcompute.Operation_POST, name, sg)
	response, err := c.VirtualMachineScaleSetAgentClient.Invoke(ctx, request)
	log.Infof("[VirtualMachineScaleSet][Create] [%v]", response)
	if err != nil {
		return nil, err
	}
	vmsss := getVirtualMachineScaleSetFromResponse(response)
	return &((*vmsss)[0]), nil
}

// Delete methods invokes create or update on the client
func (c *client) Delete(ctx context.Context, name string, id string) error {
	vmss, err := c.Get(ctx, name)
	if err != nil {
		return err
	}

	request := getVirtualMachineScaleSetRequest(wssdcompute.Operation_DELETE, name, &(*vmss)[0])
	response, err := c.VirtualMachineScaleSetAgentClient.Invoke(ctx, request)
	log.Infof("[VirtualMachineScaleSet][Delete] [%v]", response)
	return err
}

func getVirtualMachineScaleSetFromResponse(response *wssdcompute.VirtualMachineScaleSetResponse) *[]compute.VirtualMachineScaleSet {
	vmsss := []compute.VirtualMachineScaleSet{}
	for _, vmss := range response.GetVirtualMachineScaleSetSystems() {
		vmsss = append(vmsss, getVirtualMachineScaleSet(vmss))
	}

	return &vmsss

}

func getVirtualMachineScaleSetRequest(opType wssdcompute.Operation, name string, vmss *compute.VirtualMachineScaleSet) *wssdcompute.VirtualMachineScaleSetRequest {
	request := &wssdcompute.VirtualMachineScaleSetRequest{
		OperationType:                 opType,
		VirtualMachineScaleSetSystems: []*wssdcompute.VirtualMachineScaleSet{},
	}
	if vmss != nil {
		request.VirtualMachineScaleSetSystems = append(request.VirtualMachineScaleSetSystems, getWssdVirtualMachineScaleSet(vmss))
	} else if len(name) > 0 {
		request.VirtualMachineScaleSetSystems = append(request.VirtualMachineScaleSetSystems,
			&wssdcompute.VirtualMachineScaleSet{
				Name: name,
			})
	}

	return request

}

func getVirtualMachineScaleSet(vmss *wssdcompute.VirtualMachineScaleSet) compute.VirtualMachineScaleSet {
	return compute.VirtualMachineScaleSet{
		BaseProperties: compute.BaseProperties{
			Name: &vmss.Name,
			ID:   &vmss.Id,
		},
		Sku: &compute.Sku{
			Name:     &vmss.Sku.Name,
			Capacity: &vmss.Sku.Capacity,
		},
		VirtualMachineProfile: virtualmachine.GetVirtualMachine(vmss.Virtualmachineprofile),
	}
}

func getWssdVirtualMachineScaleSet(vmss *compute.VirtualMachineScaleSet) *wssdcompute.VirtualMachineScaleSet {
	return &wssdcompute.VirtualMachineScaleSet{
		Name: *(vmss.Name),
		Id:   *(vmss.ID),
		Sku: &wssdcompute.Sku{
			Name:     *(vmss.Sku.Name),
			Capacity: *(vmss.Sku.Capacity),
		},
		Virtualmachineprofile: virtualmachine.GetWssdVirtualMachine(vmss.VirtualMachineProfile),
	}
}
