// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package virtualmachinescaleset

import (
	"sync"
)

var (
	providerCache *VirtualMachineScaleSetProvider
	mux           sync.Mutex
)

// GetVirtualMachineScaleSetProvider
func GetVirtualMachineScaleSetProvider() *VirtualMachineScaleSetProvider {
	mux.Lock()
	defer mux.Unlock()

	if providerCache == nil {
		providerCache = newVirtualMachineScaleSetProvider()
	}
	return providerCache

}

func newVirtualMachineScaleSetProvider() *VirtualMachineScaleSetProvider {
	return NewVirtualMachineScaleSetProvider()
}
