// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license

package security

// KeyVaultProperties defines the structure of a Security Item
type KeyVaultProperties struct {
	SecretMap map[string]*string `json:"secretmap"`
}

// KeyVault defines the structure of a keyvault
type KeyVault struct {
	// ID
	ID *string `json:"ID,omitempty"`
	// Name
	Name *string `json:"name,omitempty"`
	// Type
	Type *string `json:"type,omitempty"`
	// Tags - Custom resource tags
	Tags map[string]*string `json:"tags"`
	// Properties
	*KeyVaultProperties `json:"properties,omitempty"`
}

// IdentityProperties defines the structure of a Security Item
type IdentityProperties struct {
}

// Identity defines the structure of a identity
type Identity struct {
	// ID
	ID *string `json:"ID,omitempty"`
	// Name
	Name *string `json:"name,omitempty"`
	// Type
	Type *string `json:"type,omitempty"`
	// Tags - Custom resource tags
	Certificate *[]byte            `json:"certificate,omitempty"`
	Tags        map[string]*string `json:"tags"`
	// Properties
	*IdentityProperties `json:"properties,omitempty"`
}
