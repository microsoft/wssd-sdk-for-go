// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package hcn

import (
	errors "errors"
	log "k8s.io/klog"
	"reflect"

	"github.com/Microsoft/hcsshim/hcn"

	"github.com/microsoft/wssdagent/pkg/wssdagent/apis/config"
	"github.com/microsoft/wssdagent/pkg/wssdagent/store"
	pb "github.com/microsoft/wssdagent/rpc/network"
	"github.com/microsoft/wssdagent/services/network/virtualnetwork/internal"
)

type client struct {
	config *config.ChildAgentConfiguration
	store  *store.ConfigStore
}

func newClient() *client {
	cConfig := config.GetChildAgentConfiguration("VirtualNetwork")
	return &client{
		store:  store.NewConfigStore(cConfig.DataStorePath, reflect.TypeOf(internal.VirtualNetworkInternal{})),
		config: cConfig,
	}
}

func (c *client) newVirtualNetwork(id string) *internal.VirtualNetworkInternal {
	return internal.NewVirtualNetworkInternal(id, c.config.DataStorePath)
}

// Create a Virtual Network
func (c *client) Create(network *pb.VirtualNetwork) (*pb.VirtualNetwork, error) {
	log.Infof("[VirtualNetwork][Create] spec[%v]", network)
	vnetinternal := c.newVirtualNetwork(network.Id)

	// Create network
	networkSchema, err := getHostComputeNetworkConfig(network)
	if err != nil {
		log.Errorf("Unable to create the specified network, error: %v", err)
		return nil, err
	}

	hcnNetwork, err := networkSchema.Create()
	if err != nil {
		log.Errorf("Unable to create the specified network, error: %v", err)
		return nil, err
	}

	// 3. Save the config to the store
	c.store.Add(network.Id, vnetinternal)

	newvnet, err := getVirtualNetworkConfig(hcnNetwork)
	vnetinternal.VNet = newvnet
	return newvnet, err
}

// Get a Virtual Network specified by name
func (c *client) Get(networkDef *pb.VirtualNetwork) ([]*pb.VirtualNetwork, error) {
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
func (c *client) Delete(network *pb.VirtualNetwork) error {
	log.Infof("[VirtualNetwork][Delete] spec[%v]", network)

	hcnNetwork, err := hcn.GetNetworkByName(network.Name)
	if err != nil {
		log.Errorf("Unable to get network with Id %s, error: %v", network.Id, err)
		return err
	}

	if hcn.IsNotFoundError(err) {
		return nil
	}

	err = hcnNetwork.Delete()
	if err != nil && !hcn.IsNotFoundError(err) {
		log.Errorf("Unable to delete network with Id %s, error: %v", network.Id, err)
		return err
	}

	return c.store.Delete(network.Id)
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
