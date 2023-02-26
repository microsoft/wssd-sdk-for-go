// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license

package internal

import (
	"context"
	"errors"
	"io"

	"github.com/microsoft/moc/pkg/auth"
	loggingHelpers "github.com/microsoft/moc/pkg/logging"
	wssdadmin "github.com/microsoft/moc/rpc/common/admin"
	wssdclient "github.com/microsoft/wssd-sdk-for-go/pkg/client"
)

type client struct {
	wssdadmin.LogAgentClient
}

// NewLoggingClient - creates a client session with the backend wssd agent
func NewLoggingClient(subID string, authorizer auth.Authorizer) (*client, error) {
	c, err := wssdclient.GetLogClient(&subID, authorizer)
	if err != nil {
		return nil, err
	}
	return &client{c}, nil
}

func (c *client) ForwardLogFile(ctx context.Context, forwardFunc func([]byte, error) error) error {
	request := getLoggingRequest()
	fileStreamClient, err := c.LogAgentClient.Get(ctx, request)
	if err != nil {
		return err
	}

	recFunc := func() ([]byte, error) {
		getLogFileResponse, innerErr := fileStreamClient.Recv()
		if innerErr != nil {
			return []byte{}, innerErr
		}
		if getLogFileResponse.Error == io.EOF.Error() {
			return getLogFileResponse.File, io.EOF
		}
		return getLogFileResponse.File, nil

	}
	return loggingHelpers.Forward(ctx, forwardFunc, recFunc)

	return nil
}

// Get
func (c *client) GetLogFile(ctx context.Context, filename string) error {
	request := getLoggingRequest()
	fileStreamClient, err := c.LogAgentClient.Get(ctx, request)
	if err != nil {
		return err
	}

	recFunc := func() ([]byte, error) {
		getLogFileResponse, innerErr := fileStreamClient.Recv()
		if innerErr != nil {
			return []byte{}, innerErr
		}
		if getLogFileResponse.Error == io.EOF.Error() {
			return getLogFileResponse.File, io.EOF
		}
		return getLogFileResponse.File, nil

	}
	return loggingHelpers.ReceiveFile(ctx, filename, recFunc)
}

func (c *client) SetTraceLevel(ctx context.Context, location string, settracelevel uint32) error {

	request := setTraceLevelRequest(settracelevel, location)
	fileStreamClient, err := c.LogAgentClient.Set(ctx, request)
	if err != nil {
		return err
	}
	getLogFileResponse, innerErr := fileStreamClient.Recv()

	if innerErr != nil {
		return innerErr
	}

	if getLogFileResponse.Done {
		return nil
	}

	Err := errors.New("error setting tracelevel")
	return Err
}

func getLoggingRequest() *wssdadmin.LogRequest {
	return &wssdadmin.LogRequest{}
}

func setTraceLevelRequest(settracelevel uint32, location string) *wssdadmin.SetRequest {
	return &wssdadmin.SetRequest{
		Settracelevel: settracelevel,
		Location:      location,
	}
}
