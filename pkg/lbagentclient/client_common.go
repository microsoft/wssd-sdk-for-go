// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

//NOTE: getClientConnection chould be refactored into moc repo.
package lbagentclient

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	admin_pb "github.com/microsoft/moc/rpc/common/admin"
	lbagent_pb "github.com/microsoft/moc/rpc/lbagent"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	log "k8s.io/klog"

	"github.com/microsoft/moc/pkg/auth"
)

//Note: This is the only thing that differs between the various client stuff.
const (
	debugModeTLS     = "WSSD_DEBUG_MODE"
	ServerPort   int = 46000
	AuthPort     int = 65000
)

var (
	mux             sync.Mutex
	connectionCache map[string]*grpc.ClientConn
)

func init() {
	connectionCache = map[string]*grpc.ClientConn{}
}

// Returns nil if debug mode is on; err if it is not
func isDebugMode() error {
	debugEnv := strings.ToLower(os.Getenv(debugModeTLS))
	if debugEnv == "on" {
		return nil
	}
	if viper.GetBool("Debug") {
		return nil
	}
	return fmt.Errorf("Debug Mode not set")
}

func getServerEndpoint(serverAddress *string) string {
	return fmt.Sprintf("%s:%d", *serverAddress, ServerPort)
}

func getAuthServerEndpoint(serverAddress *string) string {
	return fmt.Sprintf("%s:%d", *serverAddress, AuthPort)
}

func getDefaultDialOption(authorizer auth.Authorizer) []grpc.DialOption {
	var opts []grpc.DialOption

	// Debug Mode allows us to talk to wssdagent without a proper handshake
	// This means we can debug and test wssdagent without generating certs
	// and having proper tokens

	if ok := isDebugMode(); ok == nil {
		opts = append(opts, grpc.WithInsecure())
	} else {
		opts = append(opts, grpc.WithTransportCredentials(authorizer.WithTransportAuthorization()))
	}

	opts = append(opts, grpc.WithKeepaliveParams(
		keepalive.ClientParameters{
			Time:                1 * time.Minute,
			Timeout:             20 * time.Second,
			PermitWithoutStream: true,
		}))

	return opts
}

func getClientConnection(serverAddress *string, authorizer auth.Authorizer) (*grpc.ClientConn, error) {
	mux.Lock()
	defer mux.Unlock()
	endpoint := getServerEndpoint(serverAddress)

	conn, ok := connectionCache[endpoint]
	if ok {
		return conn, nil
	}

	opts := getDefaultDialOption(authorizer)
	conn, err := grpc.Dial(endpoint, opts...)
	if err != nil {
		log.Fatalf("Failed to dial: %v", err)
	}

	connectionCache[endpoint] = conn

	return conn, nil
}

// CloseClientConnectionByEndpoint allows a caller to close the current clientconn
// for a particular endpoint
func CloseClientConnectionByEndpoint(serverAddress *string) error {
	mux.Lock()
	defer mux.Unlock()
	endpoint := getServerEndpoint(serverAddress)

	conn, ok := connectionCache[endpoint]
	if ok {
		err := conn.Close()
		if err != nil {
			return err
		}

		delete(connectionCache, endpoint)
	}
	return nil
}

// GetLoadBalancerAgentClient returns the client to communicate with the lbagent
func GetLoadBalancerAgentClient(serverAddress *string, authorizer auth.Authorizer) (lbagent_pb.LoadBalancerAgentClient, error) {
	conn, err := getClientConnection(serverAddress, authorizer)
	if err != nil {
		log.Fatalf("Unable to get LoadBalancerAgentClient. Failed to dial: %v", err)
	}

	return lbagent_pb.NewLoadBalancerAgentClient(conn), nil
}

// GetHealthClient returns the health client to communicate with the lbagent
func GetHealthClient(serverAddress *string, authorizer auth.Authorizer) (admin_pb.HealthAgentClient, error) {
	conn, err := getClientConnection(serverAddress, authorizer)
	if err != nil {
		log.Fatalf("Unable to get HealthClient. Failed to dial: %v", err)
	}

	return admin_pb.NewHealthAgentClient(conn), nil
}
