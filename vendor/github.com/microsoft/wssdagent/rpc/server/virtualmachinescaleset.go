// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package server

import (
	context "context"
	"github.com/microsoft/wssdagent/pkg/errors"
	pb "github.com/microsoft/wssdagent/rpc/compute"
	"github.com/microsoft/wssdagent/services/compute/virtualmachinescaleset"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	log "k8s.io/klog"
	"github.com/microsoft/wssdagent/pkg/auth"
	"sync"
)

type virtualmachinescalesetAgentServer struct {
	provider *virtualmachinescaleset.VirtualMachineScaleSetProvider
	jwtAuthorizer *auth.JwtAuthorizer
	mu       sync.Mutex
}

func newVirtualMachineScaleSetAgentServer(authorizer *auth.JwtAuthorizer) *virtualmachinescalesetAgentServer {
	s := &virtualmachinescalesetAgentServer{
		provider: virtualmachinescaleset.GetVirtualMachineScaleSetProvider(),
		jwtAuthorizer: authorizer,
	}
	return s
}

func (s *virtualmachinescalesetAgentServer) Invoke(context context.Context, req *pb.VirtualMachineScaleSetRequest) (*pb.VirtualMachineScaleSetResponse, error) {
	_, err := s.jwtAuthorizer.ValidateTokenFromContext(context)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "Valid Token Required.")
	}

	log.Infof("[VirtualMachineScaleSetAgent] [Invoke] Request[%v]", req)
	res := new(pb.VirtualMachineScaleSetResponse)

	switch req.GetOperationType() {
	case pb.Operation_GET:
		res.VirtualMachineScaleSetSystems, err = s.provider.Get(context, req.GetVirtualMachineScaleSetSystems())
	case pb.Operation_POST:
		res.VirtualMachineScaleSetSystems, err = s.provider.CreateOrUpdate(context, req.GetVirtualMachineScaleSetSystems())
	case pb.Operation_DELETE:
		err = s.provider.Delete(context, req.GetVirtualMachineScaleSetSystems())
	default:
		return nil, status.Errorf(codes.Unavailable, "[VirtualMachineScaleSetAgent] Invalid operation type specified")
	}
	log.Infof("[VirtualMachineScaleSetAgent] [Invoke] Request[%v] Response[%v], Error [%v]", req, res, err)

	if err == errors.NotFound {
		err = status.Errorf(codes.NotFound, err.Error())
	}
	return res, err

}
