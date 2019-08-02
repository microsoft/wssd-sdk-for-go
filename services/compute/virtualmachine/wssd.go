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

package virtualmachine

import (
	"context"
	"fmt"
	"github.com/microsoft/wssd-sdk-for-go/services/compute"

	wssdclient "github.com/microsoft/wssdagent/rpc/client"
	wssdcompute "github.com/microsoft/wssdagent/rpc/compute"
	log "k8s.io/klog"
)

type client struct {
	wssdcompute.VirtualMachineAgentClient
}

// newVirtualMachineClient - creates a client session with the backend wssd agent
func newVirtualMachineClient(subID string) (*client, error) {
	c, err := wssdclient.GetVirtualMachineClient(&subID)
	if err != nil {
		return nil, err
	}
	return &client{c}, nil
}

// Get
func (c *client) Get(ctx context.Context, name string) (*[]compute.VirtualMachine, error) {
	request := getVirtualMachineRequest(wssdcompute.Operation_GET, name, nil)
	response, err := c.VirtualMachineAgentClient.Invoke(ctx, request)
	if err != nil {
		return nil, err
	}
	log.Infof("[VirtualMachine][Get] [%v]", response)
	return getVirtualMachineFromResponse(response), nil

}

// CreateOrUpdate
func (c *client) CreateOrUpdate(ctx context.Context, name string, id string, sg *compute.VirtualMachine) (*compute.VirtualMachine, error) {
	request := getVirtualMachineRequest(wssdcompute.Operation_POST, name, sg)
	response, err := c.VirtualMachineAgentClient.Invoke(ctx, request)
	if err != nil {
		return nil, err
	}
	log.Infof("[VirtualMachine][Create] [%v]", response)
	vms := getVirtualMachineFromResponse(response)
	if len(*vms) == 0 {
		return nil, fmt.Errorf("Creation of Virtual Machine failed to unknown reason.")
	}

	return &(*vms)[0], nil
}

// Delete methods invokes create or update on the client
func (c *client) Delete(ctx context.Context, name string, id string) error {
	vm, err := c.Get(ctx, name)
	if err != nil {
		return err
	}
	if len(*vm) == 0 {
		return fmt.Errorf("Virtual Machine [%s] not found", name)
	}

	request := getVirtualMachineRequest(wssdcompute.Operation_DELETE, name, &(*vm)[0])
	response, err := c.VirtualMachineAgentClient.Invoke(ctx, request)
	log.Infof("[VirtualMachine][Delete] [%v]", response)

	return err
}

func getVirtualMachineFromResponse(response *wssdcompute.VirtualMachineResponse) *[]compute.VirtualMachine {
	vms := []compute.VirtualMachine{}
	for _, vm := range response.GetVirtualMachineSystems() {
		vms = append(vms, *(GetVirtualMachine(vm)))
	}

	return &vms
}

func getVirtualMachineRequest(opType wssdcompute.Operation, name string, vmss *compute.VirtualMachine) *wssdcompute.VirtualMachineRequest {
	request := &wssdcompute.VirtualMachineRequest{
		OperationType:         opType,
		VirtualMachineSystems: []*wssdcompute.VirtualMachine{},
	}
	if vmss != nil {
		request.VirtualMachineSystems = append(request.VirtualMachineSystems, GetWssdVirtualMachine(vmss))
	} else if len(name) > 0 {
		request.VirtualMachineSystems = append(request.VirtualMachineSystems,
			&wssdcompute.VirtualMachine{
				Name: name,
			})
	}
	return request
}

// Conversion functions from compute to wssdcompute
func GetWssdVirtualMachine(c *compute.VirtualMachine) *wssdcompute.VirtualMachine {
	return &wssdcompute.VirtualMachine{
		Name:    *c.Name,
		Id:      *c.ID,
		Storage: getWssdVirtualMachineStorageConfiguration(c.StorageProfile),
		Os:      getWssdVirtualMachineOSConfiguration(c.OsProfile),
		Network: getWssdVirtualMachineNetworkConfiguration(c.NetworkProfile),
	}

}

func getWssdVirtualMachineStorageConfiguration(s *compute.StorageProfile) *wssdcompute.StorageConfiguration {
	return &wssdcompute.StorageConfiguration{
		Osdisk:    getWssdVirtualMachineStorageConfigurationOsDisk(s.OsDisk),
		Datadisks: getWssdVirtualMachineStorageConfigurationDataDisks(s.DataDisks),
	}
}

func getWssdVirtualMachineStorageConfigurationOsDisk(s *compute.OSDisk) *wssdcompute.Disk {
	return &wssdcompute.Disk{
		Diskid: *s.VhdId,
	}
}

func getWssdVirtualMachineStorageConfigurationDataDisks(s *[]compute.DataDisk) []*wssdcompute.Disk {
	datadisks := []*wssdcompute.Disk{}
	for _, d := range *s {
		datadisks = append(datadisks, &wssdcompute.Disk{Diskid: *d.VhdId})
	}

	return datadisks

}

func getWssdVirtualMachineNetworkConfiguration(s *compute.NetworkProfile) *wssdcompute.NetworkConfiguration {
	nc := &wssdcompute.NetworkConfiguration{
		Interfaces: []*wssdcompute.NetworkInterface{},
	}
	for _, nic := range *s.NetworkInterfaces {
		if nic.VirtualNetworkInterfaceID == nil {
			continue
		}
		nc.Interfaces = append(nc.Interfaces, &wssdcompute.NetworkInterface{NetworkInterfaceId: *nic.VirtualNetworkInterfaceID})
	}

	return nc
}

func getWssdVirtualMachineOSSSHPublicKeys(ssh *compute.SSHConfiguration) []*wssdcompute.SSHPublicKey {
	keys := []*wssdcompute.SSHPublicKey{}
	if ssh == nil {
		return keys
	}
	for _, key := range *ssh.PublicKeys {
		keys = append(keys, &wssdcompute.SSHPublicKey{Keydata: *key.KeyData})
	}
	return keys

}

func getWssdVirtualMachineOSConfiguration(s *compute.OSProfile) *wssdcompute.OperatingSystemConfiguration {
	publickeys := []*wssdcompute.SSHPublicKey{}
	if s.LinuxConfiguration != nil {
		publickeys = getWssdVirtualMachineOSSSHPublicKeys(s.LinuxConfiguration.SSH)
	}

	adminuser := &wssdcompute.UserConfiguration{}
	if s.AdminUsername != nil && s.AdminPassword != nil {
		adminuser.Username = *s.AdminUsername
		adminuser.Password = *s.AdminPassword
	}

	osconfig := wssdcompute.OperatingSystemConfiguration{
		ComputerName:  *s.ComputerName,
		Administrator: adminuser,
		Users:         []*wssdcompute.UserConfiguration{},
		Publickeys:    publickeys,
	}
	if s.CustomData != nil {
		osconfig.StartupScript = *s.CustomData
	}
	return &osconfig
}

// Conversion functions from wssdcompute to compute

func GetVirtualMachine(c *wssdcompute.VirtualMachine) *compute.VirtualMachine {
	return &compute.VirtualMachine{
		BaseProperties: compute.BaseProperties{
			Name: &c.Name,
			ID:   &c.Id,
		},
		StorageProfile: getVirtualMachineStorageProfile(c.Storage),
		OsProfile:      getVirtualMachineOSProfile(c.Os),
		NetworkProfile: getVirtualMachineNetworkProfile(c.Network),
	}
}

func getVirtualMachineStorageProfile(s *wssdcompute.StorageConfiguration) *compute.StorageProfile {
	return &compute.StorageProfile{
		OsDisk:    getVirtualMachineStorageProfileOsDisk(s.Osdisk),
		DataDisks: getVirtualMachineStorageProfileDataDisks(s.Datadisks),
	}
}

func getVirtualMachineStorageProfileOsDisk(d *wssdcompute.Disk) *compute.OSDisk {
	return &compute.OSDisk{
		VhdId: &d.Diskid,
	}
}

func getVirtualMachineStorageProfileDataDisks(dd []*wssdcompute.Disk) *[]compute.DataDisk {
	cdd := []compute.DataDisk{}

	for _, i := range dd {
		cdd = append(cdd, compute.DataDisk{VhdId: &(i.Diskid)})
	}

	return &cdd

}

func getVirtualMachineNetworkProfile(n *wssdcompute.NetworkConfiguration) *compute.NetworkProfile {
	np := &compute.NetworkProfile{
		NetworkInterfaces: &[]compute.NetworkInterfaceReference{},
	}

	for _, nic := range n.Interfaces {
		if nic == nil {
			continue
		}
		*np.NetworkInterfaces = append(*np.NetworkInterfaces, compute.NetworkInterfaceReference{VirtualNetworkInterfaceID: &((*nic).NetworkInterfaceId)})
	}
	return np
}

func getVirtualMachineOSProfile(o *wssdcompute.OperatingSystemConfiguration) *compute.OSProfile {
	return &compute.OSProfile{
		ComputerName: &o.ComputerName,
		// AdminUsername: &o.Administrator.Username,
		// AdminPassword: &o.Administrator.Password,
	}
}
