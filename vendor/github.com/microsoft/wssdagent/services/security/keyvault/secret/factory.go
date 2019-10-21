// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package secret

import (
	"sync"
)

var (
	providerCache *SecretProvider
	mux           sync.Mutex
)

// GetSecretProvider
func GetSecretProvider() *SecretProvider {
	mux.Lock()
	defer mux.Unlock()

	if providerCache == nil {
		providerCache = newSecretProvider()
	}
	return providerCache

}

func newSecretProvider() *SecretProvider {
	return NewSecretProvider()
}
