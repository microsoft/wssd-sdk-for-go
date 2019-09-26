// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package config

type BaseConfiguration struct {
	// DataStorePath
	DataStorePath string
	// ConfigStorePath
	ConfigStorePath string
	// LogPath
	LogPath string
	// featureGates is a map of feature names to bools that enable or disable alpha/experimental
	FeatureGates map[string]bool
}

// ChildAgentConfiguration
type ChildAgentConfiguration struct {
	BaseConfiguration
}

// WSSDAgentConfiguration
type WSSDAgentConfiguration struct {
	// BaseConfiguration
	BaseConfiguration
	// ImageStore
	ImageStorePath string

	// Address to listen to
	Address string

	// Port
	Port int

	// tlsCertFile is the file containing x509 Certificate for HTTPS.  (CA cert,
	// if any, concatenated after server cert). If tlsCertFile and
	// tlsPrivateKeyFile are not provided, a self-signed certificate
	// and key are generated for the public address and saved to the directory
	TLSCertificateFile string
	TLSPrivateKeyFile  string
	TLSCipherSuites    []string
	TLSMinVersion      string

	ProviderConfigurations map[string]ChildAgentConfiguration
}
