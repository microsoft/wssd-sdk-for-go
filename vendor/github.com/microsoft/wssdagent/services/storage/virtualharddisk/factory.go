// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package virtualharddisk

import (
	"sync"
)

var (
	providerCache *VirtualHardDiskProvider
	mux           sync.Mutex
)

// GetVirtualHardDiskProvider
func GetVirtualHardDiskProvider() *VirtualHardDiskProvider {
	mux.Lock()
	defer mux.Unlock()

	if providerCache == nil {
		providerCache = newVirtualHardDiskProvider()
	}
	return providerCache

}

func newVirtualHardDiskProvider() *VirtualHardDiskProvider {
	return NewVirtualHardDiskProvider()
}
