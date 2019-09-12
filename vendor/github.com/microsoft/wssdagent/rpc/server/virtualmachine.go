// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package server

import (
	context "context"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	log "k8s.io/klog"
	"sync"

	pb "github.com/microsoft/wssdagent/rpc/compute"
	"github.com/microsoft/wssdagent/services/compute/virtualmachine"
)

type virtualmachineAgentServer struct {
	provider virtualmachine.VirtualMachineProvider
	mu       sync.Mutex
}

func newVirtualMachineAgentServer() *virtualmachineAgentServer {
	s := &virtualmachineAgentServer{
		provider: virtualmachine.GetVirtualMachineProvider(virtualmachine.HCSSpec),
	}
	return s
}

func (s *virtualmachineAgentServer) Invoke(context context.Context, req *pb.VirtualMachineRequest) (*pb.VirtualMachineResponse, error) {
	log.Infof("[VirtualMachineAgent] [Invoke] Request[%v]", req)
	res := new(pb.VirtualMachineResponse)
	var err error
	switch req.GetOperationType() {
	case pb.Operation_GET:
		res.VirtualMachineSystems, err = s.provider.Get(req.GetVirtualMachineSystems())
	case pb.Operation_POST:
		res.VirtualMachineSystems, err = s.provider.CreateOrUpdate(req.GetVirtualMachineSystems())
	case pb.Operation_DELETE:
		err = s.provider.Delete(req.GetVirtualMachineSystems())
	default:
		return nil, status.Errorf(codes.Unavailable, "[VirtualMachineAgent] Invalid operation type specified")
	}

	log.Infof("[VirtualMachineAgent] [Invoke] Request[%v] Response[%v], Error [%v]", req, res, err)

	return res, err
}
