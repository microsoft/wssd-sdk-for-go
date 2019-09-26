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
	"sync"
)

type secretAgentServer struct {
	provider secret.SecretProvider
	mu       sync.Mutex
}

func newSecretAgentServer() *secretAgentServer {
	s := &secretAgentServer{
		provider: secret.GetSecretProvider(secret.FSSpec),
	}
	return s
}

func (s *secretAgentServer) Invoke(context context.Context, req *pb.SecretRequest) (*pb.SecretResponse, error) {
	res := new(pb.SecretResponse)
	var err error
	switch req.GetOperationType() {
	case pb.Operation_GET:
		res.Secrets, err = s.provider.Get(req.GetSecrets())
	case pb.Operation_POST:
		res.Secrets, err = s.provider.CreateOrUpdate(req.GetSecrets())
	case pb.Operation_DELETE:
		err = s.provider.Delete(req.GetSecrets())
	default:
		return nil, status.Errorf(codes.Unavailable, "[SecretAgent] Invalid operation type specified")
	}
	log.Infof("[SecretAgent] [Invoke] Request[%v] Response[%v], Error [%v]", req, res, err)

	return res, err

}
