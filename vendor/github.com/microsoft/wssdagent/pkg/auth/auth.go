// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package auth 

import (
	"encoding/pem"
	context "context"
	jwt "github.com/dgrijalva/jwt-go"
	"google.golang.org/grpc/metadata"
	"crypto/rsa"
	"crypto/rand"
	"crypto/x509"
	"fmt"
	"time"
	"io/ioutil"
)

const (
	RSAPublicKeyPemType = "RSA PUBLIC KEY"
	RSAPrivateKeyPemType = "RSA PRIVATE KEY"
	JWTIssuerName = "wssdagentsvc"
)

type JwtAuthorizer struct {
	jwtPublicKey *rsa.PublicKey
}

type JwtSigner struct {
	jwtPrivateKey *rsa.PrivateKey
}

type claims struct {
	Name string `json:"name"`
	jwt.StandardClaims
}

// NewJwtAuthorizer Creates Authorizer which will validate tokens
func NewJwtAuthorizer(keyLocation string) (*JwtAuthorizer, error)  {
	data, err := ioutil.ReadFile(keyLocation)
	if err != nil {

		// This means there is no public key saved.
		// Generate a random Public Key to protect the agent
		// and once a user logins in it will be replaced with a useful key
		reader := rand.Reader
		bitSize := 2048
		key, err := rsa.GenerateKey(reader, bitSize)
		if err != nil {
			return nil, err
		}
		
		return &JwtAuthorizer{&key.PublicKey}, fmt.Errorf("Error reading public key: %v", err)
	}
	
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(data)
	if err != nil {
		return nil, fmt.Errorf("Error parsing public key: %v", err)
	}
	
	return &JwtAuthorizer{publicKey}, nil
}

// NewJwtSigner Creates Signer which will sign tokens
func NewJwtSigner(privateKeyByte []byte) (*JwtSigner, error)  {
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyByte)
	if err != nil {
		return nil, fmt.Errorf("Error parsing private key: %v", err)
	}
	
	return &JwtSigner{privateKey}, nil
}

func validateToken(tokenString string, publicKey *rsa.PublicKey) (*jwt.Token, error) {
	parsedToken, err := jwt.Parse(tokenString, func (token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return publicKey, nil
	})
	if err != nil {
		return nil, err
	}
	if !parsedToken.Valid {
		return nil, fmt.Errorf("Valid Token Required") 
	}

	return parsedToken, nil
}

// ValidateTokenFromContext obtains the token from the context of the call
func (ja *JwtAuthorizer) ValidateTokenFromContext(context context.Context) (*jwt.Token, error) {
	var token *jwt.Token
	var err error

	md, ok := metadata.FromIncomingContext(context)
	if !ok {
		return nil, err
	}

	jwtToken, ok := md["authorization"]
	if !ok {
		return nil, err
	}

	token, err = validateToken(jwtToken[0], ja.jwtPublicKey)
	if err != nil {
		return nil, err
	}
	return token, nil
}

// UpdatePublicKey updates an Authorizer's Public Key
func (ja *JwtAuthorizer) UpdatePublicKey(publicKey *rsa.PublicKey) {
	ja.jwtPublicKey = publicKey
}

// WritePublicKeyToPem creates a pem for a public key and writes it to disk 
func (ja *JwtAuthorizer) WritePublicKeyToPem(keyLocation string) error {

	publicKeyBytes, err := x509.MarshalPKIXPublicKey(ja.jwtPublicKey)
	if err != nil {
		return err
	}

	pubkey := &pem.Block{
		Type:  RSAPublicKeyPemType,
		Bytes: publicKeyBytes,
	}

	outByte := pem.EncodeToMemory(pubkey)

	if outByte == nil {
		return fmt.Errorf("Could not Encode Public Key")
	}

	err = ioutil.WriteFile(
		keyLocation,
		outByte,
		0644)

	return err
}

// IssueJWT issues a JWT from the name and guid
func (js *JwtSigner) IssueJWT(name string, guid string) (string, error) {
	cl := &claims{
		Name: name,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
			Id: guid,
			IssuedAt: time.Now().Unix(),
			Issuer: JWTIssuerName,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, cl)

	return token.SignedString(js.jwtPrivateKey)
}

// GetPublicKey returns the Signer's PublicKey
func (js *JwtSigner) GetPublicKey() *rsa.PublicKey {
	return &js.jwtPrivateKey.PublicKey
}


// GeneratePrivateKey generates a private key for JWT
func GeneratePrivateKey() ([]byte, error) {
	reader := rand.Reader
	bitSize := 2048
	key, err := rsa.GenerateKey(reader, bitSize)
	if err != nil {
		return nil, err
	}

	privateKey := &pem.Block{
		Type: RSAPrivateKeyPemType,
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	}

	outByte := pem.EncodeToMemory(privateKey)

	if outByte == nil {
		return nil, fmt.Errorf("Could not Encode Secret")
	}

	return outByte, nil
}