// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package internal

import (
	"os"
	"path"

	pb "github.com/microsoft/wssdagent/rpc/security"
)

type KeyVaultInternal struct {
	Entity     *pb.KeyVault
	Id         string
	Name       string
	ConfigPath string
}

func NewKeyVaultInternal(id, basepath string, vault *pb.KeyVault) *KeyVaultInternal {
	basepath = path.Join(basepath, id)
	os.MkdirAll(basepath, os.ModePerm)
	vault.Id = id
	return &KeyVaultInternal{
		Id:         id,
		ConfigPath: basepath,
		Entity:     vault,
		Name:       vault.Name,
	}
}
