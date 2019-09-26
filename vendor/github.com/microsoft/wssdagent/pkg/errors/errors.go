// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package errors

import (
	"errors"
	perrors "github.com/pkg/errors"
)

var (
	NotFound             error = errors.New("Not Found")
	InvalidConfiguration error = errors.New("Invalid Configuration")
	InvalidInput         error = errors.New("Invalid Input")
	NotSupported         error = errors.New("Not Supported")
)

func Wrap(cause error, message string) error {
	return perrors.Wrap(cause, message)
}

func Wrapf(err error, format string, args ...interface{}) error {
	return perrors.Wrapf(err, format, args)
}
