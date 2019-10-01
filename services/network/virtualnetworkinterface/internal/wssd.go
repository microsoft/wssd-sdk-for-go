// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license

package internal

import (
	"context"
	"fmt"
	"github.com/microsoft/wssd-sdk-for-go/services/network"

	virtualnetwork "github.com/microsoft/wssd-sdk-for-go/services/network/virtualnetwork"
	wssdclient "github.com/microsoft/wssdagent/rpc/client"
	wssdnetwork "github.com/microsoft/wssdagent/rpc/network"

	wssdcommon "github.com/microsoft/wssd-sdk-for-go/common"
)

const ()

type client struct {
	subID string
	wssdnetwork.VirtualNetworkInterfaceAgentClient
}

// NewVirtualNetworkInterfaceClientN- creates a client session with the backend wssd agent
func NewVirtualNetworkInterfaceClient(subID string) (*client, error) {
	c, err := wssdclient.GetVirtualNetworkInterfaceClient(&subID)
	if err != nil {
		return nil, err
	}
	return &client{subID, c}, nil
}

// Get
func (c *client) Get(ctx context.Context, group, name string) (*[]network.VirtualNetworkInterface, error) {
	request := c.getVirtualNetworkInterfaceRequest(wssdnetwork.Operation_GET, name, nil)
	response, err := c.VirtualNetworkInterfaceAgentClient.Invoke(ctx, request)
	if err != nil {
		return nil, err
	}
	vnetInt, err := c.getVirtualNetworkInterfacesFromResponse(group, response)
	if err != nil {
		return nil, err
	}

	return vnetInt, nil
}

// CreateOrUpdate
func (c *client) CreateOrUpdate(ctx context.Context, group, name string, vnetInterface *network.VirtualNetworkInterface) (*network.VirtualNetworkInterface, error) {
	request := c.getVirtualNetworkInterfaceRequest(wssdnetwork.Operation_POST, name, vnetInterface)
	response, err := c.VirtualNetworkInterfaceAgentClient.Invoke(ctx, request)
	if err != nil {
		return nil, err
	}
	vnets, err := c.getVirtualNetworkInterfacesFromResponse(group, response)
	if err != nil {
		return nil, err
	}

	return &(*vnets)[0], nil
}

// Delete methods invokes create or update on the client
func (c *client) Delete(ctx context.Context, group, name string) error {
	vnetInterface, err := c.Get(ctx, group, name)
	if err != nil {
		return err
	}
	if len(*vnetInterface) == 0 {
		return fmt.Errorf("Virtual Network Interface [%s] not found", name)
	}

	request := c.getVirtualNetworkInterfaceRequest(wssdnetwork.Operation_DELETE, name, &(*vnetInterface)[0])
	_, err = c.VirtualNetworkInterfaceAgentClient.Invoke(ctx, request)

	if err != nil {
		return err
	}

	return err
}

/////////////// private methods  ///////////////
func (c *client) getVirtualNetworkInterfaceRequest(opType wssdnetwork.Operation, name string, networkInterface *network.VirtualNetworkInterface) *wssdnetwork.VirtualNetworkInterfaceRequest {
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

func (c *client) getVirtualNetworkInterfacesFromResponse(group string, response *wssdnetwork.VirtualNetworkInterfaceResponse) (*[]network.VirtualNetworkInterface, error) {
	virtualNetworkInterfaces := []network.VirtualNetworkInterface{}

	for _, vnetInterface := range response.GetVirtualNetworkInterfaces() {
		vnetIntf, err := GetVirtualNetworkInterface(c.subID, group, vnetInterface)
		if err != nil {
			return nil, err
		}

		virtualNetworkInterfaces = append(virtualNetworkInterfaces, *vnetIntf)
	}

	return &virtualNetworkInterfaces, nil
}

// Conversion functions from network interface to wssd network interface
func GetWssdVirtualNetworkInterface(c *network.VirtualNetworkInterface) *wssdnetwork.VirtualNetworkInterface {

	vnic := &wssdnetwork.VirtualNetworkInterface{
		Name:        *c.Name,
		Id:          *c.ID,
		Networkname: *c.VirtualNetworkName,
		// TODO: Type
		Ipconfigs: getWssdNetworkInterfaceIPConfig(c.IPConfigurations),
	}

	if c.MACAddress != nil {
		vnic.Macaddress = *c.MACAddress
	}
	return vnic
}

func getWssdNetworkInterfaceIPConfig(ipConfigs *[]network.IPConfiguration) []*wssdnetwork.IpConfiguration {
	wssdIpConfigs := []*wssdnetwork.IpConfiguration{}
	if ipConfigs == nil {
		return wssdIpConfigs
	}

	for _, ipConfig := range *ipConfigs {
		if ipConfig.IPAddress == nil {
			continue
		}
		wssdIpConfigs = append(wssdIpConfigs, &wssdnetwork.IpConfiguration{
			Ipaddress:    *ipConfig.IPAddress,
			Prefixlength: *ipConfig.PrefixLength,
			Subnetid:     *ipConfig.SubnetID,
		})
	}

	return wssdIpConfigs
}

// Conversion function from wssd network interface to network interface
func GetVirtualNetworkInterface(server, group string, c *wssdnetwork.VirtualNetworkInterface) (*network.VirtualNetworkInterface, error) {

	//vnet, err := getVirtualNetwork(server, group, c.Networkname)
	//if err != nil {
	//	return nil, fmt.Errorf("Virtual Network Interface [%s] is not on a supported network type.\n Inner error: %v", c.Name, err)
	//}

	vnetIntf := &network.VirtualNetworkInterface{
		Name: &c.Name,
		ID:   &c.Id,
		VirtualNetworkInterfaceProperties: &network.VirtualNetworkInterfaceProperties{
			VirtualNetworkName: &c.Networkname,
			MACAddress:         &c.Macaddress,
			// TODO: Type
			IPConfigurations: getNetworkIpConfigs(c.Ipconfigs),
		},
	}

	return vnetIntf, nil
}

func getVirtualNetwork(server, group, networkName string) (*network.VirtualNetwork, error) {
	vnetclient, err := virtualnetwork.NewVirtualNetworkClient(server)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), wssdcommon.DefaultServerContextTimeout)
	defer cancel()

	networks, err := vnetclient.Get(ctx, group, networkName)
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
			IPConfigurationProperties: &network.IPConfigurationProperties{
				IPAddress:    &wssdipconfig.Ipaddress,
				PrefixLength: &wssdipconfig.Prefixlength,
				SubnetID:     &wssdipconfig.Subnetid,
			},
		})
	}

	return &ipconfigs
}
