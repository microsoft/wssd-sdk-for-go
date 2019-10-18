// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package server

import (
	context "context"
	"github.com/microsoft/wssdagent/pkg/errors"
	pb "github.com/microsoft/wssdagent/rpc/network"
	"github.com/microsoft/wssdagent/services/network/virtualnetworkinterface"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	"github.com/microsoft/wssdagent/pkg/auth"

	log "k8s.io/klog"
	"sync"
)

type virtualnetworkinterfaceAgentServer struct {
	mu       sync.Mutex
	provider *virtualnetworkinterface.VirtualNetworkInterfaceProvider
	jwtAuthorizer *auth.JwtAuthorizer
}

func newVirtualNetworkInterfaceAgentServer(authorizer *auth.JwtAuthorizer) *virtualnetworkinterfaceAgentServer {
	s := &virtualnetworkinterfaceAgentServer{
		provider: virtualnetworkinterface.GetVirtualNetworkInterfaceProvider(),
		jwtAuthorizer: authorizer,
	}
	return s
}

func (s *virtualnetworkinterfaceAgentServer) Invoke(context context.Context, req *pb.VirtualNetworkInterfaceRequest) (*pb.VirtualNetworkInterfaceResponse, error) {
	_, err := s.jwtAuthorizer.ValidateTokenFromContext(context)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "Valid Token Required.")
	}

	res := new(pb.VirtualNetworkInterfaceResponse)
	switch req.GetOperationType() {
	case pb.Operation_GET:
		res.VirtualNetworkInterfaces, err = s.provider.Get(context, req.GetVirtualNetworkInterfaces())
	case pb.Operation_POST:
		res.VirtualNetworkInterfaces, err = s.provider.CreateOrUpdate(context, req.GetVirtualNetworkInterfaces())
	case pb.Operation_DELETE:
		err = s.provider.Delete(context, req.GetVirtualNetworkInterfaces())
	default:
		return nil, status.Errorf(codes.Unavailable, "[VirtualNetworkInterfaceAgent] Invalid operation type specified")
	}
	log.Infof("[VirtualNetworkInterfaceAgent] [Invoke] Request[%v] Response[%v], Error [%v]", req, res, err)
	if err == errors.NotFound {
		err = status.Errorf(codes.NotFound, err.Error())
	}

	return res, err
}
