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

	// Get (possible) ActiveImage
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

	// Get (possible) Visibility
	hidden := false
	var visibilityDetails = mapItemComponent.Entity.GetAllComponentsOfType(ecs.VisibilityComponentType)
	if visibilityDetails != nil && len(visibilityDetails) == 1 {
		hidden = visibilityDetails[0].(*ecs_components.VisibilityComponent).Hidden
	}

	model := models.CampaignScreenMapItemElement{
		Id:          mapItemComponent.GetId().String(),
		EntityName:  mapItemComponent.Entity.GetName(),
		EntityId:    mapItemComponent.Entity.GetId().String(),
		MapId:       mapId,
		Hidden:      hidden,
		Controllers: controllingPlayers,
		Image: models.CampaignImage{
			Name: image.Name,
			Url:  image.Url,
		},
	}

	// Get (Possible Health Info
	var healthDetails = mapItemComponent.Entity.GetAllComponentsOfType(ecs.HealthComponentType)
	if healthDetails != nil && len(healthDetails) > 0 {
		healthComponent := healthDetails[0].(*ecs_components.HealthComponent)
		model.Health = models.CampaignScreenMapItemHealth{
			Total:   healthComponent.Temporary + healthComponent.Maximum,
			Current: int(healthComponent.Temporary) + int(healthComponent.Maximum) - int(healthComponent.Damage),
		}
	}

	if mapItemComponent.Position != nil {
		model.Position = models.CampaignScreenMapPosition{
			X: strconv.Itoa(int(mapItemComponent.Position.X)),
			Y: strconv.Itoa(int(mapItemComponent.Position.Y)),
		}
	}

	return model
}
