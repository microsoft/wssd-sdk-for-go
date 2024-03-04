// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license

package internal

import (
	"context"
	"fmt"

	"github.com/microsoft/moc/pkg/auth"
	"github.com/microsoft/moc/pkg/status"
	prototags "github.com/microsoft/moc/pkg/tags"
	wssdcommonproto "github.com/microsoft/moc/rpc/common"
	wssdnetwork "github.com/microsoft/moc/rpc/nodeagent/network"
	wssdclient "github.com/microsoft/wssd-sdk-for-go/pkg/client"
	"github.com/microsoft/wssd-sdk-for-go/services/network"
)

type client struct {
	wssdnetwork.LogicalNetworkAgentClient
}

// NewLogicalNetworkClient - creates a client session with the backend wssd agent
func NewLogicalNetworkClient(subID string, authorizer auth.Authorizer) (*client, error) {

	c, err := wssdclient.GetLogicalNetworkClient(&subID, authorizer)
	if err != nil {
		return nil, err
	}
	return &client{c}, nil
}

// Get
func (c *client) Get(ctx context.Context, name string) (*[]network.LogicalNetwork, error) {
	request := getLogicalNetworkRequest(wssdcommonproto.Operation_GET, name, nil)
	response, err := c.LogicalNetworkAgentClient.Invoke(ctx, request)
	if err != nil {
		return nil, err
	}
	return getLogicalNetworksFromResponse(response), nil
}

// CreateOrUpdate
func (c *client) CreateOrUpdate(ctx context.Context, name string, lnet *network.LogicalNetwork) (*network.LogicalNetwork, error) {
	err := c.validate(ctx, name, lnet)
	if err != nil {
		return nil, err
	}
	request := getLogicalNetworkRequest(wssdcommonproto.Operation_POST, name, lnet)
	response, err := c.LogicalNetworkAgentClient.Invoke(ctx, request)
	if err != nil {
		return nil, err
	}
	lnets := getLogicalNetworksFromResponse(response)

	if len(*lnets) == 0 {
		return nil, fmt.Errorf("[LogicalNetwork][Create] Unexpected error: Creating a logical network returned no result")
	}

	return &((*lnets)[0]), nil
}

// Delete methods invokes create or update on the client
func (c *client) Delete(ctx context.Context, name string) error {
	lnet, err := c.Get(ctx, name)
	if err != nil {
		return err
	}
	if len(*lnet) == 0 {
		return fmt.Errorf("Logical Network [%s] not found", name)
	}

	request := getLogicalNetworkRequest(wssdcommonproto.Operation_DELETE, name, &(*lnet)[0])
	_, err = c.LogicalNetworkAgentClient.Invoke(ctx, request)

	return err
}

// validate
func (c *client) validate(ctx context.Context, name string, lnet *network.LogicalNetwork) error {
	// Validate
	return nil
}

func getLogicalNetworkRequest(opType wssdcommonproto.Operation, name string, network *network.LogicalNetwork) *wssdnetwork.LogicalNetworkRequest {
	request := &wssdnetwork.LogicalNetworkRequest{
		OperationType:   opType,
		LogicalNetworks: []*wssdnetwork.LogicalNetwork{},
	}
	if network != nil {
		request.LogicalNetworks = append(request.LogicalNetworks, getWssdLogicalNetwork(network))
	} else if len(name) > 0 {
		request.LogicalNetworks = append(request.LogicalNetworks,
			&wssdnetwork.LogicalNetwork{
				Name: name,
			})
	}

	return request
}

func getLogicalNetworksFromResponse(response *wssdnetwork.LogicalNetworkResponse) *[]network.LogicalNetwork {
	logicalNetworks := []network.LogicalNetwork{}
	for _, lnet := range response.GetLogicalNetworks() {
		logicalNetworks = append(logicalNetworks, *(GetLogicalNetwork(lnet)))
	}

	return &logicalNetworks
}

// Conversion functions from network to wssdnetwork
func getWssdLogicalNetwork(c *network.LogicalNetwork) *wssdnetwork.LogicalNetwork {

	wssdlnet := &wssdnetwork.LogicalNetwork{
		Name: *c.Name,
		Tags: getWssdTags(c.Tags),
	}
	if c.LogicalNetworkProperties == nil {
		return wssdlnet
	}

	wssdlnet.Ipams = getWssdNetworkIpams(c.LogicalNetworkProperties.Subnets)

	return wssdlnet
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

func getWssdNetworkIpams(subnets *[]network.LogicalSubnet) []*wssdnetwork.LogicalNetworkIpam {
	ipam := wssdnetwork.LogicalNetworkIpam{}
	if subnets == nil {
		return []*wssdnetwork.LogicalNetworkIpam{}
	}

	for _, subnet := range *subnets {
		wssdsubnet := &wssdnetwork.LogicalSubnet{
			Name: *subnet.Name,
			// TODO: implement something for IPConfigurationReferences
		}
		if subnet.Vlan == nil {
			wssdsubnet.Vlan = 0
		} else {
			wssdsubnet.Vlan = uint32(*subnet.Vlan)
		}
		if subnet.LogicalSubnetProperties != nil {
			if subnet.LogicalSubnetProperties.AddressPrefix != nil {
				wssdsubnet.AddressPrefix = *subnet.LogicalSubnetProperties.AddressPrefix
			}
			wssdsubnet.Routes = getWssdNetworkRoutes(subnet.LogicalSubnetProperties.Routes)
		}
		wssdsubnet.Allocation = ipAllocationMethodSdkToProtobuf(subnet.LogicalSubnetProperties.IPAllocationMethod)

		if subnet.DNSSettings != nil {
			dns := &wssdcommonproto.Dns{}
			if subnet.DNSSettings.Domain != nil {
				dns.Domain = *subnet.DNSSettings.Domain
			}
			if subnet.DNSSettings.Search != nil {
				dns.Search = *subnet.DNSSettings.Search
			}
			if subnet.DNSSettings.Servers != nil {
				dns.Servers = *subnet.DNSSettings.Servers
			}
			if subnet.DNSSettings.Options != nil {
				dns.Options = *subnet.DNSSettings.Options
			}
			wssdsubnet.Dns = dns
		}

		ipam.Subnets = append(ipam.Subnets, wssdsubnet)
	}

	return []*wssdnetwork.LogicalNetworkIpam{&ipam}
}

func getWssdNetworkRoutes(routes *[]network.Route) []*wssdcommonproto.Route {
	wssdroutes := []*wssdcommonproto.Route{}
	if routes == nil {
		return wssdroutes
	}

	for _, route := range *routes {
		if route.RouteProperties == nil {
			continue
		}
		wssdroutes = append(wssdroutes, &wssdcommonproto.Route{
			NextHop:           *route.RouteProperties.NextHop,
			DestinationPrefix: *route.RouteProperties.DestinationPrefix,
			Metric:            route.RouteProperties.Metric,
		})
	}

	return wssdroutes
}

// Conversion function from wssdnetwork to network
func GetLogicalNetwork(c *wssdnetwork.LogicalNetwork) *network.LogicalNetwork {

	lnet := &network.LogicalNetwork{
		Name: &c.Name,
		ID:   &c.Id,
		Tags: getNetworkTags(c.GetTags()),
		LogicalNetworkProperties: &network.LogicalNetworkProperties{
			Subnets:           getNetworkSubnets(c.Ipams),
			ProvisioningState: status.GetProvisioningState(c.Status.GetProvisioningStatus()),
			Statuses:          status.GetStatuses(c.Status),
		},
	}

	return lnet
}

func getNetworkSubnets(ipams []*wssdnetwork.LogicalNetworkIpam) *[]network.LogicalSubnet {
	subnets := []network.LogicalSubnet{}

	for _, ipam := range ipams {
		for _, subnet := range ipam.Subnets {

			dnsSettings := &network.DNSSetting{}
			if subnet.Dns != nil {
				dnsSettings = &network.DNSSetting{
					Domain:  &subnet.Dns.Domain,
					Search:  &subnet.Dns.Search,
					Servers: &subnet.Dns.Servers,
					Options: &subnet.Dns.Options,
				}
			}

			subnets = append(subnets, network.LogicalSubnet{
				Name: &subnet.Name,
				ID:   &subnet.Id,
				LogicalSubnetProperties: &network.LogicalSubnetProperties{
					AddressPrefix: &subnet.AddressPrefix,
					Routes:        getNetworkRoutes(subnet.Routes),
					// TODO: implement something for IPConfigurationReferences
					IPAllocationMethod: ipAllocationMethodProtobufToSdk(subnet.Allocation),
					Vlan:               getVlan(subnet.Vlan),
					DNSSettings:        dnsSettings,
				},
			})
		}
	}

	return &subnets
}

func getNetworkRoutes(wssdroutes []*wssdcommonproto.Route) *[]network.Route {
	routes := []network.Route{}

	for _, route := range wssdroutes {
		routes = append(routes, network.Route{
			RouteProperties: &network.RouteProperties{
				NextHop:           &route.NextHop,
				DestinationPrefix: &route.DestinationPrefix,
				Metric:            route.Metric,
			},
		})
	}

	return &routes
}

func getVlan(wssdvlan uint32) *uint16 {
	vlan := uint16(wssdvlan)
	return &vlan
}

func getNetworkTags(tags *wssdcommonproto.Tags) map[string]*string {
	return prototags.ProtoToMap(tags)
}

func getWssdTags(tags map[string]*string) *wssdcommonproto.Tags {
	return prototags.MapToProto(tags)
}
