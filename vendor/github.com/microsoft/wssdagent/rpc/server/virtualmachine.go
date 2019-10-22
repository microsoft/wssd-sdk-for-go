// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package server

import (
	context "context"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	log "k8s.io/klog"
	"sync"
	"github.com/microsoft/wssdagent/pkg/auth"
	pb "github.com/microsoft/wssdagent/rpc/compute"
	"github.com/microsoft/wssdagent/services/compute/virtualmachine"
)

type virtualmachineAgentServer struct {
	provider *virtualmachine.VirtualMachineProvider
	jwtAuthorizer *auth.JwtAuthorizer
	mu       sync.Mutex
}

func newVirtualMachineAgentServer(authorizer *auth.JwtAuthorizer) *virtualmachineAgentServer {
	s := &virtualmachineAgentServer{
		provider: virtualmachine.GetVirtualMachineProvider(),
		jwtAuthorizer: authorizer,
	}
	return s
}

func (s *virtualmachineAgentServer) Invoke(context context.Context, req *pb.VirtualMachineRequest) (*pb.VirtualMachineResponse, error) {
	_, err := s.jwtAuthorizer.ValidateTokenFromContext(context)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "Valid Token Required.")
	}
	log.Infof("[VirtualMachineAgent] [Invoke] Request[%v]", req)
	res := new(pb.VirtualMachineResponse)
	switch req.GetOperationType() {
	case pb.Operation_GET:
		res.VirtualMachineSystems, err = s.provider.Get(context, req.GetVirtualMachineSystems())
	case pb.Operation_POST:
		res.VirtualMachineSystems, err = s.provider.CreateOrUpdate(context, req.GetVirtualMachineSystems())
	case pb.Operation_DELETE:
		err = s.provider.Delete(context, req.GetVirtualMachineSystems())
	default:
		return nil, status.Errorf(codes.Unavailable, "[VirtualMachineAgent] Invalid operation type specified")
	}

	log.Infof("[VirtualMachineAgent] [Invoke] Request[%v] Response[%v], Error [%v]", req, res, err)

	return res, GetGRPCError(err)
}
