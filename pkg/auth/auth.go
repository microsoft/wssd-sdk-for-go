// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package auth

import (
	"os"
	"path"
	"crypto/tls"
	"google.golang.org/grpc/credentials"
	"io/ioutil"
	context "context"
)

const (
	ClientTokenName = ".token"
	ClientCertName = "wssd.pem"
	ClientTokenPath = "WSSD_CLIENT_TOKEN"
	ClientCertPath = "WSSD_CLIENT_CERT"
	DefaultWSSDFolder = ".wssd"

)

type Authorizer interface {
	WithTransportAuthorization() credentials.TransportCredentials
	WithRPCAuthorization() credentials.PerRPCCredentials
}

type ManagedIdentityConfig struct {
	ClientTokenPath string
	ClientCertPath string
}

func (ba *BearerAuthorizer) WithRPCAuthorization() credentials.PerRPCCredentials {
	return ba.tokenProvider
}

func (ba *BearerAuthorizer) WithTransportAuthorization() credentials.TransportCredentials {
	return ba.transportCredentials
}

type JwtTokenProvider struct {
	RawData string `json:"rawdata"`
}

// BearerAuthorizer implements the bearer authorization
type BearerAuthorizer struct {
	tokenProvider JwtTokenProvider
	transportCredentials credentials.TransportCredentials
}

// NewBearerAuthorizer crates a BearerAuthorizer using the given token provider
func NewBearerAuthorizer(tp JwtTokenProvider, tc credentials.TransportCredentials) *BearerAuthorizer {
	return &BearerAuthorizer{
		tokenProvider: tp,
		transportCredentials: tc,
	}
}

// EnvironmentSettings contains the available authentication settings.
type EnvironmentSettings struct {
	Values      map[string]string
}

func NewAuthorizerFromEnvironment() (Authorizer, error) {
	settings, err := GetSettingsFromEnvironment()
	if err != nil {
		return nil, err
	}
	return settings.GetAuthorizer()
}

func GetSettingsFromEnvironment() (s EnvironmentSettings, err error) {
	s = EnvironmentSettings{
		Values: map[string]string{},
	}
	s.Values[ClientTokenPath] = getClientTokenLocation()
	s.Values[ClientCertPath] = getClientCertLocation()

	return
}

func (settings EnvironmentSettings) GetAuthorizer() (Authorizer, error) {
	return settings.GetManagedIdentityConfig().Authorizer()
}

func (settings EnvironmentSettings) GetManagedIdentityConfig() ManagedIdentityConfig {
	return ManagedIdentityConfig{
		settings.Values[ClientTokenPath],
		settings.Values[ClientCertPath],
	}
}

func (mc ManagedIdentityConfig) Authorizer() (Authorizer, error) {

	jwtCreds := TokenProviderFromFile(mc.ClientTokenPath)
	transportCreds := TransportCredentialsFromFile(mc.ClientCertPath)

	return NewBearerAuthorizer(jwtCreds, transportCreds), nil
}

func TokenProviderFromFile(tokenLocation string) JwtTokenProvider {
	data, err := ioutil.ReadFile(tokenLocation)
	if err != nil {
		// Call to open the token file most likely failed do to
		// token not being set. This is expected when the an identity is not yet
		// set. Log and continue
		return JwtTokenProvider{}
	}

	return JwtTokenProvider{string(data)}
}

func TransportCredentialsFromFile(certLocation string) credentials.TransportCredentials {
	creds, err := credentials.NewClientTLSFromFile(certLocation, "")
	if err != nil {
		// Call to open the cert file most likely failed do to 
		// cert not being set. We will just create an empty TLS
		// config, log and continue as this may be the desired outcome if
		// running in debug mode
    	return credentials.NewTLS(&tls.Config{})
	}
	return creds
}


func (c JwtTokenProvider) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		"authorization": c.RawData,
	}, nil
}

func (c JwtTokenProvider) RequireTransportSecurity() bool {
	return true
}

func getClientTokenLocation() string {
	clientTokenPath := os.Getenv(ClientTokenPath);
	if clientTokenPath == "" {
		wd, err := os.Getwd()
		if err != nil {
			panic(err)
		}

		// Create the default token path and set the 
		// env variable
		defaultPath := path.Join(wd, DefaultWSSDFolder)
		os.MkdirAll(defaultPath, os.ModePerm)
		clientTokenPath = path.Join(defaultPath, ClientTokenName)
		os.Setenv(ClientTokenPath, clientTokenPath)
	}
	return clientTokenPath
}

func getClientCertLocation() string {
	clientCertPath := os.Getenv(ClientCertPath);
	if clientCertPath == "" {
		wd, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		
		clientCertPath = path.Join(wd, ClientCertName)
	}
	return clientCertPath
}

func SaveToken(tokenStr string) error {
	return ioutil.WriteFile(
		getClientTokenLocation(),
		[]byte(tokenStr),
		0644)
}