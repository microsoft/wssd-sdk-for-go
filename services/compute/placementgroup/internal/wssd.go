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
	wssdcompute.PlacementGroupAgentClient
}

func NewPlacementGroupWssdClient(subID string, authorizer auth.Authorizer) (*wssdClient, error) {
	c, err := wssd.GetPlacementGroupClient(&subID, authorizer)
	if err != nil {
		return nil, err
	}

	return &wssdClient{c}, nil
}

func (c *wssdClient) Get(ctx context.Context, name string) (*[]compute.PlacementGroup, error) {
	request, err := c.getPlacementGroupRequest(wssdcommonproto.Operation_GET, name, nil)
	if err != nil {
		return nil, err
	}

	response, err := c.PlacementGroupAgentClient.Invoke(ctx, request)
	if err != nil {
		return nil, err
	}

	return c.getPlacementGroupFromResponse(response), nil
}

func (c *wssdClient) getPlacementGroupRequest(opType wssdcommonproto.Operation, name string, pgroup *compute.PlacementGroup) (*wssdcompute.PlacementGroupRequest, error) {
	request := &wssdcompute.PlacementGroupRequest{
		OperationType:   opType,
		PlacementGroups: []*wssdcompute.PlacementGroup{},
	}

	if pgroup != nil {
		wssdpgroup, err := getWssdPlacementGroup(pgroup)
		if err != nil {
			return nil, err
		}

		request.PlacementGroups = append(request.PlacementGroups, wssdpgroup)
	} else if len(name) > 0 {
		pgroup := &wssdcompute.PlacementGroup{
			Name: name,
		}

		request.PlacementGroups = append(request.PlacementGroups, pgroup)
	}

	return request, nil
}

func (c *wssdClient) CreateOrUpdate(ctx context.Context, name string, pgroup *compute.PlacementGroup) (*compute.PlacementGroup, error) {
	request, err := c.getPlacementGroupRequest(wssdcommonproto.Operation_POST, name, pgroup)
	if err != nil {
		return nil, err
	}

	response, err := c.PlacementGroupAgentClient.Invoke(ctx, request)
	if err != nil {
		return nil, err
	}

	pgroups := c.getPlacementGroupFromResponse(response)
	if len(*pgroups) == 0 {
		return nil, errors.Wrapf(errors.Unknown, "parsing of placement group in response failed")
	}

	return &(*pgroups)[0], nil
}

func (c *wssdClient) getPlacementGroupFromResponse(response *wssdcompute.PlacementGroupResponse) *[]compute.PlacementGroup {
	pgroups := []compute.PlacementGroup{}
	for _, pgroup := range response.GetPlacementGroups() {
		pgroups = append(pgroups, *(getPlacementGroup(pgroup)))
	}

	return &pgroups
}

func (c *wssdClient) Delete(ctx context.Context, name string) error {
	pgroup, err := c.Get(ctx, name)
	if err != nil {
		return err
	}

	if len(*pgroup) == 0 {
		// any error in the Get call should have been caught above, not expecting to hit this.
		return errors.Wrapf(errors.Unknown, "placement group response yielded no results for %s", name)
	}

	request, err := c.getPlacementGroupRequest(wssdcommonproto.Operation_DELETE, name, &(*pgroup)[0])
	if err != nil {
		return err
	}

	_, err = c.PlacementGroupAgentClient.Invoke(ctx, request)
	return err
}
