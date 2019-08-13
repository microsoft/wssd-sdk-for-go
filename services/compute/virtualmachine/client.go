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

package virtualmachine

import (
	"context"
	"github.com/microsoft/wssd-sdk-for-go/services/compute"
)

type Service interface {
	Get(context.Context, string, string) (*[]compute.VirtualMachine, error)
	CreateOrUpdate(context.Context, string, string, *compute.VirtualMachine) (*compute.VirtualMachine, error)
	Delete(context.Context, string, string) error
}

type VirtualMachineClient struct {
	compute.BaseClient
	internal Service
}

func NewVirtualMachineClient(cloudFQDN string) (*VirtualMachineClient, error) {
	c, err := newVirtualMachineClient(cloudFQDN)
	if err != nil {
		return nil, err
	}

	return &VirtualMachineClient{internal: c}, nil
}

// Get methods invokes the client Get method
func (c *VirtualMachineClient) Get(ctx context.Context, group, name string) (*[]compute.VirtualMachine, error) {
	return c.internal.Get(ctx, group, name)
}

// CreateOrUpdate methods invokes create or update on the client
func (c *VirtualMachineClient) CreateOrUpdate(ctx context.Context, group, name string, compute *compute.VirtualMachine) (*compute.VirtualMachine, error) {
	return c.internal.CreateOrUpdate(ctx, group, name, compute)
}

// Delete methods invokes delete of the compute resource
func (c *VirtualMachineClient) Delete(ctx context.Context, group string, name string) error {
	return c.internal.Delete(ctx, group, name)
}
