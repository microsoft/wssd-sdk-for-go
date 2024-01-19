package internal

import (
	wssdcompute "github.com/microsoft/moc/rpc/nodeagent/compute"
	"github.com/microsoft/wssd-sdk-for-go/services/compute"
)

func getWssdAvailabilitySet(avset *compute.AvailabilitySet) (*wssdcompute.AvailabilitySet, error) {
	// Implement the logic to convert avset to wssdavset
	return nil, nil
}

func (c *wssdClient) getAvailabilitySetFromResponse(response *wssdcompute.AvailabilitySetResponse) *[]compute.AvailabilitySet {
	// Implement the logic to convert the response to a slice of compute.AvailabilitySet
	return nil
}
