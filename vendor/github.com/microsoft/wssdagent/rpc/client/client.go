// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package client

import (
	"os"
	"strings"
	"fmt"
	compute_pb "github.com/microsoft/wssdagent/rpc/compute"
	security_pb "github.com/microsoft/wssdagent/rpc/security"
	network_pb "github.com/microsoft/wssdagent/rpc/network"
	storage_pb "github.com/microsoft/wssdagent/rpc/storage"
	"google.golang.org/grpc"
	log "k8s.io/klog"

	"github.com/microsoft/wssdagent/pkg/apis/config"
	"github.com/microsoft/wssd-sdk-for-go/pkg/auth"
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

func getServerEndpoint(serverAddress *string) string {
	return fmt.Sprintf("%s:%d", *serverAddress, config.ServerPort)
}

func getDefaultDialOption(authorizer auth.Authorizer) []grpc.DialOption {
	var opts []grpc.DialOption

	// Debug Mode allows us to talk to wssdagent without a proper handshake
	// This means we can debug and test wssdagent without generating certs
	// and having proper tokens

	// Check if debug mode is on
	if ok := isDebugMode(); ok == nil {
		// Working on a way to have debug mode still include tokens ... but not include
		// cert generation ... for now that is out of scope
		// creds := credentials.NewTLS(&tls.Config{InsecureSkipVerify: true})
		opts = append(opts, grpc.WithInsecure())
	} else {

		opts = append(opts, grpc.WithTransportCredentials(authorizer.WithTransportAuthorization()))
		opts = append(opts, grpc.WithPerRPCCredentials(authorizer.WithRPCAuthorization()))
	}

	return opts
}

// GetVirtualNetworkClient returns the virtual network client to communicate with the wssdagent
func GetVirtualNetworkClient(serverAddress *string, authorizer auth.Authorizer) (network_pb.VirtualNetworkAgentClient, error) {
	opts := getDefaultDialOption(authorizer)
	conn, err := grpc.Dial(getServerEndpoint(serverAddress), opts...)
	if err != nil {
		log.Fatalf("Unable to get VirtualNetworkClient. Failed to dial: %v", err)
	}

	return network_pb.NewVirtualNetworkAgentClient(conn), nil
}

// GetVirtualNetworkInterfaceClient returns the virtual network interface client to communicate with the wssd agent
func GetVirtualNetworkInterfaceClient(serverAddress *string, authorizer auth.Authorizer) (network_pb.VirtualNetworkInterfaceAgentClient, error) {
	opts := getDefaultDialOption(authorizer)

	conn, err := grpc.Dial(getServerEndpoint(serverAddress), opts...)
	if err != nil {
		log.Fatalf("Unable to get VirtualNetworkInterfaceClient. Failed to dial: %v", err)
	}

	return network_pb.NewVirtualNetworkInterfaceAgentClient(conn), nil
}

// GetLoadBalancerClient returns the loadbalancer client to communicate with the wssd agent
func GetLoadBalancerClient(serverAddress *string, authorizer auth.Authorizer) (network_pb.LoadBalancerAgentClient, error) {
	opts := getDefaultDialOption(authorizer)
	conn, err := grpc.Dial(getServerEndpoint(serverAddress), opts...)
	if err != nil {
		log.Fatalf("Unable to get LoadBalancerClient. Failed to dial: %v", err)
	}

	return network_pb.NewLoadBalancerAgentClient(conn), nil
}

// GetVirtualMachineClient returns the virtual machine client to comminicate with the wssd agent
func GetVirtualMachineClient(serverAddress *string, authorizer auth.Authorizer) (compute_pb.VirtualMachineAgentClient, error) {
	opts := getDefaultDialOption(authorizer)
	conn, err := grpc.Dial(getServerEndpoint(serverAddress), opts...)
	if err != nil {
		log.Fatalf("Unable to get VirtualMachineClient. Failed to dial: %v", err)
	}

	return compute_pb.NewVirtualMachineAgentClient(conn), nil
}

// GetVirtualMachineScaleSetClient returns the virtual machine client to comminicate with the wssd agent
func GetVirtualMachineScaleSetClient(serverAddress *string, authorizer auth.Authorizer) (compute_pb.VirtualMachineScaleSetAgentClient, error) {
	opts := getDefaultDialOption(authorizer)
	conn, err := grpc.Dial(getServerEndpoint(serverAddress), opts...)
	if err != nil {
		log.Fatalf("Unable to get VirtualMachineScaleSetClient. Failed to dial: %v", err)
	}

	return compute_pb.NewVirtualMachineScaleSetAgentClient(conn), nil
}

// GetVirtualHardDiskClient returns the virtual network client to communicate with the wssdagent
func GetVirtualHardDiskClient(serverAddress *string, authorizer auth.Authorizer) (storage_pb.VirtualHardDiskAgentClient, error) {
	opts := getDefaultDialOption(authorizer)
	conn, err := grpc.Dial(getServerEndpoint(serverAddress), opts...)
	if err != nil {
		log.Fatalf("Unable to get VirtualHardDiskClient. Failed to dial: %v", err)
	}

	return storage_pb.NewVirtualHardDiskAgentClient(conn), nil
}

// GetKeyVaultClient returns the keyvault client to communicate with the wssdagent
func GetKeyVaultClient(serverAddress *string, authorizer auth.Authorizer) (security_pb.KeyVaultAgentClient, error) {
	opts := getDefaultDialOption(authorizer)
	conn, err := grpc.Dial(getServerEndpoint(serverAddress), opts...)
	if err != nil {
		log.Fatalf("Unable to get KeyVaultClient. Failed to dial: %v", err)
	}

	return security_pb.NewKeyVaultAgentClient(conn), nil
}

// GetSecretClient returns the secret client to communicate with the wssdagent
func GetSecretClient(serverAddress *string, authorizer auth.Authorizer) (security_pb.SecretAgentClient, error) {
	opts := getDefaultDialOption(authorizer)
	conn, err := grpc.Dial(getServerEndpoint(serverAddress), opts...)
	if err != nil {
		log.Fatalf("Unable to get KeyVaultClient. Failed to dial: %v", err)
	}

	return security_pb.NewSecretAgentClient(conn), nil
}

// GetIdentityClient returns the secret client to communicate with the wssdagent
func GetIdentityClient(serverAddress *string, authorizer auth.Authorizer) (security_pb.IdentityAgentClient, error) {
	opts := getDefaultDialOption(authorizer)
	conn, err := grpc.Dial(getServerEndpoint(serverAddress), opts...)
	if err != nil {
		log.Fatalf("Unable to get IdentityClient. Failed to dial: %v", err)
	}

	return security_pb.NewIdentityAgentClient(conn), nil
}

// GetAuthenticationClient returns the secret client to communicate with the wssdagent
func GetAuthenticationClient(serverAddress *string, authorizer auth.Authorizer) (security_pb.AuthenticationAgentClient, error) {
	opts := getDefaultDialOption(authorizer)
	conn, err := grpc.Dial(getServerEndpoint(serverAddress), opts...)
	if err != nil {
		log.Fatalf("Unable to get AuthenticationClient. Failed to dial: %v", err)
	}

	return security_pb.NewAuthenticationAgentClient(conn), nil
}