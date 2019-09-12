// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package server

import (
	context "context"
	pb "github.com/microsoft/wssdagent/rpc/network"
	"github.com/microsoft/wssdagent/services/network/virtualnetworkinterface"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"

	log "k8s.io/klog"
	"sync"
)

type virtualnetworkinterfaceAgentServer struct {
	mu       sync.Mutex
	provider virtualnetworkinterface.VirtualNetworkInterfaceProvider
}

func newVirtualNetworkInterfaceAgentServer() *virtualnetworkinterfaceAgentServer {
	s := &virtualnetworkinterfaceAgentServer{
		provider: virtualnetworkinterface.GetVirtualNetworkInterfaceProvider(virtualnetworkinterface.HCNSpec),
	}
	return s
}

func (s *virtualnetworkinterfaceAgentServer) Invoke(context context.Context, req *pb.VirtualNetworkInterfaceRequest) (*pb.VirtualNetworkInterfaceResponse, error) {
	res := new(pb.VirtualNetworkInterfaceResponse)
	var err error
	switch req.GetOperationType() {
	case pb.Operation_GET:
		res.VirtualNetworkInterfaces, err = s.provider.Get(req.GetVirtualNetworkInterfaces())
	case pb.Operation_POST:
		res.VirtualNetworkInterfaces, err = s.provider.CreateOrUpdate(req.GetVirtualNetworkInterfaces())
	case pb.Operation_DELETE:
		err = s.provider.Delete(req.GetVirtualNetworkInterfaces())
	default:
		return nil, status.Errorf(codes.Unavailable, "[VirtualNetworkInterfaceAgent] Invalid operation type specified")
	}
	log.Infof("[VirtualNetworkInterfaceAgent] [Invoke] Request[%v] Response[%v], Error [%v]", req, res, err)

	return res, err
}
