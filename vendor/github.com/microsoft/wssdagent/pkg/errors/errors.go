// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package errors

import (
	"errors"
	perrors "github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

var (
	NotFound             error = errors.New("Not Found")
	InvalidConfiguration error = errors.New("Invalid Configuration")
	InvalidInput         error = errors.New("Invalid Input")
	NotSupported         error = errors.New("Not Supported")
	AlreadyExists        error = errors.New("Already Exists")
	AlreadyInUse         error = errors.New("Already In Use")
	Duplicates           error = errors.New("Duplicates")
	InvalidFilter        error = errors.New("Invalid Filter")
	Failed               error = errors.New("Failed")
	InvalidGroup         error = errors.New("InvalidGroup")
	InvalidVersion       error = errors.New("InvalidVersion")
	OldVersion           error = errors.New("OldVersion")
	UpdateFailed         error = errors.New("Update Failed")
	NotInitialized       error = errors.New("Not Initialized")
	Unknown              error = errors.New("Unknown Reason")
)

func Wrap(cause error, message string) error {
	return perrors.Wrap(cause, message)
}

func Wrapf(err error, format string, args ...interface{}) error {
	return perrors.Wrapf(err, format, args...)
}

func GetGRPCErrorCode(err error) codes.Code {
	if derr, ok := status.FromError(err); ok {
		return derr.Code()
	}
	return codes.Unknown
}
func IsGRPCNotFound(err error) bool {
	if derr, ok := status.FromError(err); ok {
		return derr.Code() == codes.NotFound
	}
	return false
}

func IsGRPCAlreadyExist(err error) bool {
	if derr, ok := status.FromError(err); ok {
		return derr.Code() == codes.AlreadyExists
	}
	return false
}

func GetGRPCError(err error) error {
	if err == nil {
		return err
	}
	if IsNotFound(err) {
		return status.Errorf(codes.NotFound, err.Error())
	}
	if IsAlreadyExists(err) {
		return status.Errorf(codes.AlreadyExists, err.Error())
	}
	return err
}
func IsNotFound(err error) bool {
	return checkError(err, NotFound)
}
func IsAlreadyExists(err error) bool {
	return checkError(err, AlreadyExists)
}
func IsInvalidGroup(err error) bool {
	return checkError(err, InvalidGroup)
}
func checkError(wrappedError, err error) bool {
	if wrappedError == nil {
		return false
	}
	if wrappedError == err {
		return true
	}
	cerr := perrors.Cause(wrappedError)
	if cerr != nil && cerr == err {
		return true
	}
	return false

}

func New(errString string) error {
	return errors.New(errString)
}
