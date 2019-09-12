// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package virtualharddisk

import (
	pb "github.com/microsoft/wssdagent/rpc/storage"
)

type VirtualHardDiskProvider interface {
	CreateOrUpdate([]*pb.VirtualHardDisk) ([]*pb.VirtualHardDisk, error)
	Get([]*pb.VirtualHardDisk) ([]*pb.VirtualHardDisk, error)
	Delete([]*pb.VirtualHardDisk) error
}

// DeleteVirtualHardDisk helper to delete virtual hard disks
func DeleteVirtualHardDisk(provider VirtualHardDiskProvider, vhdids []string) error {
	return provider.Delete(getVirtualHardDisk(vhdids))
}

func GetVirtualHardDiskPath(provider VirtualHardDiskProvider, vhdid string) (string, error) {
	return vhdid, nil
	//vnicsnew, err := provider.Get(getVirtualHardDisk([]string{vhdid}))
	//if err != nil {
	//		return "", err
	//	}
	//
	//	if len(vnicsnew) == 0 {
	//		return "", fmt.Errorf(vnicId + " not found")
	//	}
	//
	//	return (*(vnicsnew[0])).Path, nil
}

func CreateVirtualHardDisk(provider VirtualHardDiskProvider, vhdid, sourcePath, destinationPath string) error {
	return nil
}

func getVirtualHardDisk(vhds []string) []*pb.VirtualHardDisk {
	tmp := []*pb.VirtualHardDisk{}
	for _, vhd := range vhds {
		tmp = append(tmp, &pb.VirtualHardDisk{Id: vhd})
	}
	return tmp
}
