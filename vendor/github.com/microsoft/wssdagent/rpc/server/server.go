// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package server

import (
	"google.golang.org/grpc"
	log "k8s.io/klog"

	compute_pb "github.com/microsoft/wssdagent/rpc/compute"
	network_pb "github.com/microsoft/wssdagent/rpc/network"
	storage_pb "github.com/microsoft/wssdagent/rpc/storage"

	"net/http"
	_ "net/http/pprof"
)

func RegisterServers(opts []grpc.ServerOption) *grpc.Server {
	grpcServer := grpc.NewServer(opts...)

	log.Infof("Registering Rpc Agent Servers . . .")
	// Register compute agents
	compute_pb.RegisterVirtualMachineAgentServer(grpcServer, newVirtualMachineAgentServer())
	compute_pb.RegisterVirtualMachineScaleSetAgentServer(grpcServer, newVirtualMachineScaleSetAgentServer())

	// Register network agents
	network_pb.RegisterVirtualNetworkAgentServer(grpcServer, newVirtualNetworkAgentServer())
	network_pb.RegisterVirtualNetworkInterfaceAgentServer(grpcServer, newVirtualNetworkInterfaceAgentServer())
	network_pb.RegisterLoadBalancerAgentServer(grpcServer, newLoadBalancerAgentServer())

	// Register storage agents
	storage_pb.RegisterVirtualHardDiskAgentServer(grpcServer, newVirtualHardDiskAgentServer())

	return grpcServer
}

func RegisterDebugServer() {
	go func() {
		log.Info("Starting http/pprof. Access the Webpage via http://localhost:6060/debug/pprof/")
		log.Info(http.ListenAndServe("localhost:6060", nil))
	}()
}
