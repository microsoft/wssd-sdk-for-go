// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package try

import (
	"github.com/microsoft/wssdagent/pkg/wssdagent/errors"
)

// Func represents functions that can be retried.
type Func func(attempt int) (retry bool, err error)

// Do keeps trying the function until the second argument
// returns false, or no error is returned.
func Do(retryCount int, fn Func) error {
	var err error
	var cont bool
	attempt := 1
	for {
		cont, err = fn(attempt)

		if !cont || err == nil {
			break
		}
		attempt++
		if attempt > retryCount {
			return errors.Wrapf(err, "retry limit exceeded: %d attempts", attempt)
		}
	}
	return err
}
