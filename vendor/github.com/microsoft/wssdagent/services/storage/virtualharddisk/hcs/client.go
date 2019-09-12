// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package hcs

import (
	pb "github.com/microsoft/wssdagent/rpc/storage"
)

type Service interface {
	Create(*pb.VirtualHardDisk) (*pb.VirtualHardDisk, error)
	Get(*pb.VirtualHardDisk) ([]*pb.VirtualHardDisk, error)
	Delete(*pb.VirtualHardDisk) error
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

func (c *Client) Create(vhdList []*pb.VirtualHardDisk) ([]*pb.VirtualHardDisk, error) {
	newVhdList := []*pb.VirtualHardDisk{}
	for _, virtualHardDisk := range vhdList {
		resultVhd, err := c.internal.Create(virtualHardDisk)
		if err != nil {
			return nil, err
		}
		newVhdList = append(newVhdList, resultVhd)
	}
	return newVhdList, nil
}

func (c *Client) Get(vhdList []*pb.VirtualHardDisk) ([]*pb.VirtualHardDisk, error) {
	newVhdList := []*pb.VirtualHardDisk{}
	if len(vhdList) == 0 {
		var err error
		newVhdList, err = c.internal.Get(nil)
		if err != nil {
			return nil, err
		}
	} else {
		for _, virtualHardDisk := range vhdList {
			resultVhdList, err := c.internal.Get(virtualHardDisk)
			if err != nil {
				return nil, err
			}
			newVhdList = append(newVhdList, resultVhdList[0])
		}
	}
	return newVhdList, nil
}

func (c *Client) Delete(vhdList []*pb.VirtualHardDisk) error {
	for _, virtualHardDisk := range vhdList {
		c.internal.Delete(virtualHardDisk)
	}
	return nil
}
