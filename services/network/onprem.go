// Copyright 2019 (c) Microsoft and contributors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package network

import (
	hcn "github.com/Microsoft/hcsshim/hcn"
)

type SpecOption func(*Spec) *Spec

type Spec struct {
	internal *hcn.HostComputeNetwork
}

func Internal(spec *Spec) *hcn.HostComputeNetwork {
	return spec.internal
}

func Name(name string) SpecOption {
	return func(o *Spec) *Spec {
		o.internal.Name = &name
		return o
	}
}

func AddressSpace(cidr string) SpecOption {
}

func Subnet(name, cidr, gateway string) SpecOption {
	return func(o *Spec) *Spec {
		if o.internal.HostComputeNetwork == nil {
			o.internal.HostComputeNetwork = &hcn.HostComputeNetwork{}
		}

		if o.internal.HostComputeNetwork.Ipams == nil {
			o.internal.HostComputeNetwork.Ipams = &[]hcn.Ipam{}
			o.internal.HostComputeNetwork.Ipams[0].Subnets = &[]hcn.Subnet{}
		}

		found := false
		for _, subnet := range *o.internal.HostComputeNetwork.Ipams[0].Subnets {
			if *subnet.Name == name {
				subnet.IpAddressPrefix = &cidr
				subnet.Routes = &hcn.Route{DestinationPrefix: "0.0.0.0/0", NextHop: gateway}
				found = true
			}
		}

		if !found {
			*o.internal.HostComputeNetwork.Ipam[0].Subnets = append(
				*o.internal.HostComputeNetwork.Ipam[0].Subnets,
				hcn.Subnet{
					Name:            name,
					IpAddressPrefix: cidr,
					Routes: &hcn.Route{
						DestinationPrefix: "0.0.0.0/0",
						NextHop:           gateway,
					},
				},
			)
		}

		return o
	}
}

func (s *Spec) Set(options ...SpecOption) {
	for _, option := range options {
		s = option(s)
	}
}

func (s *Spec) ID() string {
	return *s.internal.ID
}

func (s *Spec) Subnets() map[string]string {
	subnets := map[string]string{}
	for _, subnet := range *s.internal.Subnets {
		subnets[*subnet.Name] = *subnet.ID
	}
	return subnets
}
