package ecs_model_translation

import (
	"github.com/hielkefellinger/go-dnd/app/ecs"
	"github.com/hielkefellinger/go-dnd/app/ecs_components"
	"github.com/hielkefellinger/go-dnd/app/models"
	"strconv"
)

func ItemEntityToCampaignInventoryItem(rawItemEntity ecs.Entity, count uint) models.CampaignInventoryItem {
	inventoryItem := models.GetCampaignInventoryItem()
	inventoryItem.Id = rawItemEntity.GetId().String()
	inventoryItem.Count = count

	// Check Item Details
	itemComponents := rawItemEntity.GetAllComponentsOfType(ecs.ItemComponentType)
	if len(itemComponents) >= 1 {
		itemComponent := itemComponents[0].(*ecs_components.ItemComponent)
		inventoryItem.Name = itemComponent.Name
		inventoryItem.Description = itemComponent.Description
	}

	// Check Images
	imagesComponents := rawItemEntity.GetAllComponentsOfType(ecs.ImageComponentType)
	for index, image := range imagesComponents {
		currentImage := image.(*ecs_components.ImageComponent)
		charImage := models.CampaignImage{
			Id:   currentImage.Id.String(),
			Name: currentImage.Name,
			Url:  currentImage.Url,
		}

		if index == 0 || currentImage.Active {
			inventoryItem.Image = charImage
		}
		inventoryItem.Images = append(inventoryItem.Images, charImage)
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
