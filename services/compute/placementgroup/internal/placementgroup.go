package internal

import (
	"github.com/microsoft/moc/pkg/errors"
	"github.com/microsoft/moc/pkg/status"
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
	}

	if pgroup.PlacementGroupProperties == nil {
		// nothing else to convert
		return wssdpgroup, nil
	}

	pgType := wssdcompute.PlacementGroupType_Affinity
	if pgroup.Type == compute.Affinity {
		pgType = wssdcompute.PlacementGroupType_Affinity
	} else if pgroup.Type == compute.AntiAffinity {
		pgType = wssdcompute.PlacementGroupType_AntiAffinity
	} else if pgroup.Type == compute.StrictAntiAffinity {
		pgType = wssdcompute.PlacementGroupType_StrictAntiAffinity
	}

	pgScope := wssdcompute.PlacementGroupScope_Server
	if pgroup.Scope == compute.ZoneScope {
		pgScope = wssdcompute.PlacementGroupScope_Zone
	}

	wssdpgroup = &wssdcompute.PlacementGroup{
		Name:            *pgroup.Name,
		Entity:          getWssdPlacementGroupEntity(pgroup),
		VirtualMachines: getWssdPlacementGroupVMs(pgroup),
		Status:          status.GetFromStatuses(pgroup.Statuses),
		Type:            pgType,
		Scope:           pgScope,
	}

	if pgroup.PlacementGroupProperties != nil && pgroup.PlacementGroupProperties.Zones != nil {
		wssdpgroup.Zones = &wssdcommonproto.ZoneConfiguration{
			Zones:           []*wssdcommonproto.ZoneReference{},
			StrictPlacement: *pgroup.Zones.StrictPlacement,
		}

		for _, zn := range *pgroup.PlacementGroupProperties.Zones.Zones {
			rpcZoneRef, err := getRpcZoneReference(&zn)
			if err != nil {
				return nil, err
			}
			wssdpgroup.Zones.Zones = append(wssdpgroup.Zones.Zones, rpcZoneRef)
		}
	}

	return wssdpgroup, nil
}

func getRpcZoneReference(s *compute.ZoneReference) (*wssdcommonproto.ZoneReference, error) {
	if s == nil {
		return nil, errors.Wrapf(errors.InvalidInput, "Zone Name is missing")
	}

	return &wssdcommonproto.ZoneReference{
		Name: *s.Name,
	}, nil
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

	pgZoneRef := make([]compute.ZoneReference, len(pgroup.Zones.Zones))

	for i, zn := range pgroup.Zones.Zones {
		pgZnRef := compute.ZoneReference{
			Name:  &zn.Name,
			Nodes: &[]string{},
		}
		pgZoneRef[i] = pgZnRef
	}

	pgZones := &compute.ZoneConfiguration{
		Zones:           &pgZoneRef,
		StrictPlacement: &pgroup.Zones.StrictPlacement,
	}

	pgScope := compute.ServerScope
	if pgroup.Scope == wssdcompute.PlacementGroupScope_Zone {
		pgScope = compute.ZoneScope
	}

	pgType := compute.Affinity
	if pgroup.Type == wssdcompute.PlacementGroupType_Affinity {
		pgType = compute.Affinity
	} else if pgroup.Type == wssdcompute.PlacementGroupType_AntiAffinity {
		pgType = compute.AntiAffinity
	} else if pgroup.Type == wssdcompute.PlacementGroupType_StrictAntiAffinity {
		pgType = compute.StrictAntiAffinity
	}

	return &compute.PlacementGroup{
		Name: &pgroup.Name,
		ID:   &pgroup.Id,
		Type: pgType,
		PlacementGroupProperties: &compute.PlacementGroupProperties{
			VirtualMachines: getPlacementGroupVMs(pgroup),
			Statuses:        getPlacementGroupStatuses(pgroup),
			IsPlaceholder:   getPlacementGroupIsPlaceholder(pgroup),
			Zones:           pgZones,
			Scope:           pgScope,
		},
	}
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
