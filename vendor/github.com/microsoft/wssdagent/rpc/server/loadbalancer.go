// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package server

import (
	context "context"
	"github.com/microsoft/wssdagent/pkg/errors"
	pb "github.com/microsoft/wssdagent/rpc/network"
	"github.com/microsoft/wssdagent/services/network/loadbalancer"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"

	log "k8s.io/klog"
	"github.com/microsoft/wssdagent/pkg/auth"
	"sync"
)

type loadbalancerAgentServer struct {
	mu       sync.Mutex
	jwtAuthorizer *auth.JwtAuthorizer
	provider loadbalancer.LoadBalancerProvider
}

func newLoadBalancerAgentServer(authorizer *auth.JwtAuthorizer) *loadbalancerAgentServer {
	s := &loadbalancerAgentServer{
		provider: loadbalancer.GetLoadBalancerProvider(loadbalancer.HCNSpec),
		jwtAuthorizer: authorizer,
	}
	return s
}

func (s *loadbalancerAgentServer) Invoke(context context.Context, req *pb.LoadBalancerRequest) (*pb.LoadBalancerResponse, error) {
	_, err := s.jwtAuthorizer.ValidateTokenFromContext(context)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "Valid Token Required.")
	}
	res := new(pb.LoadBalancerResponse)
	switch req.GetOperationType() {
	case pb.Operation_GET:
		res.LoadBalancers, err = s.provider.Get(req.GetLoadBalancers())
	case pb.Operation_POST:
		res.LoadBalancers, err = s.provider.CreateOrUpdate(req.GetLoadBalancers())
	case pb.Operation_DELETE:
		err = s.provider.Delete(req.GetLoadBalancers())
	default:
		return nil, status.Errorf(codes.Unavailable, "[LoadbalancerAgent] Invalid operation type specified")
	}
	log.Infof("[LoadBalancerAgent] [Invoke] Request[%v] Response[%v], Error [%v]", req, res, err)

	if err == errors.NotFound {
		err = status.Errorf(codes.NotFound, err.Error())
	}
	return res, err
}
