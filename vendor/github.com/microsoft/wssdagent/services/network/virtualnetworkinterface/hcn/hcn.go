// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package hcn

import (
	"github.com/Microsoft/hcsshim/hcn"
	log "k8s.io/klog"
	"strconv"

	pb "github.com/microsoft/wssdagent/rpc/network"
	"github.com/microsoft/wssdagent/services/network/virtualnetworkinterface/internal"
)

type Client struct {
}

func NewClient() *Client {
	return &Client{}
}

// Create a Virtual Network Interface
func (c *Client) CreateVirtualNetworkInterface(vnicInternal *internal.VirtualNetworkInterfaceInternal) (err error) {
	vnetInterfaceDef := vnicInternal.Entity
	// Create a network interface
	hcnEndpointSchema, err := getHostComputeEndpointConfig(vnetInterfaceDef)
	if err != nil {
		log.Errorf("[NetworkInterface][Create] Unable to get the endpoint config for the specified network interface, error: %v", err)
		return
	}

	hcnEndpoint, err := hcnEndpointSchema.Create()
	if err != nil {
		log.Errorf("Unable to create the specified network interface, error: %v", err)
		return
	}

	newvnic, err := getVirtualNetworkInterfaceConfig(hcnEndpoint)
	if err != nil {
		return
	}
	vnicInternal.Entity = newvnic
	return
}

// Delete a Virtual Network Interface
func (c *Client) CleanupVirtualNetworkInterface(vnicInternal *internal.VirtualNetworkInterfaceInternal) (err error) {
	vnetInterfaceDef := vnicInternal.Entity
	log.Infof("[NetworkInterface][Delete] spec[%v]", vnetInterfaceDef)

	hcnEndpoint, err := hcn.GetEndpointByName(vnetInterfaceDef.Name)
	if err != nil {
		if hcn.IsNotFoundError(err) {
			err = nil
			return
		}

		log.Errorf("[NetworkInterface][Delete] Unable to get network interface with Id %s, error: %v", vnetInterfaceDef.Id, err)
		return
	}

	err = hcnEndpoint.Delete()
	if err != nil {
		return
	}
	return
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
