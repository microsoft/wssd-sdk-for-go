// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package server

import (
	context "context"
	pb "github.com/microsoft/wssdagent/rpc/storage"
	"github.com/microsoft/wssdagent/services/storage/virtualharddisk"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	"github.com/microsoft/wssdagent/pkg/auth"
	log "k8s.io/klog"
	"sync"
)

type virtualharddiskAgentServer struct {
	provider *virtualharddisk.VirtualHardDiskProvider
	jwtAuthorizer *auth.JwtAuthorizer
	mu       sync.Mutex
}

func newVirtualHardDiskAgentServer(authorizer *auth.JwtAuthorizer) *virtualharddiskAgentServer {
	s := &virtualharddiskAgentServer{
		provider: virtualharddisk.GetVirtualHardDiskProvider(),
		jwtAuthorizer: authorizer,
	}
	return s
}

func (s *virtualharddiskAgentServer) Invoke(context context.Context, req *pb.VirtualHardDiskRequest) (*pb.VirtualHardDiskResponse, error) {

	_, err := s.jwtAuthorizer.ValidateTokenFromContext(context)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "Valid Token Required.")
	}

	res := new(pb.VirtualHardDiskResponse)

	switch req.GetOperationType() {
	case pb.Operation_GET:
		res.VirtualHardDiskSystems, err = s.provider.Get(context, req.GetVirtualHardDiskSystems())
	case pb.Operation_POST:
		res.VirtualHardDiskSystems, err = s.provider.CreateOrUpdate(context, req.GetVirtualHardDiskSystems())
	case pb.Operation_DELETE:
		err = s.provider.Delete(context, req.GetVirtualHardDiskSystems())
	default:
		return nil, status.Errorf(codes.Unavailable, "[VirtualHarddiskAgent] Invalid operation type specified")
	}
	log.Infof("[VirtualHarddiskAgent] [Invoke] Request[%v] Response[%v], Error [%v]", req, res, err)
	return res, GetGRPCError(err)

}
