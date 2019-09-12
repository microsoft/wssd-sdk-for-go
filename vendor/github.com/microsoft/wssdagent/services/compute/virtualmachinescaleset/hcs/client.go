// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package hcs

import (
	pb "github.com/microsoft/wssdagent/rpc/compute"
)

type Service interface {
	Create(*pb.VirtualMachineScaleSet) (*pb.VirtualMachineScaleSet, error)
	Get(*pb.VirtualMachineScaleSet) ([]*pb.VirtualMachineScaleSet, error)
	Delete(*pb.VirtualMachineScaleSet) error
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
func (c *Client) Create(vms []*pb.VirtualMachineScaleSet) ([]*pb.VirtualMachineScaleSet, error) {
	newvms := []*pb.VirtualMachineScaleSet{}
	for _, vm := range vms {
		newvmss, err := c.internal.Create(vm)
		if err != nil {
			return newvms, err
		}
		newvms = append(newvms, newvmss)

	}

	return newvms, nil

}

// Get all HCS Virtual Machines
func (c *Client) Get(vms []*pb.VirtualMachineScaleSet) ([]*pb.VirtualMachineScaleSet, error) {
	var err error
	newvms := []*pb.VirtualMachineScaleSet{}
	if len(vms) == 0 {
		// Get everything
		return c.internal.Get(nil)
	}
	// Get only requested vmss
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
func (c *Client) Delete(vms []*pb.VirtualMachineScaleSet) error {
	for _, vm := range vms {
		err := c.internal.Delete(vm)
		if err != nil {
			return err
		}
	}

	return nil

}
