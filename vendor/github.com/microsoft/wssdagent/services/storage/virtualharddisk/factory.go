// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package virtualharddisk

import (
	hcs "github.com/microsoft/wssdagent/services/storage/virtualharddisk/hcs"
	vmms "github.com/microsoft/wssdagent/services/storage/virtualharddisk/vmms"
	log "k8s.io/klog"
)

const (
	HCSSpec  = "hcs"
	VMMSSpec = "vmms"
)

var providerCache = map[string]VirtualHardDiskProvider{}

// GetVirtualHardDiskProvider
func GetVirtualHardDiskProvider(spec string) VirtualHardDiskProvider {
	if provider, ok := providerCache[spec]; ok {
		return provider
	}

	provider := newVirtualHardDiskProvider(spec)
	providerCache[spec] = provider
	return provider

}

func newVirtualHardDiskProvider(spec string) VirtualHardDiskProvider {
	log.Infof("Creating %s VirtualHardDisk Provider", spec)
	switch spec {
	case HCSSpec:
		return hcs.NewVirtualHardDiskProvider()
	case VMMSSpec:
		return vmms.NewVirtualHardDiskProvider()
	default:
		panic("missing storage provider")
	}
}
