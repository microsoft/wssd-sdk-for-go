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

type VirtualNetworkInterface interface {
	Name(string) SpecOption
	AddressSpace(cidr string) SpecOption
	Subnet(name, cidr, nsgID, rtID string) SpecOption
}

// IpConfiguration
type IpConfiguration struct {
	// IpAddress
	IpAddress *string `json:"ipaddress,omitempty"`
	// PrefixLength
	PrefixLength *string `json:"prefixlength,omitempty"`
	// SubnetId
	IpSubnet *IpSubnet `json:"ipsubnet,omitempty"`
	// VirtualNetworkInterface reference
	VirtualNetworkInterface *VirtualNetworkInterface `json:",omitempty"`
}

// VirtualNetwork defines the structure of a VNET
type VirtualNetworkInterface struct {
	// Id
	Id *string `json:"id,omitempty"`
	// Name
	Name *string `json:"name,omitempty"`
	// VirtualMachineId
	VirtualMachineId *string `json:"virtualMachineId,omitempty"`
	// Tags - Custom resource tags
	Tags map[string]*string `json:"tags"`
	// VirtualNetwork reference
	VirtualNetwork *VirtualNetwork `json:"virtualNetworkId,omitempty"`
	// IpConfigurations
	IpConfigurations *[]IpConfiguration `json:"ipConfigurations,omitempty"`
	// Type
	Type *string `json:"type,omitempty"`
	// Dns
	DnsSettings *Dns `json:"dnsSettings,omitempty"`
	// Routes for the subnet
	Routes *[]Route `json:"routes,omitempty"`
	// MacAddress - the macaddress of the network interface
	MacAddress *string `json:"macAddress,omitempty"`
	// EnableIPForwarding
	EnableIPForwarding *bool `json:"enableIPForwarding,omitempty"`
	// EnableMacSpoofing - enable macspoofing on this nic
	EnableMacSpoofing *bool `json:"enableMacSpoofing,omitempty"`
	// EnableDhcpGuard
	EnableDhcpGuard *bool `json:"enableDhcpGuard,omitempty"`
	// EnableRouterAdvertisementGuard
	EnableRouterAdvertisementGuard *bool `json:"enableRouterAdvertisementGuard,omitempty"`
}
