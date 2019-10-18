// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package identity

import (
	"context"
	"github.com/microsoft/wssdagent/pkg/errors"
	pb "github.com/microsoft/wssdagent/rpc/security"
)

type IdentityProvider struct {
	client *Client
}

func NewIdentityProvider() *IdentityProvider {
	return &IdentityProvider{
		client: NewClient(),
	}
}

func (identityProv *IdentityProvider) Get(ctx context.Context, identities []*pb.Identity) ([]*pb.Identity, error) {
	newidentities := []*pb.Identity{}
	if len(identities) == 0 {
		// Get Everything
		return identityProv.client.Get(ctx, nil)
	}

	// Get only requested identities
	for _, iden := range identities {
		newidentity, err := identityProv.client.Get(ctx, iden)
		if err != nil {
			return newidentities, err
		}
		newidentities = append(newidentities, newidentity[0])
	}
	return newidentities, nil
}

func (identityProv *IdentityProvider) CreateOrUpdate(ctx context.Context, identities []*pb.Identity) ([]*pb.Identity, error) {
	newidentities := []*pb.Identity{}
	for _, iden := range identities {
		newidentity, err := identityProv.client.Create(ctx, iden)
		if err != nil {
			if err != errors.AlreadyExists {
				identityProv.client.Delete(ctx, iden)
			}
			return newidentities, err
		}
		newidentities = append(newidentities, newidentity)
	}

	return newidentities, nil
}

func (identityProv *IdentityProvider) Delete(ctx context.Context, identities []*pb.Identity) error {
	for _, iden := range identities {
		err := identityProv.client.Delete(ctx, iden)
		if err != nil {
			return err
		}
	}

	return nil
}

// GetIdentityByName
func (identityProv *IdentityProvider) GetIdentityByName(ctx context.Context, identityName string) (*pb.Identity, error) {
	identity, err := identityProv.Get(ctx, []*pb.Identity{&pb.Identity{Name: identityName}})
	if err != nil {
		return nil, err
	}
	if len(identity) == 0 {
		return nil, errors.NotFound
	}

	return identity[0], nil

}

// GetDataStorePath
func (identityProv *IdentityProvider) GetDataStorePath() string {
	return identityProv.client.GetDataStorePath()
}