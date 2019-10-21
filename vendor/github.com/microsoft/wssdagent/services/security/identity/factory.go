// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package identity

import (
	"sync"
)

var (
	providerCache *IdentityProvider
	mux           sync.Mutex
)

// GetIdentityProvider
func GetIdentityProvider() *IdentityProvider {
	mux.Lock()
	defer mux.Unlock()

	if providerCache == nil {
		providerCache = newIdentityProvider()
	}
	return providerCache

}

func newIdentityProvider() *IdentityProvider {
	return NewIdentityProvider()
}
