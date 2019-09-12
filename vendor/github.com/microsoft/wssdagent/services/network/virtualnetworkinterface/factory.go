// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package virtualnetworkinterface

import (
	hcn "github.com/microsoft/wssdagent/services/network/virtualnetworkinterface/hcn"
	vmms "github.com/microsoft/wssdagent/services/network/virtualnetworkinterface/vmms"
	log "k8s.io/klog"
)

const (
	HCNSpec  = "hcn"
	VMMSSpec = "vmms"
)

var providerCache = map[string]VirtualNetworkInterfaceProvider{}

// GetVirtualNetworkInterfaceProvider
func GetVirtualNetworkInterfaceProvider(spec string) VirtualNetworkInterfaceProvider {
	if provider, ok := providerCache[spec]; ok {
		return provider
	}

	provider := newVirtualNetworkInterfaceProvider(spec)
	providerCache[spec] = provider
	return provider

}

func newVirtualNetworkInterfaceProvider(spec string) VirtualNetworkInterfaceProvider {
	log.Infof("Creating %s VirtualNetworkInterface Provider", spec)
	switch spec {
	case HCNSpec:
		return hcn.NewVirtualNetworkInterfaceProvider()
	case VMMSSpec:
		return vmms.NewVirtualNetworkInterfaceProvider()
	default:
	}
	panic("missing provider")
}
