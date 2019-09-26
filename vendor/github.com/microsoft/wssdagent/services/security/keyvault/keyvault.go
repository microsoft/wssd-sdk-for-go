// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package keyvault

import (
	pb "github.com/microsoft/wssdagent/rpc/security"
)

type KeyVaultProvider interface {
	CreateOrUpdate([]*pb.KeyVault) ([]*pb.KeyVault, error)
	Get([]*pb.KeyVault) ([]*pb.KeyVault, error)
	Delete([]*pb.KeyVault) error
}
