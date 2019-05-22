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

type LoadBalancer interface {
	Name(string) SpecOption
	AddressSpace(cidr string) SpecOption
	Subnet(name, cidr, nsgID, rtID string) SpecOption
}

type FrontendIpConfiguration struct {
	// Id
	Id *string `json:"ID,omitempty"`
	// Name
	Name *string `json:"name,omitempty"`
	// IpAddress of the Frontend Ip
	IpAddress *string `json:"ipaddress,omitempty"`
	// Id of the Subnet this frontend ip configuration belongs to
	SubnetId *string `json:"subnetId,omitempty"`
}

type BackendAddressPool struct {
	// Id
	Id *string `json:"ID,omitempty"`
	// Name
	Name *string `json:"name,omitempty"`
	// BackendIPConfigurations
	BackendIPConfigurations *[]IpConfiguration `json:"backendIPConfigurations,omitempty"`
}

type LoadBalancingRule struct {
	// FrontendIpConfigurationId
	FrontendIpConfigurationId *string `json:"frontendIpConfigurationId,omitempty"`
	// BackendAddressPoolId
	BackendAddressPoolId *string `json:"backendAddressPoolId,omitempty"`
	// Dns
	Protocol *string `json:"protocol,omitempty"`
	// FrontendPort
	FrontendPort *int32 `json:"frontendPort,omitempty"`
	// BackendPort
	BackendPort *int32 `json:"backendPort,omitempty"`
}

// LoadBalancer defines the structure of a Load Balancer
type LoadBalancer struct {
	// Id
	Id *string `json:"ID,omitempty"`
	// Name
	Name *string `json:"name,omitempty"`
	// Type
	Type *string `json:"type,omitempty"`
	// FrontendIpConfigurations
	FrontendIpConfigurations *[]FrontendIpConfiguration `json:"frontendIpConfigurations,omitempty"`
	// BackendAddressPools
	BackendAddressPools *[]BackendAddressPool `json:"backendAddressPools,omitempty"`
	// LoadBalancingRules
	LoadBalancingRules *[]LoadBalancingRule `json:"loadBalancingRules,omitempty"`
	// Tags - Custom resource tags
	Tags map[string]*string `json:"tags"`
}
