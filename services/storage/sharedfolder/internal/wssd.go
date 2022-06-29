// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license

package internal

import (
	"context"
	"fmt"

	"github.com/microsoft/moc/pkg/status"
	"github.com/microsoft/wssd-sdk-for-go/services/storage"

	"github.com/microsoft/moc/pkg/auth"
	"github.com/microsoft/moc/pkg/errors"
	wssdcommonproto "github.com/microsoft/moc/rpc/common"
	wssdstorage "github.com/microsoft/moc/rpc/nodeagent/storage"
	wssdclient "github.com/microsoft/wssd-sdk-for-go/pkg/client"
	log "k8s.io/klog"
)

type client struct {
	wssdstorage.SharedFolderAgentClient
}

// NewSharedFolderClient - creates a client session with the backend wssd agent
func NewSharedFolderClient(subID string, authorizer auth.Authorizer) (*client, error) {
	c, err := wssdclient.GetSharedFolderClient(&subID, authorizer)
	if err != nil {
		return nil, err
	}
	return &client{c}, nil
}

// Get
func (c *client) Get(ctx context.Context, name string) (*[]storage.SharedFolder, error) {
	request, err := getSharedFolderRequest(wssdcommonproto.Operation_GET, name, nil)
	if err != nil {
		return nil, err
	}
	response, err := c.SharedFolderAgentClient.Invoke(ctx, request)
	if err != nil {
		return nil, err
	}
	return getSharedFoldersFromResponse(response), nil
}

// CreateOrUpdate
func (c *client) CreateOrUpdate(ctx context.Context, name string, sg *storage.SharedFolder) (*storage.SharedFolder, error) {
	request, err := getSharedFolderRequest(wssdcommonproto.Operation_POST, name, sg)
	if err != nil {
		return nil, err
	}
	response, err := c.SharedFolderAgentClient.Invoke(ctx, request)
	if err != nil {
		log.Errorf("[SharedFolder] Create failed with error %v", err)
		return nil, err
	}

	sharedfolder := getSharedFoldersFromResponse(response)

	if len(*sharedfolder) == 0 {
		return nil, fmt.Errorf("[SharedFolder][Create] Unexpected error: Creating a sharedfolder returned no result")
	}

	return &((*sharedfolder)[0]), err
}

// Delete methods invokes create or update on the client
func (c *client) Delete(ctx context.Context, name string) error {
	request, err := getSharedFolderRequest(wssdcommonproto.Operation_DELETE, name, nil)
	if err != nil {
		return err
	}
	_, err = c.SharedFolderAgentClient.Invoke(ctx, request)
	return err
}

func getSharedFoldersFromResponse(response *wssdstorage.SharedFolderResponse) *[]storage.SharedFolder {
	sharedFolders := []storage.SharedFolder{}
	for _, sharedfolder := range response.GetSharedFolderSystems() {
		sharedFolders = append(sharedFolders, *(getSharedFolder(sharedfolder)))
	}

	return &sharedFolders
}

func getSharedFolderRequest(opType wssdcommonproto.Operation, name string, sharedfolder *storage.SharedFolder) (*wssdstorage.SharedFolderRequest, error) {
	request := &wssdstorage.SharedFolderRequest{
		OperationType:       opType,
		SharedFolderSystems: []*wssdstorage.SharedFolder{},
	}

	wssdsharedfolder := &wssdstorage.SharedFolder{
		Name: name,
	}
	var err error
	if sharedfolder != nil {
		wssdsharedfolder, err = getWssdSharedFolder(sharedfolder)
		if err != nil {
			return nil, err
		}
	}
	request.SharedFolderSystems = append(request.SharedFolderSystems, wssdsharedfolder)
	return request, nil
}

func getSharedFolder(sharedfolder *wssdstorage.SharedFolder) *storage.SharedFolder {

	return &storage.SharedFolder{
		ID:   &sharedfolder.Id,
		Name: &sharedfolder.Name,
		SharedFolderProperties: &storage.SharedFolderProperties{
			ContainerName:      &sharedfolder.ContainerName,
			FolderName:         &sharedfolder.FolderName,
			ReadOnly:           &sharedfolder.ReadOnly,
			Path:               &sharedfolder.Path,
			VirtualMachineName: &sharedfolder.VirtualmachineName,
			ProvisioningState:  status.GetProvisioningState(sharedfolder.Status.GetProvisioningStatus()),
			Statuses:           status.GetStatuses(sharedfolder.Status),
		},
	}
}

func getSharedFolderIsPlaceholder(sharedfolder *wssdstorage.SharedFolder) *bool {
	isPlaceholder := false
	entity := sharedfolder.GetEntity()
	if entity != nil {
		isPlaceholder = entity.IsPlaceholder
	}
	return &isPlaceholder
}

func getWssdSharedFolder(sharedfolder *storage.SharedFolder) (*wssdstorage.SharedFolder, error) {
	wssdSharedFolder := wssdstorage.SharedFolder{
		Name: *sharedfolder.Name,
	}

	if sharedfolder.Name == nil {
		return nil, errors.Wrapf(errors.InvalidInput, "Missing Name")
	}

	wssdSharedFolder.Name = *sharedfolder.Name
	wssdSharedFolder.Entity = getWssdSharedFolderEntity(sharedfolder)

	if sharedfolder.ContainerName != nil {
		wssdSharedFolder.ContainerName = *sharedfolder.ContainerName
	}
	if sharedfolder.FolderName != nil {
		wssdSharedFolder.FolderName = *sharedfolder.FolderName
	}
	if sharedfolder.ReadOnly != nil {
		wssdSharedFolder.ReadOnly = *sharedfolder.ReadOnly
	}
	if sharedfolder.VirtualMachineName != nil {
		wssdSharedFolder.VirtualmachineName = *sharedfolder.VirtualMachineName
	}

	return &wssdSharedFolder, nil
}

func getWssdSharedFolderEntity(sharedfolder *storage.SharedFolder) *wssdcommonproto.Entity {
	/*	isPlaceholder := false
		if sharedfolder.SharedFolderProperties != nil && sharedfolder.SharedFolderProperties.IsPlaceholder != nil {
			isPlaceholder = *sharedfolder.SharedFolderProperties.IsPlaceholder
		}
	*/
	return &wssdcommonproto.Entity{
		//IsPlaceholder: isPlaceholder,
	}
}
