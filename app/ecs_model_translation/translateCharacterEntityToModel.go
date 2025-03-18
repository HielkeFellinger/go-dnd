package ecs_model_translation

import (
	"github.com/hielkefellinger/go-dnd/app/ecs"
	"github.com/hielkefellinger/go-dnd/app/ecs_components"
	"github.com/hielkefellinger/go-dnd/app/models"
	"strconv"
)

type CharModelType int

const (
	DEFAULT   CharModelType = 0
	ALL       CharModelType = 1
	INVENTORY CharModelType = 2
	STATUS    CharModelType = 3
)

func CharacterEntityToCampaignCharacterModel(rawCharacterEntity ecs.Entity, mode CharModelType) models.CampaignCharacter {
	character := models.GetNewCampaignCharacter()

	// Test If it is a character
	if rawCharacterEntity.HasComponentType(ecs.CharacterComponentType) {

		// Set the ID
		character.Id = rawCharacterEntity.GetId().String()

		// Get Name and Description
		characterComponents := rawCharacterEntity.GetAllComponentsOfType(ecs.CharacterComponentType)
		if len(characterComponents) >= 1 {
			characterComponent := characterComponents[0].(*ecs_components.CharacterComponent)
			character.Name = characterComponent.Name
			character.Description = characterComponent.Description
		}

		// Get Possible Images
		characterImages := rawCharacterEntity.GetAllComponentsOfType(ecs.ImageComponentType)
		for index, image := range characterImages {
			currentImage := image.(*ecs_components.ImageComponent)
			charImage := models.CampaignImage{
				Id:   currentImage.Id.String(),
				Name: currentImage.Name,
				Url:  currentImage.Url,
			}

			if index == 0 || currentImage.Active {
				character.Image = charImage
			}
			character.Images = append(character.Images, charImage)
		}

		// Check Ownership
		for _, rawPlayerComponent := range rawCharacterEntity.GetAllComponentsOfType(ecs.PlayerComponentType) {
			playerComponent := rawPlayerComponent.(*ecs_components.PlayerComponent)
			character.Controllers = append(character.Controllers, playerComponent.Name)
		}

		// Check Total Level
		for _, rawLevelComponent := range rawCharacterEntity.GetAllComponentsOfType(ecs.LevelComponentType) {
			levelComponent := rawLevelComponent.(*ecs_components.LevelComponent)
			character.Level = strconv.Itoa(int(levelComponent.Level))
		}

		if mode == ALL || mode == DEFAULT {
			// Check for Health
			healthComponents := rawCharacterEntity.GetAllComponentsOfType(ecs.HealthComponentType)
			if len(healthComponents) >= 1 {
				healthComponent := healthComponents[0].(*ecs_components.HealthComponent)
				character.Health = models.CampaignCharacterHealth{
					Damage:             strconv.Itoa(int(healthComponent.Damage)),
					TemporaryHitPoints: strconv.Itoa(int(healthComponent.Temporary)),
					MaximumHitPoints:   strconv.Itoa(int(healthComponent.Maximum)),
				}
			}

			// Check Total Visibility
			for _, rawVisibilityComponent := range rawCharacterEntity.GetAllComponentsOfType(ecs.VisibilityComponentType) {
				visibilityComponent := rawVisibilityComponent.(*ecs_components.VisibilityComponent)
				character.Hidden = visibilityComponent.Hidden
			}
		}

		if mode == INVENTORY || mode == ALL {
			// Check for hasRelation to InventoryEntity
			hasRelationComponents := rawCharacterEntity.GetAllComponentsOfType(ecs.HasRelationComponentType)
			for _, rawHasRelationComponent := range hasRelationComponents {
				hasRelationComponent := rawHasRelationComponent.(*ecs_components.HasRelationComponent)
				// Test if relation is an Inventory
				if hasRelationComponent.Entity.HasComponentType(ecs.InventoryComponentType) {
					character.Inventories = append(character.Inventories,
						InventoryEntityToCampaignInventoryModel(hasRelationComponent.Entity))
				}
				// @todo stats?
			}

			// Check and update linked inventories (Used in trading)
			for id, inventory := range character.Inventories {
				linkedInv := models.CampaignLinkedInventory{
					Id:          inventory.Id,
					Name:        inventory.Name,
					Description: inventory.Description,
				}
				for index, _ := range character.Inventories {
					if index != id {
						character.Inventories[index].LinkedInventories =
							append(character.Inventories[index].LinkedInventories, linkedInv)
					}
				}
			}
		}

		if mode == STATUS {
			// @todo spells & slots?
		}
	}

	return character
}
