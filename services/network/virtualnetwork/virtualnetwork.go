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

// Route is assoicated with a subnet.
type Route struct {
	// Id
	Id *string `json:"ID,omitempty"`
	// Name
	Name *string `json:"name,omitempty"`
	// Type
	Type *string `json:"type,omitempty"`
	// NextHop
	NextHop *string `json:"nexthop,omitempty"`
	// DestinationPrefix in cidr format
	DestinationPrefix *string `json:"destinationprefix,omitempty"`
}

// Subnet is assoicated with a Ipam.
type Subnet struct {
	// Id
	Id *string `json:"ID,omitempty"`
	// Name
	Name *string `json:"name,omitempty"`
	// Type
	Type *string `json:"type,omitempty"`
	// IpSubnets for the subnet
	IpSubnets *[]IpSubnet `json:ipsubnets",omitempty"`
	// Routes for the subnet
	Routes *[]Route `json:"routes,omitempty"`
}

// Subnet is assoicated with a Ipam.
type IpSubnet struct {
	// Id
	Id *string `json:"ID,omitempty"`
	// Name
	Name *string `json:"name,omitempty"`
	// Type
	Type *string `json:"type,omitempty"`
	// Cidr for this subnet - Ipv4, Ipv6
	Cidr *string `json:"cidr,omitempty"`
	// Routes for the subnet
	Routes *[]Route `json:"routes,omitempty"`
	// IpConfigurations references that are on this IpSubnet
	IpConfigurations *[]IpConfiguration `json:"ipconfigurations,omitempty"`
}

// Ipam (Internet Protocol Address Management) is assoicated with a network
// and represents the address space(s) of a network.
type Ipam struct {
	// Type of the subnet - Static/Dhcp
	Type *string `json:"type,omitempty"`
	// Subnets that could hold ipv4 and ipv6 subnets
	Subnets *[]Subnet `json:"subnets,omitempty"`
}

// MacRange is associated with MacPool and respresents the start and end addresses.
type MacRange struct {
	// StartMacAddress
	StartMacAddress *string `json:"startmacaddress,omitempty"`
	// EndMacAddress
	EndMacAddress *string `json:"endmacaddress,omitempty"`
}

// MacPool is assoicated with a network and represents pool of MacRanges.
type MacPool struct {
	// Ranges of mac
	Ranges *[]MacRange `json:"ranges,omitempty"`
}

// Dns (Domain Name System is associated with a network.
type Dns struct {
	// Domain
	Domain *string `json:"domain,omitempty"`
	// Search strings
	Search *[]string `json:"search,omitempty"`
	// Servers is list of nameservers
	Servers *[]string `json:"servers,omitempty"`
	// Options for Dns
	Options *[]string `json:"options,omitempty"`
}

// VirtualNetwork defines the structure of a VNET
type VirtualNetwork struct {
	// Id
	Id *string `json:"ID,omitempty"`
	// Name
	Name *string `json:"name,omitempty"`
	// Type
	Type *string `json:"type,omitempty"`
	// AddressSpace
	AddressSpace *string `json:"addressspace,omitempty"`
	// MacPool
	MacPool *MacPool `json:"macpool,omitempty"`
	// Dns
	Dns Dns `json:"dns,omitempty"`
	// Ipams
	Ipams *[]Ipam `json:"ipams,omitempty"`
	// Tags - Custom resource tags
	Tags map[string]*string `json:"tags"`
}
