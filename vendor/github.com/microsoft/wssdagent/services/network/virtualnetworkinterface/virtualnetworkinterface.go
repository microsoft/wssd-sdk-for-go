// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package virtualnetworkinterface

import (
	"context"
	"fmt"
	"time"

	"github.com/microsoft/wssdagent/pkg/errors"
	"github.com/microsoft/wssdagent/pkg/guid"
	pb "github.com/microsoft/wssdagent/rpc/network"
	log "k8s.io/klog"
)

type VirtualNetworkInterfaceProvider struct {
	client *Client
}

func NewVirtualNetworkInterfaceProvider() *VirtualNetworkInterfaceProvider {
	return &VirtualNetworkInterfaceProvider{
		client: NewClient(),
	}
}

func (vnicProv *VirtualNetworkInterfaceProvider) Get(ctx context.Context, vnics []*pb.VirtualNetworkInterface) ([]*pb.VirtualNetworkInterface, error) {
	newvnics := []*pb.VirtualNetworkInterface{}
	if len(vnics) == 0 {
		// Get Everything
		return vnicProv.client.Get(ctx, nil)
	}

	// Get only requested vnics
	for _, vnic := range vnics {
		newvnic, err := vnicProv.client.Get(ctx, vnic)
		if err != nil {
			return newvnics, err
		}
		newvnics = append(newvnics, newvnic[0])
	}
	return newvnics, nil
}

func (vnicProv *VirtualNetworkInterfaceProvider) CreateOrUpdate(ctx context.Context, vnics []*pb.VirtualNetworkInterface) ([]*pb.VirtualNetworkInterface, error) {
	newvnics := []*pb.VirtualNetworkInterface{}
	for _, vnic := range vnics {
		newvnic, err := vnicProv.client.Create(ctx, vnic)
		if err != nil {
			if err != errors.AlreadyExists {
				vnicProv.client.Delete(ctx, vnic)
			}
			return newvnics, err
		}
		newvnics = append(newvnics, newvnic)
	}

	return newvnics, nil
}

func (vnicProv *VirtualNetworkInterfaceProvider) Delete(ctx context.Context, vnics []*pb.VirtualNetworkInterface) error {
	for _, vnic := range vnics {
		err := vnicProv.client.Delete(ctx, vnic)
		if err != nil {
			return err
		}
	}

	return nil
}

// CreateVirtualNetworkInterface
func (vnicProv *VirtualNetworkInterfaceProvider) CreateVirtualNetworkInterface(ctx context.Context, name, vnetName string) error {
	vnic := &pb.VirtualNetworkInterface{Name: name, Id: guid.NewGuid(), Networkname: vnetName}
	_, err := vnicProv.CreateOrUpdate(ctx, []*pb.VirtualNetworkInterface{vnic})
	if err != nil {
		return err
	}
	return nil
}

// DeleteVirtualNetworkInterface helper to delete network interfaces
func (vnicProv *VirtualNetworkInterfaceProvider) DeleteVirtualNetworkInterface(ctx context.Context, vnics []string) error {
	return vnicProv.Delete(ctx, getVirtualNetworkInterfaceByName(vnics))
}

// GetVirtualNetworkInterfaceByName
func (vnicProv *VirtualNetworkInterfaceProvider) GetVirtualNetworkInterfaceByName(ctx context.Context, name string) (*pb.VirtualNetworkInterface, error) {
	if len(name) == 0 {
		return nil, errors.Wrapf(errors.InvalidInput, "GetVirtualNetworkInterfaceByName cannot query empty name")
	}
	vnicsnew, err := vnicProv.Get(ctx, getVirtualNetworkInterfaceByName([]string{name}))
	if err != nil {
		return nil, err
	}

	if len(vnicsnew) == 0 {
		return nil, errors.NotFound
	}

	return vnicsnew[0], nil
}
func (vnicProv *VirtualNetworkInterfaceProvider) GetVirtualNetworkInterfaceById(ctx context.Context, Id string) (*pb.VirtualNetworkInterface, error) {
	vnicsnew, err := vnicProv.Get(ctx, getVirtualNetworkInterfaceById([]string{Id}))
	if err != nil {
		return nil, err
	}

	if len(vnicsnew) == 0 {
		return nil, fmt.Errorf(Id + " not found")
	}

	return vnicsnew[0], nil
}

func (vnicProv *VirtualNetworkInterfaceProvider) WaitForIPAddress(ctx context.Context, name string) (string, error) {
	log.Infof("[VirtualNetworkInterface][WaitForIPAddress] vnic[%s]", name)
	for i := 0; i < 100; i++ {
		vmnic, err := vnicProv.GetVirtualNetworkInterfaceByName(ctx, name)
		log.Infof("[VirtualNetworkInterface][WaitForIPAddress] vnic[%v]", vmnic)
		if err != nil {
			time.Sleep(1 * time.Second)
			continue
		}

		if len(vmnic.Ipconfigs) == 0 || len(vmnic.Ipconfigs[0].GetIpaddress()) == 0 {
			time.Sleep(5 * time.Second)
			continue
		}
		return vmnic.Ipconfigs[0].GetIpaddress(), nil
	}

	return "", fmt.Errorf("Unable to get IPAddress")
}

func getVirtualNetworkInterfaceById(vnics []string) []*pb.VirtualNetworkInterface {
	tmp := []*pb.VirtualNetworkInterface{}
	for _, vnic := range vnics {
		tmp = append(tmp, &pb.VirtualNetworkInterface{Id: vnic})
	}
	return tmp
}
func getVirtualNetworkInterfaceByName(vnics []string) []*pb.VirtualNetworkInterface {
	tmp := []*pb.VirtualNetworkInterface{}
	for _, vnic := range vnics {
		tmp = append(tmp, &pb.VirtualNetworkInterface{Name: vnic})
	}
	return tmp
}
