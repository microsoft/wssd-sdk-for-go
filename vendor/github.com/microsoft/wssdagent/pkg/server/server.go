// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package server

import (
	"os"
	"strings"
	"fmt"
	"google.golang.org/grpc"
	log "k8s.io/klog"
	"net"

	"github.com/microsoft/wssdagent/pkg/apis/config"
	"github.com/microsoft/wssdagent/rpc/server"
	"google.golang.org/grpc/credentials"
)

const debugModeTLS = "WSSD_DEBUG_MODE"

// Returns nil if debug mode is on; err if it is not
func isDebugMode() error {
	debugEnv := strings.ToLower(os.Getenv(debugModeTLS))
	if debugEnv == "on" {
		return nil
	}
	return fmt.Errorf("Debug Mode not set")
}

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
	
	// Check if debug mode is on
	if ok := isDebugMode(); ok != nil {
		// Create the TLS credentials
		creds, err := credentials.NewServerTLSFromFile(
			config.GetTLSServerCertConfiguration(), 
			config.GetTLSServerKeyConfiguration())
		if err != nil {
		  log.Fatalf("could not load TLS keys: %s", err)
		}
		
	   opts = append(opts, grpc.Creds(creds))
	}

	grpcServer := server.RegisterServers(opts)
	server.RegisterDebugServer()
	return grpcServer.Serve(lis)
}
