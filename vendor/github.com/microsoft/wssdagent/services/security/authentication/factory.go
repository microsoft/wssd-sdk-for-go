// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package authentication

import (
	"sync"
)

var (
	providerCache *AuthenticationProvider
	mux           sync.Mutex
)

// GetAuthenticationProvider
func GetAuthenticationProvider() *AuthenticationProvider {
	mux.Lock()
	defer mux.Unlock()

	if providerCache == nil {
		providerCache = newAuthenticationProvider()
	}
	return providerCache

}

func newAuthenticationProvider() *AuthenticationProvider {
	return NewAuthenticationProvider()
}