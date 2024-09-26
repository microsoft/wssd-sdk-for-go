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
func getWssdPlacementGroup(pgroup *compute.PlacementGroup) (*wssdcompute.PlacementGroup, error) {
	errorPrefix := "Error converting PlacementGroup to WssdPlacementGroup"

	if pgroup == nil {
		return nil, errors.Wrapf(errors.InvalidInput, "%s, PlacementGroup cannot be nil", errorPrefix)
	}

	if pgroup.Name == nil {
		return nil, errors.Wrapf(errors.InvalidInput, "%s, Name is missing", errorPrefix)
	}

	wssdpgroup := &wssdcompute.PlacementGroup{
		Name: *pgroup.Name,
		Tags: getWssdTags(pgroup.Tags),
	}

	if pgroup.PlacementGroupProperties == nil {
		// nothing else to convert
		return wssdpgroup, nil
	}

	wssdpgroup = &wssdcompute.PlacementGroup{
		Name:                     *pgroup.Name,
		Tags:                     getWssdTags(pgroup.Tags),
		Entity:                   getWssdPlacementGroupEntity(pgroup),
		PlatformFaultDomainCount: getWssdPlatformFaultDomainCount(pgroup),
		VirtualMachines:          getWssdPlacementGroupVMs(pgroup),
		Status:                   status.GetFromStatuses(pgroup.Statuses),
	}

	return wssdpgroup, nil
}

func getWssdTags(tags map[string]*string) *wssdcommonproto.Tags {
	return prototags.MapToProto(tags)
}

func getWssdPlacementGroupEntity(pgroup *compute.PlacementGroup) *wssdcommonproto.Entity {
	isPlaceholder := false
	if pgroup.IsPlaceholder != nil {
		isPlaceholder = *pgroup.IsPlaceholder
	}

	return &wssdcommonproto.Entity{
		IsPlaceholder: isPlaceholder,
	}
}

func getWssdPlatformFaultDomainCount(pgroup *compute.PlacementGroup) int32 {
	var faultDomainCount int32 = 0
	if pgroup.PlatformFaultDomainCount != nil {
		faultDomainCount = *pgroup.PlatformFaultDomainCount
	}

	return faultDomainCount
}

func getWssdPlacementGroupVMs(pgroup *compute.PlacementGroup) []*wssdcompute.VirtualMachineReference {
	var vms []*wssdcompute.VirtualMachineReference
	for _, vm := range pgroup.VirtualMachines {
		if vm != nil && vm.Name != nil {
			vms = append(vms, &wssdcompute.VirtualMachineReference{
				Name: *vm.Name,
			})
		}
	}

	return vms
}

// Conversion functions from wssdcompute to compute
func getPlacementGroup(pgroup *wssdcompute.PlacementGroup) *compute.PlacementGroup {
	return &compute.PlacementGroup{
		Name: &pgroup.Name,
		ID:   &pgroup.Id,
		Tags: getComputeTags(pgroup.GetTags()),
		PlacementGroupProperties: &compute.PlacementGroupProperties{
			PlatformFaultDomainCount: getPlacementGroupPlatformFaultDomainCount(pgroup),
			VirtualMachines:          getPlacementGroupVMs(pgroup),
			Statuses:                 getPlacementGroupStatuses(pgroup),
			IsPlaceholder:            getPlacementGroupIsPlaceholder(pgroup),
		},
	}
}

func getComputeTags(tags *wssdcommonproto.Tags) map[string]*string {
	return prototags.ProtoToMap(tags)
}

func getPlacementGroupPlatformFaultDomainCount(pgroup *wssdcompute.PlacementGroup) *int32 {
	return &pgroup.PlatformFaultDomainCount
}

func getPlacementGroupVMs(pgroup *wssdcompute.PlacementGroup) []*compute.SubResource {
	var vms []*compute.SubResource
	for _, vm := range pgroup.VirtualMachines {
		sr := compute.SubResource{
			Name: &vm.Name,
		}

		vms = append(vms, &sr)
	}

	return vms
}

func getPlacementGroupStatuses(pgroup *wssdcompute.PlacementGroup) map[string]*string {
	return status.GetStatuses(pgroup.Status)
}

func getPlacementGroupIsPlaceholder(pgroup *wssdcompute.PlacementGroup) *bool {
	isPlaceholder := false
	entity := pgroup.GetEntity()
	if entity != nil {
		isPlaceholder = entity.IsPlaceholder
	}
	return &isPlaceholder
}
