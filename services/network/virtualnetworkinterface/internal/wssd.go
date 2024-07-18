// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license

package internal

import (
	"context"
	"fmt"

	wssdcommon "github.com/microsoft/moc/common"
	"github.com/microsoft/moc/pkg/auth"
	"github.com/microsoft/moc/pkg/errors"
	"github.com/microsoft/moc/pkg/status"
	prototags "github.com/microsoft/moc/pkg/tags"
	wssdcommonproto "github.com/microsoft/moc/rpc/common"
	wssdnetwork "github.com/microsoft/moc/rpc/nodeagent/network"
	wssdclient "github.com/microsoft/wssd-sdk-for-go/pkg/client"
	"github.com/microsoft/wssd-sdk-for-go/services/network"
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
	request, err := c.getVirtualNetworkInterfaceRequest(wssdcommonproto.Operation_GET, name, nil, "", "")
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
	request, err := c.getVirtualNetworkInterfaceRequest(wssdcommonproto.Operation_POST, name, vnetInterface, "", "")
	if err != nil {
		return nil, err
	}
	response, err := c.VirtualNetworkInterfaceAgentClient.Invoke(ctx, request)
	if err != nil {
		return nil, err
	}
	vnics, err := c.getVirtualNetworkInterfacesFromResponse(group, response)
	if err != nil {
		return nil, err
	}

	return &(*vnics)[0], nil
}

// Hydrate
func (c *client) Hydrate(ctx context.Context, group, name string, subnetId string, macAddress string) (*network.VirtualNetworkInterface, error) {
	// do we need to make a request to the logicalnetwork client for the subnetId?
	request, err := c.getVirtualNetworkInterfaceRequest(wssdcommonproto.Operation_HYDRATE, name, nil, subnetId, macAddress)
	if err != nil {
		return nil, err
	}
	response, err := c.VirtualNetworkInterfaceAgentClient.Invoke(ctx, request)
	if err != nil {
		return nil, err
	}
	vnics, err := c.getVirtualNetworkInterfacesFromResponse(group, response)
	if len(*vnics) == 0 {
		return nil, fmt.Errorf("Hydration of Virtual Network Interface failed with error: %v", err)
	}

	return &(*vnics)[0], nil
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

	request, err := c.getVirtualNetworkInterfaceRequest(wssdcommonproto.Operation_DELETE, name, &(*vnetInterface)[0], "", "")
	if err != nil {
		return err
	}
	_, err = c.VirtualNetworkInterfaceAgentClient.Invoke(ctx, request)

	if err != nil {
		return err
	}

	return err
}

// Update
func (c *client) Update(ctx context.Context, group, name string, vnetInterface *network.VirtualNetworkInterface) (*network.VirtualNetworkInterface, error) {
	request, err := c.getVirtualNetworkInterfaceRequest(wssdcommonproto.Operation_UPDATE, name, vnetInterface, "", "")
	if err != nil {
		return nil, err
	}
	response, err := c.VirtualNetworkInterfaceAgentClient.Invoke(ctx, request)
	if err != nil {
		return nil, err
	}
	vnics, err := c.getVirtualNetworkInterfacesFromResponse(group, response)
	if err != nil {
		return nil, err
	}

	return &(*vnics)[0], nil
}

// ///////////// private methods  ///////////////
func (c *client) getVirtualNetworkInterfaceRequest(opType wssdcommonproto.Operation, name string, networkInterface *network.VirtualNetworkInterface, subnetId string, macAddress string) (*wssdnetwork.VirtualNetworkInterfaceRequest, error) {
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
	} else if len(name) > 0 && len(subnetId) > 0 && len(macAddress) > 0 {
		ipconfig := &wssdnetwork.IpConfiguration{
			Subnetid: subnetId,
		}
		ipconfigs := []*wssdnetwork.IpConfiguration{ipconfig}

		request.VirtualNetworkInterfaces = append(request.VirtualNetworkInterfaces,
			&wssdnetwork.VirtualNetworkInterface{
				Name:       name,
				Macaddress: macAddress,
				Ipconfigs:  ipconfigs,
			})
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
		Name:        *c.Name,
		Ipconfigs:   wssdipconfigs,
		DnsSettings: cc.getDns(c.DNSSettings),
		Tags:        prototags.MapToProto(c.Tags),
	}

	if c.MACAddress != nil {
		vnic.Macaddress = *c.MACAddress
	}

	if c.VirtualMachineID != nil {
		vnic.VirtualMachineName = *c.VirtualMachineID
	}

	if c.EnableAcceleratedNetworking != nil {
		if *c.EnableAcceleratedNetworking {
			vnic.IovWeight = uint32(100)
		} else {
			vnic.IovWeight = uint32(0)
		}
	}

	vnic.Entity = cc.getWssdVirtualMachineEntity(c)
	return vnic, nil
}

func (c *client) getWssdVirtualMachineEntity(vnic *network.VirtualNetworkInterface) *wssdcommonproto.Entity {
	isPlaceholder := false
	if vnic.VirtualNetworkInterfaceProperties != nil && vnic.VirtualNetworkInterfaceProperties.IsPlaceholder != nil {
		isPlaceholder = *vnic.VirtualNetworkInterfaceProperties.IsPlaceholder
	}

	return &wssdcommonproto.Entity{
		IsPlaceholder: isPlaceholder,
	}
}

func networkTypeProtobufToSdk(networkType wssdnetwork.NetworkType) network.NetworkType {
	switch networkType {
	case wssdnetwork.NetworkType_LOGICAL_NETWORK:
		return network.Logical
	case wssdnetwork.NetworkType_VIRTUAL_NETWORK:
		return network.Virtual
	}
	return network.Virtual
}

func networkTypeSdkToProtobuf(networkType network.NetworkType) wssdnetwork.NetworkType {
	switch networkType {
	case network.Logical:
		return wssdnetwork.NetworkType_LOGICAL_NETWORK
	case network.Virtual:
		return wssdnetwork.NetworkType_VIRTUAL_NETWORK
	}
	return wssdnetwork.NetworkType_VIRTUAL_NETWORK
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
		Subnetid:    *ipconfig.SubnetID,
		NetworkType: networkTypeSdkToProtobuf(ipconfig.NetworkType),
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
			VirtualMachineID:            &c.VirtualMachineName,
			MACAddress:                  &c.Macaddress,
			DNSSettings:                 cc.getWssdDNSSettings(c.DnsSettings),
			IPConfigurations:            cc.getNetworkIpConfigs(c.Ipconfigs),
			ProvisioningState:           status.GetProvisioningState(c.Status.GetProvisioningStatus()),
			Statuses:                    status.GetStatuses(c.Status),
			IsPlaceholder:               cc.getVirtualNetworkIsPlaceholder(c),
			EnableAcceleratedNetworking: cc.getIovSetting(c),
		},
		Tags: prototags.ProtoToMap(c.Tags),
	}

	return vnetIntf, nil
}

func (cc *client) getDns(dnssetting *network.DNSSetting) *wssdcommonproto.Dns {
	if dnssetting == nil {
		return nil
	}
	var dns wssdcommonproto.Dns
	if dnssetting.Servers != nil {
		dns.Servers = *dnssetting.Servers
	}
	if dnssetting.Domain != nil {
		dns.Domain = *dnssetting.Domain
	}
	if dnssetting.Domain != nil {
		dns.Domain = *dnssetting.Domain
	}
	if dnssetting.Search != nil {
		dns.Search = *dnssetting.Search
	}
	if dnssetting.Options != nil {
		dns.Options = *dnssetting.Options
	}
	return &dns
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
				NetworkType:        networkTypeProtobufToSdk(wssdipconfig.NetworkType),
			},
		})
	}

	return &ipconfigs
}

func (cc *client) getWssdDNSSettings(dnssetting *wssdcommonproto.Dns) *network.DNSSetting {
	if dnssetting == nil {
		return nil
	}
	return &network.DNSSetting{
		Servers: &dnssetting.Servers,
		Domain:  &dnssetting.Domain,
		Search:  &dnssetting.Search,
		Options: &dnssetting.Options,
	}
}

func (c *client) getVirtualNetworkIsPlaceholder(vnic *wssdnetwork.VirtualNetworkInterface) *bool {
	isPlaceholder := false
	entity := vnic.GetEntity()
	if entity != nil {
		isPlaceholder = entity.IsPlaceholder
	}
	return &isPlaceholder
}

func (c *client) getIovSetting(vnic *wssdnetwork.VirtualNetworkInterface) *bool {
	isAcceleratedNetworkingEnabled := false
	if vnic.IovWeight > 0 {
		isAcceleratedNetworkingEnabled = true
	}
	return &isAcceleratedNetworkingEnabled
}
