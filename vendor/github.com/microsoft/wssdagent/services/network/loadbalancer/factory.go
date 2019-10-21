// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package loadbalancer

import (
	hcn "github.com/microsoft/wssdagent/services/network/loadbalancer/hcn"
	vmms "github.com/microsoft/wssdagent/services/network/loadbalancer/vmms"
	log "k8s.io/klog"
)

const (
	// HcnSpec
	HCNSpec = "hcn"
	// VMMSSpec
	VMMSSpec = "vmms"
)

var providerCache = map[string]LoadBalancerProvider{}

// GetLoadBalancerProvider
func GetLoadBalancerProvider(spec string) LoadBalancerProvider {
	if provider, ok := providerCache[spec]; ok {
		return provider
	}

	provider := newLoadBalancerProvider(spec)
	providerCache[spec] = provider
	return provider

}

func newLoadBalancerProvider(spec string) LoadBalancerProvider {
	log.Infof("Creating %s LoadBalancer Provider", spec)
	switch spec {
	case HCNSpec:
		return hcn.NewLoadBalancerProvider()
	case VMMSSpec:
		return vmms.NewLoadBalancerProvider()
	default:
		panic("missing loadbalancer provider")
	}
}
