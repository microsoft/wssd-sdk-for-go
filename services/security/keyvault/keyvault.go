// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license

package keyvault

import "time"

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

// JSONWebKeyEncryptionAlgorithm enumerates the values for json web key encryption algorithm.
type JSONWebKeyEncryptionAlgorithm string

const (
	// A256KW AES Key Wrap with 256 bit key-encryption key
	A256KW JSONWebKeyEncryptionAlgorithm = "A256KW"
)

// JSONWebKeyType enumerates the values for json web key type.
type JSONWebKeyType string

const (
	// AES Advanced Encrytion Standard.
	AES JSONWebKeyType = "AES"
)

type Key struct {
	// ID
	ID *string `json:"ID,omitempty"`
	// Name
	Name *string `json:"name,omitempty"`
	// VaultName
	VaultName *string `json:"vaultname"`
	// Algorithm
	Type *JSONWebKeyType `json:"keytype,omitempty"`
	// CreationTime
	CreationTime *time.Time `json:"ct,omitempty"`
	// KeyVersion
	KeyVersion *uint32 `json:"keyversion,omitempty"`
	// ProvisioningState - READ-ONLY; The provisioning state
	ProvisioningState *string `json:"provisioningState,omitempty"`
}

// KeyOperationResult the key operation result.
type KeyOperationResult struct {
	// Key
	*Key `json:"key,omitempty"`
	// Result - READ-ONLY; a URL-encoded base64 string
	Result *string `json:"result,omitempty"`
}

// KeyOperationResult the key operation result.
type KeyOperationRequest struct {
	// Key
	*Key `json:"key,omitempty"`
	// Algorithm
	Algorithm *JSONWebKeyEncryptionAlgorithm `json:"Algorithm,omitempty"`
	//Data
	Data *string `json:"data,omitempty"`
}
