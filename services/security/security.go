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
	Tags map[string]*string `json:"tags"`
	// Certificates
	Certificate *[]byte `json:"certificate,omitempty"`
	// Properties
	*IdentityProperties `json:"properties,omitempty"`
}

// CertificateAttributes the certificate management attributes
type CertificateAttributes struct {
	// Enabled - Determines whether the object is enabled
	Enabled *bool `json:"enabled,omitempty"`
	// NotBefore - Not before date in seconds since 1970-01-01T00:00:00Z
	NotBefore *int64 `json:"nbf,omitempty"`
	// Expires - Expiry date in seconds since 1970-01-01T00:00:00Z
	Expires *int64 `json:"exp,omitempty"`
	// Created - READ-ONLY; Creation time in seconds since 1970-01-01T00:00:00Z
	Created *int64 `json:"created,omitempty"`
	// Updated - READ-ONLY; Last updated time in seconds since 1970-01-01T00:00:00Z
	Updated *int64 `json:"updated,omitempty"`
}

// Certificate a certificate consists of a certificate (X509) plus its attributes.
type Certificate struct {
	// ID - READ-ONLY; The certificate id
	ID *string `json:"id,omitempty"`
	// Name
	Name *string `json:"name,omitempty"`
	// Cer - CER contents of x509 certificate.
	Cer *[]byte `json:"cer,omitempty"`
	// Type - The content type of the certificate
	Type *string `json:"contentType,omitempty"`
	// Attributes - The certificate attributes.
	Attributes *CertificateAttributes `json:"attributes,omitempty"`
	// Tags - Application-specific metadata in the form of key-value pairs
	Tags map[string]*string `json:"tags"`
}
