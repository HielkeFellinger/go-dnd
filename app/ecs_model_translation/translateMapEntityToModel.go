package ecs_model_translation

import (
	"github.com/hielkefellinger/go-dnd/app/ecs"
	"github.com/hielkefellinger/go-dnd/app/ecs_components"
	"github.com/hielkefellinger/go-dnd/app/models"
)

func MapEntityToCampaignMapModel(mapEntity ecs.Entity) models.CampaignMap {
	// Get Image
	var image *ecs_components.ImageComponent
	var imageDetails = mapEntity.GetAllComponentsOfType(ecs.ImageComponentType)
	if imageDetails != nil && len(imageDetails) == 1 {
		image = imageDetails[0].(any).(*ecs_components.ImageComponent)
	}

	// Get Area
	var area *ecs_components.AreaComponent
	var areaDetails = mapEntity.GetAllComponentsOfType(ecs.AreaComponentType)
	if areaDetails != nil && len(areaDetails) == 1 {
		area = areaDetails[0].(any).(*ecs_components.AreaComponent)
	}

	return models.CampaignMap{
		Id:          mapEntity.GetId().String(),
		Name:        mapEntity.GetName(),
		Description: mapEntity.GetDescription(),
		X:           area.Width,
		Y:           area.Length,
		Image: models.CampaignMapImage{
			Name: image.Name,
			Url:  image.Url,
		},
	}
}
