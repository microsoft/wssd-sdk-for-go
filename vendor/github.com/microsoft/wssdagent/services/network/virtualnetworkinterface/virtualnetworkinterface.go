// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package virtualnetworkinterface

import (
	"fmt"
	"time"

	"github.com/microsoft/wssdagent/common"
	pb "github.com/microsoft/wssdagent/rpc/network"
	log "k8s.io/klog"
)

type VirtualNetworkInterfaceProvider interface {
	CreateOrUpdate([]*pb.VirtualNetworkInterface) ([]*pb.VirtualNetworkInterface, error)
	Get([]*pb.VirtualNetworkInterface) ([]*pb.VirtualNetworkInterface, error)
	Delete([]*pb.VirtualNetworkInterface) error
}

// CreateVirtualNetworkInterface
func CreateVirtualNetworkInterface(provider VirtualNetworkInterfaceProvider, name, vnetName string) error {
	vnic := &pb.VirtualNetworkInterface{Name: name, Id: common.NewGuid(), Networkname: vnetName}
	_, err := provider.CreateOrUpdate([]*pb.VirtualNetworkInterface{vnic})
	if err != nil {
		return err
	}
	return nil
}

// DeleteVirtualNetworkInterface helper to delete network interfaces
func DeleteVirtualNetworkInterface(provider VirtualNetworkInterfaceProvider, vnics []string) error {
	return provider.Delete(getVirtualNetworkInterfaceByName(vnics))
}

// GetVirtualNetworkInterfaceByName
func GetVirtualNetworkInterfaceByName(provider VirtualNetworkInterfaceProvider, name string) (*pb.VirtualNetworkInterface, error) {
	if len(name) == 0 {
		return nil, fmt.Errorf("GetVirtualNetworkInterfaceByName cannot query empty name")
	}
	vnicsnew, err := provider.Get(getVirtualNetworkInterfaceByName([]string{name}))
	if err != nil {
		return nil, err
	}

	if len(vnicsnew) == 0 {
		return nil, fmt.Errorf(name + " not found")
	}

	return vnicsnew[0], nil
}
func GetVirtualNetworkInterfaceById(provider VirtualNetworkInterfaceProvider, Id string) (*pb.VirtualNetworkInterface, error) {
	vnicsnew, err := provider.Get(getVirtualNetworkInterfaceById([]string{Id}))
	if err != nil {
		return nil, err
	}

	if len(vnicsnew) == 0 {
		return nil, fmt.Errorf(Id + " not found")
	}

	return vnicsnew[0], nil
}

func WaitForIPAddress(provider VirtualNetworkInterfaceProvider, name string) (string, error) {
	log.Infof("[VirtualNetworkInterface][WaitForIPAddress] vnic[%s]", name)
	for i := 0; i < 100; i++ {
		vmnic, err := GetVirtualNetworkInterfaceByName(provider, name)
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
