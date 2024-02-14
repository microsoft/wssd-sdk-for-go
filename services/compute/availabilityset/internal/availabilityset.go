package internal

import (
	"github.com/microsoft/moc/pkg/errors"
	"github.com/microsoft/moc/pkg/status"
	prototags "github.com/microsoft/moc/pkg/tags"
	wssdcommonproto "github.com/microsoft/moc/rpc/common"
	wssdcompute "github.com/microsoft/moc/rpc/nodeagent/compute"
	"github.com/microsoft/wssd-sdk-for-go/services/compute"
)

// Conversion functions from client to rpc
func getWssdAvailabilitySet(avset *compute.AvailabilitySet) (*wssdcompute.AvailabilitySet, error) {
	errorPrefix := "Error converting AvailabilitySet to WssdAvailabilitySet"

	if avset == nil {
		return nil, errors.Wrapf(errors.InvalidInput, "%s, AvailabilitySet cannot be nil", errorPrefix)
	}

	if avset.Name == nil || len(*avset.Name) == 0 || len(*avset.Name) > 200 {
		return nil, errors.Wrapf(errors.InvalidInput, "%s, Name cannot be empty or more than 200 characters", errorPrefix)
	}

	wssdavset := &wssdcompute.AvailabilitySet{
		Name: *avset.Name,
		Tags: getWssdTags(avset.Tags),
	}

	if avset.AvailabilitySetProperties == nil {
		return wssdavset, nil
	}

	fd, err := getWssdAvailabilitySetFaultDomain(avset)
	if err != nil {
		return nil, errors.Wrapf(err, "%s, failed to get fault domain", errorPrefix)
	}

	entity, err := getWssdAvailabilitySetEntity(avset)
	if err != nil {
		return nil, errors.Wrapf(err, "%s, failed to get entity", errorPrefix)
	}

	vms, err := getWssdAvailabilitySetVMs(avset)
	if err != nil {
		return nil, errors.Wrapf(err, "%s, failed to get vms", errorPrefix)
	}

	wssdavset = &wssdcompute.AvailabilitySet{
		Name:                     *avset.Name,
		Tags:                     getWssdTags(avset.Tags),
		Entity:                   entity,
		PlatformFaultDomainCount: *fd,
		VirtualMachines:          vms,
		Status:                   status.GetFromStatuses(avset.Statuses),
	}

	return wssdavset, nil
}

func getWssdTags(tags map[string]*string) *wssdcommonproto.Tags {
	return prototags.MapToProto(tags)
}

func getWssdAvailabilitySetFaultDomain(avset *compute.AvailabilitySet) (*int32, error) {
	// PlatformFaultDomainCount upperbound should be nodecount (number of nodes managed by cloud agent) inclusive
	// We don't apply upper bound validation here because nodeagent is not aware of this information
	// this validation needs to take place in the cloud agent.
	if avset.PlatformFaultDomainCount == nil || *avset.PlatformFaultDomainCount < 2 {
		return nil, errors.Wrapf(errors.InvalidInput, "PlatformFaultDomainCount cannot be less than 2")
	}

	return avset.PlatformFaultDomainCount, nil
}

func getWssdAvailabilitySetEntity(avset *compute.AvailabilitySet) (*wssdcommonproto.Entity, error) {
	isPlaceholder := false
	if avset.IsPlaceholder != nil {
		isPlaceholder = *avset.IsPlaceholder
	}

	return &wssdcommonproto.Entity{
		IsPlaceholder: isPlaceholder,
	}, nil
}

func getWssdAvailabilitySetVMs(avset *compute.AvailabilitySet) ([]*wssdcommonproto.NodeSubResource, error) {
	var vms []*wssdcommonproto.NodeSubResource
	for _, vm := range avset.VirtualMachines {
		err := validateVM(vm)
		if err != nil {
			return nil, err
		}

		vms = append(vms, &wssdcommonproto.NodeSubResource{
			Name: *vm.Name,
		})
	}

	return vms, nil
}

func validateVM(vm *compute.SubResource) error {
	if vm.Name == nil || len(*vm.Name) == 0 {
		return errors.Wrapf(errors.InvalidInput, "avset member VM name cannot be empty")
	}

	// TODO: should we validate the VM exists in the system?
	return nil
}

// Conversion functions from wssdcompute to compute
func getAvailabilitySet(avset *wssdcompute.AvailabilitySet) *compute.AvailabilitySet {
	return &compute.AvailabilitySet{
		Name: &avset.Name,
		ID:   &avset.Id,
		Tags: getComputeTags(avset.GetTags()),
		AvailabilitySetProperties: &compute.AvailabilitySetProperties{
			PlatformFaultDomainCount: getAvailabilitySetPlatformFaultDomainCount(avset),
			VirtualMachines:          getAvailabilitySetVMs(avset),
			Statuses:                 getAvailabilitySetStatuses(avset),
			IsPlaceholder:            getAvailabilitySetIsPlaceholder(avset),
		},
	}
}

func getAvailabilitySetPlatformFaultDomainCount(avset *wssdcompute.AvailabilitySet) *int32 {
	return &avset.PlatformFaultDomainCount
}

func getAvailabilitySetVMs(avset *wssdcompute.AvailabilitySet) []*compute.SubResource {
	var vms []*compute.SubResource
	for _, vm := range avset.VirtualMachines {
		sr := compute.SubResource{
			Name: &vm.Name,
		}

		vms = append(vms, &sr)
	}

	return vms
}

func getAvailabilitySetStatuses(avset *wssdcompute.AvailabilitySet) map[string]*string {
	return status.GetStatuses(avset.Status)
}

func getComputeTags(tags *wssdcommonproto.Tags) map[string]*string {
	return prototags.ProtoToMap(tags)
}

func getAvailabilitySetIsPlaceholder(avset *wssdcompute.AvailabilitySet) *bool {
	isPlaceholder := false
	entity := avset.GetEntity()
	if entity != nil {
		isPlaceholder = entity.IsPlaceholder
	}
	return &isPlaceholder
}
