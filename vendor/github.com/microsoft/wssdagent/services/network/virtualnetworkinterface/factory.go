// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package virtualnetworkinterface

import (
	"sync"
)

var (
	providerCache *VirtualNetworkInterfaceProvider
	mux           sync.Mutex
)

// GetVirtualNetworkInterfaceProvider
func GetVirtualNetworkInterfaceProvider() *VirtualNetworkInterfaceProvider {
	mux.Lock()
	defer mux.Unlock()

	if providerCache == nil {
		providerCache = newVirtualNetworkInterfaceProvider()
	}
	return providerCache

}

func newVirtualNetworkInterfaceProvider() *VirtualNetworkInterfaceProvider {
	return NewVirtualNetworkInterfaceProvider()
}
