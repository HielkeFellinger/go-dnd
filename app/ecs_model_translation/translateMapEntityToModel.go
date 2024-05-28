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

	// Get CampaignImage(s)
	images := make([]models.CampaignImage, 0)
	var image *ecs_components.ImageComponent
	var imageDetails = rawMapEntity.GetAllComponentsOfType(ecs.ImageComponentType)
	if imageDetails != nil && len(imageDetails) > 0 {
		// Loop to find active image and build image option; but set the first as default
		for index, imageDetail := range imageDetails {
			castedImage := imageDetail.(*ecs_components.ImageComponent)
			images = append(images, models.CampaignImage{
				Name:   castedImage.Name,
				Url:    castedImage.Url,
				Id:     castedImage.Id.String(),
				Active: castedImage.Active,
			})
			if index == 0 || castedImage.Active {
				image = castedImage
			}
		}
	}
	if image == nil {
		image = ecs_components.NewMissingImageComponent()
	}

	// Get Area
	var area *ecs_components.AreaComponent
	var areaDetails = rawMapEntity.GetAllComponentsOfType(ecs.AreaComponentType)
	if areaDetails != nil && len(areaDetails) == 1 {
		area = areaDetails[0].(*ecs_components.AreaComponent)
	}

	model := models.CampaignMap{
		Id:          rawMapEntity.GetId().String(),
		Name:        rawMapEntity.GetName(),
		Description: rawMapEntity.GetDescription(),
		ActiveImage: models.CampaignImage{Id: image.Id.String(), Name: image.Name, Url: image.Url},
		Images:      images,
	}

	if mapEntity != nil {
		model.Active = mapEntity.Active
	}

	if area != nil {
		model.X = area.Width
		model.Y = area.Length
	}

	return model
}
