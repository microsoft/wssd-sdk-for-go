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
	"errors"
	"github.com/microsoft/wssd-sdk-for-go/services/compute"
	"github.com/microsoft/wssd-sdk-for-go/services/network"

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

// GetAll
func (c *client) List(ctx context.Context) (*[]compute.VirtualMachine, error) {
	request := &wssdcompute.VirtualMachineRequest{
		OperationType: wssdcompute.Operation_GET,
	}

	response, err := c.VirtualMachineAgentClient.Invoke(ctx, request)

	if err != nil {
		return nil, err
	}

	// if response.GetResult().Value == false {
	//	return nil, errors.New(response.GetError())
	// }

	if len(response.GetVirtualMachineSystems()) == 0 {
		return nil, nil
	}
	log.Infof("[VirtualMachine][List] [%v]", response)

	vms := []compute.VirtualMachine{}
	for _, v := range response.GetVirtualMachineSystems() {
		vms = append(vms, *getVirtualMachine(v))
	}

	// Pick the first virtual machine returned
	return &vms, nil
}

// Get
func (c *client) Get(ctx context.Context, name string) (*compute.VirtualMachine, error) {
	request := &wssdcompute.VirtualMachineRequest{
		OperationType:         wssdcompute.Operation_GET,
		VirtualMachineSystems: []*wssdcompute.VirtualMachine{},
	}
	vm := &wssdcompute.VirtualMachine{
		Name: name,
	}
	request.VirtualMachineSystems = append(request.VirtualMachineSystems, vm)
	response, err := c.VirtualMachineAgentClient.Invoke(ctx, request)

	if err != nil {
		return nil, err
	}

	if len(response.GetVirtualMachineSystems()) == 0 {
		return nil, errors.New(response.GetError())
	}
	log.Infof("[VirtualMachine][Get] [%v]", response)

	// Pick the first virtual machine returned
	return getVirtualMachine(response.GetVirtualMachineSystems()[0]), nil
}

// CreateOrUpdate
func (c *client) CreateOrUpdate(ctx context.Context, name string, id string, sg *compute.VirtualMachine) (*compute.VirtualMachine, error) {
	request := &wssdcompute.VirtualMachineRequest{
		OperationType:         wssdcompute.Operation_POST,
		VirtualMachineSystems: make([]*wssdcompute.VirtualMachine, 0),
	}
	request.VirtualMachineSystems = append(request.VirtualMachineSystems, getWssdVirtualMachine(sg))
	response, err := c.VirtualMachineAgentClient.Invoke(ctx, request)
	if err != nil {
		return nil, err
	}

	if len(response.GetVirtualMachineSystems()) == 0 {
		return nil, errors.New(response.GetError())
	}

	// Pick the first virtual machine returned
	return getVirtualMachine(response.GetVirtualMachineSystems()[0]), nil
}

// Delete methods invokes create or update on the client
func (c *client) Delete(ctx context.Context, name string, id string) error {
	request := &wssdcompute.VirtualMachineRequest{
		OperationType:         wssdcompute.Operation_DELETE,
		VirtualMachineSystems: []*wssdcompute.VirtualMachine{},
	}
	vm, err := c.Get(ctx, name)
	if err != nil {
		return err
	}

	request.VirtualMachineSystems = append(request.VirtualMachineSystems, getWssdVirtualMachine(vm))
	_, err = c.VirtualMachineAgentClient.Invoke(ctx, request)

	if err != nil {
		return err
	}

	return nil
}

// Conversion functions from compute to wssdcompute

func getWssdVirtualMachine(c *compute.VirtualMachine) *wssdcompute.VirtualMachine {
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
	for _, _ = range *s.NetworkInterfaceConfigurations {
		nc.Interfaces = append(nc.Interfaces, &wssdcompute.NetworkInterface{}) // FixMe
	}

	return nc
}

func getWssdVirtualMachineOSConfiguration(s *compute.OSProfile) *wssdcompute.OperatingSystemConfiguration {
	return &wssdcompute.OperatingSystemConfiguration{
		ComputerName:  *s.ComputerName,
		Administrator: &wssdcompute.UserConfiguration{},
		Users:         []*wssdcompute.UserConfiguration{},
		Publickeys:    []*wssdcompute.SSHPublicKey{},
	}
}

// Conversion functions from wssdcompute to compute

func getVirtualMachine(c *wssdcompute.VirtualMachine) *compute.VirtualMachine {
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
		NetworkInterfaceConfigurations: &[]network.VirtualNetworkInterface{},
	}

	for _, _ = range n.Interfaces {
		*np.NetworkInterfaceConfigurations = append(*np.NetworkInterfaceConfigurations, network.VirtualNetworkInterface{}) // FixMe
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
