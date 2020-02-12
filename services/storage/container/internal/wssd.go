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
	wssdstorage.ContainerAgentClient
}

// NewContainerClient - creates a client session with the backend wssd agent
func NewContainerClient(subID string, authorizer auth.Authorizer) (*client, error) {
	c, err := wssdclient.GetContainerClient(&subID, authorizer)
	if err != nil {
		return nil, err
	}
	return &client{c}, nil
}

// Get
func (c *client) Get(ctx context.Context, group, name string) (*[]storage.Container, error) {
	request := getContainerRequest(wssdcommonproto.Operation_GET, name, nil)
	response, err := c.ContainerAgentClient.Invoke(ctx, request)
	if err != nil {
		return nil, err
	}
	return getContainersFromResponse(response), nil
}

// CreateOrUpdate
func (c *client) CreateOrUpdate(ctx context.Context, group, name string, sg *storage.Container) (*storage.Container, error) {
	request := getContainerRequest(wssdcommonproto.Operation_POST, name, sg)
	response, err := c.ContainerAgentClient.Invoke(ctx, request)
	if err != nil {
		log.Errorf("[Container] Create failed with error %v", err)
		return nil, err
	}

	ctainer := getContainersFromResponse(response)

	if len(*ctainer) == 0 {
		return nil, fmt.Errorf("[Container][Create] Unexpected error: Creating a network interface returned no result")
	}

	return &((*ctainer)[0]), err
}

// Delete methods invokes create or update on the client
func (c *client) Delete(ctx context.Context, group, name string) error {
	request := getContainerRequest(wssdcommonproto.Operation_DELETE, name, nil)
	_, err := c.ContainerAgentClient.Invoke(ctx, request)
	return err
}

func getContainersFromResponse(response *wssdstorage.ContainerResponse) *[]storage.Container {
	containers := []storage.Container{}
	for _, ctainer := range response.GetContainers() {
		containers = append(containers, *(getContainer(ctainer)))
	}

	return &containers
}

func getContainerRequest(opType wssdcommonproto.Operation, name string, ctainer *storage.Container) *wssdstorage.ContainerRequest {
	request := &wssdstorage.ContainerRequest{
		OperationType: opType,
		Containers:    []*wssdstorage.Container{},
	}
	if ctainer != nil {
		request.Containers = append(request.Containers, getWssdContainer(ctainer))
	} else if len(name) > 0 {
		request.Containers = append(request.Containers,
			&wssdstorage.Container{
				Name: name,
			})
	}
	return request
}

func getContainer(ctainer *wssdstorage.Container) *storage.Container {
	return &storage.Container{
		ID:   &ctainer.Id,
		Name: &ctainer.Name,
		ContainerProperties: &storage.ContainerProperties{
			Path:              &ctainer.Path,
			ProvisioningState: getContainerProvisioningState(ctainer.Status.ProvisioningStatus),
		},
	}
}

func getContainerProvisioningState(status *wssdcommonproto.ProvisionStatus) *string {
	provisionState := wssdcommonproto.ProvisionState_UNKNOWN
	if status != nil {
		provisionState = status.CurrentState
	}
	stateString := provisionState.String()
	return &stateString
}

func getWssdContainer(ctainer *storage.Container) *wssdstorage.Container {

	var disk wssdstorage.Container

	if ctainer.Name != nil {
		disk.Name = *ctainer.Name
	}
	if ctainer.Path != nil {
		disk.Path = *ctainer.Path
	}

	return &disk
}
