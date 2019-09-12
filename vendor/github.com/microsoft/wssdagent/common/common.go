// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package common

import (
	"os/exec"

	"github.com/Microsoft/go-winio/pkg/guid"
)

// ExecuteCommand
func ExecuteCommand(name string, arg ...string) (output string, err error) {
	outbytes, err := exec.Command(name, arg...).Output()
	if err != nil {
		return "", err
	}
	return string(outbytes), err
}

// NewGuid
func NewGuid() string {
	g, err := guid.NewV4()
	if err != nil {
		return "00000000-0000-0000-0000-000000000000"
	}
	return g.String()
}
