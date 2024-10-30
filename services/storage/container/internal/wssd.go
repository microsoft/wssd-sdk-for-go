// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license

package internal

import (
	"context"
	"fmt"

	"code.cloudfoundry.org/bytefmt"

	"github.com/microsoft/moc/pkg/status"
	"github.com/microsoft/wssd-sdk-for-go/services/storage"

	"github.com/microsoft/moc/pkg/auth"
	prototags "github.com/microsoft/moc/pkg/tags"
	wssdcommonproto "github.com/microsoft/moc/rpc/common"
	wssdstorage "github.com/microsoft/moc/rpc/nodeagent/storage"
	wssdclient "github.com/microsoft/wssd-sdk-for-go/pkg/client"
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
	var totalSize, availSize string
	if ctainer.Info != nil {
		totalSize = bytefmt.ByteSize(ctainer.Info.Capacity.TotalBytes)
		availSize = bytefmt.ByteSize(ctainer.Info.Capacity.AvailableBytes)
	}
	return &storage.Container{
		ID:   &ctainer.Id,
		Name: &ctainer.Name,
		ContainerProperties: &storage.ContainerProperties{
			Path:              &ctainer.Path,
			ProvisioningState: status.GetProvisioningState(ctainer.Status.GetProvisioningStatus()),
			Statuses:          status.GetStatuses(ctainer.Status),
			ContainerInfo: &storage.ContainerInfo{
				AvailableSize: availSize,
				TotalSize:     totalSize,
				NodeName:      ctainer.Info.NodeName,
			},
			IsPlaceholder: getContainerPlaceHolder(ctainer),
		},
		Tags: prototags.ProtoToMap(ctainer.Tags),
	}
}

func getWssdContainer(ctainer *storage.Container) *wssdstorage.Container {

	wssdctainer := &wssdstorage.Container{
		Tags: prototags.MapToProto(ctainer.Tags),
	}

	if ctainer.Name != nil {
		wssdctainer.Name = *ctainer.Name
	}
	if ctainer.Path != nil {
		wssdctainer.Path = *ctainer.Path
	}

	isPlaceholder := false
	if ctainer.ContainerProperties != nil && ctainer.ContainerProperties.IsPlaceholder != nil {
		isPlaceholder = *ctainer.ContainerProperties.IsPlaceholder
	}

	wssdctainer.Entity = &wssdcommonproto.Entity{
		IsPlaceholder: isPlaceholder,
	}

	if ctainer.ContainerInfo != nil {
		if wssdctainer.Info != nil && wssdctainer.Info.Capacity != nil {
			wssdctainer.Info.Capacity.AvailableBytes, _ = bytefmt.ToBytes(ctainer.ContainerInfo.AvailableSize)
			wssdctainer.Info.Capacity.TotalBytes, _ = bytefmt.ToBytes(ctainer.ContainerInfo.TotalSize)
			wssdctainer.Info.NodeName = ctainer.ContainerInfo.NodeName
		}
	}
	return wssdctainer
}

func getContainerPlaceHolder(ctainer *wssdstorage.Container) *bool {
	isPlaceholder := false
	entity := ctainer.GetEntity()
	if entity != nil {
		isPlaceholder = entity.IsPlaceholder
	}
	return &isPlaceholder
}
