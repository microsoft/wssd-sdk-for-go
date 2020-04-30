// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license

package internal

import (
	"context"
	"fmt"
	"github.com/microsoft/moc/pkg/status"
	"github.com/microsoft/wssd-sdk-for-go/services/network"

	wssdcommon "github.com/microsoft/moc/common"
	"github.com/microsoft/moc/pkg/auth"
	"github.com/microsoft/moc/pkg/errors"
	wssdcommonproto "github.com/microsoft/moc/rpc/common"
	wssdnetwork "github.com/microsoft/moc/rpc/nodeagent/network"
	wssdclient "github.com/microsoft/wssd-sdk-for-go/pkg/client"
	virtualnetwork "github.com/microsoft/wssd-sdk-for-go/services/network/virtualnetwork"
)

type client struct {
	subID string
	wssdnetwork.VirtualNetworkInterfaceAgentClient
}

// NewVirtualNetworkInterfaceClientN- creates a client session with the backend wssd agent
func NewVirtualNetworkInterfaceClient(subID string, authorizer auth.Authorizer) (*client, error) {
	c, err := wssdclient.GetVirtualNetworkInterfaceClient(&subID, authorizer)
	if err != nil {
		return nil, err
	}
	return &client{subID, c}, nil
}

// Get
func (c *client) Get(ctx context.Context, group, name string) (*[]network.VirtualNetworkInterface, error) {
	request, err := c.getVirtualNetworkInterfaceRequest(wssdcommonproto.Operation_GET, name, nil)
	if err != nil {
		return nil, err
	}
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
	request, err := c.getVirtualNetworkInterfaceRequest(wssdcommonproto.Operation_POST, name, vnetInterface)
	if err != nil {
		return nil, err
	}
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

	request, err := c.getVirtualNetworkInterfaceRequest(wssdcommonproto.Operation_DELETE, name, &(*vnetInterface)[0])
	if err != nil {
		return err
	}
	_, err = c.VirtualNetworkInterfaceAgentClient.Invoke(ctx, request)

	if err != nil {
		return err
	}

	return err
}

/////////////// private methods  ///////////////
func (c *client) getVirtualNetworkInterfaceRequest(opType wssdcommonproto.Operation, name string, networkInterface *network.VirtualNetworkInterface) (*wssdnetwork.VirtualNetworkInterfaceRequest, error) {
	request := &wssdnetwork.VirtualNetworkInterfaceRequest{
		OperationType:            opType,
		VirtualNetworkInterfaces: []*wssdnetwork.VirtualNetworkInterface{},
	}
	if networkInterface != nil {
		wssdnetworkinterface, err := c.getWssdVirtualNetworkInterface(networkInterface)
		if err != nil {
			return nil, err
		}
		request.VirtualNetworkInterfaces = append(request.VirtualNetworkInterfaces, wssdnetworkinterface)
	} else if len(name) > 0 {
		request.VirtualNetworkInterfaces = append(request.VirtualNetworkInterfaces,
			&wssdnetwork.VirtualNetworkInterface{
				Name: name,
			})
	}
	return request, nil
}

func (c *client) getVirtualNetworkInterfacesFromResponse(group string, response *wssdnetwork.VirtualNetworkInterfaceResponse) (*[]network.VirtualNetworkInterface, error) {
	virtualNetworkInterfaces := []network.VirtualNetworkInterface{}

	for _, vnetInterface := range response.GetVirtualNetworkInterfaces() {
		vnetIntf, err := c.getVirtualNetworkInterface(c.subID, group, vnetInterface)
		if err != nil {
			return nil, err
		}

		virtualNetworkInterfaces = append(virtualNetworkInterfaces, *vnetIntf)
	}

	return &virtualNetworkInterfaces, nil
}

// Conversion functions from network interface to wssd network interface
func (cc *client) getWssdVirtualNetworkInterface(c *network.VirtualNetworkInterface) (*wssdnetwork.VirtualNetworkInterface, error) {
	if c.VirtualNetworkInterfaceProperties == nil {
		return nil, errors.Wrapf(errors.InvalidInput, "Missing Network Interface Properties")
	}

	wssdipconfigs := []*wssdnetwork.IpConfiguration{}
	for _, ipconfig := range *c.IPConfigurations {
		wssdipconfig, err := cc.getWssdNetworkInterfaceIPConfig(&ipconfig)
		if err != nil {
			return nil, err
		}
		wssdipconfigs = append(wssdipconfigs, wssdipconfig)
	}

	vnic := &wssdnetwork.VirtualNetworkInterface{
		Name:      *c.Name,
		Ipconfigs: wssdipconfigs,
	}

	if c.MACAddress != nil {
		vnic.Macaddress = *c.MACAddress
	}
	return vnic, nil
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

func (c *client) getWssdNetworkInterfaceIPConfig(ipconfig *network.IPConfiguration) (*wssdnetwork.IpConfiguration, error) {
	if ipconfig.IPConfigurationProperties == nil {
		return nil, errors.Wrapf(errors.InvalidInput, "Missing IPConfiguration Properties")
	}
	if ipconfig.IPConfigurationProperties.SubnetID == nil ||
		len(*ipconfig.IPConfigurationProperties.SubnetID) == 0 {
		return nil, errors.Wrapf(errors.InvalidInput, "Missing IPConfiguration Properties")
	}

	wssdipconfig := &wssdnetwork.IpConfiguration{
		Subnetid: *ipconfig.SubnetID,
	}
	if ipconfig.IPAddress != nil {
		wssdipconfig.Ipaddress = *ipconfig.IPAddress
	}
	if ipconfig.PrefixLength != nil {
		wssdipconfig.Prefixlength = *ipconfig.PrefixLength
	}
	if ipconfig.Gateway != nil {
		wssdipconfig.Gateway = *ipconfig.Gateway
	}
	wssdipconfig.Allocation = ipAllocationMethodSdkToProtobuf(ipconfig.IPAllocationMethod)

	return wssdipconfig, nil
}

// Conversion function from wssd network interface to network interface
func (cc *client) getVirtualNetworkInterface(server, group string, c *wssdnetwork.VirtualNetworkInterface) (*network.VirtualNetworkInterface, error) {
	vnetIntf := &network.VirtualNetworkInterface{
		Name: &c.Name,
		ID:   &c.Id,
		VirtualNetworkInterfaceProperties: &network.VirtualNetworkInterfaceProperties{
			MACAddress:        &c.Macaddress,
			IPConfigurations:  cc.getNetworkIpConfigs(c.Ipconfigs),
			ProvisioningState: status.GetProvisioningState(c.Status.GetProvisioningStatus()),
			Statuses:          status.GetStatuses(c.Status),
		},
	}

	return vnetIntf, nil
}

func (c *client) getVirtualNetwork(server, group, networkName string) (*network.VirtualNetwork, error) {

	authorizer, err := auth.NewAuthorizerFromEnvironment(server)
	if err != nil {
		return nil, err
	}

	vnetclient, err := virtualnetwork.NewVirtualNetworkClient(server, authorizer)
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

func (c *client) getNetworkIpConfigs(wssdipconfigs []*wssdnetwork.IpConfiguration) *[]network.IPConfiguration {
	ipconfigs := []network.IPConfiguration{}

	for _, wssdipconfig := range wssdipconfigs {
		ipconfigs = append(ipconfigs, network.IPConfiguration{
			IPConfigurationProperties: &network.IPConfigurationProperties{
				IPAddress:          &wssdipconfig.Ipaddress,
				PrefixLength:       &wssdipconfig.Prefixlength,
				SubnetID:           &wssdipconfig.Subnetid,
				Gateway:            &wssdipconfig.Gateway,
				IPAllocationMethod: ipAllocationMethodProtobufToSdk(wssdipconfig.Allocation),
			},
		})
	}

	return &ipconfigs
}
