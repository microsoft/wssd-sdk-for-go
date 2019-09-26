// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package fs

import (
	//	"fmt"
	pb "github.com/microsoft/wssdagent/rpc/security"
)

type SecretProvider struct {
	client *Client
}

func NewSecretProvider() *SecretProvider {
	return &SecretProvider{
		client: NewClient(),
	}
}

func (svp *SecretProvider) Get(secrets []*pb.Secret) ([]*pb.Secret, error) {
	return svp.client.Get(secrets)
}

func (svp *SecretProvider) CreateOrUpdate(secrets []*pb.Secret) ([]*pb.Secret, error) {
	return svp.client.Create(secrets)
}

func (svp *SecretProvider) Delete(secrets []*pb.Secret) error {
	return svp.client.Delete(secrets)
}
