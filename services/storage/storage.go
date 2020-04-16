// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license

package storage

// VirtualHardDiskProperties defines the structure of a Virtual HardDisk
type VirtualHardDiskProperties struct {
	// Path - READONLY
	Path *string `json:"path,omitempty"`
	// Source
	Source *string `json:"source,omitempty"`
	// DiskSizeBytes
	DiskSizeBytes *int64 `json:"disksizebytes,omitempty"`
	// Dynamic
	Dynamic *bool `json:"dynamic,omitempty"`
	// Blocksizebytes
	Blocksizebytes *int32 `json:"blocksizebytes,omitempty"`
	//Logicalsectorbytes
	Logicalsectorbytes *int32 `json:"logicalsectorbytes,omitempty"`
	//Physicalsectorbytes
	Physicalsectorbytes *int32 `json:"physicalsectorbytes,omitempty"`
	//Controllernumber - READONLY
	Controllernumber *int64 `json:"controllernumber,omitempty"`
	//Controllerlocation - READONLY
	Controllerlocation *int64 `json:"controllerlocation,omitempty"`
	//Disknumber - READONLY
	Disknumber *int64 `json:"disknumber,omitempty"`
	// VirtualMachineName
	VirtualMachineName *string `json:"virtualmachinename,omitempty"`
	//Scsipath - READONLY
	Scsipath *string `json:"scsipath,omitempty"`
	//Virtualharddisktype
	Virtualharddisktype string `json:"virtualharddisktype,omitempty"`
	// ProvisioningState - READ-ONLY; The provisioning state, which only appears in the response.
	ProvisioningState *string `json:"provisioningState,omitempty"`
	// Statuses - Status
	Statuses map[string]*string `json:"statuses"`
}

// VirtualHardDisk defines the structure of a VHD
type VirtualHardDisk struct {
	// ID
	ID *string `json:"ID,omitempty"`
	// Name
	Name *string `json:"name,omitempty"`
	// Type
	Type *string `json:"type,omitempty"`
	// Tags - Custom resource tags
	Tags map[string]*string `json:"tags"`
	// Properties
	*VirtualHardDiskProperties `json:"properties,omitempty"`
}

// ContainerProperties defines the structure of a ContainerProperties
type ContainerProperties struct {
	// Path
	Path *string `json:"path,omitempty"`
	// ProvisioningState - READ-ONLY; The provisioning state, which only appears in the response.
	ProvisioningState *string `json:"provisioningState,omitempty"`
	// Statuses - Status
	Statuses map[string]*string `json:"statuses"`
}

// VirtualHardDisk defines the structure of a VHD
type Container struct {
	// ID
	ID *string `json:"ID,omitempty"`
	// Name
	Name *string `json:"name,omitempty"`
	// Type
	Type *string `json:"type,omitempty"`
	// Tags - Custom resource tags
	Tags map[string]*string `json:"tags"`
	// Properties
	*ContainerProperties `json:"properties,omitempty"`
}
