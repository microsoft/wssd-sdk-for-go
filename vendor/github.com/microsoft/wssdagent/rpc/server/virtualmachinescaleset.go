// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package server

import (
	context "context"
	pb "github.com/microsoft/wssdagent/rpc/compute"
	"github.com/microsoft/wssdagent/services/compute/virtualmachinescaleset"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	log "k8s.io/klog"
	"sync"
)

type virtualmachinescalesetAgentServer struct {
	provider virtualmachinescaleset.VirtualMachineScaleSetProvider
	mu       sync.Mutex
}

func newVirtualMachineScaleSetAgentServer() *virtualmachinescalesetAgentServer {
	s := &virtualmachinescalesetAgentServer{
		provider: virtualmachinescaleset.GetVirtualMachineScaleSetProvider(virtualmachinescaleset.HCSSpec),
	}
	return s
}

func (s *virtualmachinescalesetAgentServer) Invoke(context context.Context, req *pb.VirtualMachineScaleSetRequest) (*pb.VirtualMachineScaleSetResponse, error) {
	log.Infof("[VirtualMachineScaleSetAgent] [Invoke] Request[%v]", req)
	res := new(pb.VirtualMachineScaleSetResponse)
	var err error
	switch req.GetOperationType() {
	case pb.Operation_GET:
		res.VirtualMachineScaleSetSystems, err = s.provider.Get(req.GetVirtualMachineScaleSetSystems())
	case pb.Operation_POST:
		res.VirtualMachineScaleSetSystems, err = s.provider.CreateOrUpdate(req.GetVirtualMachineScaleSetSystems())
	case pb.Operation_DELETE:
		err = s.provider.Delete(req.GetVirtualMachineScaleSetSystems())
	default:
		return nil, status.Errorf(codes.Unavailable, "[VirtualMachineScaleSetAgent] Invalid operation type specified")
	}
	log.Infof("[VirtualMachineScaleSetAgent] [Invoke] Request[%v] Response[%v], Error [%v]", req, res, err)

	return res, err

}
