// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package server

import (
	"fmt"
	"google.golang.org/grpc"
	log "k8s.io/klog"
	"net"

	"github.com/microsoft/wssdagent/pkg/wssdagent/apis/config"
	"github.com/microsoft/wssdagent/rpc/server"
)

func NewWssdAgentServer() error {
	log.Info("Starting wssdagent...")
	agentConfig := config.GetAgentConfiguration()
	log.Infof("AgentConfiguration [%v]", agentConfig)
	listenerStr := fmt.Sprintf("%s:%d", agentConfig.Address, agentConfig.Port)
	log.Infof("Starting a listener on [%s]", listenerStr)

	lis, err := net.Listen("tcp", listenerStr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
		return err
	}

	var opts []grpc.ServerOption
	grpcServer := server.RegisterServers(opts)
	server.RegisterDebugServer()
	return grpcServer.Serve(lis)
}
