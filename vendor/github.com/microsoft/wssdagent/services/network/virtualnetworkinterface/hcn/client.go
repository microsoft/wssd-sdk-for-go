// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package hcn

import (
	"github.com/microsoft/wssdagent/pkg/errors"
	pb "github.com/microsoft/wssdagent/rpc/network"
)

type Service interface {
	Create(*pb.VirtualNetworkInterface) (*pb.VirtualNetworkInterface, error)
	Get(*pb.VirtualNetworkInterface) ([]*pb.VirtualNetworkInterface, error)
	Delete(*pb.VirtualNetworkInterface) error
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

// Create or Update the specified virtual machine(s)
func (c *Client) Create(vnetInterfaces []*pb.VirtualNetworkInterface) ([]*pb.VirtualNetworkInterface, error) {
	newVnetInterfaces := []*pb.VirtualNetworkInterface{}
	for _, vnetInterface := range vnetInterfaces {
		newVnetInterface, err := c.internal.Create(vnetInterface)
		if err != nil {
			c.Delete(newVnetInterfaces)
			return nil, err
		}
		newVnetInterfaces = append(newVnetInterfaces, newVnetInterface)
	}

	return newVnetInterfaces, nil

}

// Get all/selected HCS Virtual Networks
func (c *Client) Get(vnetInterfaces []*pb.VirtualNetworkInterface) ([]*pb.VirtualNetworkInterface, error) {
	newVnetInterfaces := []*pb.VirtualNetworkInterface{}
	if len(vnetInterfaces) == 0 {
		// Get Everything
		return c.internal.Get(nil)
	}

	// Get only requested vnetInterfaces
	for _, vnetInterface := range vnetInterfaces {
		newVnetInterface, err := c.internal.Get(vnetInterface)
		if err != nil {
			return newVnetInterfaces, errors.Wrap(err, "Error finding the interface")
		}
		newVnetInterfaces = append(newVnetInterfaces, newVnetInterface[0])
	}
	return newVnetInterfaces, nil
}

// Delete the specified Virtual Network Interface
func (c *Client) Delete(vnetInterfaces []*pb.VirtualNetworkInterface) error {
	for _, vnetInterface := range vnetInterfaces {
		err := c.internal.Delete(vnetInterface)
		if err != nil {
			return err
		}
	}

	return nil

}
