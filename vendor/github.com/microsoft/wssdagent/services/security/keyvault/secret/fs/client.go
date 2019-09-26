// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package fs

import (
	pb "github.com/microsoft/wssdagent/rpc/security"
)

type Service interface {
	Create(*pb.Secret) (*pb.Secret, error)
	Get(*pb.Secret) ([]*pb.Secret, error)
	Delete(*pb.Secret) error
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

func (c *Client) Create(secrets []*pb.Secret) ([]*pb.Secret, error) {
	newSecrets := []*pb.Secret{}
	for _, secret := range secrets {
		resultSecret, err := c.internal.Create(secret)
		if err != nil {
			return nil, err
		}
		newSecrets = append(newSecrets, resultSecret)
	}
	return newSecrets, nil
}

func (c *Client) Get(secrets []*pb.Secret) ([]*pb.Secret, error) {
	newSecrets := []*pb.Secret{}
	if len(secrets) == 0 {
		var err error
		newSecrets, err = c.internal.Get(nil)
		if err != nil {
			return nil, err
		}
	} else {
		for _, secret := range secrets {
			resultSecretList, err := c.internal.Get(secret)
			if err != nil {
				return nil, err
			}
			newSecrets = append(newSecrets, resultSecretList[0])
		}
	}
	return newSecrets, nil
}

func (c *Client) Delete(secrets []*pb.Secret) error {
	for _, secret := range secrets {
		c.internal.Delete(secret)
	}
	return nil
}
