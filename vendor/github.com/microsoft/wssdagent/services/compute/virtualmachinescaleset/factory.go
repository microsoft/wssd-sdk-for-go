// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package virtualmachinescaleset

import (
	hcs "github.com/microsoft/wssdagent/services/compute/virtualmachinescaleset/hcs"
	vmms "github.com/microsoft/wssdagent/services/compute/virtualmachinescaleset/vmms"
	log "k8s.io/klog"
)

const (
	HCSSpec  = "hcs"
	VMMSSpec = "vmms"
)

var providerCache = map[string]VirtualMachineScaleSetProvider{}

// GetVirtualMachineScaleSetProvider
func GetVirtualMachineScaleSetProvider(spec string) VirtualMachineScaleSetProvider {
	if provider, ok := providerCache[spec]; ok {
		return provider
	}

	provider := newVirtualMachineScaleSetProvider(spec)
	providerCache[spec] = provider
	return provider

}

func newVirtualMachineScaleSetProvider(spec string) VirtualMachineScaleSetProvider {
	log.Infof("Creating %s VirtualMachineScaleSet Provider", spec)
	switch spec {
	case HCSSpec:
		return hcs.NewVirtualMachineScaleSetProvider()
	case VMMSSpec:
		return vmms.NewVirtualMachineScaleSetProvider()
	default:
	}
	panic("missing provider")
}
