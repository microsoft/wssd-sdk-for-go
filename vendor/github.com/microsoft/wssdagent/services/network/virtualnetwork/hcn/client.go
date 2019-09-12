// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package hcn

import (
	pb "github.com/microsoft/wssdagent/rpc/network"
)

type Service interface {
	Create(*pb.VirtualNetwork) (*pb.VirtualNetwork, error)
	Get(*pb.VirtualNetwork) ([]*pb.VirtualNetwork, error)
	Delete(*pb.VirtualNetwork) error
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

// Create or Update the specified virtual network(s)
func (c *Client) Create(vnets []*pb.VirtualNetwork) ([]*pb.VirtualNetwork, error) {
	newvnets := []*pb.VirtualNetwork{}
	for _, vnet := range vnets {
		newvnet, err := c.internal.Create(vnet)
		if err != nil {
			c.Delete(newvnets)
			return nil, err
		}
		newvnets = append(newvnets, newvnet)
	}

	return newvnets, nil

}

// Get all/selected HCS virtual network(s)
func (c *Client) Get(vnets []*pb.VirtualNetwork) ([]*pb.VirtualNetwork, error) {
	newvnets := []*pb.VirtualNetwork{}
	if len(vnets) == 0 {
		// Get Everything
		return c.internal.Get(nil)
	}

	// Get only requested vnets
	for _, vnet := range vnets {
		newvnet, err := c.internal.Get(vnet)
		if err != nil {
			return newvnets, err
		}
		newvnets = append(newvnets, newvnet[0])
	}
	return newvnets, nil
}

// Delete the specified virtual network(s)
func (c *Client) Delete(vnets []*pb.VirtualNetwork) error {
	for _, vnet := range vnets {
		err := c.internal.Delete(vnet)
		if err != nil {
			return err
		}
	}

	return nil

}
