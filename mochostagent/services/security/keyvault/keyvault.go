// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license

package keyvault

type SecretProperties struct {
	// VaultName
	VaultName *string `json:"vaultname"`
	// FileName
	FileName *string `json:"filename"`
	// State - State would container ProvisioningState-SubState
	Statuses map[string]*string `json:"statuses"`
	// ProvisioningState - READ-ONLY; The provisioning state, which only appears in the response.
	ProvisioningState *string `json:"provisioningState,omitempty"`
}

// Secret defines the structure of a secret
type Secret struct {
	// ID
	ID *string `json:"ID,omitempty"`
	// Name
	Name *string `json:"name,omitempty"`
	// Type
	Type *string `json:"type,omitempty"`
	// Tags - Custom resource tags
	Tags map[string]*string `json:"tags"`
	// Value
	Value *string `json:"value"`
	// Properties
	*SecretProperties `json:"properties,omitempty"`
}
