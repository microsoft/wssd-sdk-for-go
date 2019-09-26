// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package keyvault

import (
	fs "github.com/microsoft/wssdagent/services/security/keyvault/fs"
	log "k8s.io/klog"
)

const (
	FSSpec = "fs"
)

var providerCache = map[string]KeyVaultProvider{}

// GetKeyVaultProvider
func GetKeyVaultProvider(spec string) KeyVaultProvider {
	if provider, ok := providerCache[spec]; ok {
		return provider
	}

	provider := newKeyVaultProvider(spec)
	providerCache[spec] = provider
	return provider

}

func newKeyVaultProvider(spec string) KeyVaultProvider {
	log.Infof("Creating %s KeyVault Provider", spec)
	switch spec {
	case FSSpec:
		return fs.NewKeyVaultProvider()
	default:
	}
	panic("missing provider")
}
