// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package internal

import (
	"os"
	"path"

	pb "github.com/microsoft/wssdagent/rpc/security"
)

type SecretInternal struct {
	Srt        *pb.Secret
	Id         string
	ConfigPath string
}

func NewSecretInternal(id, basepath string) *SecretInternal {
	basepath = path.Join(basepath, id)
	os.MkdirAll(basepath, os.ModePerm)
	return &SecretInternal{
		Id:         id,
		ConfigPath: basepath,
	}
}
