// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package virtualnetwork

import (
	hcn "github.com/microsoft/wssdagent/services/network/virtualnetwork/hcn"
	vmms "github.com/microsoft/wssdagent/services/network/virtualnetwork/vmms"
	log "k8s.io/klog"
)

const (
	HCNSpec  = "hcn"
	VMMSSpec = "vmms"
)

var providerCache = map[string]VirtualNetworkProvider{}

// GetVirtualNetworkProvider
func GetVirtualNetworkProvider(spec string) VirtualNetworkProvider {
	if provider, ok := providerCache[spec]; ok {
		return provider
	}

	provider := newVirtualNetworkProvider(spec)
	providerCache[spec] = provider
	return provider

}

func newVirtualNetworkProvider(spec string) VirtualNetworkProvider {
	log.Infof("Creating %s VirtualNetwork Provider", spec)
	switch spec {
	case HCNSpec:
		return hcn.NewVirtualNetworkProvider()
	case VMMSSpec:
		return vmms.NewVirtualNetworkProvider()
	default:
	}
	panic("missing provider")
}
