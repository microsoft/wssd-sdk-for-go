// Copyright 2019 (c) Microsoft and contributors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package virtualnetwork

import (
	"context"
	"fmt"
	"github.com/microsoft/wssd-sdk-for-go/services/network"

	wssdclient "github.com/microsoft/wssdagent/rpc/client"
	wssdnetwork "github.com/microsoft/wssdagent/rpc/network"
	log "k8s.io/klog"
)

type client struct {
	wssdnetwork.VirtualNetworkAgentClient
}

// newClient - creates a client session with the backend wssd agent
func newVirtualNetworkClient(subID string) (*client, error) {
	c, err := wssdclient.GetVirtualNetworkClient(&subID)
	if err != nil {
		return nil, err
	}
	return &client{c}, nil
}

// Get
func (c *client) Get(ctx context.Context, name string) (*[]network.VirtualNetwork, error) {
	request := getVirtualNetworkRequest(wssdnetwork.Operation_GET, name, nil)
	response, err := c.VirtualNetworkAgentClient.Invoke(ctx, request)
	if err != nil {
		return nil, err
	}

	PrintListWssd(response.VirtualNetworks)
	log.Infof("[VirtualNetwork][Get] [%v]", response)
	return getVirtualNetworksFromResponse(response), nil
}

// CreateOrUpdate
func (c *client) CreateOrUpdate(ctx context.Context, name string, id string, vnet *network.VirtualNetwork) (*network.VirtualNetwork, error) {
	request := getVirtualNetworkRequest(wssdnetwork.Operation_POST, name, vnet)
	response, err := c.VirtualNetworkAgentClient.Invoke(ctx, request)
	if err != nil {
		log.Errorf("[Virtual Network] Create failed with error", err)
		return nil, err
	}
	log.Infof("[VirtualNetwork][Create] [%v]", response)
	vnets := getVirtualNetworksFromResponse(response)
	return &(*vnets)[0], nil
}

// Delete methods invokes create or update on the client
func (c *client) Delete(ctx context.Context, name string, id string) error {
	vnet, err := c.Get(ctx, name)
	if err != nil {
		return err
	}
	if len(*vnet) == 0 {
		return fmt.Errorf("Virtual Network [%s] not found", name)
	}

	request := getVirtualNetworkRequest(wssdnetwork.Operation_DELETE, name, &(*vnet)[0])
	response, err := c.VirtualNetworkAgentClient.Invoke(ctx, request)
	log.Infof("[Virtual Network][Delete] [%v]", response)

	return err
}

func getVirtualNetworkRequest(opType wssdnetwork.Operation, name string, network *network.VirtualNetwork) *wssdnetwork.VirtualNetworkRequest {
	request := &wssdnetwork.VirtualNetworkRequest{
		OperationType:         opType,
		VirtualNetworks: []*wssdnetwork.VirtualNetwork{},
	}
	if network != nil {
		request.VirtualNetworks = append(request.VirtualNetworks, GetWssdVirtualNetwork(network))
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
func GetWssdVirtualNetwork(c *network.VirtualNetwork) *wssdnetwork.VirtualNetwork {

	return &wssdnetwork.VirtualNetwork{
		Name: *c.BaseProperties.Name,
		Id:   *c.BaseProperties.ID,
		// TODO: MACPool (it is currently missing from network.VirtualNetwork)
		Ipams: getWssdNetworkIpams(c.Subnets),
		Dns:  &wssdnetwork.Dns{
			Domain: *c.DNSSettings.Domain,
			Search: *c.DNSSettings.Search,
			Servers: *c.DNSSettings.Servers,
			Options: *c.DNSSettings.Options,
		},
		Type: wssdnetwork.VirtualNetworkType_Transparent, // TODO: we should make this a parameter instead of hardcoding it here
	}
}

func getWssdNetworkIpams(subnets *[]network.Subnet) []*wssdnetwork.Ipam {
	ipam := wssdnetwork.Ipam{}

	for _, subnet := range *subnets {
		ipam.Subnets = append(ipam.Subnets, &wssdnetwork.Subnet {
			Name: *subnet.BaseProperties.Name,
			Id:   *subnet.BaseProperties.ID,
			Cidr: *subnet.AddressPrefix,
			Routes: getWssdNetworkRoutes(subnet.Routes),
			// TODO: implement something for IPConfigurationReferences
		})
	}

	return []*wssdnetwork.Ipam{&ipam}
}

func getWssdNetworkRoutes(routes *[]network.Route) []*wssdnetwork.Route {
	wssdroutes := []*wssdnetwork.Route{}

	for _, route := range *routes {
		wssdroutes = append(wssdroutes, &wssdnetwork.Route {
			Nexthop: *route.NextHop,
			Destinationprefix: *route.DestinationPrefix,
			Metric: route.Metric,
		})
	}

	return wssdroutes
}


// Conversion function from wssdnetwork to network
func GetVirtualNetwork(c *wssdnetwork.VirtualNetwork) *network.VirtualNetwork {

	return &network.VirtualNetwork{
		BaseProperties: network.BaseProperties{
			Name: &c.Name,
			ID:   &c.Id,
		},
		// TODO: MACPool (it is currently missing from network.VirtualNetwork)
		Subnets: getNetworkSubnets(c.Ipams),
		DNSSettings:  network.DNS{
			Domain: &c.Dns.Domain,
			Search: &c.Dns.Search,
			Servers: &c.Dns.Servers,
			Options: &c.Dns.Options,
		},
		// TODO: do something with c.VirtualNetworkType
	}
}

func getNetworkSubnets(ipams []*wssdnetwork.Ipam) *[]network.Subnet {
	subnets := []network.Subnet{}

	for _, ipam := range ipams {
		for _, subnet := range ipam.Subnets {
			subnets = append(subnets, network.Subnet {
				BaseProperties: network.BaseProperties{
					Name: &subnet.Name,
					ID:   &subnet.Id,
				},
				AddressPrefix: &subnet.Cidr,
				Routes: getNetworkRoutes(subnet.Routes),
				// TODO: implement something for IPConfigurationReferences
			})
		}
	}

	return &subnets
}

func getNetworkRoutes(wssdroutes []*wssdnetwork.Route) *[]network.Route {
	routes := []network.Route{}

	for _, route := range wssdroutes {
		routes = append(routes, network.Route {
			NextHop: &route.Nexthop,
			DestinationPrefix: &route.Destinationprefix,
			Metric: route.Metric,
		})
	}

	return &routes
}
