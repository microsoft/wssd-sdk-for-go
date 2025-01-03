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
	prototags "github.com/microsoft/moc/pkg/tags"
	"github.com/microsoft/moc/rpc/common"
	wssdcommonproto "github.com/microsoft/moc/rpc/common"
	wssdstorage "github.com/microsoft/moc/rpc/nodeagent/storage"
	wssdclient "github.com/microsoft/wssd-sdk-for-go/pkg/client"
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

// Get
func (c *client) get(ctx context.Context, containerName, name string) ([]*wssdstorage.VirtualHardDisk, error) {
	request, err := getVirtualHardDiskRequest(wssdcommonproto.Operation_GET, name, containerName, nil)
	if err != nil {
		return nil, err
	}
	response, err := c.VirtualHardDiskAgentClient.Invoke(ctx, request)
	if err != nil {
		return nil, err
	}
	return response.GetVirtualHardDiskSystems(), nil
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

// Upload
func (c *client) Upload(ctx context.Context, containerName, name string, targeturl string) error {
	fmt.Printf("wssd-sdk-for-go: wssd.go: Creating operation request %s to %s\n", name, targeturl)
	if targeturl == "" {
		return errors.Wrapf(errors.InvalidInput, "Target URL cannot be empty")
	}
	request, err := c.getVirtualHardDiskOperationRequest(ctx, wssdcommonproto.VirtualHardDiskOperation_UPLOAD, name, containerName, targeturl)
	if err != nil {
		return err
	}
	_, err = c.VirtualHardDiskAgentClient.Operate(ctx, request)
	if err != nil {
		log.Errorf("[VirtualHardDisk] Upload failed with error %v", err)
		return err
	}
	fmt.Printf("wssd-sdk-for-go: wssd.go: Done with Operate\n")

	return nil
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

// The interface for the hydrate call takes the container name and the name of the disk file.
// Ultimately, we need the full path on disk to the disk file which we assemble from the path of the container plus the file name of the disk.
// (e.g. "C:\ClusterStorage\Userdata_1\abc123" for the container path and "my_disk.vhd" for the disk name)
func (c *client) Hydrate(ctx context.Context, containerName, name string, vhdDef *storage.VirtualHardDisk) (*storage.VirtualHardDisk, error) {
	request, err := getVirtualHardDiskRequest(wssdcommonproto.Operation_HYDRATE, name, containerName, vhdDef)
	if err != nil {
		return nil, err
	}
	response, err := c.VirtualHardDiskAgentClient.Invoke(ctx, request)
	if err != nil {
		return nil, err
	}

	vhd := getVirtualHardDisksFromResponse(response)
	if len(*vhd) == 0 {
		return nil, fmt.Errorf("[VirtualHardDisk][Hydrate] Unexpected error: Hydrating a disk returned no result")
	}

	return &((*vhd)[0]), err
}

func (c *client) getVirtualHardDiskOperationRequest(ctx context.Context, opType wssdcommonproto.VirtualHardDiskOperation, name, containerName string, targeturl string) (*wssdstorage.VirtualHardDiskOperationRequest, error) {
	var err error
	vhd, err := c.get(ctx, containerName, name)

	if len(vhd) != 1 {
		err = errors.Wrapf(errors.InvalidInput, "Multiple or No Virtual Hard Disks found in container %s with name %s", containerName, name)
		return nil, err
	}

	fmt.Printf("wssd-sdk-for-go: wssd.go: getVirtualHardDiskOperationRequest: Get worked on the vhd\n")

	vhd[0].TargetUrl = targeturl

	request := &wssdstorage.VirtualHardDiskOperationRequest{
		OperationType:    opType,
		VirtualHardDisks: vhd,
	}

	return request, nil
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
	fmt.Printf("wssd-sdk-for-go: wssd.go: getVirtualHardDisk: Enter\n")

	return &storage.VirtualHardDisk{
		ID:   &vhd.Id,
		Name: &vhd.Name,
		Tags: getComputeTags(vhd.GetTags()),
		VirtualHardDiskProperties: &storage.VirtualHardDiskProperties{
			Source:              &vhd.Source,
			SourceType:          vhd.SourceType,
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
			HyperVGeneration:    vhd.HyperVGeneration,
			ProvisioningState:   status.GetProvisioningState(vhd.Status.GetProvisioningStatus()),
			Statuses:            status.GetStatuses(vhd.Status),
			IsPlaceholder:       getVirtualHardDiskIsPlaceholder(vhd),
			CloudInitDataSource: vhd.CloudInitDataSource,
			DiskFileFormat:      common.DiskFileFormat(vhd.DiskFileFormat),
			ContainerName:       &vhd.ContainerName,
		},
	}
}

func getVirtualHardDiskIsPlaceholder(vhd *wssdstorage.VirtualHardDisk) *bool {
	isPlaceholder := false
	entity := vhd.GetEntity()
	if entity != nil {
		isPlaceholder = entity.IsPlaceholder
	}
	return &isPlaceholder
}

func getWssdVirtualHardDisk(containerName string, vhd *storage.VirtualHardDisk) (*wssdstorage.VirtualHardDisk, error) {
	disk := wssdstorage.VirtualHardDisk{
		ContainerName: containerName,
		Tags:          getWssdTags(vhd.Tags),
	}

	if vhd.Name == nil {
		return nil, errors.Wrapf(errors.InvalidInput, "Missing Name")
	}

	disk.Name = *vhd.Name
	disk.Entity = getWssdVirtualHardDiskEntity(vhd)

	if vhd.VirtualHardDiskProperties == nil {
		return &disk, nil
	}

	disk.Virtualharddisktype = getVirtualharddisktype(vhd.Virtualharddisktype)
	disk.HyperVGeneration = vhd.HyperVGeneration
	disk.DiskFileFormat = vhd.DiskFileFormat
	disk.SourceType = vhd.SourceType

	if vhd.Path != nil {
		disk.Path = *vhd.Path
	}

	if disk.Virtualharddisktype == wssdstorage.VirtualHardDiskType_OS_VIRTUALHARDDISK {
		if vhd.Source == nil {
			return nil, errors.Wrapf(errors.InvalidInput, "Missing Source")
		}
		disk.Source = *vhd.Source
		disk.CloudInitDataSource = vhd.CloudInitDataSource

	} else {
		if vhd.DiskSizeBytes == nil && vhd.Source == nil {
			return nil, errors.Wrapf(errors.InvalidInput, "Need to define atleast one of: DiskSize, Source")
		}
		if vhd.SourceType == common.ImageSource_HTTP_SOURCE && vhd.Source == nil {
			return nil, errors.Wrapf(errors.InvalidInput, "Source cannot be empty if source type is HTTP")
		}
		if vhd.Source != nil {
			disk.Source = *vhd.Source
		}
		if vhd.DiskSizeBytes != nil {
			disk.Size = *vhd.DiskSizeBytes
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
		if vhd.VirtualMachineName != nil {
			disk.VirtualmachineName = *vhd.VirtualMachineName
		}
	}

	return &disk, nil
}

func getWssdVirtualHardDiskEntity(vhd *storage.VirtualHardDisk) *wssdcommonproto.Entity {
	isPlaceholder := false
	if vhd.VirtualHardDiskProperties != nil && vhd.VirtualHardDiskProperties.IsPlaceholder != nil {
		isPlaceholder = *vhd.VirtualHardDiskProperties.IsPlaceholder
	}

	return &wssdcommonproto.Entity{
		IsPlaceholder: isPlaceholder,
	}
}

func getVirtualharddisktype(enum string) wssdstorage.VirtualHardDiskType {
	typevalue := wssdstorage.VirtualHardDiskType(0)
	typevTmp, ok := wssdstorage.VirtualHardDiskType_value[enum]
	if ok {
		typevalue = wssdstorage.VirtualHardDiskType(typevTmp)
	}
	return typevalue
}

func getComputeTags(tags *wssdcommonproto.Tags) map[string]*string {
	return prototags.ProtoToMap(tags)
}

func getWssdTags(tags map[string]*string) *wssdcommonproto.Tags {
	return prototags.MapToProto(tags)
}
