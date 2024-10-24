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
			Zones: []*wssdcommonproto.ZoneReference{},
            StrictPlacement: pgroup.PlacementGroupProperties.StrictPlacement,
		}

		for _, zn := range *pgroup.PlacementGroupProperties.Zones {
            rpcZoneRef, err := getRpcZoneReference(&zn)
			if err != nil {
				return nil, err
			}
			wssdpgroup.Zones.Zones = append(wssdpgroup.Zones.Zones, rpcZoneRef)
		}
	}

	return wssdpgroup, nil
}

func getRpcZoneReference(s *string) (*wssdcommonproto.ZoneReference, error) {
	if s == nil {
		return nil, errors.Wrapf(errors.InvalidInput, "Zone Name is missing")
	}

	return &wssdcommonproto.ZoneReference{
        Name: *s,
	},nil
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

	pgZone := []string{}
	for _, zn := range pgroup.Zones.Zones {
		pgZone = append(pgZone, zn.Name) 
	} 
	
	pgScope := compute.ServerScope
    if pgroup.Scope == wssdcompute.PlacementGroupScope_Zone {
		pgScope = compute.ZoneScope
	}

	return &compute.PlacementGroup{
		Name: &pgroup.Name,
		ID:   &pgroup.Id,
		PlacementGroupProperties: &compute.PlacementGroupProperties{
			VirtualMachines: getPlacementGroupVMs(pgroup),
			Statuses:        getPlacementGroupStatuses(pgroup),
			IsPlaceholder:   getPlacementGroupIsPlaceholder(pgroup),
			Zones:           &pgZone,
			Scope: 			 pgScope,
			StrictPlacement: pgroup.Zones.StrictPlacement,
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
