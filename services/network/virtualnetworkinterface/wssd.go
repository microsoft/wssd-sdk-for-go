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

package virtualnetworkinterface

import (
	"context"
	"fmt"
	"github.com/microsoft/wssd-sdk-for-go/services/network"

	virtualnetwork "github.com/microsoft/wssd-sdk-for-go/services/network/virtualnetwork"
	wssdclient "github.com/microsoft/wssdagent/rpc/client"
	wssdnetwork "github.com/microsoft/wssdagent/rpc/network"
	log "k8s.io/klog"

	wssdcommon "github.com/microsoft/wssd-sdk-for-go/common"
)

const (

)

type client struct {
	subID string
	wssdnetwork.VirtualNetworkInterfaceAgentClient
}

// newVirtualNetworkInterfaceClient - creates a client session with the backend wssd agent
func newVirtualNetworkInterfaceClient(subID string) (*client, error) {
	c, err := wssdclient.GetVirtualNetworkInterfaceClient(&subID)
	if err != nil {
		return nil, err
	}
	return &client{subID, c}, nil
}

// Get
func (c *client) Get(ctx context.Context, name string) (*[]network.VirtualNetworkInterface, error) {
	request := getVirtualNetworkInterfaceRequest(wssdnetwork.Operation_GET, name, nil)
	response, err := c.VirtualNetworkInterfaceAgentClient.Invoke(ctx, request)
	if err != nil {
		return nil, err
	}
	vnetInt, err := getVirtualNetworkInterfacesFromResponse(c.subID, response)
	if err != nil {
		log.Errorf("[VirtualNetworkInterface][Get] getVirtualNetworkInterfacesFromResponse failed with error %v", err)
		return nil, err
	}

	return vnetInt, nil
}

// CreateOrUpdate
func (c *client) CreateOrUpdate(ctx context.Context, name string, id string, vnetInterface *network.VirtualNetworkInterface) (*network.VirtualNetworkInterface, error) {
	request := getVirtualNetworkInterfaceRequest(wssdnetwork.Operation_POST, name, vnetInterface)
	response, err := c.VirtualNetworkInterfaceAgentClient.Invoke(ctx, request)
	if err != nil {
		log.Errorf("[Virtual Network Interface] Create failed with error %v", err)
		return nil, err
	}
	log.Infof("[VirtualNetworkInterface][Create] [%v]", response)
	vnets, err := getVirtualNetworkInterfacesFromResponse(c.subID, response)
	if err != nil {
		log.Errorf("[VirtualNetworkInterface][Create] getVirtualNetworkInterfacesFromResponse failed with error %v", err)
		return nil, err
	}

	return &(*vnets)[0], nil
}

// Delete methods invokes create or update on the client
func (c *client) Delete(ctx context.Context, name string, id string) error {
	vnetInterface, err := c.Get(ctx, name)
	if err != nil {
		return err
	}
	if len(*vnetInterface) == 0 {
		return fmt.Errorf("Virtual Network Interface [%s] not found", name)
	}

	request := getVirtualNetworkInterfaceRequest(wssdnetwork.Operation_DELETE, name, &(*vnetInterface)[0])
	response, err := c.VirtualNetworkInterfaceAgentClient.Invoke(ctx, request)
	log.Infof("[Virtual Network Interface][Delete] [%v]", response)

	if err != nil {
		log.Errorf("[VirtualNetworkInterface][Delete] failed with error %v", err)
		return err
	}

	return err
}

func getVirtualNetworkInterfaceRequest(opType wssdnetwork.Operation, name string, networkInterface *network.VirtualNetworkInterface) *wssdnetwork.VirtualNetworkInterfaceRequest {
	request := &wssdnetwork.VirtualNetworkInterfaceRequest{
		OperationType:            opType,
		VirtualNetworkInterfaces: []*wssdnetwork.VirtualNetworkInterface{},
	}
	if networkInterface != nil {
		request.VirtualNetworkInterfaces = append(request.VirtualNetworkInterfaces, GetWssdVirtualNetworkInterface(networkInterface))
	} else if len(name) > 0 {
		request.VirtualNetworkInterfaces = append(request.VirtualNetworkInterfaces,
			&wssdnetwork.VirtualNetworkInterface{
				Name: name,
			})
	}
	return request
}

func getVirtualNetworkInterfacesFromResponse(server string, response *wssdnetwork.VirtualNetworkInterfaceResponse) (*[]network.VirtualNetworkInterface, error) {
	virtualNetworkInterfaces := []network.VirtualNetworkInterface{}

	for _, vnetInterface := range response.GetVirtualNetworkInterfaces() {
		vnetIntf, err := GetVirtualNetworkInterface(server, vnetInterface)
		if err != nil {
			return nil, err
		}

		virtualNetworkInterfaces = append(virtualNetworkInterfaces, *vnetIntf)
	}

	return &virtualNetworkInterfaces, nil
}

// Conversion functions from network interface to wssd network interface
func GetWssdVirtualNetworkInterface(c *network.VirtualNetworkInterface) *wssdnetwork.VirtualNetworkInterface {

	return &wssdnetwork.VirtualNetworkInterface{
		Name:        *c.BaseProperties.Name,
		Id:          *c.BaseProperties.ID,
		Networkname: *c.VirtualNetwork.BaseProperties.Name,
		// TODO: Type
		Ipconfigs: getWssdNetworkInterfaceIPConfig(c.IPConfigurations),
	}
}

func getWssdNetworkInterfaceIPConfig(ipConfigs *[]network.IPConfiguration) []*wssdnetwork.IpConfiguration {
	wssdIpConfigs := []*wssdnetwork.IpConfiguration{}

	for _, ipConfig := range *ipConfigs {
		wssdIpConfigs = append(wssdIpConfigs, &wssdnetwork.IpConfiguration{
			Ipaddress:    *ipConfig.IPAddress,
			Prefixlength: *ipConfig.PrefixLength,
			Subnetid:     *ipConfig.SubnetID,
		})
	}

	return wssdIpConfigs
}

// Conversion function from wssd network interface to network interface
func GetVirtualNetworkInterface(server string, c *wssdnetwork.VirtualNetworkInterface) (*network.VirtualNetworkInterface, error) {

	vnet, err := getVirtualNetwork(server, c.Networkname)
	if err != nil {
		return nil, fmt.Errorf("Virtual Network Interface [%s] is not on a supported network type.\n Inner error: %v", c.Name, err)
	}

	vnetIntf := &network.VirtualNetworkInterface{

		BaseProperties: network.BaseProperties{
			Name: &c.Name,
			ID:   &c.Id,
		},
		VirtualNetwork: vnet,
		// TODO: Type
		IPConfigurations: getNetworkIpConfigs(c.Ipconfigs),
	}

	return vnetIntf, nil
}

func getVirtualNetwork(server string, networkName string) (*network.VirtualNetwork, error) {
	vnetclient, err := virtualnetwork.NewVirtualNetworkClient(server)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), wssdcommon.DefaultServerContextTimeout)
	defer cancel()

	networks, err := vnetclient.Get(ctx, networkName)
	if err != nil {
		return nil, err
	}

	if len(*networks) > 0 {
		return &(*networks)[0], nil
	}

	return nil, fmt.Errorf("Virtual Network [%s] not found or network type not supported", networkName)
}

func getNetworkIpConfigs(wssdipconfigs []*wssdnetwork.IpConfiguration) *[]network.IPConfiguration {
	ipconfigs := []network.IPConfiguration{}

	for _, wssdipconfig := range wssdipconfigs {
		ipconfigs = append(ipconfigs, network.IPConfiguration{
			IPAddress:    &wssdipconfig.Ipaddress,
			PrefixLength: &wssdipconfig.Prefixlength,
			SubnetID:     &wssdipconfig.Subnetid,
		})
	}

	return &ipconfigs
}
