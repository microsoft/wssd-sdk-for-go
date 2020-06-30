// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license

package network

type TransportProtocol string

const (
	// TransportProtocolAll
	TransportProtocolAll TransportProtocol = "All"
	// TransportProtocolTCP
	TransportProtocolTCP TransportProtocol = "Tcp"
	// TransportProtocolUDP
	TransportProtocolUDP TransportProtocol = "Udp"
)

// RouteProperties
type RouteProperties struct {
	// NextHop
	NextHop *string `json:"nexthop,omitempty"`
	// DestinationPrefix in cidr format
	DestinationPrefix *string `json:"destinationprefix,omitempty"`
	// Metric
	Metric uint32 `json:"metric,omitempty"`
}

// Route is assoicated with a subnet.
type Route struct {
	// ID
	ID *string `json:"ID,omitempty"`
	// Name
	Name *string `json:"name,omitempty"`
	// Type
	Type *string `json:"type,omitempty"`
	// Tags - Custom resource tags
	Tags map[string]*string `json:"tags"`
	// Properties
	*RouteProperties `json:"properties,omitempty"`
}

// IPConfigurationReference
type IPConfigurationReference struct {
	// IPConfigurationID
	IPConfigurationID *string `json:"ID,omitempty"`
}

// IPAllocationMethod enumerates the values for ip allocation method.
type IPAllocationMethod string

const (
	// Dynamic ...
	Dynamic IPAllocationMethod = "Dynamic"
	// Static ...
	Static IPAllocationMethod = "Static"
)

// SubnetProperties
type SubnetProperties struct {
	// Cidr for this subnet - IPv4, IPv6
	AddressPrefix *string `json:"cidr,omitempty"`
	// Routes for the subnet
	Routes *[]Route `json:"routes,omitempty"`
	// IPConfigurationReferences
	IPConfigurationReferences *[]IPConfigurationReference `json:"ipConfigurationReferences,omitempty"`
	// IPAllocationMethod - The IP address allocation method. Possible values include: 'Static', 'Dynamic'
	IPAllocationMethod IPAllocationMethod `json:"ipAllocationMethod,omitempty"`
}

// Subnet is assoicated with a Virtual Network.
type Subnet struct {
	// ID
	ID *string `json:"ID,omitempty"`
	// Name
	Name *string `json:"name,omitempty"`
	// Type
	Type *string `json:"type,omitempty"`
	// Tags - Custom resource tags
	Tags map[string]*string `json:"tags"`
	// Properties
	*SubnetProperties `json:"properties,omitempty"`
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

// DNSSetting (Domain Name System is associated with a network.
type DNSSetting struct {
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

// FrontendIPConfigurationProperties
type FrontendIPConfigurationProperties struct {
	// IPAddress of the Frontend IP
	IPAddress *string `json:"ipaddress,omitempty"`
	// ID of the Subnet this frontend ip configuration belongs to
	SubnetID *string `json:"subnetID,omitempty"`
}

// FrontendIPConfiguration
type FrontendIPConfiguration struct {
	// ID
	ID *string `json:"ID,omitempty"`
	// Name
	Name *string `json:"name,omitempty"`
	// Type
	Type *string `json:"type,omitempty"`
	// Tags - Custom resource tags
	Tags map[string]*string `json:"tags"`
	// Properties
	*FrontendIPConfigurationProperties `json:"properties,omitempty"`
}

// BackendAddressPoolProperties
type BackendAddressPoolProperties struct {
	// BackendIPConfigurations
	BackendIPConfigurations *[]IPConfiguration `json:"backendIPConfigurations,omitempty"`
}

// BackendAddressPool
type BackendAddressPool struct {
	// ID
	ID *string `json:"ID,omitempty"`
	// Name
	Name *string `json:"name,omitempty"`
	// Type
	Type *string `json:"type,omitempty"`
	// Tags - Custom resource tags
	Tags map[string]*string `json:"tags"`

	// Properties
	*BackendAddressPoolProperties `json:"properties,omitempty"`
}

// LoadBalancingRuleProperties
type LoadBalancingRuleProperties struct {
	// FrontendIPConfigurationID
	FrontendIPConfigurationID *string `json:"frontendIPConfigurationID,omitempty"`
	// BackendAddressPoolID
	BackendAddressPoolID *string `json:"backendAddressPoolID,omitempty"`
	// TransportProtocol
	Protocol TransportProtocol `json:"protocol,omitempty"`
	// FrontendPort
	FrontendPort *int32 `json:"frontendPort,omitempty"`
	// BackendPort
	BackendPort *int32 `json:"backendPort,omitempty"`
}

// LoadBalancingRule
type LoadBalancingRule struct {
	// ID
	ID *string `json:"ID,omitempty"`
	// Name
	Name *string `json:"name,omitempty"`
	// Type
	Type *string `json:"type,omitempty"`
	// Tags - Custom resource tags
	Tags map[string]*string `json:"tags"`
	// Properties
	*LoadBalancingRuleProperties `json:"properties,omitempty"`
}

// LoadBalancerProperties
type LoadBalancerProperties struct {
	// FrontendIPConfigurations
	FrontendIPConfigurations *[]FrontendIPConfiguration `json:"frontendIPConfigurations,omitempty"`
	// BackendAddressPools
	BackendAddressPools *[]BackendAddressPool `json:"backendAddressPools,omitempty"`
	// LoadBalancingRules
	LoadBalancingRules *[]LoadBalancingRule `json:"loadBalancingRules,omitempty"`
	// ProvisioningState - READ-ONLY; The provisioning state, which only appears in the response.
	ProvisioningState *string `json:"provisioningState,omitempty"`
	// Statuses - Status
	Statuses map[string]*string `json:"statuses"`
}

// LoadBalancer defines the structure of a Load Balancer
type LoadBalancer struct {
	// ID
	ID *string `json:"ID,omitempty"`
	// Name
	Name *string `json:"name,omitempty"`
	// Type
	Type *string `json:"type,omitempty"`
	// Tags - Custom resource tags
	Tags map[string]*string `json:"tags"`
	// Properties
	*LoadBalancerProperties `json:"properties,omitempty"`
}

type VirtualNetworkProperties struct {
	// AddressSpace
	AddressSpace *AddressSpace `json:"addressSpace,omitempty"`
	// MACPool
	MACPool *MACPool `json:"macPool,omitempty"`
	// DNS
	DNSSettings *DNSSetting `json:"dnsSettings,omitempty"`
	// ProvisioningState - READ-ONLY; The provisioning state, which only appears in the response.
	ProvisioningState *string `json:"provisioningState,omitempty"`
	// Subnets that could hold ipv4 and ipv6 subnets
	Subnets *[]Subnet `json:"subnets,omitempty"`
	// Vlan
	Vlan *int32 `json:"vlan,omitempty"`
	// Statuses - Status
	Statuses map[string]*string `json:"statuses"`
}

// VirtualNetwork defines the structure of a VNET
type VirtualNetwork struct {
	// ID
	ID *string `json:"ID,omitempty"`
	// Name
	Name *string `json:"name,omitempty"`
	// Type
	Type *string `json:"type,omitempty"`
	// Tags - Custom resource tags
	Tags map[string]*string `json:"tags"`
	// Properties
	*VirtualNetworkProperties `json:"properties,omitempty"`
}

// IPConfigurationProperties
type IPConfigurationProperties struct {
	// IPAddress
	IPAddress *string `json:"ipaddress,omitempty"`
	// PrefixLength
	PrefixLength *string `json:"prefixlength,omitempty"`
	// SubnetID
	SubnetID *string `json:"subnetId,omitempty"`
	// Gateway
	Gateway *string `json:"gateway,omitempty"`
	// Primary indicates that this is the primary IPaddress of the Nic
	Primary *bool `json:"primary,omitempty"`
	// VirtualNetworkInterface reference
	VirtualNetworkInterfaceID *string `json:"virtualNetworkInterfaceID,omitempty"`
	// LoadBalancerBackendAddressPoolIDs
	LoadBalancerBackendAddressPoolIDs *[]string `json:"loadBalancerBackendAddressPools,omitempty"`
	// LoadBalancerInboundNatPools
	LoadBalancerInboundNatPoolIDs *[]string `json:"loadBalancerInboundNatPools,omitempty"`
	// IPAllocationMethod - The IP address allocation method. Possible values include: 'Static', 'Dynamic'
	IPAllocationMethod IPAllocationMethod `json:"ipAllocationMethod,omitempty"`
}

// IPConfiguration
type IPConfiguration struct {
	// ID
	ID *string `json:"ID,omitempty"`
	// Name
	Name *string `json:"name,omitempty"`
	// Type
	Type *string `json:"type,omitempty"`
	// Tags - Custom resource tags
	Tags map[string]*string `json:"tags"`
	// Properties
	*IPConfigurationProperties `json:"properties,omitempty"`
}

type VirtualNetworkInterfaceProperties struct {
	// VirtualMAChineID
	VirtualMachineID *string `json:"virtualMAChineID,omitempty"`
	// VirtualNetwork reference
	VirtualNetwork *VirtualNetwork `json:"virtualNetworkID,omitempty"`
	// IPConfigurations
	IPConfigurations *[]IPConfiguration `json:"ipConfigurations,omitempty"`
	// DNS
	DNSSettings *DNSSetting `json:"dnsSettings,omitempty"`
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
	// ProvisioningState - READ-ONLY; The provisioning state, which only appears in the response.
	ProvisioningState *string `json:"provisioningState,omitempty"`
	// Statuses - Status
	Statuses map[string]*string `json:"statuses"`
}

// VirtualNetwork defines the structure of a VNET
type VirtualNetworkInterface struct {
	// ID
	ID *string `json:"ID,omitempty"`
	// Name
	Name *string `json:"name,omitempty"`
	// Type
	Type *string `json:"type,omitempty"`
	// Tags - Custom resource tags
	Tags map[string]*string `json:"tags"`
	// Properties
	*VirtualNetworkInterfaceProperties `json:"properties,omitempty"`
}
