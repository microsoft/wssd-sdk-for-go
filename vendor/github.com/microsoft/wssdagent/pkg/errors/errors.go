// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package errors

import (
	"errors"
	perrors "github.com/pkg/errors"
)

var (
	NotFound             error = errors.New("Not Found")
	Failed               error = errors.New("Failed")
	InvalidConfiguration error = errors.New("Invalid Configuration")
	InvalidInput         error = errors.New("Invalid Input")
	InvalidFilter        error = errors.New("Invalid Filter")
	NotSupported         error = errors.New("Not Supported")
	AlreadyExists        error = errors.New("Already Exists")
)

func Wrap(cause error, message string) error {
	return perrors.Wrap(cause, message)
}

func Wrapf(err error, format string, args ...interface{}) error {
	return perrors.Wrapf(err, format, args)
}

func IsNotFound(err error) bool {
	return checkError(err, NotFound)
}
func IsAlreadyExists(err error) bool {
	return checkError(err, AlreadyExists)
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
