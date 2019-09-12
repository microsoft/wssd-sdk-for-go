// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package client

import (
	"fmt"
	compute_pb "github.com/microsoft/wssdagent/rpc/compute"
	network_pb "github.com/microsoft/wssdagent/rpc/network"
	storage_pb "github.com/microsoft/wssdagent/rpc/storage"
	"google.golang.org/grpc"
	log "k8s.io/klog"

	"github.com/microsoft/wssdagent/pkg/wssdagent/apis/config"
)

func getServerEndpoint(serverAddress *string) string {
	return fmt.Sprintf("%s:%d", *serverAddress, config.ServerPort)
}

func getDefaultDialOption() []grpc.DialOption {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())

	return opts
}

// GetVirtualNetworkClient returns the virtual network client to communicate with the wssdagent
func GetVirtualNetworkClient(serverAddress *string) (network_pb.VirtualNetworkAgentClient, error) {
	opts := getDefaultDialOption()
	conn, err := grpc.Dial(getServerEndpoint(serverAddress), opts...)
	if err != nil {
		log.Fatalf("Unable to get VirtualNetworkClient. Failed to dial: %v", err)
	}

	return network_pb.NewVirtualNetworkAgentClient(conn), nil
}

// GetVirtualNetworkInterfaceClient returns the virtual network interface client to communicate with the wssd agent
func GetVirtualNetworkInterfaceClient(serverAddress *string) (network_pb.VirtualNetworkInterfaceAgentClient, error) {
	opts := getDefaultDialOption()
	conn, err := grpc.Dial(getServerEndpoint(serverAddress), opts...)
	if err != nil {
		log.Fatalf("Unable to get VirtualNetworkInterfaceClient. Failed to dial: %v", err)
	}

	return network_pb.NewVirtualNetworkInterfaceAgentClient(conn), nil
}

// GetLoadBalancerClient returns the loadbalancer client to communicate with the wssd agent
func GetLoadBalancerClient(serverAddress *string) (network_pb.LoadBalancerAgentClient, error) {
	opts := getDefaultDialOption()
	conn, err := grpc.Dial(getServerEndpoint(serverAddress), opts...)
	if err != nil {
		log.Fatalf("Unable to get LoadBalancerClient. Failed to dial: %v", err)
	}

	return network_pb.NewLoadBalancerAgentClient(conn), nil
}

// GetVirtualMachineClient returns the virtual machine client to comminicate with the wssd agent
func GetVirtualMachineClient(serverAddress *string) (compute_pb.VirtualMachineAgentClient, error) {
	opts := getDefaultDialOption()
	conn, err := grpc.Dial(getServerEndpoint(serverAddress), opts...)
	if err != nil {
		log.Fatalf("Unable to get VirtualMachineClient. Failed to dial: %v", err)
	}

	return compute_pb.NewVirtualMachineAgentClient(conn), nil
}

// GetVirtualMachineScaleSetClient returns the virtual machine client to comminicate with the wssd agent
func GetVirtualMachineScaleSetClient(serverAddress *string) (compute_pb.VirtualMachineScaleSetAgentClient, error) {
	opts := getDefaultDialOption()
	conn, err := grpc.Dial(getServerEndpoint(serverAddress), opts...)
	if err != nil {
		log.Fatalf("Unable to get VirtualMachineScaleSetClient. Failed to dial: %v", err)
	}

	return compute_pb.NewVirtualMachineScaleSetAgentClient(conn), nil
}

// GetVirtualHardDiskClient returns the virtual network client to communicate with the wssdagent
func GetVirtualHardDiskClient(serverAddress *string) (storage_pb.VirtualHardDiskAgentClient, error) {
	opts := getDefaultDialOption()
	conn, err := grpc.Dial(getServerEndpoint(serverAddress), opts...)
	if err != nil {
		log.Fatalf("Unable to get VirtualHardDiskClient. Failed to dial: %v", err)
	}

	return storage_pb.NewVirtualHardDiskAgentClient(conn), nil
}
