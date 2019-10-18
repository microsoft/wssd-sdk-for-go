// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package virtualnetwork

import (
	"sync"
)

var (
	providerCache *VirtualNetworkProvider
	mux           sync.Mutex
)

// GetVirtualNetworkProvider
func GetVirtualNetworkProvider() *VirtualNetworkProvider {
	mux.Lock()
	defer mux.Unlock()

	if providerCache == nil {
		providerCache = newVirtualNetworkProvider()
	}
	return providerCache

}

func newVirtualNetworkProvider() *VirtualNetworkProvider {
	return NewVirtualNetworkProvider()
}
