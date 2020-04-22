// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license

package internal

import (
	"context"
	"fmt"
	"github.com/microsoft/moc/pkg/status"
	"github.com/microsoft/wssd-sdk-for-go/services/security"

	"github.com/microsoft/moc/pkg/auth"
	wssdcommonproto "github.com/microsoft/moc/rpc/common"
	wssdsecurity "github.com/microsoft/moc/rpc/nodeagent/security"
	wssdclient "github.com/microsoft/wssd-sdk-for-go/pkg/client"
	log "k8s.io/klog"
)

type client struct {
	wssdsecurity.IdentityAgentClient
}

// NewIdentityClientN- creates a client session with the backend wssd agent
func NewIdentityClient(subID string, authorizer auth.Authorizer) (*client, error) {
	c, err := wssdclient.GetIdentityClient(&subID, authorizer)
	if err != nil {
		return nil, err
	}
	return &client{c}, nil
}

// Get
func (c *client) Get(ctx context.Context, group, name string) (*[]security.Identity, error) {
	request := getIdentityRequest(wssdcommonproto.Operation_GET, name, nil)
	response, err := c.IdentityAgentClient.Invoke(ctx, request)
	if err != nil {
		return nil, err
	}
	return getIdentitiesFromResponse(response), nil
}

// CreateOrUpdate
func (c *client) CreateOrUpdate(ctx context.Context, group, name string, sg *security.Identity) (*security.Identity, error) {
	request := getIdentityRequest(wssdcommonproto.Operation_POST, name, sg)
	response, err := c.IdentityAgentClient.Invoke(ctx, request)
	if err != nil {
		log.Errorf("[Identity] Create failed with error %v", err)
		return nil, err
	}

	identity := getIdentitiesFromResponse(response)

	if len(*identity) == 0 {
		return nil, fmt.Errorf("[Identity][Create] Unexpected error: Creating a identity returned no result")
	}

	return &((*identity)[0]), err
}

// Delete methods invokes create or update on the client
func (c *client) Delete(ctx context.Context, group, name string) error {
	identity, err := c.Get(ctx, group, name)
	if err != nil {
		return err
	}
	if len(*identity) == 0 {
		return fmt.Errorf("Identity [%s] not found", name)
	}

	request := getIdentityRequest(wssdcommonproto.Operation_DELETE, name, &(*identity)[0])
	_, err = c.IdentityAgentClient.Invoke(ctx, request)
	return err
}

func getIdentitiesFromResponse(response *wssdsecurity.IdentityResponse) *[]security.Identity {
	identities := []security.Identity{}
	for _, resIdentities := range response.GetIdentitys() {
		identities = append(identities, *(getIdentity(resIdentities)))
	}

	return &identities
}

func getIdentityRequest(opType wssdcommonproto.Operation, name string, identity *security.Identity) *wssdsecurity.IdentityRequest {
	request := &wssdsecurity.IdentityRequest{
		OperationType: opType,
		Identitys:     []*wssdsecurity.Identity{},
	}
	if identity != nil {
		request.Identitys = append(request.Identitys, getWssdIdentity(identity))
	} else if len(name) > 0 {
		request.Identitys = append(request.Identitys,
			&wssdsecurity.Identity{
				Name: name,
			})
	}
	return request
}

func getIdentity(identity *wssdsecurity.Identity) *security.Identity {
	return &security.Identity{
		ID:   &identity.Id,
		Name: &identity.Name,
		IdentityProperties: &security.IdentityProperties{
			ProvisioningState: status.GetProvisioningState(identity.GetStatus().GetProvisioningStatus()),
			Statuses:          status.GetStatuses(identity.GetStatus()),
		},
	}
}

func getWssdIdentity(identity *security.Identity) *wssdsecurity.Identity {
	return &wssdsecurity.Identity{
		Name: *identity.Name,
	}
}
