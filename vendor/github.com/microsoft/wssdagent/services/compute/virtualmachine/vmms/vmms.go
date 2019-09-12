// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package vmms

type BaseProperties struct {
	ID   *string
	Name *string
}

type VirtualMachineType int

const (
	RealizedVirtualMachine VirtualMachineType = 0
)

type VirtualMachineState int

const (
	Off VirtualMachineState = 0
)

type VirtualMachineHardDiskDrive struct {
}

type VirtualNetworkAdapter struct {
}
type VirtualMachine struct {
	BaseProperties
	Generation      *string
	Path            *string
	ProcessorCount  int
	Memory          int
	State           VirtualMachineState
	Type            VirtualMachineType
	HardDrives      *[]VirtualMachineHardDiskDrive
	NetworkAdapters *[]VirtualNetworkAdapter
}
