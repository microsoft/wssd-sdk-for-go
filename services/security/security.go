// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license

package security

import ()

// BaseProperties defines the structure of a Security Item
type BaseProperties struct {
	// ID
	ID *string `json:"ID,omitempty"`
	// Name
	Name *string `json:"name,omitempty"`
	// Type
	Type *string `json:"type,omitempty"`
	// Tags - Custom resource tags
	Tags map[string]*string `json:"tags"`
}

// KeyVault defines the structure of a keyvault
type KeyVault struct {
	BaseProperties
	// KeyValues
	//SecretMap map[string]*string `json:"secretmap"`
}
