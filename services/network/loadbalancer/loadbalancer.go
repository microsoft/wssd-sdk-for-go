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

type FrontendIPConfiguration struct {
	// ID
	ID *string `json:"ID,omitempty"`
	// Name
	Name *string `json:"name,omitempty"`
	// IPAddress of the Frontend IP
	IPAddress *string `json:"ipaddress,omitempty"`
	// ID of the Subnet this frontend ip configuration belongs to
	SubnetID *string `json:"subnetID,omitempty"`
}

type BackendAddressPool struct {
	// ID
	ID *string `json:"ID,omitempty"`
	// Name
	Name *string `json:"name,omitempty"`
	// BackendIPConfigurations
	BackendIPConfigurations *[]IPConfiguration `json:"backendIPConfigurations,omitempty"`
}

type LoadBalancingRule struct {
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
	// ID
	ID *string `json:"ID,omitempty"`
	// Name
	Name *string `json:"name,omitempty"`
	// Type
	Type *string `json:"type,omitempty"`
	// FrontendIPConfigurations
	FrontendIPConfigurations *[]FrontendIPConfiguration `json:"frontendIPConfigurations,omitempty"`
	// BackendAddressPools
	BackendAddressPools *[]BackendAddressPool `json:"backendAddressPools,omitempty"`
	// LoadBalancingRules
	LoadBalancingRules *[]LoadBalancingRule `json:"loadBalancingRules,omitempty"`
	// Tags - Custom resource tags
	Tags map[string]*string `json:"tags"`
}
