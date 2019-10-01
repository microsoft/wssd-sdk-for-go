// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license

package storage

import ()

// BaseProperties defines the structure of a Load Balancer
type VirtualHardDiskProperties struct {
	// Path
	Source *string `json:"source,omitempty"`
	// DiskSizeGB
	DiskSizeGB *int32 `json:"diskSizeGB,omitempty"`
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
