// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license

package compute

import (
	"github.com/microsoft/wssd-sdk-for-go/services/network"
)

// ref: github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2018-10-01/compute/models.go

// VirtualMachineCustomSize Specifies cpu/memory information for custom VMSize types.
type VirtualMachineCustomSize struct {
	// CpuCount - Specifies number of CPUs for VM
	CpuCount *int32 `json:"cpucount,omitempty"`
	// MemoryMB - Specifies memory for VM in MB
	MemoryMB *int32 `json:"memorymb,omitempty"`
}

// DynamicMemoryConfiguration Specifies the dynamic memory configuration for a VM.
type DynamicMemoryConfiguration struct {
	// MaximumMemoryMB - Specifies the maximum amount of memory the VM is allowed to use.
	MaximumMemoryMB *uint64 `json:"maximummemorymb,omitempty"`
	// MinimumMemoryMB - Specifies the minimum amount of memory the VM is allocated.
	MinimumMemoryMB *uint64 `json:"minimummemorymb,omitempty"`
	// TargetMemoryBuffer - Specifies the size of the VMs memory buffer as a percentage of the current memory usage.
	TargetMemoryBuffer *uint32 `json:"targetmemorybuffer,omitempty"`
}

// HardwareProfile
type HardwareProfile struct {
	// VMSize - Specifies the size of the virtual machine.
	VMSize VirtualMachineSizeTypes `json:"vmSize,omitempty"`
	// CustomSize - Specifies cpu/memory information for custom VMSize types.
	CustomSize *VirtualMachineCustomSize `json:"customsize,omitempty"`
	// DynamicMemoryConfig - Specifies the dynamic memory configuration for a VM, dynamic memory will be enabled if this field is present.
	DynamicMemoryConfig *DynamicMemoryConfiguration `json:"dynamicmemoryconfig,omitempty"`
}

type OperatingSystemTypes string

const (
	// Linux
	Linux OperatingSystemTypes = "Linux"
	// Windows
	Windows OperatingSystemTypes = "Windows"
)

type OperatingSystemBootstrapEngine string

const (
	CloudInit          OperatingSystemBootstrapEngine = "CloudInit"
	WindowsAnswerFiles OperatingSystemBootstrapEngine = "WindowsAnswerFiles"
)

// ImageReference specifies information about the image to use. You can specify information about platform
// images, marketplace images, or virtual machine images. This element is required when you want to use a
// platform image, marketplace image, or virtual machine image, but is not used in other creation
// operations.
type ImageReference struct {
	// Publisher - The image publisher.
	Publisher *string `json:"publisher,omitempty"`
	// Offer - Specifies the offer of the platform image or marketplace image used to create the virtual machine.
	Offer *string `json:"offer,omitempty"`
	// Sku - The image SKU.
	Sku *string `json:"sku,omitempty"`
	// Version - Specifies the version of the platform image or marketplace image used to create the virtual machine. The allowed formats are Major.Minor.Build or 'latest'. Major, Minor, and Build are decimal numbers. Specify 'latest' to use the latest version of an image available at deploy time. Even if you use 'latest', the VM image will not automatically update after deploy time even if a new version becomes available.
	Version *string `json:"version,omitempty"`
	// ID - Resource Id
	ID *string `json:"id,omitempty"`
}

type OSDisk struct {
	// Name
	Name *string `json:"name,omitempty"`
	// OsType
	OsType OperatingSystemTypes `json:"osType,omitempty"`
	// VhdName reference to virtual hard disk
	VhdName *string `json:"vhd,omitempty"`
}

type DataDisk struct {
	// Name
	Name *string `json:"name,omitempty"`
	// VhdName reference to VirtualHardDisk
	VhdName        *string         `json:"vhd,omitempty"`
	ImageReference *ImageReference `json:"imageReference,omitempty"`
}

type StorageProfile struct {
	// ImageReference - Specifies information about the image to use. You can specify information about platform images, marketplace images, or virtual machine images. This element is required when you want to use a platform image, marketplace image, or virtual machine image, but is not used in other creation operations.
	ImageReference *ImageReference `json:"imageReference,omitempty"`
	// OSDisk
	OsDisk *OSDisk `json:"osDisk,omitempty"`
	// DataDisks
	DataDisks *[]DataDisk `json:"dataDisks,omitempty"`
	// VmConfigContainerName
	VmConfigContainerName *string `json:"vmConfigContainerName,omitempty"`
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

type RDPConfiguration struct {
	// Set to 'true' to disable Remote Desktop
	DisableRDP *bool
}

// ProtocolTypes enumerates the values for protocol types.
type ProtocolTypes string

const (
	// HTTP ...
	HTTP ProtocolTypes = "Http"
	// HTTPS ...
	HTTPS ProtocolTypes = "Https"
)

// WinRMConfiguration describes Windows Remote Management configuration of the VM
type WinRMConfiguration struct {
	// Listeners - The list of Windows Remote Management listeners
	Listeners *[]WinRMListener `json:"listeners,omitempty"`
}

// WinRMListener describes Protocol and thumbprint of Windows Remote Management listener
type WinRMListener struct {
	// Protocol - Specifies the protocol of WinRM listener. Possible values include: 'HTTP', 'HTTPS'
	Protocol ProtocolTypes `json:"protocol,omitempty"`
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
	// RDP
	RDP *RDPConfiguration `json:"rdp,omitempty"`
	// WinRM - Specifies the Windows Remote Management listeners. This enables remote Windows PowerShell.
	WinRM *WinRMConfiguration `json:"winRM,omitempty"`
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
	// CustomData Specifies a base-64 encoded string of custom data. The base-64 encoded string is decoded to a binary array that is saved as a file on the Virtual Machine. The maximum length of the binary array is 65535 bytes. <br><br> For using cloud-init for your VM, see [Using cloud-init to customize a Linux VM during creation](https://docs.microsoft.com/azure/virtual-machines/virtual-machines-linux-using-cloud-init?toc=%2fazure%2fvirtual-machines%2flinux%2ftoc.json)
	CustomData *string `json:"customData,omitempty"`
	// WindowsConfiguration
	WindowsConfiguration *WindowsConfiguration `json:"windowsConfiguration,omitempty"`
	// LinuxConfiguration
	LinuxConfiguration *LinuxConfiguration `json:"linuxConfiguration,omitempty"`
	// Bootstrap engine
	OsBootstrapEngine OperatingSystemBootstrapEngine `json:"osbootstrapengine,omitempty"`
}

type NetworkInterfaceReference struct {
	// VirtualNetworkReference
	VirtualNetworkReference *string `json:"virtualNetworkReference,omitempty"`
	// VirtualNetworkInterfaceReference
	VirtualNetworkInterfaceReference *string `json:"virtualNetworkInterfaceReference,omitempty"`
}
type NetworkProfile struct {
	// NetworkInterfaces
	NetworkInterfaces *[]NetworkInterfaceReference `json:"networkInterfaces,omitempty"`
}

type UefiSettings struct {
	// SecureBootEnabled - Specifies whether secure boot should be enabled on the virtual machine.
	SecureBootEnabled *bool `json:"secureBootEnabled,omitempty"`
}

// SecurityProfile
type SecurityProfile struct {
	EnableTPM *bool `json:"enableTPM,omitempty"`
	// Security related configuration used while creating the virtual machine.
	UefiSettings *UefiSettings `json:"uefiSettings,omitempty"`
}

type VirtualMachineProperties struct {
	// SecurityProfile - Specifies the security settings for the virtual machine.
	SecurityProfile *SecurityProfile `json:"securityProfile,omitempty"`
	// HardwareProfile - Specifies the hardware settings for the virtual machine.
	HardwareProfile *HardwareProfile `json:"hardwareProfile,omitempty"`
	// StorageProfile - Specifies the storage settings for the virtual machine disks.
	StorageProfile *StorageProfile `json:"storageProfile,omitempty"`
	// OsProfile
	OsProfile *OSProfile `json:"osProfile,omitempty"`
	// NetworkProfile
	NetworkProfile *NetworkProfile `json:"networkProfile,omitempty"`
	// ProvisioningState - READ-ONLY; The provisioning state, which only appears in the response.
	ProvisioningState *string `json:"provisioningState,omitempty"`
	// DisableHighAvailability
	DisableHighAvailability *bool `json:"disableHighAvailability,omitempty"`
	// State - State would container PowerState/ProvisioningState-SubState
	// https://docs.microsoft.com/en-us/azure/virtual-machines/windows/states-lifecycle
	Statuses map[string]*string `json:"statuses"`
	// IsPlaceholder - On a multi-node system, the entity (such as a VM) is created on a node where
	// IsPlacehoder is false. On all the other nodes, IsPlaceholder is set to true.
	// When an entity moves among these nodes (such as when a VM migrates), the
	// IsPlacehoder property is updated accordingly on all the nodes.
	// IsPlacehoder therefore defines where the entity (VM) is *not* located.
	// This property is the exact inverse of the node agent's SystemOwned property.
	IsPlaceholder *bool `json:"isPlaceholder,omitempty"`
	// HighAvailabilityState
	HighAvailabilityState *string `json:"HighAvailabilityState,omitempty"`
}

type VirtualMachine struct {
	// ID
	ID *string `json:"ID,omitempty"`
	// Name
	Name *string `json:"name,omitempty"`
	// Type
	Type *string `json:"type,omitempty"`
	// Tags - Custom resource tags
	Tags map[string]*string `json:"tags"`
	// Properties
	*VirtualMachineProperties `json:"properties,omitempty"`
}

type Sku struct {
	// Name
	Name *string `json:"name,omitempty"`
	// Capacity
	Capacity *int64 `json:"capacity,omitempty"`
}

type VirtualMachineScaleSetNetworkConfigurationProperties struct {
	// IPConfigurations
	IPConfigurations *[]network.IPConfiguration `json:"ipConfigurations,omitempty"`
	// DNS
	DNSSettings *network.DNSSetting `json:"dnsSettings,omitempty"`
	// EnableIPForwarding
	EnableIPForwarding *bool `json:"enableIPForwarding,omitempty"`
}

// VirtualNetwork defines the structure of a VNET
type VirtualMachineScaleSetNetworkConfiguration struct {
	// ID
	ID *string `json:"ID,omitempty"`
	// Name
	Name *string `json:"name,omitempty"`
	// Type
	Type *string `json:"type,omitempty"`
	// Tags - Custom resource tags
	Tags map[string]*string `json:"tags"`
	// Properties
	*VirtualMachineScaleSetNetworkConfigurationProperties `json:"properties,omitempty"`
}

type VirtualMachineScaleSetNetworkProfile struct {
	// NetworkInterfaceConfigurations
	NetworkInterfaceConfigurations *[]VirtualMachineScaleSetNetworkConfiguration `json:"networkInterfaceConfigurations,omitempty"`
}

// BootDiagnostics boot Diagnostics is a debugging feature which allows you to view Console Output and
// Screenshot to diagnose VM status. <br><br> You can easily view the output of your console log. <br><br>
// Azure also enables you to see a screenshot of the VM from the hypervisor.
type BootDiagnostics struct {
	// Enabled - Whether boot diagnostics should be enabled on the Virtual Machine.
	Enabled *bool `json:"enabled,omitempty"`
	// StorageURI - Uri of the storage account to use for placing the console output and screenshot.
	StorageURI *string `json:"storageUri,omitempty"`
}

type DiagnosticsProfile struct {
	// BootDiagnostics - Boot Diagnostics is a debugging feature which allows you to view Console Output and Screenshot to diagnose VM status. <br><br> You can easily view the output of your console log. <br><br> Azure also enables you to see a screenshot of the VM from the hypervisor.
	BootDiagnostics *BootDiagnostics `json:"bootDiagnostics,omitempty"`
}

// VirtualMachinePriorityTypes enumerates the values for virtual machine priority types.
type VirtualMachinePriorityTypes string

const (
	Low     VirtualMachinePriorityTypes = "Low"
	Regular VirtualMachinePriorityTypes = "Regular"
)

// VirtualMachineEvictionPolicyTypes enumerates the values for virtual machine eviction policy types.
type VirtualMachineEvictionPolicyTypes string

const (
	Deallocate VirtualMachineEvictionPolicyTypes = "Deallocate"
	Delete     VirtualMachineEvictionPolicyTypes = "Delete"
)

// VirtualMachineScaleSetVMProfileProperties
type VirtualMachineScaleSetVMProfileProperties struct {
	// SecurityProfile - Specifies the security settings for the virtual machine.
	SecurityProfile *SecurityProfile `json:"securityProfile,omitempty"`
	// HardwareProfile - Specifies the hardware settings for the virtual machine.
	HardwareProfile *HardwareProfile `json:"hardwareProfile,omitempty"`
	// StorageProfile - Specifies the storage settings for the virtual machine disks.
	StorageProfile *StorageProfile `json:"storageProfile,omitempty"`
	// OsProfile
	OsProfile *OSProfile `json:"osProfile,omitempty"`
	// NetworkProfile
	NetworkProfile *VirtualMachineScaleSetNetworkProfile `json:"networkProfile,omitempty"`
	// DiagnosticsProfile - Specifies the boot diagnostic settings state
	DiagnosticsProfile *DiagnosticsProfile `json:"diagnosticsProfile,omitempty"`
	// Priority - Specifies the priority for the virtual machines in the scale set. <br><br>Minimum api-version: 2017-10-30-preview. Possible values include: 'Regular', 'Low'
	Priority VirtualMachinePriorityTypes `json:"priority,omitempty"`
	// EvictionPolicy - Specifies the eviction policy for virtual machines in a low priority scale set. <br><br>Minimum api-version: 2017-10-30-preview. Possible values include: 'Deallocate', 'Delete'
	EvictionPolicy VirtualMachineEvictionPolicyTypes `json:"evictionPolicy,omitempty"`
}

// VirtualMachineScaleSetVMProfile
type VirtualMachineScaleSetVMProfile struct {
	// ID
	ID *string `json:"ID,omitempty"`
	// Name
	Name *string `json:"name,omitempty"`
	// Type
	Type *string `json:"type,omitempty"`
	// Tags - Custom resource tags
	Tags map[string]*string `json:"tags"`
	// Properties
	*VirtualMachineScaleSetVMProfileProperties `json:"properties,omitempty"`
}

// ResourceIdentityType enumerates the values for resource identity type.
type ResourceIdentityType string

const (
	// ResourceIdentityTypeNone ...
	ResourceIdentityTypeNone ResourceIdentityType = "None"
	// ResourceIdentityTypeSystemAssigned ...
	ResourceIdentityTypeSystemAssigned ResourceIdentityType = "SystemAssigned"
	// ResourceIdentityTypeSystemAssignedUserAssigned ...
	ResourceIdentityTypeSystemAssignedUserAssigned ResourceIdentityType = "SystemAssigned, UserAssigned"
	// ResourceIdentityTypeUserAssigned ...
	ResourceIdentityTypeUserAssigned ResourceIdentityType = "UserAssigned"
)

// VirtualMachineScaleSetIdentityUserAssignedIdentitiesValue ...
type VirtualMachineScaleSetIdentityUserAssignedIdentitiesValue struct {
	// PrincipalID - READ-ONLY; The principal id of user assigned identity.
	PrincipalID *string `json:"principalId,omitempty"`
	// ClientID - READ-ONLY; The client id of user assigned identity.
	ClientID *string `json:"clientId,omitempty"`
}

// VirtualMachineScaleSetIdentity identity for the virtual machine scale set.
type VirtualMachineScaleSetIdentity struct {
	// PrincipalID - READ-ONLY; The principal id of virtual machine scale set identity. This property will only be provided for a system assigned identity.
	PrincipalID *string `json:"principalId,omitempty"`
	// TenantID - READ-ONLY; The tenant id associated with the virtual machine scale set. This property will only be provided for a system assigned identity.
	TenantID *string `json:"tenantId,omitempty"`
	// Type - The type of identity used for the virtual machine scale set. The type 'SystemAssigned, UserAssigned' includes both an implicitly created identity and a set of user assigned identities. The type 'None' will remove any identities from the virtual machine scale set. Possible values include: 'ResourceIdentityTypeSystemAssigned', 'ResourceIdentityTypeUserAssigned', 'ResourceIdentityTypeSystemAssignedUserAssigned', 'ResourceIdentityTypeNone'
	Type ResourceIdentityType `json:"type,omitempty"`
	// UserAssignedIdentities - The list of user identities associated with the virtual machine scale set. The user identity dictionary key references will be ARM resource ids in the form: '/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ManagedIdentity/userAssignedIdentities/{identityName}'.
	UserAssignedIdentities map[string]*VirtualMachineScaleSetIdentityUserAssignedIdentitiesValue `json:"userAssignedIdentities"`
}

// VirtualMachineScaleSetProperties
type VirtualMachineScaleSetProperties struct {
	// VirtualMachineProfile
	VirtualMachineProfile *VirtualMachineScaleSetVMProfile `json:"virtualMachineProfile,omitempty"`
	// ProvisioningState - READ-ONLY; The provisioning state, which only appears in the response.
	ProvisioningState *string `json:"provisioningState,omitempty"`
	// State - State would container PowerState/ProvisioningState-SubState
	Statuses map[string]*string `json:"statuses"`
}

// VirtualMachineScaleSet
type VirtualMachineScaleSet struct {
	// ID
	ID *string `json:"ID,omitempty"`
	// Name
	Name *string `json:"name,omitempty"`
	// Type
	Type *string `json:"type,omitempty"`
	// Tags - Custom resource tags
	Tags map[string]*string `json:"tags"`
	// Sku
	Sku *Sku `json:"sku,omitempty"`
	// Identity - The identity of the virtual machine scale set, if configured.
	Identity *VirtualMachineScaleSetIdentity `json:"identity,omitempty"`
	// Properties
	*VirtualMachineScaleSetProperties `json:"properties,omitempty"`
	// DisableHighAvailability
	DisableHighAvailability *bool `json:"disableHighAvailability,omitempty"`
	// Statuses - Status
	Statuses map[string]*string `json:"statuses"`
	// IsPlaceholder - On a multi-node system, the entity (such as a VM) is created on a node where
	// IsPlacehoder is false. On all the other nodes, IsPlaceholder is set to true.
	// When an entity moves among these nodes (such as when a VM migrates), the
	// IsPlacehoder property is updated accordingly on all the nodes.
	// IsPlacehoder therefore defines where the entity (VM) is *not* located.
	// This property is the exact inverse of the node agent's SystemOwned property.
	IsPlaceholder *bool `json:"isPlaceholder,omitempty"`
	// HighAvailabilityState
	HighAvailabilityState *string `json:"HighAvailabilityState,omitempty"`
}

type ExecutionState string

const (
	// ExecutionStateFailed ...
	ExecutionStateFailed ExecutionState = "Failed"
	// ExecutionStateSucceeded ...
	ExecutionStateSucceeded ExecutionState = "Succeeded"
	// ExecutionStateUnknown ...
	ExecutionStateUnknown ExecutionState = "Unknown"
)

// VirtualMachineRunCommandScriptSource describes the script sources for run command.
type VirtualMachineRunCommandScriptSource struct {
	// Script - Specifies the script content to be executed on the VM.
	Script *string `json:"script,omitempty"`
	// ScriptURI - Specifies the script download location.
	ScriptURI *string `json:"scriptUri,omitempty"`
	// CommandID - Specifies a commandId of predefined built-in script.
	CommandID *string `json:"commandId,omitempty"`
}

// RunCommandInputParameter describes the properties of a run command parameter.
type RunCommandInputParameter struct {
	// Name - The run command parameter name.
	Name *string `json:"name,omitempty"`
	// Value - The run command parameter value.
	Value *string `json:"value,omitempty"`
}

// VirtualMachineRunCommandInstanceView the instance view of a virtual machine run command.
type VirtualMachineRunCommandInstanceView struct {
	// ExecutionState - Script execution status. Possible values include: 'ExecutionStateUnknown', 'ExecutionStateFailed', 'ExecutionStateSucceeded'
	ExecutionState ExecutionState `json:"executionState,omitempty"`
	// ExitCode - Exit code returned from script execution.
	ExitCode *int32 `json:"exitCode,omitempty"`
	// Output - Script output stream.
	Output *string `json:"output,omitempty"`
	// Error - Script error stream.
	Error *string `json:"error,omitempty"`
}

// VirtualMachineRunCommandRequest describes the properties of a Virtual Machine run command.
type VirtualMachineRunCommandRequest struct {
	// Source - The source of the run command script.
	Source *VirtualMachineRunCommandScriptSource `json:"source,omitempty"`
	// Parameters - The parameters used by the script.
	Parameters    *[]RunCommandInputParameter `json:"parameters,omitempty"`
	RunAsUser     *string                     `json:"runasuser,omitempty"`
	RunAsPassword *string                     `json:"runaspassword,omitempty"`
}

// VirtualMachineRunCommandResponse
type VirtualMachineRunCommandResponse struct {
	// InstanceView - The virtual machine run command instance view.
	InstanceView *VirtualMachineRunCommandInstanceView `json:"instanceView,omitempty"`
}
