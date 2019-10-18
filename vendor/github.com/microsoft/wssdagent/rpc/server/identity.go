// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package server

import (
	context "context"
	pb "github.com/microsoft/wssdagent/rpc/security"
	"github.com/microsoft/wssdagent/services/security/identity"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	log "k8s.io/klog"
	"sync"
	"github.com/microsoft/wssdagent/pkg/auth"
)

type identityAgentServer struct {
	provider *identity.IdentityProvider
	jwtAuthorizer *auth.JwtAuthorizer
	mu       sync.Mutex
}

func newIdentityAgentServer(authorizer *auth.JwtAuthorizer) *identityAgentServer {
	s := &identityAgentServer{
		provider: identity.GetIdentityProvider(),
		jwtAuthorizer: authorizer,
	}
	return s
}

func (s *identityAgentServer) Invoke(context context.Context, req *pb.IdentityRequest) (*pb.IdentityResponse, error) {
	res := new(pb.IdentityResponse)
	var err error
	switch req.GetOperationType() {
	case pb.Operation_GET:
		res.Identitys, err = s.provider.Get(context, req.GetIdentitys())
	case pb.Operation_POST:
		res.Identitys, err = s.provider.CreateOrUpdate(context, req.GetIdentitys())
	case pb.Operation_DELETE:
		err = s.provider.Delete(context, req.GetIdentitys())
	default:
		return nil, status.Errorf(codes.Unavailable, "[IdentityAgent] Invalid operation type specified")
	}
	log.Infof("[IdentityAgent] [Invoke] Request[%v] Response[%v], Error [%v]", req, res, err)

	return res, err

}