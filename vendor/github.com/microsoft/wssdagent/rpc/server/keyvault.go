// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package server

import (
	context "context"
	pb "github.com/microsoft/wssdagent/rpc/security"
	"github.com/microsoft/wssdagent/services/security/keyvault"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	log "k8s.io/klog"
	"sync"
)

type keyvaultAgentServer struct {
	provider keyvault.KeyVaultProvider
	mu       sync.Mutex
}

func newKeyVaultAgentServer() *keyvaultAgentServer {
	s := &keyvaultAgentServer{
		provider: keyvault.GetKeyVaultProvider(keyvault.FSSpec),
	}
	return s
}

func (s *keyvaultAgentServer) Invoke(context context.Context, req *pb.KeyVaultRequest) (*pb.KeyVaultResponse, error) {
	res := new(pb.KeyVaultResponse)
	var err error
	switch req.GetOperationType() {
	case pb.Operation_GET:
		res.KeyVaults, err = s.provider.Get(req.GetKeyVaults())
	case pb.Operation_POST:
		res.KeyVaults, err = s.provider.CreateOrUpdate(req.GetKeyVaults())
	case pb.Operation_DELETE:
		err = s.provider.Delete(req.GetKeyVaults())
	default:
		return nil, status.Errorf(codes.Unavailable, "[KeyVaultAgent] Invalid operation type specified")
	}
	log.Infof("[KeyVaultAgent] [Invoke] Request[%v] Response[%v], Error [%v]", req, res, err)

	return res, err

}
