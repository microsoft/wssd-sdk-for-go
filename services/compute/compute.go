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

package compute

import (
	"github.com/microsoft/wssd-sdk-for-go/services/network"
)

// BaseProperties defines the structure of
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

// ref: github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2018-10-01/compute/models.go

// HardwareProfile
type HardwareProfile struct {
	VMSize *string `json:"vmSize,omitempty"`
}

type OperatingSystemTypes string

const (
	// Linux
	Linux OperatingSystemTypes = "Linux"
	// Windows
	Windows OperatingSystemTypes = "Windows"
)

type OSDisk struct {
	// Name
	Name *string `json:"name,omitempty"`
	// OsType
	OsType OperatingSystemTypes `json:"osType,omitempty"`
	// VhdId reference to virtual hard disk
	VhdId *string `json:"vhd,omitempty"`
}

type DataDisk struct {
	// Name
	Name *string `json:"name,omitempty"`
	// VhdId reference to VirtualHardDisk
	VhdId *string `json:"vhd,omitempty"`
}

type StorageProfile struct {
	// OSDisk
	OsDisk *OSDisk `json:"osDisk,omitempty"`
	// DataDisks
	DataDisks *[]DataDisk `json:"dataDisks,omitempty"`
}
type SSHPublicKey struct {
	// Path - Specifies the full path on the created VM where ssh public key is stored. If the file already exists, the specified key is appended to the file. Example: /home/user/.ssh/authorized_keys
	Path *string `json:"path,omitempty"`
	// KeyData - SSH public key certificate used to authenticate with the VM through ssh. The key needs to be at least 2048-bit and in ssh-rsa format. <br><br> For creating ssh keys, see [Create SSH keys on Linux and Mac for Li      nux VMs in Azure](https://docs.microsoft.com/azure/virtual-machines/virtual-machines-linux-mac-create-ssh-keys?toc=%2fazure%2fvirtual-machines%2flinux%2ftoc.json).
	KeyData *string `json:"keyData,omitempty"`
}

type SSHConfiguration struct {
	// PublicKeys - The list of SSH public keys used to authenticate with linux based VMs.
	PublicKeys *[]SSHPublicKey `json:"publicKeys,omitempty"`
}

type WindowsConfiguration struct {
	// EnableAutomaticUpdates
	EnableAutomaticUpdates *bool `json:"enableAutomaticUpdates,omitempty"`
	// TimeZone
	TimeZone *string `json:"timeZone,omitempty"`
	// AdditionalUnattendContent
	// AdditionalUnattendContent *[]AdditionalUnattendContent `json:"additionalUnattendContent,omitempty"`
	// SSH
	SSH *SSHConfiguration `json:"ssh,omitempty"`
}

type LinuxConfiguration struct {
	// SSH
	SSH *SSHConfiguration `json:"ssh,omitempty"`
	// DisablePasswordAuthentication
	DisablePasswordAuthentication *bool `json:"disablePasswordAuthentication,omitempty"`
}

type OSProfile struct {
	// ComputerName
	ComputerName *string `json:"computerName,omitempty"`
	// AdminUsername
	AdminUsername *string `json:"adminUsername,omitempty"`
	// AdminPassword
	AdminPassword *string `json:"adminPassword,omitempty"`
	// CustomData
	CustomData *string `json:"customData,omitempty"`
	// WindowsConfiguration
	WindowsConfiguration *WindowsConfiguration `json:"windowsConfiguration,omitempty"`
	// LinuxConfiguration
	LinuxConfiguration *LinuxConfiguration `json:"linuxConfiguration,omitempty"`
}

type NetworkInterfaceReference struct {
	// VirtualNetworkID
	VirtualNetworkID *string `json:"id,omitempty"`
	// VirtualNetworkInterfaceID
	VirtualNetworkInterfaceID *string `json:"id,omitempty"`
}
type NetworkProfile struct {
	// NetworkInterfaces
	NetworkInterfaceConfigurations *[]network.VirtualNetworkInterface `json:"networkInterfaceConfigurations,omitempty"`
}

type VirtualMachine struct {
	BaseProperties
	// StorageProfile
	StorageProfile *StorageProfile `json:"storageProfile,omitempty"`
	// OsProfile
	OsProfile *OSProfile `json:"osProfile,omitempty"`
	// NetworkProfile
	NetworkProfile *NetworkProfile `json:"networkProfile,omitempty"`
}

type Sku struct {
	// Name
	Name *string `json:"name,omitempty"`
	// Capacity
	Capacity *int64 `json:"capacity,omitempty"`
}

// VirtualMachineScaleSet
type VirtualMachineScaleSet struct {
	BaseProperties
	// Sku
	Sku *Sku `json:"sku,omitempty"`
	// VirtualMachineProfile
	VirtualMachineProfile *VirtualMachine `json:"virtualMachineProfile,omitempty"`
}
