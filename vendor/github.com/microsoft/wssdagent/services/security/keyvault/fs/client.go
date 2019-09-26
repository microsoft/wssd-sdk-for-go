// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package fs

import (
	pb "github.com/microsoft/wssdagent/rpc/security"
)

type Service interface {
	Create(*pb.KeyVault) (*pb.KeyVault, error)
	Get(*pb.KeyVault) ([]*pb.KeyVault, error)
	Delete(*pb.KeyVault) error
}

type Client struct {
	internal Service
}

func NewClient() *Client {
	c := newClient()
	return &Client{
		internal: c,
	}
}

func (c *Client) Create(keyvaults []*pb.KeyVault) ([]*pb.KeyVault, error) {
	newKeyVaults := []*pb.KeyVault{}
	for _, keyvault := range keyvaults {
		resultVault, err := c.internal.Create(keyvault)
		if err != nil {
			return nil, err
		}
		newKeyVaults = append(newKeyVaults, resultVault)
	}
	return newKeyVaults, nil
}

func (c *Client) Get(keyvaults []*pb.KeyVault) ([]*pb.KeyVault, error) {
	newKeyVaults := []*pb.KeyVault{}
	if len(keyvaults) == 0 {
		var err error
		newKeyVaults, err = c.internal.Get(nil)
		if err != nil {
			return nil, err
		}
	} else {
		for _, keyvault := range keyvaults {
			resultVaultList, err := c.internal.Get(keyvault)
			if err != nil {
				return nil, err
			}
			newKeyVaults = append(newKeyVaults, resultVaultList[0])
		}
	}
	return newKeyVaults, nil
}

func (c *Client) Delete(keyvaults []*pb.KeyVault) error {
	for _, keyvault := range keyvaults {
		c.internal.Delete(keyvault)
	}
	return nil
}
