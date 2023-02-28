// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license

package validation

import (
	"context"

	"github.com/microsoft/moc/pkg/auth"
	"github.com/microsoft/wssd-sdk-for-go/services/admin/validation/internal"
)

// Service interfacetype Service interface {
type Service interface {
	Validate(context.Context) error
}

// Client structure
type ValidationClient struct {
	internal Service
}

// NewClient method returns new client
func NewValidationClient(cloudFQDN string, authorizer auth.Authorizer) (*ValidationClient, error) {
	c, err := internal.NewValidationClient(cloudFQDN, authorizer)
	return &ValidationClient{c}, err
}

// validate
func (c *ValidationClient) Validate(ctx context.Context) error {
	return c.internal.Validate(ctx)
}
