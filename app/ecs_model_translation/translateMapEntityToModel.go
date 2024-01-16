package ecs_model_translation

import (
	"github.com/hielkefellinger/go-dnd/app/ecs"
	"github.com/hielkefellinger/go-dnd/app/ecs_components"
	"github.com/hielkefellinger/go-dnd/app/models"
)

func MapEntityToCampaignMapModel(rawMapEntity ecs.Entity) models.CampaignMap {

	var mapEntity *ecs_components.MapComponent
	var mapDetails = rawMapEntity.GetAllComponentsOfType(ecs.MapComponentType)
	if mapDetails != nil && len(mapDetails) == 1 {
		mapEntity = mapDetails[0].(*ecs_components.MapComponent)
	}

	// Get CampaignImage
	var image *ecs_components.ImageComponent
	var imageDetails = rawMapEntity.GetAllComponentsOfType(ecs.ImageComponentType)
	if imageDetails != nil && len(imageDetails) == 1 {
		image = imageDetails[0].(*ecs_components.ImageComponent)
	}

	// Get Area
	var area *ecs_components.AreaComponent
	var areaDetails = rawMapEntity.GetAllComponentsOfType(ecs.AreaComponentType)
	if areaDetails != nil && len(areaDetails) == 1 {
		area = areaDetails[0].(*ecs_components.AreaComponent)
	}

	return models.CampaignMap{
		Id:          rawMapEntity.GetId().String(),
		Name:        rawMapEntity.GetName(),
		Description: rawMapEntity.GetDescription(),
		Enabled:     mapEntity.Active,
		X:           area.Width,
		Y:           area.Length,
		Image: models.CampaignImage{
			Name: image.Name,
			Url:  image.Url,
		},
	}
}
