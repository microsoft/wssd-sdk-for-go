// Copyright 2019 (c) Microsoft and contributors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package storage

import ()

// BaseProperties defines the structure of a Load Balancer
type BaseProperties struct {
	// ID
	ID *string `json:"ID,omitempty"`
	// Name
	Name *string `json:"name,omitempty"`
	// Type
	Type *string `json:"type,omitempty"`
	// Tags - Custom resource tags
	Tags map[string]*string `json:"tags"`
}

// VirtualHardDisk defines the structure of a VHD
type VirtualHardDisk struct {
	BaseProperties
	// Path
	Source *string `json:"source,omitempty"`
	// DiskSizeGB
	DiskSizeGB *int32 `json:"diskSizeGB,omitempty"`
}
