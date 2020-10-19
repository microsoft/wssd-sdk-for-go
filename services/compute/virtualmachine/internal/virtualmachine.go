// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.
package internal

import (
	"github.com/microsoft/moc/pkg/errors"
	"github.com/microsoft/moc/pkg/status"
	"github.com/microsoft/wssd-sdk-for-go/services/compute"

	wssdcommonproto "github.com/microsoft/moc/rpc/common"
	wssdcompute "github.com/microsoft/moc/rpc/nodeagent/compute"
)

// Conversion functions from compute to wssdcompute
func (c *client) getWssdVirtualMachine(vm *compute.VirtualMachine) (*wssdcompute.VirtualMachine, error) {
	if vm.Name == nil {
		return nil, errors.Wrapf(errors.InvalidInput, "Virtual Machine name is missing")
	}

	wssdvm := &wssdcompute.VirtualMachine{
		Name: *vm.Name,
	}

	if vm.VirtualMachineProperties == nil {
		return wssdvm, nil
	}

	storageConfig, err := c.getWssdVirtualMachineStorageConfiguration(vm.StorageProfile)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to get Storage Configuration")
	}
	hardwareConfig, err := c.getWssdVirtualMachineHardwareConfiguration(vm)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to get Hardware Configuration")
	}
	securityConfig, err := c.getWssdVirtualMachineSecurityConfiguration(vm)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to get Security Configuration")
	}
	osconfig, err := c.getWssdVirtualMachineOSConfiguration(vm.OsProfile)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to get OS Configuration")
	}
	networkConfig, err := c.getWssdVirtualMachineNetworkConfiguration(vm.NetworkProfile)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to get Network Configuration")
	}
	entity, err := c.getWssdVirtualMachineEntity(vm)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to get Entity")
	}

	wssdvm = &wssdcompute.VirtualMachine{
		Name:     *vm.Name,
		Storage:  storageConfig,
		Hardware: hardwareConfig,
		Security: securityConfig,
		Os:       osconfig,
		Network:  networkConfig,
		Entity:   entity,
	}

	if vm.DisableHighAvailability != nil {
		wssdvm.DisableHighAvailability = *vm.DisableHighAvailability
	}

	return wssdvm, nil
}

func (c *client) getWssdVirtualMachineEntity(vm *compute.VirtualMachine) (*wssdcommonproto.Entity, error) {
	isPlaceholder := false
	if vm.IsPlaceholder != nil {
		isPlaceholder = *vm.IsPlaceholder
	}

	return &wssdcommonproto.Entity{
		IsPlaceholder: isPlaceholder,
	}, nil
}

func (c *client) getWssdVirtualMachineHardwareConfiguration(vm *compute.VirtualMachine) (*wssdcompute.HardwareConfiguration, error) {
	sizeType := wssdcommonproto.VirtualMachineSizeType_Default
	if vm.HardwareProfile != nil {
		sizeType = compute.GetWssdVirtualMachineSizeFromVirtualMachineSize(vm.HardwareProfile.VMSize)
	}
	return &wssdcompute.HardwareConfiguration{
		VMSize: sizeType,
	}, nil
}

func (c *client) getWssdVirtualMachineSecurityConfiguration(vm *compute.VirtualMachine) (*wssdcompute.SecurityConfiguration, error) {
	enableTPM := false
	if vm.SecurityProfile != nil {
		enableTPM = *vm.SecurityProfile.EnableTPM
	}
	return &wssdcompute.SecurityConfiguration{
		EnableTPM: enableTPM,
	}, nil
}

func (c *client) getWssdVirtualMachineStorageConfiguration(s *compute.StorageProfile) (*wssdcompute.StorageConfiguration, error) {
	wssdstorage := &wssdcompute.StorageConfiguration{
		Osdisk:    &wssdcompute.Disk{},
		Datadisks: []*wssdcompute.Disk{},
	}

	if s == nil {
		return wssdstorage, nil
	}

	vmConfigContainerName := ""
	if s.VmConfigContainerName != nil {
		vmConfigContainerName = *s.VmConfigContainerName
	}
	wssdstorage.VmConfigContainerName = vmConfigContainerName

	if s.OsDisk != nil {
		osdisk, err := c.getWssdVirtualMachineStorageConfigurationOsDisk(s.OsDisk)
		if err != nil {
			return nil, err
		}
		wssdstorage.Osdisk = osdisk
	}

	if s.DataDisks != nil {
		datadisks, err := c.getWssdVirtualMachineStorageConfigurationDataDisks(s.DataDisks)
		if err != nil {
			return nil, err
		}
		wssdstorage.Datadisks = datadisks
	}

	return wssdstorage, nil
}

func (c *client) getWssdVirtualMachineStorageConfigurationOsDisk(s *compute.OSDisk) (*wssdcompute.Disk, error) {
	if s.VhdName == nil {
		return nil, errors.Wrapf(errors.InvalidInput, "Vhd Name is missing in OSDisk")
	}
	return &wssdcompute.Disk{
		Diskname: *s.VhdName,
	}, nil
}

func (c *client) getWssdVirtualMachineStorageConfigurationDataDisks(s *[]compute.DataDisk) ([]*wssdcompute.Disk, error) {
	datadisks := []*wssdcompute.Disk{}
	for _, d := range *s {
		if d.VhdName == nil {
			return nil, errors.Wrapf(errors.InvalidInput, "Vhd Name is missing in DataDisk ")
		}
		datadisk := &wssdcompute.Disk{
			Diskname: *d.VhdName,
		}
		datadisks = append(datadisks, datadisk)
	}

	return datadisks, nil

}

func (c *client) getWssdVirtualMachineNetworkConfiguration(s *compute.NetworkProfile) (*wssdcompute.NetworkConfiguration, error) {
	nc := &wssdcompute.NetworkConfiguration{
		Interfaces: []*wssdcompute.NetworkInterface{},
	}
	if s == nil || s.NetworkInterfaces == nil {
		return nc, nil
	}
	for _, nic := range *s.NetworkInterfaces {
		if nic.VirtualNetworkInterfaceReference == nil {
			continue
		}
		nc.Interfaces = append(nc.Interfaces, &wssdcompute.NetworkInterface{NetworkInterfaceName: *nic.VirtualNetworkInterfaceReference})
	}

	return nc, nil
}

func (c *client) getWssdVirtualMachineOSSSHPublicKeys(ssh *compute.SSHConfiguration) ([]*wssdcompute.SSHPublicKey, error) {
	keys := []*wssdcompute.SSHPublicKey{}
	if ssh == nil {
		return keys, nil
	}
	for _, key := range *ssh.PublicKeys {
		if key.KeyData == nil {
			return nil, errors.Wrapf(errors.InvalidInput, "SSH KeyData is missing")
		}
		keys = append(keys, &wssdcompute.SSHPublicKey{Keydata: *key.KeyData})
	}
	return keys, nil

}

func (c *client) getWssdVirtualMachineOSConfiguration(s *compute.OSProfile) (*wssdcompute.OperatingSystemConfiguration, error) {
	publickeys := []*wssdcompute.SSHPublicKey{}
	osconfig := wssdcompute.OperatingSystemConfiguration{ // should Publickeys be here??
		Users:  []*wssdcompute.UserConfiguration{},
		Ostype: wssdcommonproto.OperatingSystemType_WINDOWS,
	}

	if s == nil {
		return &osconfig, nil
	}

	var err error

	if s.LinuxConfiguration != nil {
		publickeys, err = c.getWssdVirtualMachineOSSSHPublicKeys(s.LinuxConfiguration.SSH)
	} else if s.WindowsConfiguration != nil {
		publickeys, err = c.getWssdVirtualMachineOSSSHPublicKeys(s.WindowsConfiguration.SSH)
	}
	if err != nil {
		return nil, errors.Wrapf(err, "SSH Configuration Invalid")
	}

	adminuser := &wssdcompute.UserConfiguration{}
	if s.AdminUsername != nil {
		adminuser.Username = *s.AdminUsername
	}

	if s.AdminPassword != nil {
		adminuser.Password = *s.AdminPassword
	}

	if s.ComputerName == nil {
		return nil, errors.Wrapf(errors.InvalidInput, "ComputerName is missing")
	}

	osconfig.ComputerName = *s.ComputerName
	osconfig.Administrator = adminuser
	osconfig.Publickeys = publickeys

	if s.LinuxConfiguration != nil {
		osconfig.Ostype = wssdcommonproto.OperatingSystemType_LINUX
	}

	if s.CustomData != nil {
		osconfig.CustomData = *s.CustomData
	}
	return &osconfig, nil
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
