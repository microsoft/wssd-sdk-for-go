// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package server

import (
	context "context"
	pb "github.com/microsoft/wssdagent/rpc/security"
	"github.com/microsoft/wssdagent/services/security/keyvault/secret"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	log "k8s.io/klog"
	"github.com/microsoft/wssdagent/pkg/auth"
	"sync"
)

type secretAgentServer struct {
	provider *secret.SecretProvider
	jwtAuthorizer *auth.JwtAuthorizer
	mu       sync.Mutex
}

func newSecretAgentServer(authorizer *auth.JwtAuthorizer) *secretAgentServer {
	s := &secretAgentServer{
		provider: secret.GetSecretProvider(),
		jwtAuthorizer: authorizer,
	}
	return s
}

func (s *secretAgentServer) Invoke(context context.Context, req *pb.SecretRequest) (*pb.SecretResponse, error) {

	_, err := s.jwtAuthorizer.ValidateTokenFromContext(context)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "Valid Token Required.")
	}

	res := new(pb.SecretResponse)

	switch req.GetOperationType() {
	case pb.Operation_GET:
		res.Secrets, err = s.provider.Get(context, req.GetSecrets())
	case pb.Operation_POST:
		res.Secrets, err = s.provider.CreateOrUpdate(context, req.GetSecrets())
	case pb.Operation_DELETE:
		err = s.provider.Delete(context, req.GetSecrets())
	default:
		return nil, status.Errorf(codes.Unavailable, "[SecretAgent] Invalid operation type specified")
	}
	log.Infof("[SecretAgent] [Invoke] Request[%v] Response[%v], Error [%v]", req, res, err)
	return res, GetGRPCError(err)
}
