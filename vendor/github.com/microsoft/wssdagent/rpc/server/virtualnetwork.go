// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package server

import (
	context "context"
	pb "github.com/microsoft/wssdagent/rpc/network"
	"github.com/microsoft/wssdagent/services/network/virtualnetwork"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"

	log "k8s.io/klog"
	"sync"
)

type virtualnetworkAgentServer struct {
	mu       sync.Mutex
	provider virtualnetwork.VirtualNetworkProvider
}

func newVirtualNetworkAgentServer() *virtualnetworkAgentServer {
	s := &virtualnetworkAgentServer{
		provider: virtualnetwork.GetVirtualNetworkProvider(virtualnetwork.HCNSpec),
	}
	return s
}

func (s *virtualnetworkAgentServer) Invoke(context context.Context, req *pb.VirtualNetworkRequest) (*pb.VirtualNetworkResponse, error) {
	res := new(pb.VirtualNetworkResponse)
	var err error
	switch req.GetOperationType() {
	case pb.Operation_GET:
		res.VirtualNetworks, err = s.provider.Get(req.GetVirtualNetworks())
	case pb.Operation_POST:
		res.VirtualNetworks, err = s.provider.CreateOrUpdate(req.GetVirtualNetworks())
	case pb.Operation_DELETE:
		err = s.provider.Delete(req.GetVirtualNetworks())
	default:
		return nil, status.Errorf(codes.Unavailable, "[VirtualNetworkAgent] Invalid operation type specified")
	}
	log.Infof("[VirtualNetworkAgent] [Invoke] Request[%v] Response[%v], Error [%v]", req, res, err)

	return res, err
}
