// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license

package storage

import (
	wssdstorage "github.com/microsoft/wssdagent/rpc/storage"
)

// BaseProperties defines the structure of a Load Balancer
type VirtualHardDiskProperties struct {
	// Path
	Path *string `json:"path,omitempty"`
	// Source
	Source *string `json:"source,omitempty"`
	// DiskSizeGB
	DiskSizeGB *int64 `json:"diskSizeGB,omitempty"`
	// Dynamic
	Dynamic *bool `json:"dynamic,omitempty"`
	// Blocksizebytes
	Blocksizebytes *int32 `json:"blocksizebytes,omitempty"`
	//Logicalsectorbytes
	Logicalsectorbytes *int32 `json:"logicalsectorbytes,omitempty"`
	//Physicalsectorbytes
	Physicalsectorbytes *int32 `json:"physicalsectorbytes,omitempty"`
	//Controllernumber
	Controllernumber *int64 `json:"controllernumber,omitempty"`
	//Controllerlocation
	Controllerlocation *int64 `json:"controllerlocation,omitempty"`
	//Disknumber
	Disknumber *int64 `json:"disknumber,omitempty"`
	//Vmname
	Vmname *string `json:"vmname,omitempty"`
	//Vmid
	Vmid *string `json:"vmid,omitempty"`
	//Scsipath
	Scsipath *string `json:"scsipath,omitempty"`
	//Virtualharddisktype
	Virtualharddisktype *wssdstorage.VirtualHardDiskType `json:"virtualharddisktype,omitempty"`
	// ProvisioningState - READ-ONLY; The provisioning state, which only appears in the response.
	ProvisioningState *string `json:"provisioningState,omitempty"`
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
