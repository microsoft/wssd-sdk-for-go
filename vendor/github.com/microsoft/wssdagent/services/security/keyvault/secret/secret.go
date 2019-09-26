// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package secret

import (
	pb "github.com/microsoft/wssdagent/rpc/security"
)

type SecretProvider interface {
	CreateOrUpdate([]*pb.Secret) ([]*pb.Secret, error)
	Get([]*pb.Secret) ([]*pb.Secret, error)
	Delete([]*pb.Secret) error
}
