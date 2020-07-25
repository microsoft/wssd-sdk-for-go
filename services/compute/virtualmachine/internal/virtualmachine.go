// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.
package internal

import (
	"github.com/microsoft/moc/pkg/status"
	"github.com/microsoft/wssd-sdk-for-go/services/compute"

	wssdcommonproto "github.com/microsoft/moc/rpc/common"
	wssdcompute "github.com/microsoft/moc/rpc/nodeagent/compute"
)

// Conversion functions from compute to wssdcompute
func (c *client) getWssdVirtualMachine(vm *compute.VirtualMachine) *wssdcompute.VirtualMachine {
	wssdvm := &wssdcompute.VirtualMachine{
		Name: *vm.Name,
	}

	if vm.VirtualMachineProperties == nil {
		return wssdvm
	}
	wssdvm.Hardware = c.getWssdVirtualMachineHardwareConfiguration(vm)
	wssdvm.Security = c.getWssdVirtualMachineSecurityConfiguration(vm)
	wssdvm.Storage = c.getWssdVirtualMachineStorageConfiguration(vm.StorageProfile)
	wssdvm.Os = c.getWssdVirtualMachineOSConfiguration(vm.OsProfile)
	wssdvm.Network = c.getWssdVirtualMachineNetworkConfiguration(vm.NetworkProfile)
	wssdvm.Entity = c.getWssdVirtualMachineEntity(vm)

	if vm.DisableHighAvailability != nil {
		wssdvm.DisableHighAvailability = *vm.DisableHighAvailability
	}

	return wssdvm

}

func (c *client) getWssdVirtualMachineEntity(vm *compute.VirtualMachine) *wssdcommonproto.Entity {
	isPlaceholder := false
	if vm.IsPlaceholder != nil {
		isPlaceholder = *vm.IsPlaceholder
	}

	return &wssdcommonproto.Entity{
		IsPlaceholder: isPlaceholder,
	}
}

func (c *client) getWssdVirtualMachineHardwareConfiguration(vm *compute.VirtualMachine) *wssdcompute.HardwareConfiguration {
	sizeType := wssdcommonproto.VirtualMachineSizeType_Default
	if vm.HardwareProfile != nil {
		sizeType = compute.GetWssdVirtualMachineSizeFromVirtualMachineSize(vm.HardwareProfile.VMSize)
	}
	return &wssdcompute.HardwareConfiguration{
		VMSize: sizeType,
	}
}

func (c *client) getWssdVirtualMachineSecurityConfiguration(vm *compute.VirtualMachine) *wssdcompute.SecurityConfiguration {
	enableTPM := false
	if vm.SecurityProfile != nil {
		enableTPM = *vm.SecurityProfile.EnableTPM
	}
	return &wssdcompute.SecurityConfiguration{
		EnableTPM: enableTPM,
	}
}

func (c *client) getWssdVirtualMachineStorageConfiguration(s *compute.StorageProfile) *wssdcompute.StorageConfiguration {
	vmConfigContainerName := ""
	if s.VmConfigContainerName != nil {
		vmConfigContainerName = *s.VmConfigContainerName
	}
	return &wssdcompute.StorageConfiguration{
		Osdisk:                c.getWssdVirtualMachineStorageConfigurationOsDisk(s.OsDisk),
		Datadisks:             c.getWssdVirtualMachineStorageConfigurationDataDisks(s.DataDisks),
		VmConfigContainerName: vmConfigContainerName,
	}
}

func (c *client) getWssdVirtualMachineStorageConfigurationOsDisk(s *compute.OSDisk) *wssdcompute.Disk {
	return &wssdcompute.Disk{
		Diskname: *s.VhdName,
	}
}

func (c *client) getWssdVirtualMachineStorageConfigurationDataDisks(s *[]compute.DataDisk) []*wssdcompute.Disk {
	datadisks := []*wssdcompute.Disk{}
	for _, d := range *s {
		datadisks = append(datadisks, &wssdcompute.Disk{Diskname: *d.VhdName})
	}

	return datadisks

}

func (c *client) getWssdVirtualMachineNetworkConfiguration(s *compute.NetworkProfile) *wssdcompute.NetworkConfiguration {
	nc := &wssdcompute.NetworkConfiguration{
		Interfaces: []*wssdcompute.NetworkInterface{},
	}
	if s.NetworkInterfaces == nil {
		return nc
	}
	for _, nic := range *s.NetworkInterfaces {
		if nic.VirtualNetworkInterfaceReference == nil {
			continue
		}
		nc.Interfaces = append(nc.Interfaces, &wssdcompute.NetworkInterface{NetworkInterfaceName: *nic.VirtualNetworkInterfaceReference})
	}

	return nc
}

func (c *client) getWssdVirtualMachineOSSSHPublicKeys(ssh *compute.SSHConfiguration) []*wssdcompute.SSHPublicKey {
	keys := []*wssdcompute.SSHPublicKey{}
	if ssh == nil {
		return keys
	}
	for _, key := range *ssh.PublicKeys {
		keys = append(keys, &wssdcompute.SSHPublicKey{Keydata: *key.KeyData})
	}
	return keys

}

func (c *client) getWssdVirtualMachineOSConfiguration(s *compute.OSProfile) *wssdcompute.OperatingSystemConfiguration {
	publickeys := []*wssdcompute.SSHPublicKey{}
	if s.LinuxConfiguration != nil {
		publickeys = c.getWssdVirtualMachineOSSSHPublicKeys(s.LinuxConfiguration.SSH)
	} else if s.WindowsConfiguration != nil {
		publickeys = c.getWssdVirtualMachineOSSSHPublicKeys(s.WindowsConfiguration.SSH)
	}

	adminuser := &wssdcompute.UserConfiguration{}
	if s.AdminUsername != nil {
		adminuser.Username = *s.AdminUsername
	}

	if s.AdminPassword != nil {
		adminuser.Password = *s.AdminPassword
	}

	osconfig := wssdcompute.OperatingSystemConfiguration{
		ComputerName:  *s.ComputerName,
		Administrator: adminuser,
		Users:         []*wssdcompute.UserConfiguration{},
		Publickeys:    publickeys,
		Ostype:        wssdcommonproto.OperatingSystemType_WINDOWS,
	}

	if s.LinuxConfiguration != nil {
		osconfig.Ostype = wssdcommonproto.OperatingSystemType_LINUX
	}

	if s.CustomData != nil {
		osconfig.CustomData = *s.CustomData
	}
	return &osconfig
}

// Conversion functions from wssdcompute to compute

func (c *client) getVirtualMachine(vm *wssdcompute.VirtualMachine) *compute.VirtualMachine {
	return &compute.VirtualMachine{
		Name: &vm.Name,
		ID:   &vm.Id,
		VirtualMachineProperties: &compute.VirtualMachineProperties{
			HardwareProfile:         c.getVirtualMachineHardwareProfile(vm),
			SecurityProfile:         c.getVirtualMachineSecurityProfile(vm),
			StorageProfile:          c.getVirtualMachineStorageProfile(vm.Storage),
			OsProfile:               c.getVirtualMachineOSProfile(vm.Os),
			NetworkProfile:          c.getVirtualMachineNetworkProfile(vm.Network),
			DisableHighAvailability: &vm.DisableHighAvailability,
			ProvisioningState:       status.GetProvisioningState(vm.Status.GetProvisioningStatus()),
			Statuses:                c.getVirtualMachineStatuses(vm),
			IsPlaceholder:           c.getVirtualMachineIsPlaceholder(vm),
			HighAvailabilityState:   c.getVirtualMachineScaleSetHighAvailabilityState(vm),
		},
	}
}

func (c *client) getVirtualMachinePowerState(status wssdcommonproto.PowerState) *string {
	stateString := status.String()
	return &stateString
}

func (c *client) getVirtualMachineStatuses(vm *wssdcompute.VirtualMachine) map[string]*string {
	statuses := status.GetStatuses(vm.GetStatus())
	statuses["PowerState"] = c.getVirtualMachinePowerState(vm.GetPowerState())
	return statuses
}

func (c *client) getVirtualMachineHardwareProfile(vm *wssdcompute.VirtualMachine) *compute.HardwareProfile {
	sizeType := compute.VirtualMachineSizeTypesDefault
	if vm.Hardware != nil {
		sizeType = compute.GetVirtualMachineSizeFromWssdVirtualMachineSize(vm.Hardware.VMSize)
	}
	return &compute.HardwareProfile{
		VMSize: sizeType,
	}
}

func (c *client) getVirtualMachineIsPlaceholder(vm *wssdcompute.VirtualMachine) *bool {
	isPlaceholder := false
	entity := vm.GetEntity()
	if entity != nil {
		isPlaceholder = entity.IsPlaceholder
	}
	return &isPlaceholder
}

func (c *client) getVirtualMachineScaleSetHighAvailabilityState(vm *wssdcompute.VirtualMachine) *string {
	haState := wssdcommonproto.HighAvailabilityState_UNKNOWN_HA_STATE
	if vm != nil {
		haState = vm.HighAvailabilityState
	}
	stateString := haState.String()
	return &stateString
}

func (c *client) getVirtualMachineSecurityProfile(vm *wssdcompute.VirtualMachine) *compute.SecurityProfile {
	enableTPM := false
	if vm.Security != nil {
		enableTPM = vm.Security.EnableTPM
	}
	return &compute.SecurityProfile{
		EnableTPM: &enableTPM,
	}
}

func (c *client) getVirtualMachineStorageProfile(s *wssdcompute.StorageConfiguration) *compute.StorageProfile {
	return &compute.StorageProfile{
		OsDisk:                c.getVirtualMachineStorageProfileOsDisk(s.Osdisk),
		DataDisks:             c.getVirtualMachineStorageProfileDataDisks(s.Datadisks),
		VmConfigContainerName: &s.VmConfigContainerName,
	}
}

func (c *client) getVirtualMachineStorageProfileOsDisk(d *wssdcompute.Disk) *compute.OSDisk {
	return &compute.OSDisk{
		VhdName: &d.Diskname,
	}
}

func (c *client) getVirtualMachineStorageProfileDataDisks(dd []*wssdcompute.Disk) *[]compute.DataDisk {
	cdd := []compute.DataDisk{}

	for _, i := range dd {
		cdd = append(cdd, compute.DataDisk{VhdName: &(i.Diskname)})
	}

	return &cdd

}

func (c *client) getVirtualMachineNetworkProfile(n *wssdcompute.NetworkConfiguration) *compute.NetworkProfile {
	np := &compute.NetworkProfile{
		NetworkInterfaces: &[]compute.NetworkInterfaceReference{},
	}

	for _, nic := range n.Interfaces {
		if nic == nil {
			continue
		}
		*np.NetworkInterfaces = append(*np.NetworkInterfaces, compute.NetworkInterfaceReference{VirtualNetworkInterfaceReference: &((*nic).NetworkInterfaceName)})
	}
	return np
}

func (c *client) getVirtualMachineOSProfile(o *wssdcompute.OperatingSystemConfiguration) *compute.OSProfile {
	return &compute.OSProfile{
		ComputerName: &o.ComputerName,
		// AdminUsername: &o.Administrator.Username,
		// AdminPassword: &o.Administrator.Password,
	}
}
