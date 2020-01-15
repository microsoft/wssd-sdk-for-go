// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package internal

import (
	"context"
	"fmt"
	"github.com/microsoft/wssd-sdk-for-go/pkg/auth"
	"github.com/microsoft/wssd-sdk-for-go/services/compute"

	wssdclient "github.com/microsoft/wssd-sdk-for-go/pkg/client"
	wssdcommonproto "github.com/microsoft/wssdagent/rpc/common"
	wssdcompute "github.com/microsoft/wssdagent/rpc/compute"
)

type client struct {
	wssdcompute.VirtualMachineAgentClient
}

// newVirtualMachineClient - creates a client session with the backend wssd agent
func NewVirtualMachineClient(subID string, authorizer auth.Authorizer) (*client, error) {
	c, err := wssdclient.GetVirtualMachineClient(&subID, authorizer)
	if err != nil {
		return nil, err
	}
	return &client{c}, nil
}

// Get
func (c *client) Get(ctx context.Context, group, name string) (*[]compute.VirtualMachine, error) {
	request := c.getVirtualMachineRequest(wssdcommonproto.Operation_GET, name, nil)
	response, err := c.VirtualMachineAgentClient.Invoke(ctx, request)
	if err != nil {
		return nil, err
	}
	return c.getVirtualMachineFromResponse(response), nil

}

// CreateOrUpdate
func (c *client) CreateOrUpdate(ctx context.Context, group, name string, sg *compute.VirtualMachine) (*compute.VirtualMachine, error) {
	request := c.getVirtualMachineRequest(wssdcommonproto.Operation_POST, name, sg)
	response, err := c.VirtualMachineAgentClient.Invoke(ctx, request)
	if err != nil {
		return nil, err
	}
	vms := c.getVirtualMachineFromResponse(response)
	if len(*vms) == 0 {
		return nil, fmt.Errorf("Creation of Virtual Machine failed to unknown reason.")
	}

	return &(*vms)[0], nil
}

// Delete methods invokes create or update on the client
func (c *client) Delete(ctx context.Context, group, name string) error {
	vm, err := c.Get(ctx, group, name)
	if err != nil {
		return err
	}
	if len(*vm) == 0 {
		return fmt.Errorf("Virtual Machine [%s] not found", name)
	}

	request := c.getVirtualMachineRequest(wssdcommonproto.Operation_DELETE, name, &(*vm)[0])
	_, err = c.VirtualMachineAgentClient.Invoke(ctx, request)

	return err
}

func (c *client) getVirtualMachineFromResponse(response *wssdcompute.VirtualMachineResponse) *[]compute.VirtualMachine {
	vms := []compute.VirtualMachine{}
	for _, vm := range response.GetVirtualMachineSystems() {
		vms = append(vms, *(c.getVirtualMachine(vm)))
	}

	return &vms
}

func (c *client) getVirtualMachineRequest(opType wssdcommonproto.Operation, name string, vmss *compute.VirtualMachine) *wssdcompute.VirtualMachineRequest {
	request := &wssdcompute.VirtualMachineRequest{
		OperationType:         opType,
		VirtualMachineSystems: []*wssdcompute.VirtualMachine{},
	}
	if vmss != nil {
		request.VirtualMachineSystems = append(request.VirtualMachineSystems, c.getWssdVirtualMachine(vmss))
	} else if len(name) > 0 {
		request.VirtualMachineSystems = append(request.VirtualMachineSystems,
			&wssdcompute.VirtualMachine{
				Name: name,
			})
	}
	return request
}

// Conversion functions from compute to wssdcompute
func (c *client) getWssdVirtualMachine(vm *compute.VirtualMachine) *wssdcompute.VirtualMachine {
	wssdvm := &wssdcompute.VirtualMachine{
		Name: *vm.Name,
	}

	if vm.VirtualMachineProperties == nil {
		return wssdvm
	}
	wssdvm.Storage = c.getWssdVirtualMachineStorageConfiguration(vm.StorageProfile)
	wssdvm.Os = c.getWssdVirtualMachineOSConfiguration(vm.OsProfile)
	wssdvm.Network = c.getWssdVirtualMachineNetworkConfiguration(vm.NetworkProfile)
	return wssdvm

}

func (c *client) getWssdVirtualMachineStorageConfiguration(s *compute.StorageProfile) *wssdcompute.StorageConfiguration {
	return &wssdcompute.StorageConfiguration{
		Osdisk:    c.getWssdVirtualMachineStorageConfigurationOsDisk(s.OsDisk),
		Datadisks: c.getWssdVirtualMachineStorageConfigurationDataDisks(s.DataDisks),
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

// Conversion functions from wssdcompute to compute

func (c *client) getVirtualMachine(vm *wssdcompute.VirtualMachine) *compute.VirtualMachine {
	return &compute.VirtualMachine{
		Name: &vm.Name,
		ID:   &vm.Id,
		VirtualMachineProperties: &compute.VirtualMachineProperties{
			StorageProfile:    c.getVirtualMachineStorageProfile(vm.Storage),
			OsProfile:         c.getVirtualMachineOSProfile(vm.Os),
			NetworkProfile:    c.getVirtualMachineNetworkProfile(vm.Network),
			ProvisioningState: c.getVirtualMachineProvisioningState(vm.ProvisionStatus),
		},
	}
}

func (c *client) getVirtualMachineProvisioningState(status *wssdcommonproto.ProvisionStatus) *string {
	provisionState := wssdcommonproto.ProvisionState_UNKNOWN
	if status != nil {
		provisionState = status.CurrentState
	}
	stateString := provisionState.String()
	return &stateString
}

func (c *client) getVirtualMachineStorageProfile(s *wssdcompute.StorageConfiguration) *compute.StorageProfile {
	return &compute.StorageProfile{
		OsDisk:    c.getVirtualMachineStorageProfileOsDisk(s.Osdisk),
		DataDisks: c.getVirtualMachineStorageProfileDataDisks(s.Datadisks),
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
