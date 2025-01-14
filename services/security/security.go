// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license

package security

type ClientType string

const (
	ControlPlane   ClientType = "ControlPlane"
	ExternalClient ClientType = "ExternalClient"
	Node           ClientType = "Node"
)

// KeyVaultProperties defines the structure of a Security Item
type KeyVaultProperties struct {
	SecretMap map[string]*string `json:"secretmap"`
	// State - State would container ProvisioningState-SubState
	Statuses map[string]*string `json:"statuses"`
	// ProvisioningState - READ-ONLY; The provisioning state, which only appears in the response.
	ProvisioningState *string `json:"provisioningState,omitempty"`
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
	// State - State would container ProvisioningState-SubState
	Statuses map[string]*string `json:"statuses"`
	// ProvisioningState - READ-ONLY; The provisioning state, which only appears in the response.
	ProvisioningState *string `json:"provisioningState,omitempty"`
	// Client type
	ClientType ClientType `json:"clienttype,omitempty"`
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
	// Certificate string encoded in base64
	Certificate *string `json:"certificate,omitempty"`
	// Token Expiry
	TokenExpiry *int64 `json:"tokenexpiry,omitempty"`
	// Token
	Token *string `json:"token,omitempty"`
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
	// State - State would container ProvisioningState-SubState
	Statuses map[string]*string `json:"statuses"`
	// ProvisioningState - READ-ONLY; The provisioning state, which only appears in the response.
	ProvisioningState *string `json:"provisioningState,omitempty"`
}

// Certificate a certificate consists of a certificate (X509) plus its attributes.
type Certificate struct {
	// ID - READ-ONLY; The certificate id
	ID *string `json:"id,omitempty"`
	// Name
	Name *string `json:"name,omitempty"`
	// Cer - CER contents of x509 certificate string encoded in base64
	Cer *string `json:"cer,omitempty"`
	// Type - The content type of the certificate
	Type *string `json:"contentType,omitempty"`
	// Attributes - The certificate attributes.
	Attributes *CertificateAttributes `json:"attributes,omitempty"`
	// Tags - Application-specific metadata in the form of key-value pairs
	Tags map[string]*string `json:"tags"`
}

// CertificateAttributes the certificate management attributes
type CertificateRequestAttributes struct {
	// DNSNames - DNS names to be added to the certificate
	DNSNames *[]string `json:"DNSNames,omitempty"`
	// IPs - IPs to be added to the certificate
	IPs *[]string `json:"IPs,omitempty"`
	// State - State
	Statuses map[string]*string `json:"statuses"`
}

// Certificate a certificate consists of a certificate (X509) plus its attributes.
type CertificateRequest struct {
	// Name - The certificate name
	Name *string `json:"name,omitempty"`
	// CaName - The ca certificate name to sign the certificate
	CaName *string `json:"caname,omitempty"`
	// PrivateKey Key contents of RSA Private Key string encoded in base64
	PrivateKey *string `json:"privatekey,omitempty"`
	// OldCertificate Certificate contents of x509 certificate string to be renewed encoded in base64
	OldCertificate *string `json:"oldcert,omitempty"`
	// ServerAuth - If the certificate to have ServerAuth for mTLS
	ServerAuth *bool `json:"serverauth,omitempty"`
	// Attributes - The certificate attributes.
	Attributes *CertificateRequestAttributes `json:"attributes,omitempty"`
	// Tags - Application-specific metadata in the form of key-value pairs
	Tags map[string]*string `json:"tags"`
}
