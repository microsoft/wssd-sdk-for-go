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

package network

import ()

type TransportProtocol string

const (
	// TransportProtocolAll
	TransportProtocolAll TransportProtocol = "All"
	// TransportProtocolTCP
	TransportProtocolTCP TransportProtocol = "Tcp"
	// TransportProtocolUDP
	TransportProtocolUDP TransportProtocol = "Udp"
)

// BaseProperties defines the structure of a Load Balancer
type BaseProperties struct {
	// ID
	ID *string `json:"ID,omitempty"`
	// Name
	Name *string `json:"name,omitempty"`
	// Type
	Type *string `json:"type,omitempty"`
	// Tags - Custom resource tags
	Tags map[string]*string `json:"tags"`
}

// Route is assoicated with a subnet.
type Route struct {
	BaseProperties
	// NextHop
	NextHop *string `json:"nexthop,omitempty"`
	// DestinationPrefix in cidr format
	DestinationPrefix *string `json:"destinationprefix,omitempty"`
}

// IPConfigurationReference
type IPConfigurationReference struct {
	// IPConfigurationID
	IPConfigurationID *string `json:"ID,omitempty"`
}

// Subnet is assoicated with a Virtual Network.
type Subnet struct {
	BaseProperties
	// Cidr for this subnet - IPv4, IPv6
	AddressPrefix *string `json:"cidr,omitempty"`
	// Routes for the subnet
	Routes *[]Route `json:"routes,omitempty"`
	// IPConfigurationReferences
	IPConfigurationReferences *[]IPConfigurationReference `json:"ipConfigurationReferences,omitempty"`
}

// MACRange is associated with MACPool and respresents the start and end addresses.
type MACRange struct {
	// StartMACAddress
	StartMACAddress *string `json:"startmacaddress,omitempty"`
	// EndMACAddress
	EndMACAddress *string `json:"endmacaddress,omitempty"`
}

// MACPool is assoicated with a network and represents pool of MACRanges.
type MACPool struct {
	// Ranges of mac
	Ranges *[]MACRange `json:"ranges,omitempty"`
}

// DNS (Domain Name System is associated with a network.
type DNS struct {
	// Domain
	Domain *string `json:"domain,omitempty"`
	// Search strings
	Search *[]string `json:"search,omitempty"`
	// Servers is list of nameservers
	Servers *[]string `json:"servers,omitempty"`
	// Options for DNS
	Options *[]string `json:"options,omitempty"`
}

// AddressSpace addressSpace contains an array of IP address ranges that can be used by subnets of the
// virtual network.
type AddressSpace struct {
	// AddressPrefixes - A list of address blocks reserved for this virtual network in CIDR notation.
	AddressPrefixes *[]string `json:"addressPrefixes,omitempty"`
}

type FrontendIPConfiguration struct {
	BaseProperties
	// IPAddress of the Frontend IP
	IPAddress *string `json:"ipaddress,omitempty"`
	// ID of the Subnet this frontend ip configuration belongs to
	SubnetID *string `json:"subnetID,omitempty"`
}

type BackendAddressPool struct {
	BaseProperties
	// BackendIPConfigurations
	BackendIPConfigurations *[]IPConfiguration `json:"backendIPConfigurations,omitempty"`
}

type LoadBalancingRule struct {
	BaseProperties
	// FrontendIPConfigurationID
	FrontendIPConfigurationID *string `json:"frontendIPConfigurationID,omitempty"`
	// BackendAddressPoolID
	BackendAddressPoolID *string `json:"backendAddressPoolID,omitempty"`
	// Dns
	Protocol *string `json:"protocol,omitempty"`
	// FrontendPort
	FrontendPort *int32 `json:"frontendPort,omitempty"`
	// BackendPort
	BackendPort *int32 `json:"backendPort,omitempty"`
}

// LoadBalancer defines the structure of a Load Balancer
type LoadBalancer struct {
	BaseProperties
	// FrontendIPConfigurations
	FrontendIPConfigurations *[]FrontendIPConfiguration `json:"frontendIPConfigurations,omitempty"`
	// BackendAddressPools
	BackendAddressPools *[]BackendAddressPool `json:"backendAddressPools,omitempty"`
	// LoadBalancingRules
	LoadBalancingRules *[]LoadBalancingRule `json:"loadBalancingRules,omitempty"`
}

// VirtualNetwork defines the structure of a VNET
type VirtualNetwork struct {
	BaseProperties
	// AddressSpace
	AddressSpace *AddressSpace `json:"addressSpace,omitempty"`
	// MACPool
	MACPool *MACPool `json:"macPool,omitempty"`
	// DNS
	DNSSettings DNS `json:"dnsSettings,omitempty"`
	// Subnets that could hold ipv4 and ipv6 subnets
	Subnets *[]Subnet `json:"subnets,omitempty"`
}

// IPConfiguration
type IPConfiguration struct {
	BaseProperties
	// IPAddress
	IPAddress *string `json:"ipaddress,omitempty"`
	// PrefixLength
	PrefixLength *string `json:"prefixlength,omitempty"`
	// SubnetID
	SubnetID *string `json:"subnetId,omitempty"`
	// VirtualNetworkInterface reference
	VirtualNetworkInterfaceID *string `json:"virtualNetworkInterfaceID,omitempty"`
	// LoadBalancerBackendAddressPoolIDs
	LoadBalancerBackendAddressPoolIDs *[]string `json:"loadBalancerBackendAddressPools,omitempty"`
	// LoadBalancerInboundNatPools
	LoadBalancerInboundNatPoolIDs *[]string `json:"loadBalancerInboundNatPools,omitempty"`
}

// VirtualNetwork defines the structure of a VNET
type VirtualNetworkInterface struct {
	BaseProperties
	// VirtualMAChineID
	VirtualMachineID *string `json:"virtualMAChineID,omitempty"`
	// Tags - Custom resource tags
	Tags map[string]*string `json:"tags"`
	// VirtualNetwork reference
	VirtualNetwork *VirtualNetwork `json:"virtualNetworkID,omitempty"`
	// IPConfigurations
	IPConfigurations *[]IPConfiguration `json:"ipConfigurations,omitempty"`
	// DNS
	DNSSettings *DNS `json:"dnsSettings,omitempty"`
	// Routes for the subnet
	Routes *[]Route `json:"routes,omitempty"`
	// MACAddress - the macaddress of the network interface
	MACAddress *string `json:"macAddress,omitempty"`
	// EnableIPForwarding
	EnableIPForwarding *bool `json:"enableIPForwarding,omitempty"`
	// EnableMACSpoofing - enable macspoofing on this nic
	EnableMACSpoofing *bool `json:"enableMACSpoofing,omitempty"`
	// EnableDHCPGuard
	EnableDHCPGuard *bool `json:"enableDHCPGuard,omitempty"`
	// EnableRouterAdvertisementGuard
	EnableRouterAdvertisementGuard *bool `json:"enableRouterAdvertisementGuard,omitempty"`
}
