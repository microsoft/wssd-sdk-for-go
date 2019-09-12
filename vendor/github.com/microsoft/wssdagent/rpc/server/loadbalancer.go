// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package server

import (
	context "context"
	pb "github.com/microsoft/wssdagent/rpc/network"
	"github.com/microsoft/wssdagent/services/network/loadbalancer"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"

	log "k8s.io/klog"
	"sync"
)

type loadbalancerAgentServer struct {
	mu       sync.Mutex
	provider loadbalancer.LoadBalancerProvider
}

func newLoadBalancerAgentServer() *loadbalancerAgentServer {
	s := &loadbalancerAgentServer{
		provider: loadbalancer.GetLoadBalancerProvider(loadbalancer.HCNSpec),
	}
	return s
}

func (s *loadbalancerAgentServer) Invoke(context context.Context, req *pb.LoadBalancerRequest) (*pb.LoadBalancerResponse, error) {
	res := new(pb.LoadBalancerResponse)
	var err error
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

	return res, err
}
