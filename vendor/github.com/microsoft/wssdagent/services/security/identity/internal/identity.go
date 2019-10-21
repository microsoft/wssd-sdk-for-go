// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package internal

import (
	"os"
	"path"

	pb "github.com/microsoft/wssdagent/rpc/security"
)

type IdentityInternal struct {
	Entity     *pb.Identity
	Id         string
	Name       string
	ConfigPath string
}

func NewIdentityInternal(id, basepath string, ident *pb.Identity) *IdentityInternal {
	basepath = path.Join(basepath, id)
	os.MkdirAll(basepath, os.ModePerm)
	ident.Id = id
	return &IdentityInternal{
		Id:         id,
		ConfigPath: basepath,
		Entity: 	ident,
		Name:		ident.Name,
	}
}