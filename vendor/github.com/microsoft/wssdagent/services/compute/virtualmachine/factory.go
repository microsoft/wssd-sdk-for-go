// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package virtualmachine

import (
	"sync"
)

var (
	providerCache *VirtualMachineProvider
	mux           sync.Mutex
)

// GetVirtualMachineProvider
func GetVirtualMachineProvider() *VirtualMachineProvider {
	mux.Lock()
	defer mux.Unlock()

	if providerCache == nil {
		providerCache = newVirtualMachineProvider()
	}
	return providerCache

}

func newVirtualMachineProvider() *VirtualMachineProvider {
	return NewVirtualMachineProvider()
}
