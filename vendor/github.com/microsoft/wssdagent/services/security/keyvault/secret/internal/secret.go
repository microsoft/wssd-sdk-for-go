// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package internal

import (
	"os"
	"path"

	pb "github.com/microsoft/wssdagent/rpc/security"
)

type SecretInternal struct {
	Entity     *pb.Secret
	Id         string
	Name       string
	VaultName  string
	ConfigPath string
}

func NewSecretInternal(id, basepath string, sec *pb.Secret) *SecretInternal {
	basepath = path.Join(basepath, id)
	os.MkdirAll(basepath, os.ModePerm)
	sec.Id = id
	return &SecretInternal{
		Id:         id,
		ConfigPath: basepath,
		Entity:     sec,
		Name:       sec.Name,
		VaultName:  sec.VaultName,
	}
}
