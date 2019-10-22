// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package server

import (
	"google.golang.org/grpc"
	log "k8s.io/klog"

	compute_pb "github.com/microsoft/wssdagent/rpc/compute"
	network_pb "github.com/microsoft/wssdagent/rpc/network"
	security_pb "github.com/microsoft/wssdagent/rpc/security"
	storage_pb "github.com/microsoft/wssdagent/rpc/storage"
	"github.com/microsoft/wssdagent/pkg/apis/config"
	"github.com/microsoft/wssdagent/pkg/auth"
	"net/http"
	_ "net/http/pprof"
)



func RegisterServers(opts []grpc.ServerOption) *grpc.Server {
	grpcServer := grpc.NewServer(opts...)

	keyLocation := config.GetPublicKeyConfiguration()
	jwtAuthorizer, err := auth.NewJwtAuthorizer(keyLocation)

	if err != nil {
		// This is most likely not found error
		// On Service start up we try to read the public key from the last users
		// login. If that fails we create a new public key

		// Log and continue
		log.Infof("Failed reading public key pem with error: %v", err)
	}

	log.Infof("Registering Rpc Agent Servers . . .")
	// Register compute agents
	compute_pb.RegisterVirtualMachineAgentServer(grpcServer, newVirtualMachineAgentServer(jwtAuthorizer))
	compute_pb.RegisterVirtualMachineScaleSetAgentServer(grpcServer, newVirtualMachineScaleSetAgentServer(jwtAuthorizer))

	// Register network agents
	network_pb.RegisterVirtualNetworkAgentServer(grpcServer, newVirtualNetworkAgentServer(jwtAuthorizer))
	network_pb.RegisterVirtualNetworkInterfaceAgentServer(grpcServer, newVirtualNetworkInterfaceAgentServer(jwtAuthorizer))
	network_pb.RegisterLoadBalancerAgentServer(grpcServer, newLoadBalancerAgentServer(jwtAuthorizer))

	// Register storage agents
	storage_pb.RegisterVirtualHardDiskAgentServer(grpcServer, newVirtualHardDiskAgentServer(jwtAuthorizer))

	// Register security agents
	security_pb.RegisterIdentityAgentServer(grpcServer, newIdentityAgentServer(jwtAuthorizer))
	security_pb.RegisterKeyVaultAgentServer(grpcServer, newKeyVaultAgentServer(jwtAuthorizer))
	security_pb.RegisterSecretAgentServer(grpcServer, newSecretAgentServer(jwtAuthorizer))
	security_pb.RegisterAuthenticationAgentServer(grpcServer, newAuthenticationAgentServer(jwtAuthorizer))

	return grpcServer
}

func RegisterDebugServer() {
	go func() {
		log.Info("Starting http/pprof. Access the Webpage via http://localhost:6060/debug/pprof/")
		log.Info(http.ListenAndServe("localhost:6060", nil))
	}()
}
