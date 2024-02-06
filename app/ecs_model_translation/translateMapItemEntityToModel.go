package ecs_model_translation

import (
	"github.com/hielkefellinger/go-dnd/app/ecs"
	"github.com/hielkefellinger/go-dnd/app/ecs_components"
	"github.com/hielkefellinger/go-dnd/app/models"
	"strconv"
)

func MapItemEntityToCampaignMapItemElement(rawMapItemComponent ecs.Component, mapId string) models.CampaignScreenMapItemElement {

	var image = ecs_components.NewMissingImageComponent()
	var controllingPlayers = make([]string, 0)

	// Translate; nil save exit
	mapItemComponent, ok := rawMapItemComponent.(*ecs_components.MapItemRelationComponent)
	if !ok || mapItemComponent == nil || mapItemComponent.Entity == nil {
		return models.CampaignScreenMapItemElement{
			Id:          rawMapItemComponent.GetId().String(),
			MapId:       mapId,
			Controllers: controllingPlayers,
			Image: models.CampaignImage{
				Name: image.Name,
				Url:  image.Url,
			},
		}
	}

	// Get (possible) Image
	var imageDetails = mapItemComponent.Entity.GetAllComponentsOfType(ecs.ImageComponentType)
	if imageDetails != nil && len(imageDetails) == 1 {
		image = imageDetails[0].(*ecs_components.ImageComponent)
	}

	// Get (all possible) controlling players
	var playerDetails = mapItemComponent.Entity.GetAllComponentsOfType(ecs.PlayerComponentType)
	if playerDetails != nil && len(playerDetails) > 0 {
		controllingPlayers = make([]string, len(playerDetails))
		for i, detail := range playerDetails {
			controllingPlayers[i] = detail.(*ecs_components.PlayerComponent).Name
		}
	}

	model := models.CampaignScreenMapItemElement{
		Id:          mapItemComponent.GetId().String(),
		EntityName:  mapItemComponent.Entity.GetName(),
		EntityId:    mapItemComponent.Entity.GetId().String(),
		MapId:       mapId,
		Controllers: controllingPlayers,
		Image: models.CampaignImage{
			Name: image.Name,
			Url:  image.Url,
		},
	}

	if mapItemComponent.Position != nil {
		model.Position = models.CampaignScreenMapPosition{
			X: strconv.Itoa(int(mapItemComponent.Position.X)),
			Y: strconv.Itoa(int(mapItemComponent.Position.Y)),
		}
	}

	return model
}
