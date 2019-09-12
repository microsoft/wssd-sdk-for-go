// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package hcn

import (
	"github.com/Microsoft/hcsshim/hcn"
	log "k8s.io/klog"
	"reflect"
	"strconv"

	pb "github.com/microsoft/wssdagent/rpc/network"

	"github.com/microsoft/wssdagent/pkg/wssdagent/apis/config"
	"github.com/microsoft/wssdagent/pkg/wssdagent/store"
	"github.com/microsoft/wssdagent/services/network/virtualnetworkinterface/internal"
)

type client struct {
	config *config.ChildAgentConfiguration
	store  *store.ConfigStore
}

func newClient() *client {
	cConfig := config.GetChildAgentConfiguration("VirtualNetworkInterface")
	return &client{
		store:  store.NewConfigStore(cConfig.DataStorePath, reflect.TypeOf(internal.VirtualNetworkInterfaceInternal{})),
		config: cConfig,
	}
}

func (c *client) newVirtualNetworkInterface(id string) *internal.VirtualNetworkInterfaceInternal {
	return internal.NewVirtualNetworkInterfaceInternal(id, c.config.DataStorePath)
}

// Create a Virtual Network Interface
func (c *client) Create(vnetInterfaceDef *pb.VirtualNetworkInterface) (*pb.VirtualNetworkInterface, error) {
	log.Infof("[NetworkInterface][Create] spec[%v]", vnetInterfaceDef)
	vnicinternal := c.newVirtualNetworkInterface(vnetInterfaceDef.Id)

	// Create a network interface
	hcnEndpointSchema, err := getHostComputeEndpointConfig(vnetInterfaceDef)
	if err != nil {
		log.Errorf("[NetworkInterface][Create] Unable to get the endpoint config for the specified network interface, error: %v", err)
		return nil, err
	}

	hcnEndpoint, err := hcnEndpointSchema.Create()
	if err != nil {
		log.Errorf("Unable to create the specified network interface, error: %v", err)
		return nil, err
	}

	newvnic, err := getVirtualNetworkInterfaceConfig(hcnEndpoint)
	vnicinternal.VNic = newvnic

	// 3. Save the config to the store
	c.store.Add(vnetInterfaceDef.Id, vnicinternal)

	return newvnic, err
}

// Get a Virtual Network Interface specified by Id
func (c *client) Get(vnetInterfaceDef *pb.VirtualNetworkInterface) ([]*pb.VirtualNetworkInterface, error) {
	log.Infof("[NetworkInterface][Get] spec[%v]", vnetInterfaceDef)

	networkInterfaces := []*pb.VirtualNetworkInterface{}
	var err error = nil

	if vnetInterfaceDef == nil || len(vnetInterfaceDef.Name) == 0 {
		var hcnEndpoints []hcn.HostComputeEndpoint
		hcnEndpoints, err = hcn.ListEndpoints()

		if err != nil {
			log.Errorf("Unable to get network interface list, error: %v", err)
			return nil, err
		}

		for _, hcnEndpoint := range hcnEndpoints {
			var networkInterface *pb.VirtualNetworkInterface
			networkInterface, err = getVirtualNetworkInterfaceConfig(&hcnEndpoint)

			networkInterfaces = append(networkInterfaces, networkInterface)
		}

	} else {
		var hcnEndpoint *hcn.HostComputeEndpoint
		hcnEndpoint, err := hcn.GetEndpointByName(vnetInterfaceDef.Name)

		if err != nil {
			log.Errorf("[NetworkInterface][Get] Unable to get network interface with Id %s, error: %v", vnetInterfaceDef.Id, err)
			return nil, err
		}

		var networkInterface *pb.VirtualNetworkInterface
		networkInterface, err = getVirtualNetworkInterfaceConfig(hcnEndpoint)

		networkInterfaces = append(networkInterfaces, networkInterface)
	}

	return networkInterfaces, err
}

// Delete a Virtual Network Interface
func (c *client) Delete(vnetInterfaceDef *pb.VirtualNetworkInterface) error {
	log.Infof("[NetworkInterface][Delete] spec[%v]", vnetInterfaceDef)

	hcnEndpoint, err := hcn.GetEndpointByName(vnetInterfaceDef.Name)
	if err != nil {
		log.Errorf("[NetworkInterface][Delete] Unable to get network interface with Id %s, error: %v", vnetInterfaceDef.Id, err)
		return err
	}

	if hcn.IsNotFoundError(err) {
		return nil
	}

	err = hcnEndpoint.Delete()
	if err != nil && !hcn.IsNotFoundError(err) {
		log.Errorf("[NetworkInterface][Delete] Unable to delete network interface with Id %s, error: %v", vnetInterfaceDef.Id, err)
		return err
	}
	return c.store.Delete(vnetInterfaceDef.Id)
}

// getHostComputeEndpointConfig converts a protobuf VirtualNetworkInterface network to HostComputeEndpoint.
func getHostComputeEndpointConfig(protobufNetworkInterface *pb.VirtualNetworkInterface) (*hcn.HostComputeEndpoint, error) {

	hcnEndpoint := hcn.HostComputeEndpoint{
		SchemaVersion: hcn.SchemaVersion{
			Major: 2,
			Minor: 0,
		},
	}

	net, err := hcn.GetNetworkByName(protobufNetworkInterface.Networkname)

	if err != nil {
		return nil, err
	}

	hcnEndpoint.Id = protobufNetworkInterface.Id
	hcnEndpoint.Name = protobufNetworkInterface.Name
	if len(protobufNetworkInterface.Macaddress) != 0 {
		hcnEndpoint.MacAddress = protobufNetworkInterface.Macaddress
	}
	hcnEndpoint.HostComputeNetwork = net.Id
	// TODO: do something with protobufNetworkInterface.Type
	for _, ipconfig := range protobufNetworkInterface.Ipconfigs {
		prefixLength, err := strconv.ParseUint(ipconfig.Prefixlength, 10, 32)
		if err != nil {
			log.Errorf("Invalid prefix length %s (cannot be converted to a uint), error: %v", ipconfig.Prefixlength, err)
			return nil, err
		}

		hcnIpConfig := hcn.IpConfig{
			IpAddress:    ipconfig.Ipaddress,
			PrefixLength: uint8(prefixLength),
			// TODO: subnetid
		}

		hcnEndpoint.IpConfigurations = append(hcnEndpoint.IpConfigurations, hcnIpConfig)
	}

	return &hcnEndpoint, nil
}

// getVirtualNetworkInterfaceConfig converts a HostComputeEndpoint to the protobuf VirtualNetworkInterface format.
func getVirtualNetworkInterfaceConfig(hcnEndpoint *hcn.HostComputeEndpoint) (*pb.VirtualNetworkInterface, error) {

	var protobufNetworkInterface pb.VirtualNetworkInterface

	net, err := hcn.GetNetworkByID(hcnEndpoint.HostComputeNetwork)

	if err != nil {
		return nil, err
	}

	protobufNetworkInterface.Id = hcnEndpoint.Id
	protobufNetworkInterface.Name = hcnEndpoint.Name
	protobufNetworkInterface.Networkname = net.Name
	protobufNetworkInterface.Macaddress = hcnEndpoint.MacAddress
	// TODO: do something with protobufNetworkInterface.Type

	for _, ipconfig := range hcnEndpoint.IpConfigurations {

		prefixLength := uint32(ipconfig.PrefixLength)
		protobufIpConfig := pb.IpConfiguration{
			Ipaddress:    ipconfig.IpAddress,
			Prefixlength: strconv.FormatUint(uint64(prefixLength), 10),
			// TODO: subnetid
		}

		protobufNetworkInterface.Ipconfigs = append(protobufNetworkInterface.Ipconfigs, &protobufIpConfig)
	}

	return &protobufNetworkInterface, nil
}
