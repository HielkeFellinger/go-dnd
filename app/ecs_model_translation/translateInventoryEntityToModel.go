package ecs_model_translation

import (
	"github.com/hielkefellinger/go-dnd/app/ecs"
	"github.com/hielkefellinger/go-dnd/app/ecs_components"
	"github.com/hielkefellinger/go-dnd/app/models"
	"strconv"
)

func InventoryEntityToCampaignInventoryModel(rawInventoryEntity ecs.Entity) models.CampaignInventory {

	inventory := &models.CampaignInventory{
		Items: make([]models.CampaignInventoryItem, 0),
	}

	// Check if valid
	if rawInventoryEntity.HasComponentType(ecs.SlotsComponentType) {

		// Set Inv. Entity ID
		inventory.Id = rawInventoryEntity.GetId().String()

		// Loop over all the hasRelations and get the items
		rawHasRelations := rawInventoryEntity.GetAllComponentsOfType(ecs.HasRelationComponentType)
		for _, rawHasRelation := range rawHasRelations {

			hasRelation := rawHasRelation.(*ecs_components.HasRelationComponent)

			// Check if containing entity is an Item
			if hasRelation.Entity != nil && hasRelation.Entity.HasComponentType(ecs.ItemComponentType) {

				inventoryItem := &models.CampaignInventoryItem{
					Id: hasRelation.Entity.GetId().String(),
				}

				// Check Item Details
				itemComponents := hasRelation.Entity.GetAllComponentsOfType(ecs.ItemComponentType)
				if len(itemComponents) >= 1 {
					itemComponent := itemComponents[0].(*ecs_components.ItemComponent)
					inventoryItem.Name = itemComponent.Name
					inventoryItem.Description = itemComponent.Description
				}

				// Check Restore
				restoreComponents := hasRelation.Entity.GetAllComponentsOfType(ecs.RestoreComponentType)
				if len(restoreComponents) >= 1 {
					inventoryItem.Restore = restoreComponents[0].(*ecs_components.RestoreComponent).Amount
				}

				// Check Damage
				damageComponents := hasRelation.Entity.GetAllComponentsOfType(ecs.DamageComponentType)
				if len(damageComponents) >= 1 {
					inventoryItem.Damage = damageComponents[0].(*ecs_components.DamageComponent).Amount
				}

				// Check Range
				rangeComponents := hasRelation.Entity.GetAllComponentsOfType(ecs.RangeComponentType)
				if len(rangeComponents) >= 1 {
					rangeComponent := rangeComponents[0].(*ecs_components.RangeComponent)
					inventoryItem.Range = models.CampaignInventoryItemRange{
						Min: strconv.Itoa(rangeComponent.Min),
						Max: strconv.Itoa(rangeComponent.Max),
					}
				}

				inventory.Items = append(inventory.Items, *inventoryItem)
			}
		}
	}

	return *inventory
}
