package ecs_model_translation

import (
	"github.com/hielkefellinger/go-dnd/app/ecs"
	"github.com/hielkefellinger/go-dnd/app/ecs_components"
	"github.com/hielkefellinger/go-dnd/app/models"
	"strconv"
)

func InventoryEntityToCampaignInventoryModel(rawInventoryEntity ecs.Entity) models.CampaignInventory {

	inventory := models.GetCampaignInventory()
	inventory.Name = rawInventoryEntity.GetName()
	inventory.Description = rawInventoryEntity.GetDescription()

	// Check if valid
	if rawInventoryEntity.HasComponentType(ecs.InventoryComponentType) {

		// Set Inv. Entity ID
		inventory.Id = rawInventoryEntity.GetId().String()

		inventoryEntities := rawInventoryEntity.GetAllComponentsOfType(ecs.InventoryComponentType)
		for _, inventoryEntity := range inventoryEntities {
			inventory.Slots = strconv.Itoa(int(inventoryEntity.(*ecs_components.InventoryComponent).Slots))
			break
		}

		// Loop over all the hasRelations and get the items
		rawHasRelations := rawInventoryEntity.GetAllComponentsOfType(ecs.HasRelationComponentType)
		for _, rawHasRelation := range rawHasRelations {
			hasRelation := rawHasRelation.(*ecs_components.HasRelationComponent)

			// Check if containing entity is an Item
			if hasRelation.Entity != nil && hasRelation.Entity.HasComponentType(ecs.ItemComponentType) {
				inventoryItem := ItemEntityToCampaignInventoryItem(hasRelation.Entity, hasRelation.Count)
				inventory.Size += inventoryItem.Count
				inventory.Items = append(inventory.Items, inventoryItem)
			}
		}
	}

	return inventory
}
