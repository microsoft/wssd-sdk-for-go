// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.
package internal

import (
	"github.com/google/go-cmp/cmp"
	"github.com/microsoft/moc/pkg/errors"
	"github.com/microsoft/moc/pkg/status"
	"github.com/microsoft/wssd-sdk-for-go/services/compute"

	wssdcommonproto "github.com/microsoft/moc/rpc/common"
	wssdcompute "github.com/microsoft/moc/rpc/nodeagent/compute"
	wssdnetwork "github.com/microsoft/moc/rpc/nodeagent/network"
	wssdstorage "github.com/microsoft/moc/rpc/nodeagent/storage"
	"github.com/microsoft/wssd-sdk-for-go/services/network"
)

// Conversion functions from compute to wssdcompute
func (c *client) getWssdVirtualMachine(vm *compute.VirtualMachine) (*wssdcompute.VirtualMachine, error) {
	if vm.Name == nil {
		return nil, errors.Wrapf(errors.InvalidInput, "Virtual Machine name is missing")
	}

	vmID := ""
	if vm.ID != nil {
		vmID = *vm.ID
	}

	wssdvm := &wssdcompute.VirtualMachine{
		Name: *vm.Name,
		Id:   vmID,
		Tags: getWssdTags(vm.Tags),
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
	guestAgentConfig, err := c.getWssdVirtualMachineGuestAgentConfiguration(vm.GuestAgentProfile)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to get GuestAgent Configuration")
	}

	zoneConfig, err := c.getWssdVirtualMachineZoneConfiguration(vm.ZoneConfiguration)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to get Cluster Configuration")
	}

	entity, err := c.getWssdVirtualMachineEntity(vm)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to get Entity")
	}

	wssdvm = &wssdcompute.VirtualMachine{
		Name:              *vm.Name,
		Id:                vmID,
		Tags:              getWssdTags(vm.Tags),
		Storage:           storageConfig,
		Hardware:          hardwareConfig,
		Security:          securityConfig,
		Os:                osconfig,
		Network:           networkConfig,
		GuestAgent:        guestAgentConfig,
		ZoneConfiguration: zoneConfig,
		Entity:            entity,
		Priority:          vm.Priority,
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
	var customSize *wssdcommonproto.VirtualMachineCustomSize
	var dynMemConfig *wssdcommonproto.DynamicMemoryConfiguration
	var vmGPUs []*wssdcommonproto.VirtualMachineGPU
	if vm.HardwareProfile != nil {
		sizeType = compute.GetWssdVirtualMachineSizeFromVirtualMachineSize(vm.HardwareProfile.VMSize)
		if vm.HardwareProfile.CustomSize != nil {
			customSize = &wssdcommonproto.VirtualMachineCustomSize{
				CpuCount: *vm.HardwareProfile.CustomSize.CpuCount,
				MemoryMB: *vm.HardwareProfile.CustomSize.MemoryMB,
			}
			if vm.HardwareProfile.CustomSize.GpuCount != nil {
				customSize.GpuCount = *vm.HardwareProfile.CustomSize.GpuCount
			}
		}
		if vm.HardwareProfile.DynamicMemoryConfig != nil {
			dynMemConfig = &wssdcommonproto.DynamicMemoryConfiguration{}
			if vm.HardwareProfile.DynamicMemoryConfig.MaximumMemoryMB != nil {
				dynMemConfig.MaximumMemoryMB = *vm.HardwareProfile.DynamicMemoryConfig.MaximumMemoryMB
			}
			if vm.HardwareProfile.DynamicMemoryConfig.MinimumMemoryMB != nil {
				dynMemConfig.MinimumMemoryMB = *vm.HardwareProfile.DynamicMemoryConfig.MinimumMemoryMB
			}
			if vm.HardwareProfile.DynamicMemoryConfig.TargetMemoryBuffer != nil {
				dynMemConfig.TargetMemoryBuffer = *vm.HardwareProfile.DynamicMemoryConfig.TargetMemoryBuffer
			}
		}
		if vm.HardwareProfile.VirtualMachineGPUs != nil {
			for _, gpu := range vm.HardwareProfile.VirtualMachineGPUs {
				if gpu == nil {
					return nil, errors.Wrapf(errors.InvalidInput, "nil value in Hardware.VirtualMachineGPUs")
				}
				if gpu.Assignment == nil {
					return nil, errors.Wrapf(errors.InvalidInput, "GPU assignment cannot be nil")
				}
				var assignment wssdcommonproto.AssignmentType
				switch *gpu.Assignment {
				case compute.GpuDDA:
					assignment = wssdcommonproto.AssignmentType_GpuDDA
				case compute.GpuP:
					assignment = wssdcommonproto.AssignmentType_GpuP
				case compute.GpuPV:
					assignment = wssdcommonproto.AssignmentType_GpuPV
				case compute.GpuDefault:
					assignment = wssdcommonproto.AssignmentType_GpuDefault
				}
				if gpu.PartitionSizeMB == nil {
					defaultInt := uint64(0)
					gpu.PartitionSizeMB = &defaultInt
				}
				if gpu.Name == nil {
					defaultString := ""
					gpu.Name = &defaultString
				}
				vmGPU := &wssdcommonproto.VirtualMachineGPU{
					Assignment:      assignment,
					PartitionSizeMB: *gpu.PartitionSizeMB,
					Name:            *gpu.Name,
				}
				vmGPUs = append(vmGPUs, vmGPU)
			}
		}
	}
	return &wssdcompute.HardwareConfiguration{
		VMSize:                     sizeType,
		CustomSize:                 customSize,
		DynamicMemoryConfiguration: dynMemConfig,
		VirtualMachineGPUs:         vmGPUs,
	}, nil
}

func (c *client) getWssdVirtualMachineSecurityConfiguration(vm *compute.VirtualMachine) (*wssdcompute.SecurityConfiguration, error) {
	enableTPM := false
	var uefiSettings *wssdcompute.UefiSettings
	uefiSettings = nil
	securityType := wssdcommonproto.SecurityType_NOTCONFIGURED
	if vm.SecurityProfile != nil {
		if vm.SecurityProfile.EnableTPM != nil {
			enableTPM = *vm.SecurityProfile.EnableTPM
		}
		if vm.SecurityProfile.UefiSettings != nil && vm.SecurityProfile.UefiSettings.SecureBootEnabled != nil {
			uefiSettings = &wssdcompute.UefiSettings{
				SecureBootEnabled: *vm.SecurityProfile.UefiSettings.SecureBootEnabled,
			}

		}
		switch vm.SecurityProfile.SecurityType {
		case compute.TrustedLaunch:
			securityType = wssdcommonproto.SecurityType_TRUSTEDLAUNCH
		case compute.ConfidentialVM:
			securityType = wssdcommonproto.SecurityType_CONFIDENTIALVM
		}
	}

	return &wssdcompute.SecurityConfiguration{
		EnableTPM:    enableTPM,
		UefiSettings: uefiSettings,
		SecurityType: securityType,
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
	if s.Vhd == nil || s.Vhd.Path == "" {
		return nil, errors.Wrapf(errors.InvalidInput, "Vhd Path is missing in OSDisk")
	}
	var managedDisk *wssdcommonproto.VirtualMachineManagedDiskParameters
	if s.ManagedDisk != nil {
		managedDisk = &wssdcommonproto.VirtualMachineManagedDiskParameters{}
		if s.ManagedDisk.SecurityProfile != nil {
			var securityEncryptionType wssdcommonproto.SecurityEncryptionTypes
			switch s.ManagedDisk.SecurityProfile.SecurityEncryptionType {
			case compute.NonPersistedTPM:
				securityEncryptionType = wssdcommonproto.SecurityEncryptionTypes_NonPersistedTPM
			default:
				securityEncryptionType = wssdcommonproto.SecurityEncryptionTypes_SecurityEncryptionNone
			}
			managedDisk.SecurityProfile = &wssdcommonproto.VMDiskSecurityProfile{
				SecurityEncryptionType: securityEncryptionType,
			}
		}
	}

	return &wssdcompute.Disk{
		Vhd: &wssdstorage.VirtualHardDisk{
			Name:                s.Vhd.Name,
			Id:                  s.Vhd.Id,
			Source:              s.Vhd.Source,
			Path:                s.Vhd.Path,
			ContainerName:       s.Vhd.ContainerName,
			Size:                s.Vhd.Size,
			Dynamic:             s.Vhd.Dynamic,
			Blocksizebytes:      s.Vhd.Blocksizebytes,
			Logicalsectorbytes:  s.Vhd.Logicalsectorbytes,
			Physicalsectorbytes: s.Vhd.Physicalsectorbytes,
			Controllernumber:    s.Vhd.Controllernumber,
			Controllerlocation:  s.Vhd.Controllerlocation,
			Disknumber:          s.Vhd.Disknumber,
			VirtualmachineName:  s.Vhd.VirtualmachineName,
			Scsipath:            s.Vhd.Scsipath,
			Virtualharddisktype: wssdstorage.VirtualHardDiskType(s.Vhd.Virtualharddisktype),
			SourceType:          s.Vhd.SourceType,
			HyperVGeneration:    s.Vhd.HyperVGeneration,
			CloudInitDataSource: s.Vhd.CloudInitDataSource,
			DiskFileFormat:      s.Vhd.DiskFileFormat,
			TargetUrl:           s.Vhd.TargetUrl,
			Vmid:                s.Vhd.Vmid,
		},
		ManagedDisk: managedDisk,
	}, nil
}

func (c *client) getWssdVirtualMachineStorageConfigurationDataDisks(s *[]compute.DataDisk) ([]*wssdcompute.Disk, error) {
	datadisks := []*wssdcompute.Disk{}
	for _, d := range *s {
		if d.Vhd == nil || d.Vhd.Path == "" {
			return nil, errors.Wrapf(errors.InvalidInput, "Vhd path is missing ")
		}
		datadisk := &wssdcompute.Disk{
			Vhd: &wssdstorage.VirtualHardDisk{
				Name:                d.Vhd.Name,
				Id:                  d.Vhd.Id,
				Source:              d.Vhd.Source,
				Path:                d.Vhd.Path,
				ContainerName:       d.Vhd.ContainerName,
				Size:                d.Vhd.Size,
				Dynamic:             d.Vhd.Dynamic,
				Blocksizebytes:      d.Vhd.Blocksizebytes,
				Logicalsectorbytes:  d.Vhd.Logicalsectorbytes,
				Physicalsectorbytes: d.Vhd.Physicalsectorbytes,
				Controllernumber:    d.Vhd.Controllernumber,
				Controllerlocation:  d.Vhd.Controllerlocation,
				Disknumber:          d.Vhd.Disknumber,
				VirtualmachineName:  d.Vhd.VirtualmachineName,
				Scsipath:            d.Vhd.Scsipath,
				Virtualharddisktype: wssdstorage.VirtualHardDiskType(d.Vhd.Virtualharddisktype),
				SourceType:          d.Vhd.SourceType,
				HyperVGeneration:    d.Vhd.HyperVGeneration,
				CloudInitDataSource: d.Vhd.CloudInitDataSource,
				DiskFileFormat:      d.Vhd.DiskFileFormat,
				TargetUrl:           d.Vhd.TargetUrl,
				Vmid:                d.Vhd.Vmid,
			},
		}
		datadisks = append(datadisks, datadisk)
	}

	return datadisks, nil

}

func (c *client) getWssdVirtualMachineNetworkConfiguration(s *compute.NetworkProfile) (*wssdcompute.NetworkConfiguration, error) {
	nc := &wssdcompute.NetworkConfiguration{
		Vnics: []*wssdnetwork.VirtualNetworkInterface{},
	}
	if s == nil || s.NetworkInterfaces == nil {
		return nc, nil
	}
	for _, nic := range *s.NetworkInterfaces {
		if nic.Vnic == nil {
			continue
		}
		var ipconfigs []*wssdnetwork.IpConfiguration
		for _, ipconf := range *nic.Vnic.IPConfigurations {
			subnetId := ""
			switchName := ""
			ipAddress := ""
			prefixLength := ""
			gateway := ""
			ipAlloactionMethod := wssdcommonproto.IPAllocationMethod_Invalid
			if ipconf.SubnetID != nil {
				subnetId = *ipconf.SubnetID
			}
			if ipconf.SwitchName != nil {
				switchName = *ipconf.SwitchName
			}
			if ipconf.IPAddress != nil {
				ipAddress = *ipconf.IPAddress
			}
			if ipconf.PrefixLength != nil {
				prefixLength = *ipconf.PrefixLength
			}
			if ipconf.Gateway != nil {
				gateway = *ipconf.Gateway
			}
			if ipconf.IPAllocationMethod != "" {
				ipAlloactionMethod = wssdcommonproto.IPAllocationMethod(wssdcommonproto.IPAllocationMethod_value[string(ipconf.IPAllocationMethod)])
			}

			ipconfigs = append(ipconfigs, &wssdnetwork.IpConfiguration{
				Subnetid:     subnetId,
				SwitchName:   switchName,
				Ipaddress:    ipAddress,
				Prefixlength: prefixLength,
				Gateway:      gateway,
				Allocation:   ipAlloactionMethod})
		}

		dnsSettings := &wssdcommonproto.Dns{}
		if nic.Vnic.DNSSettings != nil && nic.Vnic.DNSSettings.Servers != nil {
			dnsSettings.Servers = *nic.Vnic.DNSSettings.Servers
		}

		macAddress := ""
		vnicname := ""
		vnicid := ""
		vnicvmid := ""
		if nic.Vnic.Name != nil {
			vnicname = *nic.Vnic.Name
		}
		if nic.Vnic.ID != nil {
			vnicid = *nic.Vnic.ID
		}
		if nic.Vnic.VirtualMachineID != nil {
			vnicvmid = *nic.Vnic.VirtualMachineID
		}
		if nic.Vnic.MACAddress != nil {
			macAddress = *nic.Vnic.MACAddress
		}
		nc.Vnics = append(nc.Vnics, &wssdnetwork.VirtualNetworkInterface{
			Name:        vnicname,
			Id:          vnicid,
			Vmid:        vnicvmid,
			Ipconfigs:   ipconfigs,
			DnsSettings: dnsSettings,
			Macaddress:  macAddress,
		})
	}

	return nc, nil
}

func (c *client) getWssdVirtualMachineGuestAgentConfiguration(s *compute.GuestAgentProfile) (*wssdcommonproto.GuestAgentConfiguration, error) {
	gac := &wssdcommonproto.GuestAgentConfiguration{}

	if s == nil || s.Enabled == nil {
		return gac, nil
	}
	gac.Enabled = *s.Enabled

	return gac, nil
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

func (c *client) getWssdVirtualMachineWindowsConfiguration(windowsConfiguration *compute.WindowsConfiguration) *wssdcompute.WindowsConfiguration {
	wc := &wssdcompute.WindowsConfiguration{
		RDPConfiguration: &wssdcompute.RDPConfiguration{},
	}

	if windowsConfiguration == nil {
		return wc
	}

	if windowsConfiguration.WinRM != nil && windowsConfiguration.WinRM.Listeners != nil && len(*windowsConfiguration.WinRM.Listeners) >= 1 {
		listeners := make([]*wssdcommonproto.WinRMListener, len(*windowsConfiguration.WinRM.Listeners))
		for i, listener := range *windowsConfiguration.WinRM.Listeners {
			protocol := wssdcommonproto.WinRMProtocolType_HTTP
			if listener.Protocol == compute.HTTPS {
				protocol = wssdcommonproto.WinRMProtocolType_HTTPS
			}
			listeners[i] = &wssdcommonproto.WinRMListener{
				Protocol: protocol,
			}
		}
		wc.WinRMConfiguration = &wssdcommonproto.WinRMConfiguration{
			Listeners: listeners,
		}
	}

	if windowsConfiguration.RDP != nil {
		if windowsConfiguration.RDP.DisableRDP != nil {
			wc.RDPConfiguration.DisableRDP = *windowsConfiguration.RDP.DisableRDP
		}
		if windowsConfiguration.RDP.Port != nil {
			wc.RDPConfiguration.Port = uint32(*windowsConfiguration.RDP.Port)
		}
	}
	if windowsConfiguration.EnableAutomaticUpdates != nil {
		wc.EnableAutomaticUpdates = *windowsConfiguration.EnableAutomaticUpdates
	}

	if windowsConfiguration.TimeZone != nil {
		wc.TimeZone = *windowsConfiguration.TimeZone
	}

	return wc
}

func (c *client) getWssdVirtualMachineLinuxConfiguration(linuxConfiguration *compute.LinuxConfiguration) *wssdcompute.LinuxConfiguration {
	lc := &wssdcompute.LinuxConfiguration{}

	if linuxConfiguration.DisablePasswordAuthentication != nil {
		lc.DisablePasswordAuthentication = *linuxConfiguration.DisablePasswordAuthentication
	}

	lc.CloudInitDataSource = linuxConfiguration.CloudInitDataSource

	return lc

}

func (c *client) getWssdVirtualMachineOSConfiguration(s *compute.OSProfile) (*wssdcompute.OperatingSystemConfiguration, error) {
	publickeys := []*wssdcompute.SSHPublicKey{}
	osconfig := wssdcompute.OperatingSystemConfiguration{
		Users:             []*wssdcompute.UserConfiguration{},
		Ostype:            wssdcommonproto.OperatingSystemType_WINDOWS,
		OsBootstrapEngine: wssdcommonproto.OperatingSystemBootstrapEngine_CLOUD_INIT,
	}

	if s == nil || cmp.Equal(s, compute.OSProfile{}) {
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

	switch s.OsBootstrapEngine {
	case compute.WindowsAnswerFiles:
		osconfig.OsBootstrapEngine = wssdcommonproto.OperatingSystemBootstrapEngine_WINDOWS_ANSWER_FILES
	case compute.CloudInit:
		fallthrough
	default:
		osconfig.OsBootstrapEngine = wssdcommonproto.OperatingSystemBootstrapEngine_CLOUD_INIT
	}

	if s.WindowsConfiguration != nil {
		osconfig.WindowsConfiguration = c.getWssdVirtualMachineWindowsConfiguration(s.WindowsConfiguration)
	}

	if s.LinuxConfiguration != nil {
		osconfig.LinuxConfiguration = c.getWssdVirtualMachineLinuxConfiguration(s.LinuxConfiguration)
	}

	if s.LinuxConfiguration != nil {
		osconfig.Ostype = wssdcommonproto.OperatingSystemType_LINUX
	}

	if s.CustomData != nil {
		osconfig.CustomData = *s.CustomData
	}

	osconfig.ProxyConfiguration = c.getWssdVirtualMachineProxyConfiguration(s.ProxyConfiguration)

	return &osconfig, nil
}

func (c *client) getWssdVirtualMachineProxyConfiguration(proxyConfig *compute.ProxyConfiguration) *wssdcommonproto.ProxyConfiguration {
	if proxyConfig == nil {
		return nil
	}

	proxyConfiguration := &wssdcommonproto.ProxyConfiguration{}

	if proxyConfig.HttpProxy != nil {
		proxyConfiguration.HttpProxy = *proxyConfig.HttpProxy
	}

	if proxyConfig.HttpsProxy != nil {
		proxyConfiguration.HttpsProxy = *proxyConfig.HttpsProxy
	}

	if proxyConfig.NoProxy != nil {
		proxyConfiguration.NoProxy = *proxyConfig.NoProxy
	}

	if proxyConfig.TrustedCa != nil {
		proxyConfiguration.TrustedCa = *proxyConfig.TrustedCa
	}

	return proxyConfiguration
}

// Conversion functions from wssdcompute to compute

func (c *client) getVirtualMachine(vm *wssdcompute.VirtualMachine) *compute.VirtualMachine {
	if vm == nil || cmp.Equal(vm, wssdcompute.VirtualMachine{}) {
		return &compute.VirtualMachine{}
	}

	return &compute.VirtualMachine{
		Name: &vm.Name,
		ID:   &vm.Id,
		Tags: getComputeTags(vm.GetTags()),
		VirtualMachineProperties: &compute.VirtualMachineProperties{
			HardwareProfile:         c.getVirtualMachineHardwareProfile(vm),
			SecurityProfile:         c.getVirtualMachineSecurityProfile(vm),
			StorageProfile:          c.getVirtualMachineStorageProfile(vm.Storage),
			OsProfile:               c.getVirtualMachineOSProfile(vm.Os),
			NetworkProfile:          c.getVirtualMachineNetworkProfile(vm.Network),
			GuestAgentProfile:       c.getVirtualMachineGuestProfile(vm.GuestAgent),
			GuestAgentInstanceView:  c.getVirtualMachineGuestInstanceView(vm.GuestAgentInstanceView),
			DisableHighAvailability: &vm.DisableHighAvailability,
			ProvisioningState:       status.GetProvisioningState(vm.Status.GetProvisioningStatus()),
			ValidationStatus:        status.GetValidationStatus(vm.GetStatus()),
			Statuses:                c.getVirtualMachineStatuses(vm),
			IsPlaceholder:           c.getVirtualMachineIsPlaceholder(vm),
			HighAvailabilityState:   c.getVirtualMachineScaleSetHighAvailabilityState(vm),
			ZoneConfiguration:       c.getVirtualMachineZoneConfiguration(vm),
			Priority:                vm.Priority,
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
	var customSize *compute.VirtualMachineCustomSize
	var dynamicMemoryConfig *compute.DynamicMemoryConfiguration
	var vmGPUs []*compute.VirtualMachineGPU
	if vm.Hardware != nil {
		sizeType = compute.GetVirtualMachineSizeFromWssdVirtualMachineSize(vm.Hardware.VMSize)
		if vm.Hardware.CustomSize != nil {
			customSize = &compute.VirtualMachineCustomSize{
				CpuCount: &vm.Hardware.CustomSize.CpuCount,
				MemoryMB: &vm.Hardware.CustomSize.MemoryMB,
				GpuCount: &vm.Hardware.CustomSize.GpuCount,
			}
		}
		if vm.Hardware.DynamicMemoryConfiguration != nil {
			dynamicMemoryConfig = &compute.DynamicMemoryConfiguration{
				MaximumMemoryMB:    &vm.Hardware.DynamicMemoryConfiguration.MaximumMemoryMB,
				MinimumMemoryMB:    &vm.Hardware.DynamicMemoryConfiguration.MinimumMemoryMB,
				TargetMemoryBuffer: &vm.Hardware.DynamicMemoryConfiguration.TargetMemoryBuffer,
			}
		}
		if vm.Hardware.VirtualMachineGPUs != nil {
			for _, commonVMGPU := range vm.Hardware.VirtualMachineGPUs {
				var assignment compute.Assignment
				switch commonVMGPU.Assignment {
				case wssdcommonproto.AssignmentType_GpuDDA:
					assignment = compute.GpuDDA
				case wssdcommonproto.AssignmentType_GpuP:
					assignment = compute.GpuP
				case wssdcommonproto.AssignmentType_GpuPV:
					assignment = compute.GpuPV
				case wssdcommonproto.AssignmentType_GpuDefault:
					assignment = compute.GpuDefault
				}
				vmGPU := &compute.VirtualMachineGPU{
					Assignment:      &assignment,
					PartitionSizeMB: &commonVMGPU.PartitionSizeMB,
					Name:            &commonVMGPU.Name,
				}
				vmGPUs = append(vmGPUs, vmGPU)
			}
		}
	}
	return &compute.HardwareProfile{
		VMSize:              sizeType,
		CustomSize:          customSize,
		DynamicMemoryConfig: dynamicMemoryConfig,
		VirtualMachineGPUs:  vmGPUs,
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
	var uefiSettings *compute.UefiSettings
	uefiSettings = nil
	var securityType compute.SecurityTypes = ""

	if vm.Security != nil {
		enableTPM = vm.Security.EnableTPM
		if vm.Security.UefiSettings != nil {
			uefiSettings = &compute.UefiSettings{
				SecureBootEnabled: &vm.Security.UefiSettings.SecureBootEnabled,
			}
		}

		switch vm.Security.SecurityType {
		case wssdcommonproto.SecurityType_TRUSTEDLAUNCH:
			securityType = compute.TrustedLaunch
		case wssdcommonproto.SecurityType_CONFIDENTIALVM:
			securityType = compute.ConfidentialVM
		}
	}

	return &compute.SecurityProfile{
		EnableTPM:    &enableTPM,
		UefiSettings: uefiSettings,
		SecurityType: securityType,
	}

}

func (c *client) getVirtualMachineStorageProfile(s *wssdcompute.StorageConfiguration) *compute.StorageProfile {
	if s == nil || cmp.Equal(s, wssdcompute.StorageConfiguration{}) {
		return &compute.StorageProfile{}
	}
	return &compute.StorageProfile{
		OsDisk:                c.getVirtualMachineStorageProfileOsDisk(s.Osdisk),
		DataDisks:             c.getVirtualMachineStorageProfileDataDisks(s.Datadisks),
		VmConfigContainerName: &s.VmConfigContainerName,
	}
}

func (c *client) getVirtualMachineStorageProfileOsDisk(d *wssdcompute.Disk) *compute.OSDisk {
	if d == nil {
		return nil
	}
	var managedDisk *compute.VirtualMachineManagedDiskParameters
	if d.ManagedDisk != nil {
		managedDisk = &compute.VirtualMachineManagedDiskParameters{}
		if d.ManagedDisk.SecurityProfile != nil {
			var securityEncryptionType compute.SecurityEncryptionTypes
			switch d.ManagedDisk.SecurityProfile.SecurityEncryptionType {
			case wssdcommonproto.SecurityEncryptionTypes_NonPersistedTPM:
				securityEncryptionType = compute.NonPersistedTPM
			default:
				securityEncryptionType = ""
			}
			managedDisk.SecurityProfile = &compute.VMDiskSecurityProfile{
				SecurityEncryptionType: securityEncryptionType,
			}
		}
	}
	return &compute.OSDisk{
		Vhd: &compute.VirtualHardDisk{
			Name:                d.Vhd.Name,
			Id:                  d.Vhd.Id,
			Source:              d.Vhd.Source,
			Path:                d.Vhd.Path,
			ContainerName:       d.Vhd.ContainerName,
			Size:                d.Vhd.Size,
			Dynamic:             d.Vhd.Dynamic,
			Blocksizebytes:      d.Vhd.Blocksizebytes,
			Logicalsectorbytes:  d.Vhd.Logicalsectorbytes,
			Physicalsectorbytes: d.Vhd.Physicalsectorbytes,
			Controllernumber:    d.Vhd.Controllernumber,
			Controllerlocation:  d.Vhd.Controllerlocation,
			Disknumber:          d.Vhd.Disknumber,
			VirtualmachineName:  d.Vhd.VirtualmachineName,
			Scsipath:            d.Vhd.Scsipath,
			Virtualharddisktype: compute.VirtualHardDiskType(d.Vhd.Virtualharddisktype),
			SourceType:          d.Vhd.SourceType,
			HyperVGeneration:    d.Vhd.HyperVGeneration,
			CloudInitDataSource: d.Vhd.CloudInitDataSource,
			DiskFileFormat:      d.Vhd.DiskFileFormat,
			TargetUrl:           d.Vhd.TargetUrl,
			Vmid:                d.Vhd.Vmid,
		},
		ManagedDisk: managedDisk,
	}
}

func (c *client) getVirtualMachineStorageProfileDataDisks(dd []*wssdcompute.Disk) *[]compute.DataDisk {
	cdd := []compute.DataDisk{}

	for _, i := range dd {
		cdd = append(cdd, compute.DataDisk{
			Vhd: &compute.VirtualHardDisk{
				Name:                i.Vhd.Name,
				Id:                  i.Vhd.Id,
				Source:              i.Vhd.Source,
				Path:                i.Vhd.Path,
				ContainerName:       i.Vhd.ContainerName,
				Size:                i.Vhd.Size,
				Dynamic:             i.Vhd.Dynamic,
				Blocksizebytes:      i.Vhd.Blocksizebytes,
				Logicalsectorbytes:  i.Vhd.Logicalsectorbytes,
				Physicalsectorbytes: i.Vhd.Physicalsectorbytes,
				Controllernumber:    i.Vhd.Controllernumber,
				Controllerlocation:  i.Vhd.Controllerlocation,
				Disknumber:          i.Vhd.Disknumber,
				VirtualmachineName:  i.Vhd.VirtualmachineName,
				Scsipath:            i.Vhd.Scsipath,
				Virtualharddisktype: compute.VirtualHardDiskType(i.Vhd.Virtualharddisktype),
				SourceType:          i.Vhd.SourceType,
				HyperVGeneration:    i.Vhd.HyperVGeneration,
				CloudInitDataSource: i.Vhd.CloudInitDataSource,
				DiskFileFormat:      i.Vhd.DiskFileFormat,
				TargetUrl:           i.Vhd.TargetUrl,
				Vmid:                i.Vhd.Vmid,
			}})
	}

	return &cdd

}

func (c *client) getVirtualMachineNetworkProfile(n *wssdcompute.NetworkConfiguration) *compute.NetworkProfile {
	np := &compute.NetworkProfile{
		NetworkInterfaces: &[]compute.NetworkInterfaceReference{},
	}

	if n == nil || cmp.Equal(n, wssdcompute.NetworkConfiguration{}) {
		return np
	}

	for _, nic := range n.Vnics {
		if nic == nil {
			continue
		}
		var ipconfigs []network.IPConfiguration
		for _, ipconf := range nic.Ipconfigs {
			ipconfigs = append(ipconfigs, network.IPConfiguration{
				IPConfigurationProperties: &network.IPConfigurationProperties{SubnetID: &ipconf.Subnetid, SwitchName: &ipconf.SwitchName, IPAddress: &ipconf.Ipaddress},
			})
		}

		*np.NetworkInterfaces = append(*np.NetworkInterfaces, compute.NetworkInterfaceReference{
			Vnic: &network.VirtualNetworkInterface{
				Name: &nic.Name,
				ID:   &nic.Id,
				VirtualNetworkInterfaceProperties: &network.VirtualNetworkInterfaceProperties{
					VirtualMachineID: &nic.Vmid,
					IPConfigurations: &ipconfigs,
				},
			},
		})
	}
	return np
}

func (c *client) getVirtualMachineGuestProfile(g *wssdcommonproto.GuestAgentConfiguration) *compute.GuestAgentProfile {
	if g == nil || cmp.Equal(g, wssdcommonproto.GuestAgentConfiguration{}) {
		return nil
	}

	gap := &compute.GuestAgentProfile{
		Enabled: &g.Enabled,
	}

	return gap
}

func (c *client) getVirtualMachineGuestInstanceView(g *wssdcommonproto.VirtualMachineAgentInstanceView) *compute.GuestAgentInstanceView {
	if g == nil || cmp.Equal(g, wssdcommonproto.VirtualMachineAgentInstanceView{}) {
		return nil
	}

	gap := &compute.GuestAgentInstanceView{
		AgentVersion: g.GetVmAgentVersion(),
	}

	for _, status := range g.GetStatuses() {
		gap.Statuses = append(gap.Statuses, c.getInstanceViewStatus(status))
	}

	return gap
}

func (c *client) getVirtualMachineWindowsConfiguration(windowsConfiguration *wssdcompute.WindowsConfiguration) *compute.WindowsConfiguration {
	if windowsConfiguration == nil || cmp.Equal(windowsConfiguration, wssdcompute.WindowsConfiguration{}) {
		return nil
	}

	wc := &compute.WindowsConfiguration{
		RDP: &compute.RDPConfiguration{},
	}
	if windowsConfiguration.WinRMConfiguration != nil && len(windowsConfiguration.WinRMConfiguration.Listeners) >= 1 {
		listeners := make([]compute.WinRMListener, len(windowsConfiguration.WinRMConfiguration.Listeners))
		for i, listener := range windowsConfiguration.WinRMConfiguration.Listeners {
			protocol := compute.HTTP
			if listener.Protocol == wssdcommonproto.WinRMProtocolType_HTTPS {
				protocol = compute.HTTPS
			}
			listeners[i] = compute.WinRMListener{
				Protocol: protocol,
			}
		}
		wc.WinRM = &compute.WinRMConfiguration{
			Listeners: &listeners,
		}
	}

	if windowsConfiguration.RDPConfiguration != nil {
		wc.RDP.DisableRDP = &windowsConfiguration.RDPConfiguration.DisableRDP
		rdpPort := uint16(windowsConfiguration.RDPConfiguration.Port)
		wc.RDP.Port = &rdpPort
	}

	wc.EnableAutomaticUpdates = &windowsConfiguration.EnableAutomaticUpdates
	wc.TimeZone = &windowsConfiguration.TimeZone

	return wc
}

func (c *client) getVirtualMachineLinuxConfiguration(linuxConfiguration *wssdcompute.LinuxConfiguration) *compute.LinuxConfiguration {
	if linuxConfiguration == nil || cmp.Equal(linuxConfiguration, wssdcompute.LinuxConfiguration{}) {
		return nil
	}

	return &compute.LinuxConfiguration{
		DisablePasswordAuthentication: &linuxConfiguration.DisablePasswordAuthentication,
		CloudInitDataSource:           linuxConfiguration.CloudInitDataSource,
	}
}

func (c *client) getVirtualMachineOSProfile(o *wssdcompute.OperatingSystemConfiguration) *compute.OSProfile {
	if o == nil || cmp.Equal(o, wssdcompute.OperatingSystemConfiguration{}) {
		return nil
	}
	var osBootstrapEngine compute.OperatingSystemBootstrapEngine
	switch o.OsBootstrapEngine {
	case wssdcommonproto.OperatingSystemBootstrapEngine_WINDOWS_ANSWER_FILES:
		osBootstrapEngine = compute.WindowsAnswerFiles
	case wssdcommonproto.OperatingSystemBootstrapEngine_CLOUD_INIT:
		fallthrough
	default:
		osBootstrapEngine = compute.CloudInit
	}

	return &compute.OSProfile{
		ComputerName: &o.ComputerName,
		// AdminUsername: &o.Administrator.Username,
		// AdminPassword: &o.Administrator.Password,
		// Publickeys: &o.Publickeys,
		// Users : &o.Users,
		OsBootstrapEngine:    osBootstrapEngine,
		WindowsConfiguration: c.getVirtualMachineWindowsConfiguration(o.WindowsConfiguration),
		LinuxConfiguration:   c.getVirtualMachineLinuxConfiguration(o.LinuxConfiguration),
		ProxyConfiguration:   c.getVirtualMachineProxyConfiguration(o.ProxyConfiguration),
	}
}

func (c *client) getInstanceViewStatus(status *wssdcommonproto.InstanceViewStatus) *compute.InstanceViewStatus {
	level := compute.StatusLevelUnknown
	switch status.GetLevel() {
	case wssdcommonproto.InstanceViewStatus_Info:
		level = compute.StatusLevelInfo
	case wssdcommonproto.InstanceViewStatus_Warning:
		level = compute.StatusLevelWarning
	case wssdcommonproto.InstanceViewStatus_Error:
		level = compute.StatusLevelError
	}

	return &compute.InstanceViewStatus{
		Code:          status.GetCode(),
		Level:         level,
		DisplayStatus: status.GetDisplayStatus(),
		Message:       status.GetMessage(),
		Time:          status.GetTime(),
	}
}

func (c *client) getVirtualMachineProxyConfiguration(proxyConfiguration *wssdcommonproto.ProxyConfiguration) *compute.ProxyConfiguration {

	if proxyConfiguration == nil || cmp.Equal(proxyConfiguration, compute.ProxyConfiguration{}) {
		return nil
	}

	return &compute.ProxyConfiguration{
		HttpProxy:  &proxyConfiguration.HttpProxy,
		HttpsProxy: &proxyConfiguration.HttpsProxy,
		NoProxy:    &proxyConfiguration.NoProxy,
		TrustedCa:  &proxyConfiguration.TrustedCa,
	}
}

func (c *client) getWssdVirtualMachineZoneConfiguration(zoneProfile *compute.ZoneConfiguration) (*wssdcompute.ZoneConfiguration, error) {
	if zoneProfile == nil {
		return nil, nil
	}

	wssdZones := []*wssdcompute.ZoneReference{}
	for _, computeZone := range *zoneProfile.Zones {
		nodes := []string{}
		nodes = append(nodes, *computeZone.Nodes...)
		wssdZones = append(wssdZones, &wssdcompute.ZoneReference{
			Name:  *computeZone.Name,
			Nodes: nodes,
		})
	}
	strictPlacement := false
	if zoneProfile.StrictPlacement != nil {
		strictPlacement = *zoneProfile.StrictPlacement
	}
	wssdZoneConfiguration := &wssdcompute.ZoneConfiguration{
		Zones:           wssdZones,
		StrictPlacement: strictPlacement,
	}
	return wssdZoneConfiguration, nil
}

func (c *client) getVirtualMachineZoneConfiguration(vm *wssdcompute.VirtualMachine) *compute.ZoneConfiguration {
	zones := vm.GetZoneConfiguration().GetZones()
	if zones == nil || len(zones) == 0 {
		return nil
	}

	computeZones := []compute.ZoneReference{}

	for _, avZone := range zones {
		nodes := []string{}
		nodes = append(nodes, avZone.GetNodes()...)
		zoneName := avZone.GetName()
		computeZones = append(computeZones, compute.ZoneReference{
			Name:  &zoneName,
			Nodes: &nodes,
		})
	}

	strictPlacement := false
	if vm.GetZoneConfiguration().GetStrictPlacement() {
		strictPlacement = true
	}

	return &compute.ZoneConfiguration{
		Zones:           &computeZones,
		StrictPlacement: &strictPlacement,
	}
}
