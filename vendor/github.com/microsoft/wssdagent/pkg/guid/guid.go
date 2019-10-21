// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package guid

import (
	"github.com/google/uuid"
)

// NewGuid
func NewGuid() string {
	g, err := uuid.NewUUID()
	if err != nil {
		return "00000000-0000-0000-0000-000000000000"
	}
	return g.String()
}
