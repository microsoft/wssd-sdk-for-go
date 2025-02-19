// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license

package virtualnetwork

import (
	"context"

	"github.com/microsoft/moc/pkg/auth"
	"github.com/microsoft/wssd-sdk-for-go/services/network"
	"github.com/microsoft/wssd-sdk-for-go/services/network/virtualnetwork/internal"
)

// Client structure
type VirtualNetworkClient interface {
	Get(context.Context, string, string) (*[]network.VirtualNetwork, error)
	CreateOrUpdate(context.Context, string, string, *network.VirtualNetwork) (*network.VirtualNetwork, error)
	Delete(context.Context, string, string) error
}

// NewClient method returns new client
func NewVirtualNetworkClient(cloudFQDN string, authorizer auth.Authorizer) (VirtualNetworkClient, error) {
	c, err := internal.NewVirtualNetworkClient(cloudFQDN, authorizer)
	if err != nil {
		return nil, err
	}

	return c, nil
}
