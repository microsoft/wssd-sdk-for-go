// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package internal

import (
	"testing"
)

func Test_getWssdVirtualMachine(t *testing.T) {
}

func Test_getWssdVirtualMachineStorageConfiguration(t *testing.T) {}

func Test_getWssdVirtualMachineStorageConfigurationOsDisk(t *testing.T) {}

func Test_getWssdVirtualMachineStorageConfigurationDataDisks(t *testing.T) {}

func Test_getWssdVirtualMachineNetworkConfiguration(t *testing.T) {}

func Test_getWssdVirtualMachineOSSSHPublicKeys(t *testing.T) {}
func Test_getWssdVirtualMachineOSConfiguration(t *testing.T) {}

func Test_getVirtualMachineStorageProfile(t *testing.T)          {}
func Test_getVirtualMachineStorageProfileOsDisk(t *testing.T)    {}
func Test_getVirtualMachineStorageProfileDataDisks(t *testing.T) {}
func Test_getVirtualMachineNetworkProfile(t *testing.T)          {}
func Test_getVirtualMachineOSProfile(t *testing.T)               {}

func Test_getVirtualMachineRunCommandResponse_NilResponse(t *testing.T) {
	c := &client{}
	_, err := c.getVirtualMachineRunCommandResponse(nil)
	if err == nil {
		t.Error("expected error for nil response")
	}
}
