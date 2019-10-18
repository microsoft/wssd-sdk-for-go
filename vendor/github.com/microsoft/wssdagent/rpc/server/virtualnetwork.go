// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package server

import (
	context "context"
	"github.com/microsoft/wssdagent/pkg/errors"
	pb "github.com/microsoft/wssdagent/rpc/network"
	"github.com/microsoft/wssdagent/services/network/virtualnetwork"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	"github.com/microsoft/wssdagent/pkg/auth"

	log "k8s.io/klog"
	"sync"
)

type virtualnetworkAgentServer struct {
	mu       sync.Mutex
	provider *virtualnetwork.VirtualNetworkProvider
	jwtAuthorizer *auth.JwtAuthorizer
}

func newVirtualNetworkAgentServer(authorizer *auth.JwtAuthorizer) *virtualnetworkAgentServer {
	s := &virtualnetworkAgentServer{
		provider: virtualnetwork.GetVirtualNetworkProvider(),
		jwtAuthorizer: authorizer,
	}
	return s
}

func (s *virtualnetworkAgentServer) Invoke(context context.Context, req *pb.VirtualNetworkRequest) (*pb.VirtualNetworkResponse, error) {

	_, err := s.jwtAuthorizer.ValidateTokenFromContext(context)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "Valid Token Required.")
	}

	res := new(pb.VirtualNetworkResponse)

	switch req.GetOperationType() {
	case pb.Operation_GET:
		res.VirtualNetworks, err = s.provider.Get(context, req.GetVirtualNetworks())
	case pb.Operation_POST:
		res.VirtualNetworks, err = s.provider.CreateOrUpdate(context, req.GetVirtualNetworks())
	case pb.Operation_DELETE:
		err = s.provider.Delete(context, req.GetVirtualNetworks())
	default:
		return nil, status.Errorf(codes.Unavailable, "[VirtualNetworkAgent] Invalid operation type specified")
	}
	log.Infof("[VirtualNetworkAgent] [Invoke] Request[%v] Response[%v], Error [%v]", req, res, err)
	if err == errors.NotFound {
		err = status.Errorf(codes.NotFound, err.Error())
	}

	return res, err
}
