package ecs_model_translation

import (
	"github.com/hielkefellinger/go-dnd/app/ecs"
	"github.com/hielkefellinger/go-dnd/app/ecs_components"
	"github.com/hielkefellinger/go-dnd/app/models"
	"strconv"
)

func CharacterEntityToCampaignCharacterModel(rawCharacterEntity ecs.Entity) models.CampaignCharacter {
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

		// Get Possible Image
		imagePlaceholder := ecs_components.NewMissingImageComponent()
		characterImages := rawCharacterEntity.GetAllComponentsOfType(ecs.ImageComponentType)
		if len(characterImages) >= 1 {
			imagePlaceholder = characterImages[0].(*ecs_components.ImageComponent)
		}
		character.Image = models.CampaignImage{
			Name: imagePlaceholder.Name,
			Url:  imagePlaceholder.Url,
		}

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

		// Check Ownership
		for _, rawPlayerComponent := range rawCharacterEntity.GetAllComponentsOfType(ecs.PlayerComponentType) {
			playerComponent := rawPlayerComponent.(*ecs_components.PlayerComponent)
			character.Controllers = append(character.Controllers, playerComponent.Name)
		}

		// Check for hasRelation to InventoryEntity
		hasRelationComponents := rawCharacterEntity.GetAllComponentsOfType(ecs.HasRelationComponentType)
		for _, rawHasRelationComponent := range hasRelationComponents {
			hasRelationComponent := rawHasRelationComponent.(*ecs_components.HasRelationComponent)
			// Test if relation is an Inventory
			if hasRelationComponent.Entity.HasComponentType(ecs.SlotsComponentType) {
				character.Inventories = append(character.Inventories,
					InventoryEntityToCampaignInventoryModel(hasRelationComponent.Entity))
			}

			// @todo stats?
		}

		// @todo spells & slots?
	}

	return character
}
