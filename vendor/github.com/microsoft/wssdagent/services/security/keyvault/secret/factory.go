// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package secret

import (
	fs "github.com/microsoft/wssdagent/services/security/keyvault/secret/fs"
	log "k8s.io/klog"
)

const (
	FSSpec = "fs"
)

var providerCache = map[string]SecretProvider{}

// GetSecretProvider
func GetSecretProvider(spec string) SecretProvider {
	if provider, ok := providerCache[spec]; ok {
		return provider
	}

	provider := newSecretProvider(spec)
	providerCache[spec] = provider
	return provider

}

func newSecretProvider(spec string) SecretProvider {
	log.Infof("Creating %s Secret Provider", spec)
	switch spec {
	case FSSpec:
		return fs.NewSecretProvider()
	default:
	}
	panic("missing provider")
}
