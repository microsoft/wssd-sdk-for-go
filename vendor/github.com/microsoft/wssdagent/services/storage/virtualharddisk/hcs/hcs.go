// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package hcs

import (
	"fmt"
	"github.com/microsoft/wssdagent/pkg/errors"
	"github.com/microsoft/wssdagent/services/storage/virtualharddisk/internal"
	"io"
	"os"
)

const (
	defaultReadWrite = 0666
)

type Client struct {
}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) CreateVirtualHardDisk(vhdInternal *internal.VirtualHardDiskInternal) (err error) {
	virtualHardDisk := vhdInternal.Entity

	err = copyFile(virtualHardDisk.Path, virtualHardDisk.Source)
	if err != nil {
		return
	}

	vhdInternal.Entity = virtualHardDisk
	return
}

func (c *Client) CleanupVirtualHardDisk(vhdInternal *internal.VirtualHardDiskInternal) (err error) {
	vhd := vhdInternal.Entity
	if _, err1 := os.Stat(vhd.Path); os.IsNotExist(err1) {
		err = nil
		return
	}

	err = os.Remove(vhd.Path)
	if err != nil {
		err = errors.Wrapf(err, fmt.Sprintf("error while trying to remove virtual hard disk [%s] at path [%s]", vhd.Name, vhd.Path))
		return
	}
	return
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
