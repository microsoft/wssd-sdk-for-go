// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package auth

import (
	context "context"
	"crypto/tls"
	"crypto/x509"
	"github.com/microsoft/wssdagent/pkg/marshal"
	"google.golang.org/grpc/credentials"
	"io/ioutil"
	"os"
	"path"
)

const (
	ClientTokenName   = ".token"
	ClientCertName    = "wssd.pem"
	ClientTokenPath   = "WSSD_CLIENT_TOKEN"
	AzConfigPath      = "AZCONFIG_PATH"
	DefaultWSSDFolder = ".wssd"
	ServerName        = "ServerName"
)

type azConfig struct {
	CloudCertificate  string
	ClientCertificate string
	ClientKey         string
}

type Authorizer interface {
	WithTransportAuthorization() credentials.TransportCredentials
	WithRPCAuthorization() credentials.PerRPCCredentials
}

type ManagedIdentityConfig struct {
	ClientTokenPath string
	AzConfigPath    string
	ServerName      string
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
	tokenProvider        JwtTokenProvider
	transportCredentials credentials.TransportCredentials
}

// NewBearerAuthorizer crates a BearerAuthorizer using the given token provider
func NewBearerAuthorizer(tp JwtTokenProvider, tc credentials.TransportCredentials) *BearerAuthorizer {
	return &BearerAuthorizer{
		tokenProvider:        tp,
		transportCredentials: tc,
	}
}

// EnvironmentSettings contains the available authentication settings.
type EnvironmentSettings struct {
	Values map[string]string
}

func NewAuthorizerFromEnvironment(serverName string) (Authorizer, error) {
	settings, err := GetSettingsFromEnvironment(serverName)
	if err != nil {
		return nil, err
	}
	return settings.GetAuthorizer()
}

func NewAuthorizerFromInput(tlsCert tls.Certificate, serverCertificate []byte, server string) (Authorizer, error) {
	transportCreds := TransportCredentialsFromNode(tlsCert, serverCertificate, server)
	return NewBearerAuthorizer(JwtTokenProvider{}, transportCreds), nil
}

func GetSettingsFromEnvironment(serverName string) (s EnvironmentSettings, err error) {
	s = EnvironmentSettings{
		Values: map[string]string{},
	}
	s.Values[ClientTokenPath] = getClientTokenLocation()
	s.Values[AzConfigPath] = GetAzConfigLocation()

	s.Values[ServerName] = serverName

	return
}

func (settings EnvironmentSettings) GetAuthorizer() (Authorizer, error) {
	return settings.GetManagedIdentityConfig().Authorizer()
}

func (settings EnvironmentSettings) GetManagedIdentityConfig() ManagedIdentityConfig {
	return ManagedIdentityConfig{
		settings.Values[ClientTokenPath],
		settings.Values[AzConfigPath],
		settings.Values[ServerName],
	}
}

func (mc ManagedIdentityConfig) Authorizer() (Authorizer, error) {

	jwtCreds := TokenProviderFromFile(mc.ClientTokenPath)
	transportCreds := TransportCredentialsFromFile(mc.AzConfigPath, mc.ServerName)

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

func TransportCredentialsFromFile(azConfigLocation string, server string) credentials.TransportCredentials {
	clientCerts := []tls.Certificate{}
	certPool := x509.NewCertPool()

	serverPem, tlsCert, err := readAccessFile(azConfigLocation)
	if err == nil {
		clientCerts = append(clientCerts, tlsCert)
		// Append the client certificates from the CA
		if ok := certPool.AppendCertsFromPEM(serverPem); !ok {
			return credentials.NewTLS(&tls.Config{})
		}
	}
	verifyPeerCertificate := func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
		// This is the for extra verification
		return nil
	}

	return credentials.NewTLS(&tls.Config{
		ServerName:            server,
		Certificates:          clientCerts,
		RootCAs:               certPool,
		VerifyPeerCertificate: verifyPeerCertificate,
	})

}

func readAccessFile(accessFileLocation string) ([]byte, tls.Certificate, error) {
	accessFile := azConfig{}
	err := marshal.FromJSONFile(accessFileLocation, &accessFile)
	if err != nil {
		return []byte{}, tls.Certificate{}, err
	}
	serverPem, err := marshal.FromBase64(accessFile.CloudCertificate)
	if err != nil {
		return []byte{}, tls.Certificate{}, err
	}
	clientPem, err := marshal.FromBase64(accessFile.ClientCertificate)
	if err != nil {
		return []byte{}, tls.Certificate{}, err
	}
	keyPem, err := marshal.FromBase64(accessFile.ClientKey)
	if err != nil {
		return []byte{}, tls.Certificate{}, err
	}
	tlsCert, err := tls.X509KeyPair(clientPem, keyPem)
	if err != nil {
		return []byte{}, tls.Certificate{}, err
	}

	return serverPem, tlsCert, nil
}

func TransportCredentialsFromNode(tlsCert tls.Certificate, serverCertificate []byte, server string) credentials.TransportCredentials {

	certPool := x509.NewCertPool()
	// Append the client certificates from the CA
	if ok := certPool.AppendCertsFromPEM(serverCertificate); !ok {
		return credentials.NewTLS(&tls.Config{})
	}
	verifyPeerCertificate := func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
		// This is the for extra verification
		return nil
	}

	return credentials.NewTLS(&tls.Config{
		ServerName:            server,
		Certificates:          []tls.Certificate{tlsCert},
		RootCAs:               certPool,
		VerifyPeerCertificate: verifyPeerCertificate,
	})

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
	clientTokenPath := os.Getenv(ClientTokenPath)
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

func GetAzConfigLocation() string {
	return os.Getenv(AzConfigPath)
}
func SaveToken(tokenStr string) error {
	return ioutil.WriteFile(
		getClientTokenLocation(),
		[]byte(tokenStr),
		0644)
}
