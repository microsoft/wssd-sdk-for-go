// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license

package internal

import (
	"context"
	"fmt"

	"github.com/microsoft/wssd-sdk-for-go/pkg/auth"
	wssdclient "github.com/microsoft/wssd-sdk-for-go/pkg/client"
	"github.com/microsoft/wssd-sdk-for-go/services/compute"
	"github.com/microsoft/wssd-sdk-for-go/services/compute/virtualmachine"
	"github.com/microsoft/wssd-sdk-for-go/services/network"
	"github.com/microsoft/wssdagent/pkg/errors"
	wssdcommonproto "github.com/microsoft/wssdagent/rpc/common"
	wssdcompute "github.com/microsoft/wssdagent/rpc/compute"
	wssdnetwork "github.com/microsoft/wssdagent/rpc/network"
)

type client struct {
	subID    string
	vmclient *virtualmachine.VirtualMachineClient
	wssdcompute.VirtualMachineScaleSetAgentClient
}

// NewVirtualMachineScaleSetClient - creates a client session with the backend wssd agent
func NewVirtualMachineScaleSetClient(subID string, authorizer auth.Authorizer) (*client, error) {
	c, err := wssdclient.GetVirtualMachineScaleSetClient(&subID, authorizer)
	if err != nil {
		return nil, err
	}
	vmc, err := virtualmachine.NewVirtualMachineClient(subID, authorizer)
	if err != nil {
		return nil, err
	}
	return &client{subID, vmc, c}, nil
}

// Get
func (c *client) Get(ctx context.Context, group, name string) (*[]compute.VirtualMachineScaleSet, error) {
	request, err := c.getVirtualMachineScaleSetRequest(wssdcommonproto.Operation_GET, name, nil)
	if err != nil {
		return nil, err
	}

	response, err := c.VirtualMachineScaleSetAgentClient.Invoke(ctx, request)
	if err != nil {
		return nil, err
	}
	return c.getVirtualMachineScaleSetFromResponse(response)
}

// GetVirtualMachines
func (c *client) GetVirtualMachines(ctx context.Context, group, name string) (*[]compute.VirtualMachine, error) {
	request, err := c.getVirtualMachineScaleSetRequest(wssdcommonproto.Operation_GET, name, nil)
	if err != nil {
		return nil, err
	}

	response, err := c.VirtualMachineScaleSetAgentClient.Invoke(ctx, request)
	if err != nil {
		return nil, err
	}

	vms := []compute.VirtualMachine{}
	for _, vmss := range response.GetVirtualMachineScaleSetSystems() {
		for _, vm := range vmss.GetVirtualMachineSystems() {
			tvms, err := c.vmclient.Get(ctx, group, vm.Name)
			if err != nil {
				return nil, err
			}
			if tvms == nil || len(*tvms) == 0 {
				return nil, fmt.Errorf("Vmss doesnt have any Vms")
			}
			// FIXME: Make sure Vms only on this scale set is returned.
			// If another Vm with the same name exists, that could also potentially be returned.
			vms = append(vms, (*tvms)[0])
		}
	}

	return &vms, nil
}

// CreateOrUpdate
func (c *client) CreateOrUpdate(ctx context.Context, group, name string, sg *compute.VirtualMachineScaleSet) (*compute.VirtualMachineScaleSet, error) {
	request, err := c.getVirtualMachineScaleSetRequest(wssdcommonproto.Operation_POST, name, sg)
	if err != nil {
		return nil, err
	}
	response, err := c.VirtualMachineScaleSetAgentClient.Invoke(ctx, request)
	if err != nil {
		return nil, err
	}
	vmsss, err := c.getVirtualMachineScaleSetFromResponse(response)
	if err != nil {
		return nil, err
	}
	return &((*vmsss)[0]), nil
}

// Delete methods invokes create or update on the client
func (c *client) Delete(ctx context.Context, group, name string) error {
	vmss, err := c.Get(ctx, group, name)
	if err != nil {
		return err
	}
	if len(*vmss) == 0 {
		return fmt.Errorf("Virtual Machine Scale Set [%s] not found", name)
	}

	request, err := c.getVirtualMachineScaleSetRequest(wssdcommonproto.Operation_DELETE, name, &(*vmss)[0])
	if err != nil {
		return err
	}
	_, err = c.VirtualMachineScaleSetAgentClient.Invoke(ctx, request)
	return err
}

///////// private methods ////////

// Conversion from proto to sdk
func (c *client) getVirtualMachineScaleSetFromResponse(response *wssdcompute.VirtualMachineScaleSetResponse) (*[]compute.VirtualMachineScaleSet, error) {
	vmsss := []compute.VirtualMachineScaleSet{}
	for _, vmss := range response.GetVirtualMachineScaleSetSystems() {
		cvmss, err := c.getVirtualMachineScaleSet(vmss)
		if err != nil {
			return nil, err
		}
		vmsss = append(vmsss, *cvmss)
	}

	return &vmsss, nil

}

func (c *client) getVirtualMachineScaleSetRequest(opType wssdcommonproto.Operation, name string, vmss *compute.VirtualMachineScaleSet) (*wssdcompute.VirtualMachineScaleSetRequest, error) {
	request := &wssdcompute.VirtualMachineScaleSetRequest{
		OperationType:                 opType,
		VirtualMachineScaleSetSystems: []*wssdcompute.VirtualMachineScaleSet{},
	}
	if vmss != nil {
		wssd_vmss, err := c.getWssdVirtualMachineScaleSet(vmss)
		if err != nil {
			return nil, err

		}
		request.VirtualMachineScaleSetSystems = append(request.VirtualMachineScaleSetSystems, wssd_vmss)
	} else if len(name) > 0 {
		request.VirtualMachineScaleSetSystems = append(request.VirtualMachineScaleSetSystems,
			&wssdcompute.VirtualMachineScaleSet{
				Name: name,
			})
	}

	return request, nil

}

func (c *client) getVirtualMachineScaleSet(vmss *wssdcompute.VirtualMachineScaleSet) (*compute.VirtualMachineScaleSet, error) {
	vmprofile, err := c.getVirtualMachineScaleSetVMProfile(vmss.Virtualmachineprofile)
	if err != nil {
		return nil, err
	}
	return &compute.VirtualMachineScaleSet{
		Name: &vmss.Name,
		ID:   &vmss.Id,
		Sku: &compute.Sku{
			Name:     &vmss.Sku.Name,
			Capacity: &vmss.Sku.Capacity,
		},
		VirtualMachineScaleSetProperties: &compute.VirtualMachineScaleSetProperties{
			VirtualMachineProfile: vmprofile,
			ProvisioningState:     c.getVirtualMachineScaleSetProvisioningState(vmss.Status.ProvisioningStatus),
		},
	}, nil
}

func (c *client) getVirtualMachineScaleSetProvisioningState(status *wssdcommonproto.ProvisionStatus) *string {
	provisionState := wssdcommonproto.ProvisionState_UNKNOWN
	if status != nil {
		provisionState = status.CurrentState
	}
	stateString := provisionState.String()
	return &stateString
}

func (c *client) getVirtualMachineScaleSetVMProfile(vm *wssdcompute.VirtualMachineProfile) (*compute.VirtualMachineScaleSetVMProfile, error) {
	net, err := c.getVirtualMachineScaleSetNetworkProfile(vm.Network)
	if err != nil {
		return nil, err
	}

	return &compute.VirtualMachineScaleSetVMProfile{
		Name: &vm.Vmprefix,
		VirtualMachineScaleSetVMProfileProperties: &compute.VirtualMachineScaleSetVMProfileProperties{
			HardwareProfile: c.getVirtualMachineScaleSetHardwareProfile(vm),
			StorageProfile:  c.getVirtualMachineScaleSetStorageProfile(vm.Storage),
			OsProfile:       c.getVirtualMachineScaleSetOSProfile(vm.Os),
			NetworkProfile:  net,
		},
	}, nil
}

func (c *client) getVirtualMachineScaleSetHardwareProfile(vm *wssdcompute.VirtualMachineProfile) *compute.HardwareProfile {
	sizeType := compute.VirtualMachineSizeTypesDefault
	if vm.Hardware != nil {
		sizeType = c.getVirtualMachineSize(vm.Hardware.VMSize)
	}
	return &compute.HardwareProfile{
		VMSize: sizeType,
	}
}

func (c *client) getVirtualMachineSize(size wssdcompute.VirtualMachineSizeType) compute.VirtualMachineSizeTypes {
	sizeType := compute.VirtualMachineSizeTypesDefault
	switch size {
	case wssdcompute.VirtualMachineSizeType_VirtualMachineSizeTypesStandardA1V2:
		sizeType = compute.VirtualMachineSizeTypesStandardA1V2
	case wssdcompute.VirtualMachineSizeType_VirtualMachineSizeTypesStandardA2V2:
		sizeType = compute.VirtualMachineSizeTypesStandardA2V2
	case wssdcompute.VirtualMachineSizeType_VirtualMachineSizeTypesStandardA4V2:
		sizeType = compute.VirtualMachineSizeTypesStandardA4V2
	}
	return sizeType
}

func (c *client) getVirtualMachineScaleSetStorageProfile(s *wssdcompute.StorageConfiguration) *compute.StorageProfile {
	return &compute.StorageProfile{
		OsDisk:    c.getVirtualMachineScaleSetStorageProfileOsDisk(s.Osdisk),
		DataDisks: c.getVirtualMachineScaleSetStorageProfileDataDisks(s.Datadisks),
	}
}

func (c *client) getVirtualMachineScaleSetStorageProfileOsDisk(d *wssdcompute.Disk) *compute.OSDisk {
	return &compute.OSDisk{
		VhdName: &d.Diskname,
	}
}

func (c *client) getVirtualMachineScaleSetStorageProfileDataDisks(dd []*wssdcompute.Disk) *[]compute.DataDisk {
	cdd := []compute.DataDisk{}

	for _, i := range dd {
		cdd = append(cdd, compute.DataDisk{VhdName: &(i.Diskname)})
	}

	return &cdd

}

func (c *client) getVirtualMachineScaleSetNetworkProfile(n *wssdcompute.NetworkConfigurationScaleSet) (*compute.VirtualMachineScaleSetNetworkProfile, error) {
	np := &compute.VirtualMachineScaleSetNetworkProfile{
		NetworkInterfaceConfigurations: &[]compute.VirtualMachineScaleSetNetworkConfiguration{},
	}

	for _, nic := range n.Interfaces {
		if nic == nil {
			continue
		}
		vnic, err := c.getVirtualMachineScaleSetNetworkConfiguration(nic)
		if err != nil {
			return nil, err
		}
		*np.NetworkInterfaceConfigurations = append(*np.NetworkInterfaceConfigurations, *vnic)
	}
	return np, nil
}

func (c *client) getVirtualMachineScaleSetNetworkConfiguration(nic *wssdnetwork.VirtualNetworkInterface) (*compute.VirtualMachineScaleSetNetworkConfiguration, error) {
	ipconfigs := []network.IPConfiguration{}
	for _, wssdipconfig := range nic.Ipconfigs {
		ipconfigs = append(ipconfigs, *(c.getVirtualMachineScaleSetNetworkConfigurationIPConfiguration(wssdipconfig)))
	}

	return &compute.VirtualMachineScaleSetNetworkConfiguration{
		VirtualMachineScaleSetNetworkConfigurationProperties: &compute.VirtualMachineScaleSetNetworkConfigurationProperties{
			IPConfigurations: &ipconfigs,
		},
	}, nil
}

func (c *client) getVirtualMachineScaleSetNetworkConfigurationIPConfiguration(wssdipconfig *wssdnetwork.IpConfiguration) *network.IPConfiguration {
	return &network.IPConfiguration{
		IPConfigurationProperties: &network.IPConfigurationProperties{
			SubnetID:     &wssdipconfig.Subnetid,
			PrefixLength: &wssdipconfig.Prefixlength,
			IPAddress:    &wssdipconfig.Ipaddress,
		},
	}
}

func (c *client) getVirtualMachineScaleSetOSProfile(o *wssdcompute.OperatingSystemConfiguration) *compute.OSProfile {
	return &compute.OSProfile{
		ComputerName: &o.ComputerName,
		// AdminUsername: &o.Administrator.Username,
		// AdminPassword: &o.Administrator.Password,
	}
}

// Conversion from sdk to protobuf
func (c *client) getWssdVirtualMachineScaleSet(vmss *compute.VirtualMachineScaleSet) (*wssdcompute.VirtualMachineScaleSet, error) {
	vm, err := c.getWssdVirtualMachineScaleSetVMProfile(vmss.VirtualMachineProfile)
	if err != nil {
		return nil, err
	}
	return &wssdcompute.VirtualMachineScaleSet{
		Name: *(vmss.Name),
		Sku: &wssdcompute.Sku{
			Name:     *(vmss.Sku.Name),
			Capacity: *(vmss.Sku.Capacity),
		},
		Virtualmachineprofile: vm,
	}, nil
}

func (c *client) getWssdVirtualMachineScaleSetVMProfile(vmp *compute.VirtualMachineScaleSetVMProfile) (*wssdcompute.VirtualMachineProfile, error) {
	net, err := c.getWssdVirtualMachineScaleSetNetworkConfiguration(vmp.NetworkProfile)
	if err != nil {
		return nil, err
	}
	return &wssdcompute.VirtualMachineProfile{
		Vmprefix: *vmp.Name,
		Hardware: c.getWssdVirtualMachineScaleSetHardwareConfiguration(vmp),
		Storage:  c.getWssdVirtualMachineScaleSetStorageConfiguration(vmp.StorageProfile),
		Os:       c.getWssdVirtualMachineScaleSetOSConfiguration(vmp.OsProfile),
		Network:  net,
	}, nil

}

func (c *client) getWssdVirtualMachineScaleSetHardwareConfiguration(vmp *compute.VirtualMachineScaleSetVMProfile) *wssdcompute.HardwareConfiguration {
	sizeType := wssdcompute.VirtualMachineSizeType_Default
	if vmp.HardwareProfile != nil {
		sizeType = c.getWssdVirtualMachineSize(vmp.HardwareProfile.VMSize)
	}
	return &wssdcompute.HardwareConfiguration{
		VMSize: sizeType,
	}
}

func (c *client) getWssdVirtualMachineSize(size compute.VirtualMachineSizeTypes) wssdcompute.VirtualMachineSizeType {
	sizeType := wssdcompute.VirtualMachineSizeType_Default
	switch size {
	case compute.VirtualMachineSizeTypesStandardA1V2:
		sizeType = wssdcompute.VirtualMachineSizeType_VirtualMachineSizeTypesStandardA1V2
	case compute.VirtualMachineSizeTypesStandardA2V2:
		sizeType = wssdcompute.VirtualMachineSizeType_VirtualMachineSizeTypesStandardA2V2
	case compute.VirtualMachineSizeTypesStandardA4V2:
		sizeType = wssdcompute.VirtualMachineSizeType_VirtualMachineSizeTypesStandardA4V2
	}
	return sizeType
}

func (c *client) getWssdVirtualMachineScaleSetStorageConfiguration(s *compute.StorageProfile) *wssdcompute.StorageConfiguration {
	return &wssdcompute.StorageConfiguration{
		Osdisk:    c.getWssdVirtualMachineScaleSetStorageConfigurationOsDisk(s.OsDisk),
		Datadisks: c.getWssdVirtualMachineScaleSetStorageConfigurationDataDisks(s.DataDisks),
	}
}

func (c *client) getWssdVirtualMachineScaleSetStorageConfigurationOsDisk(s *compute.OSDisk) *wssdcompute.Disk {
	return &wssdcompute.Disk{
		Diskname: *s.VhdName,
	}
}

func (c *client) getWssdVirtualMachineScaleSetStorageConfigurationDataDisks(s *[]compute.DataDisk) []*wssdcompute.Disk {
	datadisks := []*wssdcompute.Disk{}
	if s == nil {
		return datadisks
	}
	for _, d := range *s {
		datadisks = append(datadisks, &wssdcompute.Disk{Diskname: *d.VhdName})
	}

	return datadisks

}

func (c *client) getWssdVirtualMachineScaleSetNetworkConfiguration(s *compute.VirtualMachineScaleSetNetworkProfile) (*wssdcompute.NetworkConfigurationScaleSet, error) {
	nc := &wssdcompute.NetworkConfigurationScaleSet{
		Interfaces: []*wssdnetwork.VirtualNetworkInterface{},
	}
	if s == nil || s.NetworkInterfaceConfigurations == nil {
		return nc, nil
	}
	for _, nic := range *s.NetworkInterfaceConfigurations {
		vnic, err := c.getWssdVirtualMachineScaleSetNetworkConfigurationNetworkInterface(&nic)
		if err != nil {
			return nil, err
		}
		nc.Interfaces = append(nc.Interfaces, vnic)
	}

	return nc, nil
}

func (c *client) getWssdVirtualMachineScaleSetNetworkConfigurationNetworkInterface(nic *compute.VirtualMachineScaleSetNetworkConfiguration) (*wssdnetwork.VirtualNetworkInterface, error) {
	nicName := ""
	if nic.Name != nil {
		nicName = *nic.Name
	}
	wssdvnic := &wssdnetwork.VirtualNetworkInterface{
		Name: nicName,
	}
	if nic.VirtualMachineScaleSetNetworkConfigurationProperties == nil ||
		nic.IPConfigurations == nil ||
		len(*nic.IPConfigurations) == 0 {
		return nil, errors.Wrapf(errors.InvalidConfiguration, "Missing IPConfiguration Properties")
	}

	for _, ipconfig := range *nic.IPConfigurations {
		wssdipconfig, err := c.getWssdVirtualMachineScaleSetNetworkConfigurationNetworkInterfaceIPConfiguration(&ipconfig)
		if err != nil {
			return nil, err
		}
		wssdvnic.Ipconfigs = append(wssdvnic.Ipconfigs, wssdipconfig)
	}

	return wssdvnic, nil
}

func (c *client) getWssdVirtualMachineScaleSetNetworkConfigurationNetworkInterfaceIPConfiguration(ipconfig *network.IPConfiguration) (*wssdnetwork.IpConfiguration, error) {
	if ipconfig.IPConfigurationProperties == nil {
		return nil, errors.Wrapf(errors.InvalidConfiguration, "Missing IPConfiguration Properties")
	}
	if ipconfig.SubnetID == nil {
		return nil, errors.Wrapf(errors.InvalidConfiguration, "Missing Subnet Reference in IPConfiguration Properties")
	}

	wssdipconfig := &wssdnetwork.IpConfiguration{
		Subnetid: *ipconfig.SubnetID,
	}

	if ipconfig.IPAddress != nil {
		wssdipconfig.Ipaddress = *ipconfig.IPAddress
	}
	if ipconfig.PrefixLength != nil {
		wssdipconfig.Prefixlength = *ipconfig.PrefixLength
	}

	return wssdipconfig, nil
}

func (c *client) getWssdVirtualMachineScaleSetOSSSHPublicKeys(ssh *compute.SSHConfiguration) []*wssdcompute.SSHPublicKey {
	keys := []*wssdcompute.SSHPublicKey{}
	if ssh == nil {
		return keys
	}
	for _, key := range *ssh.PublicKeys {
		keys = append(keys, &wssdcompute.SSHPublicKey{Keydata: *key.KeyData})
	}
	return keys

}

func (c *client) getWssdVirtualMachineScaleSetOSConfiguration(s *compute.OSProfile) *wssdcompute.OperatingSystemConfiguration {
	publickeys := []*wssdcompute.SSHPublicKey{}
	if s.LinuxConfiguration != nil {
		publickeys = c.getWssdVirtualMachineScaleSetOSSSHPublicKeys(s.LinuxConfiguration.SSH)
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
		Ostype:        wssdcompute.OperatingSystemType_WINDOWS,
	}

	if s.LinuxConfiguration != nil {
		osconfig.Ostype = wssdcompute.OperatingSystemType_LINUX
	}

	if s.CustomData != nil {
		osconfig.CustomData = *s.CustomData
	}
	return &osconfig
}
