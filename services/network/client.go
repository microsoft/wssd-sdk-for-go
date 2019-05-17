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
	"context"
)

type Service interface {
	Get(context.Context, string, string) (network.VirtualNetwork, error)
	CreateOrUpdate(context.Context, string, string, network.VirtualNetwork) (network.VirtualNetwork, error)
	Delete(context.Context, string, string) (network.VirtualNetwork, error)
}

type Client struct {
	group    string
	internal Service
}

func NewClient(subID, group string) (*Client, error) {
	c, err := newClient(subID)
	if err != nil {
		return nil, err
	}

	return &Client{group: group, internal: c}, nil
}

func (c *Client) Get(ctx context.Context, name string) (*Spec, error) {
	id, err := c.internal.Get(ctx, c.group, name)
	if err != nil && errors.IsNotFound(err) {
		return &Spec{&network.VirtualNetwork{}}, nil
	} else if err != nil {
		return nil, err
	}

	return &Spec{&id}, nil
}

func (c *Client) Ensure(ctx context.Context, name string, spec *Spec) error {
	result, err := c.internal.CreateOrUpdate(ctx, c.group, name, *spec.internal)
	if err != nil {
		return err
	}
	spec.internal = &result
	return nil
}
