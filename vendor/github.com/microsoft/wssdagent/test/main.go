// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package main

import (
	"context"
	"flag"
	log "k8s.io/klog"
	"time"

	"github.com/microsoft/wssdagent/rpc/client"
	"github.com/microsoft/wssdagent/rpc/compute"
)

var (
	tls                = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
	caFile             = flag.String("ca_file", "", "The file containning the CA root cert file")
	serverAddr         = flag.String("server_addr", "127.0.0.1", "The server address in the format of host")
	serverHostOverride = flag.String("server_host_override", "x.test.youtube.com", "The server name use to verify the hostname returned by TLS handshake")
)

func main() {
	flag.Parse()

	vmclient, err := client.GetVirtualMachineClient(serverAddr)

	if err != nil {
		log.Fatalf("failed to create a VM Client : %v", err)
	}

	request := &compute.VirtualMachineRequest{
		OperationType: compute.Operation_GET,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	response, err := vmclient.Invoke(ctx, request)

	if err != nil {
		log.Printf("failed to list all vms : %v", err)
		return
	}

	for _, vm := range response.GetVirtualMachineSystems() {
		log.Printf("VM : %v", vm)

	}

}
