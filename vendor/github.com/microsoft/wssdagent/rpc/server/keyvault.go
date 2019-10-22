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
	"github.com/microsoft/wssdagent/pkg/auth"
	"sync"
)

type keyvaultAgentServer struct {
	provider *keyvault.KeyVaultProvider
	jwtAuthorizer *auth.JwtAuthorizer
	mu       sync.Mutex
}

func newKeyVaultAgentServer(authorizer *auth.JwtAuthorizer) *keyvaultAgentServer {
	s := &keyvaultAgentServer{
		provider: keyvault.GetKeyVaultProvider(),
		jwtAuthorizer: authorizer,
	}
	return s
}

func (s *keyvaultAgentServer) Invoke(context context.Context, req *pb.KeyVaultRequest) (*pb.KeyVaultResponse, error) {
	
	_, err := s.jwtAuthorizer.ValidateTokenFromContext(context)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "Valid Token Required.")
	}
	res := new(pb.KeyVaultResponse)
	switch req.GetOperationType() {
	case pb.Operation_GET:
		res.KeyVaults, err = s.provider.Get(context, req.GetKeyVaults())
	case pb.Operation_POST:
		res.KeyVaults, err = s.provider.CreateOrUpdate(context, req.GetKeyVaults())
	case pb.Operation_DELETE:
		err = s.provider.Delete(context, req.GetKeyVaults())
	default:
		return nil, status.Errorf(codes.Unavailable, "[KeyVaultAgent] Invalid operation type specified")
	}
	log.Infof("[KeyVaultAgent] [Invoke] Request[%v] Response[%v], Error [%v]", req, res, err)

	return res, GetGRPCError(err)

}
