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

func (c *client) SetVerbosityLevel(ctx context.Context, verbositylevel int32) error {

	if verbositylevel < int32(wssdadmin.VerboseLevel_Min_Level) || verbositylevel > int32(wssdadmin.VerboseLevel_Max_Level) {
		return errors.New(`can not set provided verbositylevel, verbositylevel should be within the range [0,9] including 0 and 9`)
	}
	request := setVerbosityLevelRequest(verbositylevel)

	_, err := c.LogAgentClient.Set(ctx, request)
	return err
}

func (c *client) GetVerbosityLevel(ctx context.Context) (string, error) {

	request := getLevelRequest()

	res, err := c.LogAgentClient.GetLevel(ctx, request)
	return res.Level, err

}

func getLoggingRequest() *wssdadmin.LogRequest {
	return &wssdadmin.LogRequest{}
}

func setVerbosityLevelRequest(verbositylevel int32) *wssdadmin.SetRequest {
	return &wssdadmin.SetRequest{
		Verbositylevel: verbositylevel,
	}
}

func getLevelRequest() *wssdadmin.GetRequest {
	return &wssdadmin.GetRequest{}
}
