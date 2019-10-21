// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package secret

import (
	"context"
	"github.com/microsoft/wssdagent/pkg/errors"
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

func (secProv *SecretProvider) Get(ctx context.Context, secs []*pb.Secret) ([]*pb.Secret, error) {
	newsecs := []*pb.Secret{}
	if len(secs) == 0 {
		// Get Everything
		return secProv.client.Get(ctx, nil)
	}

	// Get only requested secs
	for _, sec := range secs {
		newsec, err := secProv.client.Get(ctx, sec)
		if err != nil {
			return newsecs, err
		}
		newsecs = append(newsecs, newsec...)
	}
	return newsecs, nil
}

func (secProv *SecretProvider) CreateOrUpdate(ctx context.Context, secs []*pb.Secret) ([]*pb.Secret, error) {
	newsecs := []*pb.Secret{}
	for _, sec := range secs {
		newsec, err := secProv.client.Create(ctx, sec)
		if err != nil {
			if err != errors.AlreadyExists {
				secProv.client.Delete(ctx, sec)
			}
			return newsecs, err
		}
		newsecs = append(newsecs, newsec)
	}

	return newsecs, nil
}

func (secProv *SecretProvider) Delete(ctx context.Context, secs []*pb.Secret) error {
	for _, sec := range secs {
		err := secProv.client.Delete(ctx, sec)
		if err != nil {
			return err
		}
	}

	return nil
}
