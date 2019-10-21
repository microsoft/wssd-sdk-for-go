// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package auth

import (
	"os"
	"path"
//	"fmt"
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
	WithAuthorization() credentials.PerRPCCredentials
}

type ManagedIdentityConfig struct {
	DotTokenPath string
}

func (ba *BearerAuthorizer) WithAuthorization() credentials.PerRPCCredentials {
	return ba.tokenProvider
}

type JwtTokenProvider struct {
	RawData string `json:"rawdata"`
}

// BearerAuthorizer implements the bearer authorization
type BearerAuthorizer struct {
	tokenProvider JwtTokenProvider
}

// NewBearerAuthorizer crates a BearerAuthorizer using the given token provider
func NewBearerAuthorizer(tp JwtTokenProvider) *BearerAuthorizer {
	return &BearerAuthorizer{tokenProvider: tp}
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
	s.Values[ClientTokenPath] = GetClientTokenLocation()

	return
}

func (settings EnvironmentSettings) GetAuthorizer() (Authorizer, error) {
	return settings.GetManagedIdentityConfig().Authorizer()
}

func (settings EnvironmentSettings) GetManagedIdentityConfig() ManagedIdentityConfig {
	return ManagedIdentityConfig{settings.Values[ClientTokenPath]}
}

func (mc ManagedIdentityConfig) Authorizer() (Authorizer, error) {

	jwtCreds := CredsFromTokenFile(GetClientTokenLocation())
	
	return NewBearerAuthorizer(jwtCreds), nil
}

func CredsFromTokenFile(tokenLocation string) JwtTokenProvider {
	data, err := ioutil.ReadFile(tokenLocation)
	if err != nil {
		return JwtTokenProvider{}
	}

	return JwtTokenProvider{string(data)}
}

func (c JwtTokenProvider) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		"authorization": c.RawData,
	}, nil
}

func (c JwtTokenProvider) RequireTransportSecurity() bool {
	return true
}

func GetClientTokenLocation() string {
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

func GetTLSClientCertConfiguration() string {
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
		GetClientTokenLocation(),
		[]byte(tokenStr),
		0644)
}

