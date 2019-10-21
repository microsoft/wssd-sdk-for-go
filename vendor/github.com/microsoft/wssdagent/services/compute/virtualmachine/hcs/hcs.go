// +build windows
// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package hcs

import (
	"encoding/json"
	"github.com/Microsoft/hcsshim"
	log "k8s.io/klog"

	pb "github.com/microsoft/wssdagent/rpc/compute"
	networkpb "github.com/microsoft/wssdagent/rpc/network"
	storagepb "github.com/microsoft/wssdagent/rpc/storage"
	schema "github.com/microsoft/wssdagent/services/compute/virtualmachine/hcs/internal"
	"github.com/microsoft/wssdagent/services/compute/virtualmachine/internal"
)

const (
	owner    string = "wssdagent"
	memoryMB int    = 8192
	cpuCount int    = 4
)

type Client struct {
}

func NewClient() *Client {
	return &Client{}

}

// Create a Virtual Machine
func (c *Client) CreateVirtualMachine(vmInt *internal.VirtualMachineInternal, vmNic *networkpb.VirtualNetworkInterface, vhd *storagepb.VirtualHardDisk) (err error) {
	vm := vmInt.Entity
	vnicId := ""
	macAddress := ""
	vhdPath := vhd.Path

	if vmNic != nil {
		vnicId = vmNic.Id
		macAddress = vmNic.Macaddress
	}

	// Render HCS VM Spec
	vmspec, err := hcsshim.CreateVirtualMachineSpec(vm.Name, vm.Id, vhdPath, vmInt.SeedIso, owner, memoryMB, cpuCount, vnicId, macAddress)
	if err != nil {
		return
	}

	// Create the VM
	if err = vmspec.Create(); err != nil {
		return
	}

	// Start the VM
	if err = vmspec.Start(); err != nil {
		return
	}

	return
}

func (c *Client) HasVirtualMachine(vmInt *internal.VirtualMachineInternal) bool {
	return hcsshim.HasVirtualMachine(vmInt.Id)
}

// Delete a Virtual Machine
func (c *Client) CleanupVirtualMachine(vmint *internal.VirtualMachineInternal) (err error) {
	// Check with hcs
	if hcsshim.HasVirtualMachine(vmint.Id) {
		hcsvm, err1 := getVirtualMachineSpec(vmint.Entity)
		if err1 != nil {
			err = err1
			return
		}

		if err = hcsvm.Stop(); err != nil {
			log.Infof("Unable to stop the VM [%v]", err)
		}

		if err = hcsvm.Delete(); err != nil {
			return
		}
	}
	return
}

////////////////////// Private Functions //////////////////////////////////

// Conversion function
func getVirtualMachineSpec(vm *pb.VirtualMachine) (*hcsshim.VirtualMachineSpec, error) {
	vmspec, err := hcsshim.GetVirtualMachineSpec(vm.Id)
	if err != nil {
		return nil, err
	}
	vmspec.Name = vm.Name
	return vmspec, nil
}

// Conversion function
func getVirtualMachine(hcsvm *hcsshim.VirtualMachineSpec) (*pb.VirtualMachine, error) {
	vmspecString := hcsvm.String()
	log.Infof("[HCS][%s]", vmspecString)

	internalVmSchema := new(schema.ComputeSystem)
	if err := json.Unmarshal([]byte(vmspecString), internalVmSchema); err != nil {
		return nil, err
	}

	vm := &pb.VirtualMachine{
		Name: hcsvm.Name,
		Id:   hcsvm.ID,
		Storage: &pb.StorageConfiguration{
			Osdisk: &pb.Disk{
				Diskname: "", // internalVmSchema.VirtualMachine.Devices.Scsi["primary"].Attachments["0"].Path,
			},
		},
		Os:      &pb.OperatingSystemConfiguration{},
		Network: &pb.NetworkConfiguration{},
	}

	return vm, nil

}
