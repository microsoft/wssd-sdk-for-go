// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package internal

import (
	"context"
	"fmt"

	"github.com/microsoft/moc/pkg/auth"
	prototags "github.com/microsoft/moc/pkg/tags"
	wssdcommonproto "github.com/microsoft/moc/rpc/common"
	wssdcompute "github.com/microsoft/moc/rpc/nodeagent/compute"
	wssdclient "github.com/microsoft/wssd-sdk-for-go/pkg/client"
	"github.com/microsoft/wssd-sdk-for-go/services/compute"
)

type client struct {
	wssdcompute.VirtualMachineAgentClient
}

// newVirtualMachineClient - creates a client session with the backend wssd agent
func NewVirtualMachineClient(subID string, authorizer auth.Authorizer) (*client, error) {
	c, err := wssdclient.GetVirtualMachineClient(&subID, authorizer)
	if err != nil {
		return nil, err
	}
	return &client{c}, nil
}

// Get
func (c *client) Get(ctx context.Context, group, name string) (*[]compute.VirtualMachine, error) {
	request, err := c.getVirtualMachineRequest(wssdcommonproto.Operation_GET, name, nil)
	if err != nil {
		return nil, err
	}
	response, err := c.VirtualMachineAgentClient.Invoke(ctx, request)
	if err != nil {
		return nil, err
	}
	return c.getVirtualMachineFromResponse(response), nil

}

// Get
func (c *client) get(ctx context.Context, group, name string) ([]*wssdcompute.VirtualMachine, error) {
	request, err := c.getVirtualMachineRequest(wssdcommonproto.Operation_GET, name, nil)
	if err != nil {
		return nil, err
	}
	response, err := c.VirtualMachineAgentClient.Invoke(ctx, request)
	if err != nil {
		return nil, err
	}

	return response.GetVirtualMachineSystems(), nil
}

// CreateOrUpdate
func (c *client) CreateOrUpdate(ctx context.Context, group, name string, sg *compute.VirtualMachine) (*compute.VirtualMachine, error) {
	request, err := c.getVirtualMachineRequest(wssdcommonproto.Operation_POST, name, sg)
	if err != nil {
		return nil, err
	}
	response, err := c.VirtualMachineAgentClient.Invoke(ctx, request)
	if err != nil {
		return nil, err
	}
	vms := c.getVirtualMachineFromResponse(response)
	if len(*vms) == 0 {
		return nil, fmt.Errorf("Creation of Virtual Machine failed to unknown reason.")
	}

	return &(*vms)[0], nil
}

// Delete methods invokes create or update on the client
func (c *client) Delete(ctx context.Context, group, name string) error {
	vm, err := c.Get(ctx, group, name)
	if err != nil {
		return err
	}
	if len(*vm) == 0 {
		return fmt.Errorf("Virtual Machine [%s] not found", name)
	}

	request, err := c.getVirtualMachineRequest(wssdcommonproto.Operation_DELETE, name, &(*vm)[0])
	if err != nil {
		return err
	}
	_, err = c.VirtualMachineAgentClient.Invoke(ctx, request)

	return err
}

func (c *client) Start(ctx context.Context, group, name string) (err error) {
	request, err := c.getVirtualMachineOperationRequest(ctx, wssdcommonproto.VirtualMachineOperation_START, name)
	if err != nil {
		return
	}
	_, err = c.VirtualMachineAgentClient.Operate(ctx, request)
	return
}

func (c *client) Stop(ctx context.Context, group, name string) (err error) {
	request, err := c.getVirtualMachineOperationRequest(ctx, wssdcommonproto.VirtualMachineOperation_STOP, name)
	if err != nil {
		return
	}
	_, err = c.VirtualMachineAgentClient.Operate(ctx, request)
	return
}

func (c *client) getVirtualMachineFromResponse(response *wssdcompute.VirtualMachineResponse) *[]compute.VirtualMachine {
	vms := []compute.VirtualMachine{}
	for _, vm := range response.GetVirtualMachineSystems() {
		vms = append(vms, *(c.getVirtualMachine(vm)))
	}

	return &vms
}

func (c *client) getVirtualMachineRequest(opType wssdcommonproto.Operation, name string, vmss *compute.VirtualMachine) (*wssdcompute.VirtualMachineRequest, error) {
	request := &wssdcompute.VirtualMachineRequest{
		OperationType:         opType,
		VirtualMachineSystems: []*wssdcompute.VirtualMachine{},
	}
	if vmss != nil {
		wssdvm, err := c.getWssdVirtualMachine(vmss)
		if err != nil {
			return nil, err
		}
		request.VirtualMachineSystems = append(request.VirtualMachineSystems, wssdvm)
	} else if len(name) > 0 {
		wssdvm := &wssdcompute.VirtualMachine{
			Name: name,
		}
		request.VirtualMachineSystems = append(request.VirtualMachineSystems, wssdvm)
	}

	return request, nil
}

func (c *client) getVirtualMachineOperationRequest(ctx context.Context, opType wssdcommonproto.VirtualMachineOperation, name string) (request *wssdcompute.VirtualMachineOperationRequest, err error) {
	vms, err := c.get(ctx, "", name)
	if err != nil {
		return
	}

	request = &wssdcompute.VirtualMachineOperationRequest{
		OperationType:   opType,
		VirtualMachines: vms,
	}

	return
}

func getComputeTags(tags *wssdcommonproto.Tags) map[string]*string {
	return prototags.ProtoToMap(tags)
}

func getWssdTags(tags map[string]*string) *wssdcommonproto.Tags {
	return prototags.MapToProto(tags)
}
