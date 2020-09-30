// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license

package internal

import (
	"context"
	"fmt"

	"github.com/microsoft/moc/pkg/status"
	"github.com/microsoft/wssd-sdk-for-go/services/network"

	"github.com/microsoft/moc/pkg/auth"
	wssdcommonproto "github.com/microsoft/moc/rpc/common"
	wssdnetwork "github.com/microsoft/moc/rpc/nodeagent/network"
	wssdclient "github.com/microsoft/wssd-sdk-for-go/pkg/client"
)

type client struct {
	wssdnetwork.VirtualNetworkAgentClient
}

func virtualNetworkTypeToString(vnetType wssdnetwork.VirtualNetworkType) string {
	typename, ok := wssdnetwork.VirtualNetworkType_name[int32(vnetType)]
	if !ok {
		return "Unknown"
	}
	return typename

}

func virtualNetworkTypeFromString(vnNetworkString string) (wssdnetwork.VirtualNetworkType, error) {
	typevalue := wssdnetwork.VirtualNetworkType_ICS
	if len(vnNetworkString) > 0 {
		typevTmp, ok := wssdnetwork.VirtualNetworkType_value[vnNetworkString]
		if ok {
			typevalue = wssdnetwork.VirtualNetworkType(typevTmp)
		}
	}
	return typevalue, nil
}

// NewVirtualNetworkClient - creates a client session with the backend wssd agent
func NewVirtualNetworkClient(subID string, authorizer auth.Authorizer) (*client, error) {

	c, err := wssdclient.GetVirtualNetworkClient(&subID, authorizer)
	if err != nil {
		return nil, err
	}
	return &client{c}, nil
}

// Get
func (c *client) Get(ctx context.Context, group, name string) (*[]network.VirtualNetwork, error) {
	request := getVirtualNetworkRequest(wssdcommonproto.Operation_GET, name, nil)
	response, err := c.VirtualNetworkAgentClient.Invoke(ctx, request)
	if err != nil {
		return nil, err
	}
	return getVirtualNetworksFromResponse(response), nil
}

// CreateOrUpdate
func (c *client) CreateOrUpdate(ctx context.Context, group, name string, vnet *network.VirtualNetwork) (*network.VirtualNetwork, error) {
	err := c.validate(ctx, group, name, vnet)
	if err != nil {
		return nil, err
	}
	request := getVirtualNetworkRequest(wssdcommonproto.Operation_POST, name, vnet)
	response, err := c.VirtualNetworkAgentClient.Invoke(ctx, request)
	if err != nil {
		return nil, err
	}
	vnets := getVirtualNetworksFromResponse(response)

	if len(*vnets) == 0 {
		return nil, fmt.Errorf("[VirtualNetwork][Create] Unexpected error: Creating a network interface returned no result")
	}

	return &((*vnets)[0]), nil
}

// Delete methods invokes create or update on the client
func (c *client) Delete(ctx context.Context, group, name string) error {
	vnet, err := c.Get(ctx, group, name)
	if err != nil {
		return err
	}
	if len(*vnet) == 0 {
		return fmt.Errorf("Virtual Network [%s] not found", name)
	}

	request := getVirtualNetworkRequest(wssdcommonproto.Operation_DELETE, name, &(*vnet)[0])
	_, err = c.VirtualNetworkAgentClient.Invoke(ctx, request)

	return err
}

// validate
func (c *client) validate(ctx context.Context, group, name string, vnet *network.VirtualNetwork) error {
	// Validate
	return nil
}

func getVirtualNetworkRequest(opType wssdcommonproto.Operation, name string, network *network.VirtualNetwork) *wssdnetwork.VirtualNetworkRequest {
	request := &wssdnetwork.VirtualNetworkRequest{
		OperationType:   opType,
		VirtualNetworks: []*wssdnetwork.VirtualNetwork{},
	}
	if network != nil {
		request.VirtualNetworks = append(request.VirtualNetworks, getWssdVirtualNetwork(network))
	} else if len(name) > 0 {
		request.VirtualNetworks = append(request.VirtualNetworks,
			&wssdnetwork.VirtualNetwork{
				Name: name,
			})
	}

	return request
}

func getVirtualNetworksFromResponse(response *wssdnetwork.VirtualNetworkResponse) *[]network.VirtualNetwork {
	virtualNetworks := []network.VirtualNetwork{}
	for _, vnet := range response.GetVirtualNetworks() {
		virtualNetworks = append(virtualNetworks, *(GetVirtualNetwork(vnet)))
	}

	return &virtualNetworks
}

// Conversion functions from network to wssdnetwork
func getWssdVirtualNetwork(c *network.VirtualNetwork) *wssdnetwork.VirtualNetwork {
	vnetType, _ := virtualNetworkTypeFromString(*c.Type)

	wssdvnet := &wssdnetwork.VirtualNetwork{
		Name: *c.Name,
		Type: vnetType,
	}
	if c.VirtualNetworkProperties == nil {
		return wssdvnet
	}

	wssdvnet.MacPool = getWssdMacPool(c.VirtualNetworkProperties.MACPool)

	wssdvnet.Ipams = getWssdNetworkIpams(c.VirtualNetworkProperties.Subnets)

	if c.DNSSettings == nil {
		return wssdvnet
	}
	wssdvnet.Dns = &wssdcommonproto.Dns{
		Domain:  *c.DNSSettings.Domain,
		Search:  *c.DNSSettings.Search,
		Servers: *c.DNSSettings.Servers,
		Options: *c.DNSSettings.Options,
	}

	return wssdvnet
}

func ipAllocationMethodProtobufToSdk(allocation wssdcommonproto.IPAllocationMethod) network.IPAllocationMethod {
	switch allocation {
	case wssdcommonproto.IPAllocationMethod_Static:
		return network.Static
	case wssdcommonproto.IPAllocationMethod_Dynamic:
		return network.Dynamic
	}
	return network.Dynamic
}

func ipAllocationMethodSdkToProtobuf(allocation network.IPAllocationMethod) wssdcommonproto.IPAllocationMethod {
	switch allocation {
	case network.Static:
		return wssdcommonproto.IPAllocationMethod_Static
	case network.Dynamic:
		return wssdcommonproto.IPAllocationMethod_Dynamic
	}
	return wssdcommonproto.IPAllocationMethod_Dynamic
}

func getWssdMacPool(macPool *network.MACPool) *wssdnetwork.MacPool {
	wssdMacPool := wssdnetwork.MacPool{}
	if macPool == nil || macPool.Ranges == nil {
		return &wssdMacPool
	}

	for _, macRange := range *macPool.Ranges {
		wssdMacRange := &wssdnetwork.MacRange{
			StartMacAddress: *macRange.StartMACAddress,
			EndMacAddress:   *macRange.EndMACAddress,
		}

		wssdMacPool.Ranges = append(wssdMacPool.Ranges, wssdMacRange)
	}

	return &wssdMacPool
}

func getWssdNetworkIpams(subnets *[]network.Subnet) []*wssdnetwork.Ipam {
	ipam := wssdnetwork.Ipam{}
	if subnets == nil {
		return []*wssdnetwork.Ipam{}
	}

	for _, subnet := range *subnets {
		wssdsubnet := &wssdnetwork.Subnet{
			Name: *subnet.Name,
			// TODO: implement something for IPConfigurationReferences
		}
		if subnet.Vlan == nil {
			wssdsubnet.Vlan = 0
		} else {
			wssdsubnet.Vlan = uint32(*subnet.Vlan)
		}
		if subnet.SubnetProperties != nil {
			if subnet.SubnetProperties.AddressPrefix != nil {
				wssdsubnet.Cidr = *subnet.SubnetProperties.AddressPrefix
			}
			wssdsubnet.Routes = getWssdNetworkRoutes(subnet.SubnetProperties.Routes)
		}
		wssdsubnet.Allocation = ipAllocationMethodSdkToProtobuf(subnet.SubnetProperties.IPAllocationMethod)

		ipam.Subnets = append(ipam.Subnets, wssdsubnet)
	}

	return []*wssdnetwork.Ipam{&ipam}
}

func getWssdNetworkRoutes(routes *[]network.Route) []*wssdnetwork.Route {
	wssdroutes := []*wssdnetwork.Route{}
	if routes == nil {
		return wssdroutes
	}

	for _, route := range *routes {
		if route.RouteProperties == nil {
			continue
		}
		wssdroutes = append(wssdroutes, &wssdnetwork.Route{
			Nexthop:           *route.RouteProperties.NextHop,
			Destinationprefix: *route.RouteProperties.DestinationPrefix,
			Metric:            route.RouteProperties.Metric,
		})
	}

	return wssdroutes
}

// Conversion function from wssdnetwork to network
func GetVirtualNetwork(c *wssdnetwork.VirtualNetwork) *network.VirtualNetwork {

	vnetType := virtualNetworkTypeToString(c.Type)
	vnet := &network.VirtualNetwork{
		Name: &c.Name,
		ID:   &c.Id,
		Type: &vnetType,
		VirtualNetworkProperties: &network.VirtualNetworkProperties{
			// TODO: MACPool (it is currently missing from network.VirtualNetwork)
			Subnets:           getNetworkSubnets(c.Ipams),
			ProvisioningState: status.GetProvisioningState(c.Status.GetProvisioningStatus()),
			Statuses:          status.GetStatuses(c.Status),
			MACPool:           getMacPool(c.MacPool),
		},
	}

	if c.Dns == nil {
		return vnet
	}

	vnet.VirtualNetworkProperties.DNSSettings = &network.DNSSetting{
		Domain:  &c.Dns.Domain,
		Search:  &c.Dns.Search,
		Servers: &c.Dns.Servers,
		Options: &c.Dns.Options,
	}

	return vnet
}

func getNetworkSubnets(ipams []*wssdnetwork.Ipam) *[]network.Subnet {
	subnets := []network.Subnet{}

	for _, ipam := range ipams {
		for _, subnet := range ipam.Subnets {
			subnets = append(subnets, network.Subnet{
				Name: &subnet.Name,
				ID:   &subnet.Id,
				SubnetProperties: &network.SubnetProperties{
					AddressPrefix: &subnet.Cidr,
					Routes:        getNetworkRoutes(subnet.Routes),
					// TODO: implement something for IPConfigurationReferences
					IPAllocationMethod: ipAllocationMethodProtobufToSdk(subnet.Allocation),
					Vlan:               getVlan(subnet.Vlan),
				},
			})
		}
	}

	return &subnets
}

func getNetworkRoutes(wssdroutes []*wssdnetwork.Route) *[]network.Route {
	routes := []network.Route{}

	for _, route := range wssdroutes {
		routes = append(routes, network.Route{
			RouteProperties: &network.RouteProperties{
				NextHop:           &route.Nexthop,
				DestinationPrefix: &route.Destinationprefix,
				Metric:            route.Metric,
			},
		})
	}

	return &routes
}

func getMacPool(wssdMacPool *wssdnetwork.MacPool) *network.MACPool {
	macPool := network.MACPool{}
	if wssdMacPool == nil || wssdMacPool.Ranges == nil {
		return &macPool
	}

	macRanges := []network.MACRange{}
	for _, wssdMacRange := range wssdMacPool.Ranges {
		macRange := network.MACRange{
			StartMACAddress: &wssdMacRange.StartMacAddress,
			EndMACAddress:   &wssdMacRange.EndMacAddress,
		}

		macRanges = append(macRanges, macRange)
	}
	macPool.Ranges = &macRanges

	return &macPool
}

func getVlan(wssdvlan uint32) *uint16 {
	vlan := uint16(wssdvlan)
	return &vlan
}
