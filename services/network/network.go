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

// Route is assoicated with a subnet.
type Route struct {
	// ID
	ID *string `json:"ID,omitempty"`
	// Name
	Name *string `json:"name,omitempty"`
	// Type
	Type *string `json:"type,omitempty"`
	// NextHop
	NextHop *string `json:"nexthop,omitempty"`
	// DestinationPrefix in cidr format
	DestinationPrefix *string `json:"destinationprefix,omitempty"`
}

// Subnet is assoicated with a Virtual Network.
type Subnet struct {
	// ID
	ID *string `json:"ID,omitempty"`
	// Name
	Name *string `json:"name,omitempty"`
	// Type
	Type *string `json:"type,omitempty"`
	// Cidr for this subnet - IPv4, IPv6
	AddressPrefix *string `json:"cidr,omitempty"`
	// Routes for the subnet
	Routes *[]Route `json:"routes,omitempty"`
	// IPConfigurations references that are on this IPSubnet
	IPConfigurations *[]IPConfiguration `json:"ipconfigurations,omitempty"`
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
