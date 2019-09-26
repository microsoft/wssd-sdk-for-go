// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package internal

import (
	"os"
	"path"

	pb "github.com/microsoft/wssdagent/rpc/security"
)

type KeyVaultInternal struct {
	Skv        *pb.KeyVault
	Id         string
	ConfigPath string
}

func NewKeyVaultInternal(id, basepath string) *KeyVaultInternal {
	basepath = path.Join(basepath, id)
	os.MkdirAll(basepath, os.ModePerm)
	return &KeyVaultInternal{
		Id:         id,
		ConfigPath: basepath,
	}
}
