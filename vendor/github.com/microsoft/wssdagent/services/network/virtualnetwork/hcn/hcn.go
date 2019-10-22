// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package hcn

import (
	errors "errors"
	log "k8s.io/klog"

	"github.com/Microsoft/hcsshim/hcn"

	pb "github.com/microsoft/wssdagent/rpc/network"
	"github.com/microsoft/wssdagent/services/network/virtualnetwork/internal"
)

type Client struct {
}

func NewClient() *Client {
	return &Client{}
}

// Create a Virtual Network
func (c *Client) CreateVirtualNetwork(vnetInternal *internal.VirtualNetworkInternal) (err error) {
	network := vnetInternal.Entity
	// Create network
	networkSchema, err := getHostComputeNetworkConfig(network)
	if err != nil {
		log.Errorf("Unable to create the specified network, error: %v", err)
		return
	}

	hcnNetwork, err := hcn.GetNetworkByName(networkSchema.Name)
	if err == nil {
		// If network is already in the system, use it
		if networkSchema.Type == hcnNetwork.Type {
			vnetInternal.SystemOwned = true
			vnetInternal.Entity, err = getVirtualNetworkConfig(hcnNetwork)
			return
		}
	}

	hcnNetwork, err = networkSchema.Create()
	if err != nil {
		log.Errorf("Unable to create the specified network, error: %v", err)
		return
	}

	newvnet, err := getVirtualNetworkConfig(hcnNetwork)
	if err != nil {
		return
	}
	vnetInternal.Entity = newvnet
	return
}

// Get a Virtual Network specified by name
func (c *Client) Get(networkDef *pb.VirtualNetwork) ([]*pb.VirtualNetwork, error) {
	log.Infof("[VirtualNetwork][Get] spec[%v]", networkDef)

	networks := []*pb.VirtualNetwork{}
	var err error = nil

	if networkDef == nil || len(networkDef.Name) == 0 {
		var hcnNetworks []hcn.HostComputeNetwork
		hcnNetworks, err = hcn.ListNetworks()

		if err != nil {
			log.Errorf("Unable to get network list, error: %v", err)
			return nil, err
		}

		for _, hcnNetwork := range hcnNetworks {
			var network *pb.VirtualNetwork
			network, err = getVirtualNetworkConfig(&hcnNetwork)

			networks = append(networks, network)
		}

	} else {
		var hcnNetwork *hcn.HostComputeNetwork
		hcnNetwork, err = hcn.GetNetworkByName(networkDef.Name)

		if err != nil {
			log.Errorf("Unable to get network with Name %s, error: %v", networkDef.Name, err)
			return nil, err
		}

		var network *pb.VirtualNetwork
		network, err = getVirtualNetworkConfig(hcnNetwork)

		networks = append(networks, network)
	}

	return networks, err
}

// Delete a Virtual VirtualNetwork
func (c *Client) CleanupVirtualNetwork(vnetInternal *internal.VirtualNetworkInternal) error {
	network := vnetInternal.Entity
	hcnNetwork, err := hcn.GetNetworkByName(network.Name)
	if err != nil {
		if hcn.IsNotFoundError(err) {
			return nil
		}
		return err
	}

	// If owned by the system, do not attempt to delete it,
	// since we didnt create it
	if vnetInternal.SystemOwned {
		return nil
	}

	err = hcnNetwork.Delete()
	if err != nil {
		return err
	}
	return nil
}

// HasVirtualNetwork
func (c *Client) HasVirtualNetwork(vnet *pb.VirtualNetwork) bool {
	vnetName := vnet.Name
	hcnNetwork, err := hcn.GetNetworkByName(vnetName)
	if err != nil && hcn.IsNotFoundError(err) {
		return false
	}

	if err == nil {
		// Found Case
		inNetworkTypeString := VirtualNetworkTypeToString(vnet.Type)
		if string(hcnNetwork.Type) == inNetworkTypeString {
			return true
		}
	}

	return false
}

func (c *Client) getDefaultICSNetwork() (*pb.VirtualNetwork, error) {
	hcnNetwork, err := hcn.GetNetworkByID("C08CB7B8-9B3C-408E-8E30-5E16A3AEB444")
	if err != nil {
		log.Errorf("Unable to get default network error: %v", err)
		return nil, err
	}

	return getVirtualNetworkConfig(hcnNetwork)
}

////////////////////////// Private Methods //////////////////////////

// getHostComputeNetworkConfig converts a protobuf VirtualNetwork network to HCN format.
func getHostComputeNetworkConfig(protobufNetwork *pb.VirtualNetwork) (*hcn.HostComputeNetwork, error) {

	hcnNetwork := hcn.HostComputeNetwork{
		SchemaVersion: hcn.SchemaVersion{
			Major: 2,
			Minor: 0,
		},
	}

	hcnNetwork.Id = protobufNetwork.Id
	hcnNetwork.Name = protobufNetwork.Name
	hcnNetwork.Type = hcn.NetworkType(VirtualNetworkTypeToString(protobufNetwork.Type))
	for _, ipam := range protobufNetwork.Ipams {
		hcnIpam := hcn.Ipam{
			Type: ipam.Type,
		}

		for _, subnet := range ipam.Subnets {
			hcnSubnet := hcn.Subnet{
				IpAddressPrefix: subnet.Cidr,
			}

			for _, route := range subnet.Routes {
				hcnSubnet.Routes = append(hcnSubnet.Routes,
					hcn.Route{
						NextHop:           route.Nexthop,
						DestinationPrefix: route.Destinationprefix,
						Metric:            uint16(route.Metric),
					})
			}

			hcnIpam.Subnets = append(hcnIpam.Subnets, hcnSubnet)
		}
		hcnNetwork.Ipams = append(hcnNetwork.Ipams, hcnIpam)
	}

	// Process Network DNS

	if protobufNetwork.Dns == nil {
		return &hcnNetwork, nil
	}

	var dns hcn.Dns
	dns.Domain = protobufNetwork.Dns.Domain

	for _, server := range protobufNetwork.Dns.Servers {
		dns.ServerList = append(dns.ServerList, server)
	}

	for _, search := range protobufNetwork.Dns.Search {
		dns.Search = append(dns.Search, search)
	}

	for _, options := range protobufNetwork.Dns.Options {
		dns.Options = append(dns.Options, options)
	}

	hcnNetwork.Dns = dns

	return &hcnNetwork, nil
}

// getVirtualNetworkConfig converts a HostComputeNetwork to the protobuf VirtualNetwork format.
func getVirtualNetworkConfig(hcnNetwork *hcn.HostComputeNetwork) (*pb.VirtualNetwork, error) {

	var protobufNetwork pb.VirtualNetwork

	networkType, err := VirtualNetworkTypeFromString(string(hcnNetwork.Type))
	if err != nil {
		log.Errorf("Invalid or unsupported network type: '%s' for network '%s'. Error: %v", string(hcnNetwork.Type), hcnNetwork.Name, err)
		return nil, err
	}

	protobufNetwork.Id = hcnNetwork.Id
	protobufNetwork.Name = hcnNetwork.Name
	protobufNetwork.Type = networkType
	for _, ipam := range hcnNetwork.Ipams {
		protobufIpam := pb.Ipam{
			Type: ipam.Type,
		}

		for _, subnet := range ipam.Subnets {
			protobufSubnet := pb.Subnet{
				Cidr: subnet.IpAddressPrefix,
			}

			for _, route := range subnet.Routes {
				protobufSubnet.Routes = append(protobufSubnet.Routes,
					&pb.Route{
						Nexthop:           route.NextHop,
						Destinationprefix: route.DestinationPrefix,
						Metric:            uint32(route.Metric),
					})
			}

			protobufIpam.Subnets = append(protobufIpam.Subnets, &protobufSubnet)
		}
		protobufNetwork.Ipams = append(protobufNetwork.Ipams, &protobufIpam)
	}

	var dns pb.Dns

	dns.Domain = hcnNetwork.Dns.Domain

	for _, server := range hcnNetwork.Dns.ServerList {
		dns.Servers = append(dns.Servers, server)
	}

	for _, search := range hcnNetwork.Dns.Search {
		dns.Search = append(dns.Search, search)
	}

	for _, options := range hcnNetwork.Dns.Options {
		dns.Options = append(dns.Options, options)
	}

	protobufNetwork.Dns = &dns

	return &protobufNetwork, nil
}

func VirtualNetworkTypeToString(vnetType pb.VirtualNetworkType) string {
	typename, ok := pb.VirtualNetworkType_name[int32(vnetType)]

	if !ok {
		return "Unknown"
	}

	return typename
}

func VirtualNetworkTypeFromString(vnNetworkString string) (pb.VirtualNetworkType, error) {
	typevalue, ok := pb.VirtualNetworkType_value[vnNetworkString]

	if !ok {
		return -1, errors.New("Unknown network type")
	}

	return pb.VirtualNetworkType(typevalue), nil
}
