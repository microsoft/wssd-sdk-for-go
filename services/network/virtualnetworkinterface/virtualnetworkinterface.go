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

// IPConfiguration
type IPConfiguration struct {
	// IPAddress
	IPAddress *string `json:"ipaddress,omitempty"`
	// PrefixLength
	PrefixLength *string `json:"prefixlength,omitempty"`
	// SubnetID
	IPSubnet *IPSubnet `json:"ipsubnet,omitempty"`
	// VirtualNetworkInterface reference
	VirtualNetworkInterface *VirtualNetworkInterface `json:",omitempty"`
}

// VirtualNetwork defines the structure of a VNET
type VirtualNetworkInterface struct {
	// ID
	ID *string `json:"id,omitempty"`
	// Name
	Name *string `json:"name,omitempty"`
	// VirtualMAChineID
	VirtualMAChineID *string `json:"virtualMAChineID,omitempty"`
	// Tags - Custom resource tags
	Tags map[string]*string `json:"tags"`
	// VirtualNetwork reference
	VirtualNetwork *VirtualNetwork `json:"virtualNetworkID,omitempty"`
	// IPConfigurations
	IPConfigurations *[]IPConfiguration `json:"ipConfigurations,omitempty"`
	// Type
	Type *string `json:"type,omitempty"`
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
