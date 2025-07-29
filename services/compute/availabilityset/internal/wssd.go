package internal

import (
	"context"

	"github.com/microsoft/moc/pkg/auth"
	"github.com/microsoft/moc/pkg/errors"
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

	return c.getAvailabilitySetFromResponse(response), nil
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
		return nil, errors.Wrapf(errors.Unknown, "parsing of availability set in response failed")
	}

	return &(*avsets)[0], nil
}

func (c *wssdClient) getAvailabilitySetFromResponse(response *wssdcompute.AvailabilitySetResponse) *[]compute.AvailabilitySet {
	avsets := []compute.AvailabilitySet{}
	for _, avset := range response.GetAvailabilitySets() {
		avsets = append(avsets, *(getAvailabilitySet(avset)))
	}

	return &avsets
}

func (c *wssdClient) Delete(ctx context.Context, name string) error {
	avset, err := c.Get(ctx, name)
	if err != nil {
		return err
	}

	if len(*avset) == 0 {
		// any error in the Get call should have been caught above, not expecting to hit this.
		return errors.Wrapf(errors.Unknown, "availability set response yielded no results for %s", name)
	}

	request, err := c.getAvailabilitySetRequest(wssdcommonproto.Operation_DELETE, name, &(*avset)[0])
	if err != nil {
		return err
	}

	_, err = c.AvailabilitySetAgentClient.Invoke(ctx, request)
	return err
}

func (c *wssdClient) AddVmToAvailabilitySet(ctx context.Context, avset string, nodeagnetVMName, platformVMName string) error {
	request := &wssdcompute.AvailabilitySetOperationRequest{
		OperationType:   wssdcommonproto.AvailabilitySetOperation_ADD_VM,
		AvailabilitySet: avset,
		NodeagentVMName: nodeagnetVMName,
		PlatformVMName:  platformVMName,
	}

	_, err := c.AvailabilitySetAgentClient.Operate(ctx, request)
	return err
}

func (c *wssdClient) RemoveVmFromAvailabilitySet(ctx context.Context, avset string, nodeagnetVMName, platformVMName string) error {
	request := &wssdcompute.AvailabilitySetOperationRequest{
		OperationType:   wssdcommonproto.AvailabilitySetOperation_REMOVE_VM,
		AvailabilitySet: avset,
		NodeagentVMName: nodeagnetVMName,
		PlatformVMName:  platformVMName,
	}

	_, err := c.AvailabilitySetAgentClient.Operate(ctx, request)
	return err
}
