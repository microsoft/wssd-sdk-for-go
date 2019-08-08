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

package virtualmachinescaleset

import (
	"context"
	"fmt"
	log "k8s.io/klog"

	"github.com/microsoft/wssd-sdk-for-go/services/compute"
	"github.com/microsoft/wssd-sdk-for-go/services/network"
	wssdclient "github.com/microsoft/wssdagent/rpc/client"
	wssdcompute "github.com/microsoft/wssdagent/rpc/compute"
)

type client struct {
	subID string
	wssdcompute.VirtualMachineScaleSetAgentClient
}

// newClient - creates a client session with the backend wssd agent
func newVirtualMachineScaleSetClient(subID string) (*client, error) {
	c, err := wssdclient.GetVirtualMachineScaleSetClient(&subID)
	if err != nil {
		return nil, err
	}
	return &client{subID, c}, nil
}

// Get
func (c *client) Get(ctx context.Context, name string) (*[]compute.VirtualMachineScaleSet, error) {
	request, err := c.getVirtualMachineScaleSetRequest(wssdcompute.Operation_GET, name, nil)
	if err != nil {
		return nil, err
	}

	response, err := c.VirtualMachineScaleSetAgentClient.Invoke(ctx, request)
	if err != nil {
		return nil, err
	}
	log.Infof("[VirtualMachineScaleSet][Get] [%v]", response)
	return c.getVirtualMachineScaleSetFromResponse(response)
}

// CreateOrUpdate
func (c *client) CreateOrUpdate(ctx context.Context, name string, id string, sg *compute.VirtualMachineScaleSet) (*compute.VirtualMachineScaleSet, error) {
	request, err := c.getVirtualMachineScaleSetRequest(wssdcompute.Operation_POST, name, sg)
	if err != nil {
		return nil, err
	}
	response, err := c.VirtualMachineScaleSetAgentClient.Invoke(ctx, request)
	log.Infof("[VirtualMachineScaleSet][Create][Response] [%v]", response)
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
func (c *client) Delete(ctx context.Context, name string, id string) error {
	vmss, err := c.Get(ctx, name)
	if err != nil {
		return err
	}
	if len(*vmss) == 0 {
		return fmt.Errorf("Virtual Machine Scale Set [%s] not found", name)
	}

	request, err := c.getVirtualMachineScaleSetRequest(wssdcompute.Operation_DELETE, name, &(*vmss)[0])
	if err != nil {
		return err
	}
	response, err := c.VirtualMachineScaleSetAgentClient.Invoke(ctx, request)
	log.Infof("[VirtualMachineScaleSet][Delete] [%v]", response)
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

func (c *client) getVirtualMachineScaleSetRequest(opType wssdcompute.Operation, name string, vmss *compute.VirtualMachineScaleSet) (*wssdcompute.VirtualMachineScaleSetRequest, error) {
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
	log.Infof("[VirtualMachineScaleSet][Create][Request] [%v], Operation[%s]", request, opType)

	return request, nil

}

func (c *client) getVirtualMachineScaleSet(vmss *wssdcompute.VirtualMachineScaleSet) (*compute.VirtualMachineScaleSet, error) {
	vmprofile, err := c.getVirtualMachineScaleSetVMProfile(vmss.Virtualmachineprofile)
	if err != nil {
		return nil, err
	}
	return &compute.VirtualMachineScaleSet{
		BaseProperties: compute.BaseProperties{
			Name: &vmss.Name,
			ID:   &vmss.Id,
		},
		Sku: &compute.Sku{
			Name:     &vmss.Sku.Name,
			Capacity: &vmss.Sku.Capacity,
		},
		VirtualMachineProfile: vmprofile,
	}, nil
}

func (c *client) getVirtualMachineScaleSetVMProfile(vm *wssdcompute.VirtualMachineProfile) (*compute.VirtualMachineScaleSetVMProfile, error) {
	net, err := c.getVirtualMachineScaleSetNetworkProfile(vm.Network)
	if err != nil {
		return nil, err
	}

	return &compute.VirtualMachineScaleSetVMProfile{
		BaseProperties: compute.BaseProperties{
			Name: &vm.Vmprefix,
		},
		StorageProfile: c.getVirtualMachineScaleSetStorageProfile(vm.Storage),
		OsProfile:      c.getVirtualMachineScaleSetOSProfile(vm.Os),
		NetworkProfile: net,
	}, nil
}

func (c *client) getVirtualMachineScaleSetStorageProfile(s *wssdcompute.StorageConfiguration) *compute.StorageProfile {
	return &compute.StorageProfile{
		OsDisk:    c.getVirtualMachineScaleSetStorageProfileOsDisk(s.Osdisk),
		DataDisks: c.getVirtualMachineScaleSetStorageProfileDataDisks(s.Datadisks),
	}
}

func (c *client) getVirtualMachineScaleSetStorageProfileOsDisk(d *wssdcompute.Disk) *compute.OSDisk {
	return &compute.OSDisk{
		VhdId: &d.Diskid,
	}
}

func (c *client) getVirtualMachineScaleSetStorageProfileDataDisks(dd []*wssdcompute.Disk) *[]compute.DataDisk {
	cdd := []compute.DataDisk{}

	for _, i := range dd {
		cdd = append(cdd, compute.DataDisk{VhdId: &(i.Diskid)})
	}

	return &cdd

}

func (c *client) getVirtualMachineScaleSetNetworkProfile(n *wssdcompute.NetworkConfigurationScaleSet) (*compute.VirtualMachineScaleSetNetworkProfile, error) {
	np := &compute.VirtualMachineScaleSetNetworkProfile{
		NetworkInterfaceConfigurations: &[]network.VirtualNetworkInterface{},
	}

	for _, nic := range n.Interfaces {
		if nic == nil {
			continue
		}
		vnic, err := c.getVirtualMachineScaleSetNetworkInterface(nic)
		if err != nil {
			return nil, err
		}
		*np.NetworkInterfaceConfigurations = append(*np.NetworkInterfaceConfigurations, *vnic)
	}
	return np, nil
}

func (c *client) getVirtualMachineScaleSetNetworkInterface(nic *wssdcompute.VirtualNetworkInterface) (*network.VirtualNetworkInterface, error) {
	vnet := network.VirtualNetwork{
		BaseProperties: network.BaseProperties{
			Name: &nic.Networkname,
		},
	}

	vnetIntf := &network.VirtualNetworkInterface{
		BaseProperties: network.BaseProperties{
			Name: &nic.Name,
			ID:   &nic.Id,
		},
		VirtualNetwork: &vnet,
	}

	return vnetIntf, nil
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
		Id:   *(vmss.ID),
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
		Storage:  c.getWssdVirtualMachineScaleSetStorageConfiguration(vmp.StorageProfile),
		Os:       c.getWssdVirtualMachineScaleSetOSConfiguration(vmp.OsProfile),
		Network:  net,
	}, nil

}

func (c *client) getWssdVirtualMachineScaleSetStorageConfiguration(s *compute.StorageProfile) *wssdcompute.StorageConfiguration {
	return &wssdcompute.StorageConfiguration{
		Osdisk:    c.getWssdVirtualMachineScaleSetStorageConfigurationOsDisk(s.OsDisk),
		Datadisks: c.getWssdVirtualMachineScaleSetStorageConfigurationDataDisks(s.DataDisks),
	}
}

func (c *client) getWssdVirtualMachineScaleSetStorageConfigurationOsDisk(s *compute.OSDisk) *wssdcompute.Disk {
	return &wssdcompute.Disk{
		Diskid: *s.VhdId,
	}
}

func (c *client) getWssdVirtualMachineScaleSetStorageConfigurationDataDisks(s *[]compute.DataDisk) []*wssdcompute.Disk {
	datadisks := []*wssdcompute.Disk{}
	for _, d := range *s {
		datadisks = append(datadisks, &wssdcompute.Disk{Diskid: *d.VhdId})
	}

	return datadisks

}

func (c *client) getWssdVirtualMachineScaleSetNetworkConfiguration(s *compute.VirtualMachineScaleSetNetworkProfile) (*wssdcompute.NetworkConfigurationScaleSet, error) {
	nc := &wssdcompute.NetworkConfigurationScaleSet{
		Interfaces: []*wssdcompute.VirtualNetworkInterface{},
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

func (c *client) getWssdVirtualMachineScaleSetNetworkConfigurationNetworkInterface(nic *network.VirtualNetworkInterface) (*wssdcompute.VirtualNetworkInterface, error) {
	if nic.VirtualNetwork == nil {
		return nil, fmt.Errorf("Virtual Network reference required")
	}
	return &wssdcompute.VirtualNetworkInterface{
		Name:        *nic.Name,
		Networkname: *nic.VirtualNetwork.Name,
	}, nil
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
	if s.AdminUsername != nil && s.AdminPassword != nil {
		adminuser.Username = *s.AdminUsername
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
		osconfig.StartupScript = *s.CustomData
	}
	return &osconfig
}
