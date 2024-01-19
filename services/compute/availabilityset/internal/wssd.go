package internal

import (
	"context"
	"fmt"

	"github.com/microsoft/moc/pkg/auth"
	wssdcommonproto "github.com/microsoft/moc/rpc/common"
	wssdcompute "github.com/microsoft/moc/rpc/nodeagent/compute"
	wssd "github.com/microsoft/wssd-sdk-for-go/pkg/client"
	"github.com/microsoft/wssd-sdk-for-go/services/compute"
)

type wssdClient struct {
	wssdcompute.AvailabilitySetAgentClient
}

func NewAvailabilitySetWssdClient(subID string, authorizer auth.Authorizer) (*wssdClient, error) {
	c, err := wssd.GetAvailabilitySetClient(&subID, authorizer)
	if err != nil {
		return nil, err
	}

	return &wssdClient{c}, nil
}

func (c *wssdClient) Get(ctx context.Context, name string) (*[]compute.AvailabilitySet, error) {
	request, err := c.getAvailabilitySetRequest(wssdcommonproto.Operation_GET, name, nil)
	if err != nil {
		return nil, err
	}

	response, err := c.AvailabilitySetAgentClient.Invoke(ctx, request)
	if err != nil {
		return nil, err
	}

	avsets := c.getAvailabilitySetFromResponse(response)
	return avsets, nil
}

func (c *wssdClient) getAvailabilitySetRequest(opType wssdcommonproto.Operation, name string, avset *compute.AvailabilitySet) (*wssdcompute.AvailabilitySetRequest, error) {
	request := &wssdcompute.AvailabilitySetRequest{
		OperationType:    opType,
		AvailabilitySets: []*wssdcompute.AvailabilitySet{},
	}

	if avset != nil {
		wssdavset, err := getWssdAvailabilitySet(avset)
		if err != nil {
			return nil, err
		}
		request.AvailabilitySets = append(request.AvailabilitySets, wssdavset)
	} else if len(name) > 0 {
		avset := &wssdcompute.AvailabilitySet{
			Name: name,
		}
		request.AvailabilitySets = append(request.AvailabilitySets, avset)
	}

	return request, nil
}

func getWssdAvailabilitySet(avset *compute.AvailabilitySet) (*wssdcompute.AvailabilitySet, error) {
	// Implement the logic to convert avset to wssdavset
	return nil, nil
}

func (c *wssdClient) getAvailabilitySetFromResponse(response *wssdcompute.AvailabilitySetResponse) *[]compute.AvailabilitySet {
	// Implement the logic to convert the response to a slice of compute.AvailabilitySet
	return nil
}

func (c *wssdClient) CreateOrUpdate(ctx context.Context, name string, avset *compute.AvailabilitySet) (*compute.AvailabilitySet, error) {
	request, err := c.getAvailabilitySetRequest(wssdcommonproto.Operation_POST, name, avset)
	if err != nil {
		return nil, err
	}

	response, err := c.AvailabilitySetAgentClient.Invoke(ctx, request)
	if err != nil {
		return nil, err
	}

	avsets := c.getAvailabilitySetFromResponse(response)
	if len(*avsets) == 0 {
		return nil, fmt.Errorf("Creation of Virtual Machine failed to unknown reason.")
	}

	return &(*avsets)[0], nil
}

func (c *wssdClient) Delete(ctx context.Context, name string) error {
	avset, err := c.Get(ctx, name)
	if err != nil {
		return err
	}

	if len(*avset) == 0 {
		return fmt.Errorf("Availability Set [%s] not found", name)
	}

	request, err := c.getAvailabilitySetRequest(wssdcommonproto.Operation_DELETE, name, &(*avset)[0])
	if err != nil {
		return err
	}

	_, err = c.AvailabilitySetAgentClient.Invoke(ctx, request)
	return err
}
