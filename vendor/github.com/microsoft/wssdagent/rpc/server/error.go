// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package server

import (
	"github.com/microsoft/wssdagent/pkg/errors"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

func GetGRPCError(err error) error {
	if err == nil {
		return err
	}
	switch err {
	case errors.NotFound:
		return status.Errorf(codes.NotFound, err.Error())
	case errors.AlreadyExists:
		return status.Errorf(codes.AlreadyExists, err.Error())
	default:
		return err
	}
}
