// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package server

import (
	context "context"
	pb "github.com/microsoft/wssdagent/rpc/security"
	"github.com/microsoft/wssdagent/services/security/authentication"
	"github.com/microsoft/wssdagent/pkg/auth"
	log "k8s.io/klog"
	"sync"
)

type authenticationAgentServer struct {
	provider *authentication.AuthenticationProvider
	jwtAuthorizer *auth.JwtAuthorizer
	mu       sync.Mutex
}

func newAuthenticationAgentServer(authorizer *auth.JwtAuthorizer) *authenticationAgentServer {
	s := &authenticationAgentServer{
		provider: authentication.GetAuthenticationProvider(),
		jwtAuthorizer: authorizer,
	}
	return s
}

func (s *authenticationAgentServer) Login(context context.Context, req *pb.AuthenticationRequest) (*pb.AuthenticationResponse, error) {
	res := new(pb.AuthenticationResponse)
	var err error
	token, err := s.provider.Login(context, req.GetIdentity(), s.jwtAuthorizer)
	if err != nil {
		log.Infof("[AuthenticationAgent] Error [%v]", err)
		return nil, err	
	}
	res.Token = *token

	log.Infof("[AuthenticationAgent] [Invoke] Request[%v] Response[%v], Error [%v]", req, res, err)
	return res, err
}