// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package hcs

import (
	"fmt"
	"github.com/microsoft/wssdagent/common"
	"github.com/microsoft/wssdagent/pkg/wssdagent/apis/config"
	"github.com/microsoft/wssdagent/pkg/wssdagent/store"
	pb "github.com/microsoft/wssdagent/rpc/storage"
	"github.com/microsoft/wssdagent/services/storage/virtualharddisk/internal"
	"io"
	log "k8s.io/klog"
	"os"
	"path/filepath"
	"reflect"
)

const (
	defaultReadWrite       = 0666
	defaultStorageLocation = "C:/wssdstorage"
)

type client struct {
	client *Client
	config *config.ChildAgentConfiguration
	store  *store.ConfigStore
}

func newClient() *client {
	cConfig := config.GetChildAgentConfiguration("VirtualHardDisk")
	return &client{
		config: cConfig,
		store:  store.NewConfigStore(cConfig.DataStorePath, reflect.TypeOf(internal.VirtualHardDiskInternal{})),
	}
}

func (c *client) newVirtualHardDisk(id string) *internal.VirtualHardDiskInternal {
	return internal.NewVirtualHardDiskInternal(id, c.config.DataStorePath)
}

func (c *client) Get(virtualHardDisk *pb.VirtualHardDisk) ([]*pb.VirtualHardDisk, error) {
	log.Infof("[VirtualHardDisk][Get] [%v]", virtualHardDisk)
	vhdList := []*pb.VirtualHardDisk{}
	vhdName := ""
	if virtualHardDisk != nil {
		vhdName = virtualHardDisk.Name
	}
	if len(vhdName) == 0 {
		vhdInt, err := c.store.List()
		if err != nil {
			return nil, err
		}

		if *vhdInt == nil || len(*vhdInt) == 0 {
			return nil, nil
		}

		for _, vhd := range *vhdInt {
			vhdInt := vhd.(*internal.VirtualHardDiskInternal)
			vhdList = append(vhdList, vhdInt.Vhd)
		}
	} else {
		vhdInt, err := c.getVirtualHardDiskInternalByName(vhdName)
		if err != nil {
			return nil, err
		}
		vhdList = append(vhdList, vhdInt.Vhd)
	}

	return vhdList, nil
}

func (c *client) Create(virtualHardDisk *pb.VirtualHardDisk) (*pb.VirtualHardDisk, error) {
	log.Infof("[VirtualHardDisk][Create] [%v]", virtualHardDisk)
	if len(virtualHardDisk.Id) == 0 {
		virtualHardDisk.Id = common.NewGuid()
	}

	vhdInternal := c.newVirtualHardDisk(virtualHardDisk.Id)

	if virtualHardDisk.Path == "" {
		virtualHardDisk.Path = generateFilePath(virtualHardDisk.Id)
	}

	sourcePath := virtualHardDisk.Source

	// We check in case the source is actually an Id, and swap the path if it is
	// we don't care about an error in this case
	vhdSource, _ := c.getVirtualHardDiskInternalById(virtualHardDisk.Source)
	if vhdSource != nil {
		sourcePath = vhdSource.Vhd.Path
	}

	if virtualHardDisk != nil {
		copyFile(virtualHardDisk.Path, sourcePath)
	}

	vhdInternal.Vhd = virtualHardDisk

	c.store.Add(virtualHardDisk.Id, vhdInternal)

	return virtualHardDisk, nil
}

func (c *client) Delete(virtualHardDisk *pb.VirtualHardDisk) error {
	vhdList, err := c.Get(virtualHardDisk)
	if err != nil {
		return err
	}

	if len(vhdList) == 0 {
		return fmt.Errorf("Virtual Hard Disk [%s] was not found", virtualHardDisk.Name)
	}

	vhd := vhdList[0]

	log.Infof("[VirtualHardDisk][Delete] [%v]", vhd)
	err = os.Remove(vhd.Path)
	if err != nil {
		return fmt.Errorf("error while trying to remove virtual hard disk [%s] at path [%s] err: %v",
			vhd.Name, vhd.Path, err)
	}
	return c.store.Delete(vhd.Id)
}

func generateFilePath(fileName string) string {
	// Will not recreate or error if already exits
	os.MkdirAll(defaultStorageLocation, os.ModeDir)
	return filepath.Join(defaultStorageLocation, fileName)
}

func copyFile(destinationPath string, sourcePath string) error {
	sourceFile, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("failed to open source file \"%s\", err: %v", sourcePath, err)
	}
	defer sourceFile.Close()

	destinationFile, err := os.OpenFile(destinationPath, os.O_CREATE|os.O_RDWR, defaultReadWrite)
	if err != nil {
		return fmt.Errorf("failed to open destination file \"%s\", err: %v", destinationPath, err)
	}
	defer destinationFile.Close()

	_, err = io.Copy(destinationFile, sourceFile)
	if err != nil {
		return fmt.Errorf("failed to copy from file \"%s\" to file \"%s\", err: %v", sourcePath, destinationPath, err)
	}

	return nil
}

func (c *client) getVirtualHardDiskInternalByName(name string) (*internal.VirtualHardDiskInternal, error) {
	vhdInt, err := c.store.List()
	if err != nil {
		return nil, err
	}

	if *vhdInt == nil || len(*vhdInt) == 0 {
		return nil, nil
	}

	for _, vhd := range *vhdInt {
		vhdInt := vhd.(*internal.VirtualHardDiskInternal)
		if vhdInt.Vhd.Name == name {
			return vhdInt, nil
		}
	}
	return nil, fmt.Errorf("Virtual Hard Disk [%s] not found", name)
}

func (c *client) getVirtualHardDiskInternalById(id string) (*internal.VirtualHardDiskInternal, error) {
	vhdInt, err := c.store.List()
	if err != nil {
		return nil, err
	}

	if *vhdInt == nil || len(*vhdInt) == 0 {
		return nil, nil
	}

	for _, vhd := range *vhdInt {
		vhdInt := vhd.(*internal.VirtualHardDiskInternal)
		if vhdInt.Vhd.Id == id {
			return vhdInt, nil
		}
	}
	return nil, fmt.Errorf("Virtual Hard Disk with ID: [%s] not found", id)
}
