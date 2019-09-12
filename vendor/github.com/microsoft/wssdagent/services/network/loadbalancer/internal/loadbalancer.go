// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package internal

import (
	"os"
	"path"

	pb "github.com/microsoft/wssdagent/rpc/network"
)

type LoadBalancerInternal struct {
	Lb         *pb.LoadBalancer
	Id         string
	ConfigPath string
}

func NewLoadBalancerInternal(id, basepath string) *LoadBalancerInternal {
	baselbpath := path.Join(basepath, id)
	os.MkdirAll(baselbpath, os.ModePerm)
	return &LoadBalancerInternal{
		Id:         id,
		ConfigPath: baselbpath,
	}
}
