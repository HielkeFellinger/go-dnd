package ecs_model_translation

import (
	"github.com/hielkefellinger/go-dnd/app/ecs"
	"github.com/hielkefellinger/go-dnd/app/ecs_components"
	"github.com/hielkefellinger/go-dnd/app/models"
	"strconv"
)

func InventoryEntityToCampaignInventoryModel(rawInventoryEntity ecs.Entity) models.CampaignInventory {

	inventory := &models.CampaignInventory{
		Name:              rawInventoryEntity.GetName(),
		Description:       rawInventoryEntity.GetDescription(),
		Size:              0,
		Slots:             "0",
		ShowDetailButtons: true,
		Items:             make([]models.CampaignInventoryItem, 0),
	}

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
				inventory.Items = append(inventory.Items, *inventoryItem)
			}
		}
	}

	return *inventory
}

func ItemEntityToCampaignInventoryItem(rawItemEntity ecs.Entity, count uint) *models.CampaignInventoryItem {
	inventoryItem := &models.CampaignInventoryItem{
		Id:    rawItemEntity.GetId().String(),
		Count: count,
	}

	// Check Item Details
	itemComponents := rawItemEntity.GetAllComponentsOfType(ecs.ItemComponentType)
	if len(itemComponents) >= 1 {
		itemComponent := itemComponents[0].(*ecs_components.ItemComponent)
		inventoryItem.Name = itemComponent.Name
		inventoryItem.Description = itemComponent.Description
	}

	// Check Restore
	restoreComponents := rawItemEntity.GetAllComponentsOfType(ecs.RestoreComponentType)
	if len(restoreComponents) >= 1 {
		inventoryItem.Restore = restoreComponents[0].(*ecs_components.RestoreComponent).Amount
	}

	// Check Damage
	damageComponents := rawItemEntity.GetAllComponentsOfType(ecs.DamageComponentType)
	if len(damageComponents) >= 1 {
		inventoryItem.Damage = damageComponents[0].(*ecs_components.DamageComponent).Amount
	}

	// Check Weight
	weightComponents := rawItemEntity.GetAllComponentsOfType(ecs.WeightComponentType)
	if len(weightComponents) >= 1 {
		inventoryItem.Weight = weightComponents[0].(*ecs_components.WeightComponent).Amount
	}

	// Check Range
	rangeComponents := rawItemEntity.GetAllComponentsOfType(ecs.RangeComponentType)
	if len(rangeComponents) >= 1 {
		rangeComponent := rangeComponents[0].(*ecs_components.RangeComponent)
		inventoryItem.Range = models.CampaignInventoryItemRange{
			Min: strconv.Itoa(rangeComponent.Min),
			Max: strconv.Itoa(rangeComponent.Max),
		}
	}
	return inventoryItem
}
