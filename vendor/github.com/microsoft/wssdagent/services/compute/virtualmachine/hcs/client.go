// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package hcs

import (
	pb "github.com/microsoft/wssdagent/rpc/compute"
)

type Service interface {
	Create(*pb.VirtualMachine) (*pb.VirtualMachine, error)
	Get(*pb.VirtualMachine) ([]*pb.VirtualMachine, error)
	Delete(*pb.VirtualMachine) error
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
func (c *Client) Create(vms []*pb.VirtualMachine) ([]*pb.VirtualMachine, error) {
	newvms := []*pb.VirtualMachine{}
	for _, vm := range vms {
		newvm, err := c.internal.Create(vm)
		if err != nil {
			return newvms, err
		}
		newvms = append(newvms, newvm)
	}

	return newvms, nil

}

// Get all/selected HCS Virtual Machines
func (c *Client) Get(vms []*pb.VirtualMachine) ([]*pb.VirtualMachine, error) {
	var err error
	newvms := []*pb.VirtualMachine{}
	if len(vms) == 0 {
		// Get Everything
		return c.internal.Get(nil)
	}

	// Get only requested vms
	for _, vm := range vms {
		newvm, err := c.internal.Get(vm)
		if err != nil {
			return newvms, err
		}
		newvms = append(newvms, newvm[0])
	}
	return newvms, err
}

// Delete the specified virtual machine(s)
func (c *Client) Delete(vms []*pb.VirtualMachine) error {
	for _, vm := range vms {
		err := c.internal.Delete(vm)
		if err != nil {
			return err
		}
	}

	return nil

}
