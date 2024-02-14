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
// Field validations will occur in wssdagent
func getWssdAvailabilitySet(avset *compute.AvailabilitySet) (*wssdcompute.AvailabilitySet, error) {
	errorPrefix := "Error converting AvailabilitySet to WssdAvailabilitySet"

	if avset == nil {
		return nil, errors.Wrapf(errors.InvalidInput, "%s, AvailabilitySet cannot be nil", errorPrefix)
	}

	if avset.Name == nil {
		return nil, errors.Wrapf(errors.InvalidInput, "%s, Name is missing", errorPrefix)
	}

	wssdavset := &wssdcompute.AvailabilitySet{
		Name: *avset.Name,
		Tags: getWssdTags(avset.Tags),
	}

	if avset.AvailabilitySetProperties == nil {
		// nothing else to convert
		return wssdavset, nil
	}

	wssdavset = &wssdcompute.AvailabilitySet{
		Name:                     *avset.Name,
		Tags:                     getWssdTags(avset.Tags),
		Entity:                   getWssdAvailabilitySetEntity(avset),
		PlatformFaultDomainCount: getWssdPlatformFaultDomainCount(avset),
		VirtualMachines:          getWssdAvailabilitySetVMs(avset),
		Status:                   status.GetFromStatuses(avset.Statuses),
	}

	return wssdavset, nil
}

func getWssdTags(tags map[string]*string) *wssdcommonproto.Tags {
	return prototags.MapToProto(tags)
}

func getWssdAvailabilitySetEntity(avset *compute.AvailabilitySet) *wssdcommonproto.Entity {
	isPlaceholder := false
	if avset.IsPlaceholder != nil {
		isPlaceholder = *avset.IsPlaceholder
	}

	return &wssdcommonproto.Entity{
		IsPlaceholder: isPlaceholder,
	}
}

func getWssdPlatformFaultDomainCount(avset *compute.AvailabilitySet) int32 {
	var faultDomainCount int32 = 0
	if avset.PlatformFaultDomainCount != nil {
		faultDomainCount = *avset.PlatformFaultDomainCount
	}

	return faultDomainCount
}

func getWssdAvailabilitySetVMs(avset *compute.AvailabilitySet) []*wssdcommonproto.NodeSubResource {
	var vms []*wssdcommonproto.NodeSubResource
	for _, vm := range avset.VirtualMachines {
		if vm != nil && vm.Name != nil {
			vms = append(vms, &wssdcommonproto.NodeSubResource{
				Name: *vm.Name,
			})
		}
	}

	return vms
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

func getComputeTags(tags *wssdcommonproto.Tags) map[string]*string {
	return prototags.ProtoToMap(tags)
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

func getAvailabilitySetIsPlaceholder(avset *wssdcompute.AvailabilitySet) *bool {
	isPlaceholder := false
	entity := avset.GetEntity()
	if entity != nil {
		isPlaceholder = entity.IsPlaceholder
	}
	return &isPlaceholder
}
