// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package keyvault

import (
	"sync"
)

var (
	providerCache *KeyVaultProvider
	mux           sync.Mutex
)

// GetKeyVaultProvider
func GetKeyVaultProvider() *KeyVaultProvider {
	mux.Lock()
	defer mux.Unlock()

	if providerCache == nil {
		providerCache = newKeyVaultProvider()
	}
	return providerCache

}

func newKeyVaultProvider() *KeyVaultProvider {
	return NewKeyVaultProvider()
}
