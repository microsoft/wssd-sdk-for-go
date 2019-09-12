// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package server

import (
	context "context"
	pb "github.com/microsoft/wssdagent/rpc/storage"
	"github.com/microsoft/wssdagent/services/storage/virtualharddisk"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	log "k8s.io/klog"
	"sync"
)

type virtualharddiskAgentServer struct {
	provider virtualharddisk.VirtualHardDiskProvider
	mu       sync.Mutex
}

func newVirtualHardDiskAgentServer() *virtualharddiskAgentServer {
	s := &virtualharddiskAgentServer{
		provider: virtualharddisk.GetVirtualHardDiskProvider(virtualharddisk.HCSSpec),
	}
	return s
}

func (s *virtualharddiskAgentServer) Invoke(context context.Context, req *pb.VirtualHardDiskRequest) (*pb.VirtualHardDiskResponse, error) {
	res := new(pb.VirtualHardDiskResponse)
	var err error
	switch req.GetOperationType() {
	case pb.Operation_GET:
		res.VirtualHardDiskSystems, err = s.provider.Get(req.GetVirtualHardDiskSystems())
	case pb.Operation_POST:
		res.VirtualHardDiskSystems, err = s.provider.CreateOrUpdate(req.GetVirtualHardDiskSystems())
	case pb.Operation_DELETE:
		err = s.provider.Delete(req.GetVirtualHardDiskSystems())
	default:
		return nil, status.Errorf(codes.Unavailable, "[VirtualHarddiskAgent] Invalid operation type specified")
	}
	log.Infof("[VirtualHarddiskAgent] [Invoke] Request[%v] Response[%v], Error [%v]", req, res, err)

	return res, err

}
