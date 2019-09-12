// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package virtualmachine

import (
	log "k8s.io/klog"

	"github.com/microsoft/wssdagent/services/compute/virtualmachine/hcs"
	"github.com/microsoft/wssdagent/services/compute/virtualmachine/vmms"
)

const (
	Agent    = "VirtualMachine"
	HCSSpec  = "hcs"
	VMMSSpec = "vmms"
)

var providerCache = map[string]VirtualMachineProvider{}

// GetVirtualMachineProvider
func GetVirtualMachineProvider(spec string) VirtualMachineProvider {
	if provider, ok := providerCache[spec]; ok {
		return provider
	}

	provider := newVirtualMachineProvider(spec)
	providerCache[spec] = provider
	return provider

}

func newVirtualMachineProvider(spec string) VirtualMachineProvider {
	log.Infof("Creating %s VirtualMachine Provider", spec)
	switch spec {
	case HCSSpec:
		return hcs.NewVirtualMachineProvider()
	case VMMSSpec:
		return vmms.NewVirtualMachineProvider()
	default:
	}
	panic("missing provider")
}
