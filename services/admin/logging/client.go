// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license

package logging

import (
	"context"

	"github.com/microsoft/moc/pkg/auth"
	"github.com/microsoft/wssd-sdk-for-go/services/admin/logging/internal"
)

// Service interfacetype Service interface {
type Service interface {
	GetLogFile(context.Context, string) error
	ForwardLogFile(context.Context, func([]byte, error) error) error
	SetVerbosityLevel(context.Context, int32) error
	GetVerbosityLevel(context.Context) (string, error)
}

// Client structure
type LoggingClient struct {
	internal Service
}

// NewClient method returns new client
func NewLoggingClient(cloudFQDN string, authorizer auth.Authorizer) (*LoggingClient, error) {
	c, err := internal.NewLoggingClient(cloudFQDN, authorizer)
	return &LoggingClient{c}, err
}

// function not typically exposed and is used to forward files
func (c *LoggingClient) ForwardLogFile(ctx context.Context, forwardFunc func([]byte, error) error) error {
	return c.internal.ForwardLogFile(ctx, forwardFunc)
}

// gets a file from the corresponding node agent and writes it to filename
func (c *LoggingClient) GetLogFile(ctx context.Context, filename string) error {
	return c.internal.GetLogFile(ctx, filename)
}

func (c *LoggingClient) SetVerbosityLevel(ctx context.Context, verbositylevel int32) error {
	return c.internal.SetVerbosityLevel(ctx, verbositylevel)
}

func (c *LoggingClient) GetVerbosityLevel(ctx context.Context) (string, error) {
	return c.internal.GetVerbosityLevel(ctx)
}
