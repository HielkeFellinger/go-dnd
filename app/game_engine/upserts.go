package game_engine

import (
	"github.com/hielkefellinger/go-dnd/app/ecs"
	"github.com/hielkefellinger/go-dnd/app/ecs_components"
	"github.com/hielkefellinger/go-dnd/app/helpers"
	"log"
	"slices"
)

func upsertMap(mapUpdateRequest mapUpsertRequest, pool CampaignPool) (ecs.BaseEntity, error) {
	var isNewEntry = mapUpdateRequest.Id == ""
	var mapEntity ecs.BaseEntity
	var rawMapEntity ecs.Entity = nil

	if !isNewEntry {
		mapUuid, err := helpers.ParseStringToUuid(mapUpdateRequest.Id)
		if err != nil {
			return ecs.BaseEntity{}, err
		}
		mapEntityFromWorld, match := pool.GetEngine().GetWorld().GetMapEntityByUuid(mapUuid)
		if !match || mapEntityFromWorld == nil {
			return ecs.BaseEntity{}, SendManagementError("Error", "failure of loading MAP by UUID", pool)
		}
		rawMapEntity = mapEntityFromWorld
		mapEntity = *mapEntityFromWorld.(*ecs.BaseEntity)
	} else {
		mapEntity = ecs.NewEntity()
		if addErr := mapEntity.AddComponent(ecs_components.NewMapComponent()); addErr != nil {
			return ecs.BaseEntity{}, SendManagementError("Error", addErr.Error(), pool)
		}
		if addErr := mapEntity.AddComponent(ecs_components.NewAreaComponent()); addErr != nil {
			return ecs.BaseEntity{}, SendManagementError("Error", addErr.Error(), pool)
		}
	}

	// Parse Area
	mapAreaComponent := ecs_components.NewAreaComponent().(*ecs_components.AreaComponent)
	if mapUpdateRequest.X != "" {
		if areaError := mapAreaComponent.WidthFromString(mapUpdateRequest.X); areaError != nil {
			return ecs.BaseEntity{}, SendManagementError("Error", areaError.Error(), pool)
		}
		if mapAreaComponent.Width < 2 {
			return ecs.BaseEntity{}, SendManagementError("Error", "Invalid width (X). Should be > 2", pool)
		}
	}
	if mapUpdateRequest.Y != "" {
		if areaError := mapAreaComponent.LengthFromString(mapUpdateRequest.Y); areaError != nil {
			return ecs.BaseEntity{}, SendManagementError("Error", areaError.Error(), pool)
		}
		if mapAreaComponent.Length < 2 {
			return ecs.BaseEntity{}, SendManagementError("Error", "Invalid length (Y). Should be > 2", pool)
		}
	}
	if isNewEntry && (mapAreaComponent.Length < 2 || mapAreaComponent.Width < 2) {
		return ecs.BaseEntity{}, SendManagementError("Error", "Invalid length and or width set. Should be > 2", pool)
	}

	// Block editing of active maps.
	mapComponents := mapEntity.GetAllComponentsOfType(ecs.MapComponentType)
	if len(mapComponents) > 1 {
		mapComponent := mapComponents[0].(*ecs_components.MapComponent)
		if mapComponent.Active {
			return ecs.BaseEntity{}, SendManagementError("Error", "Can not modify active maps", pool)
		}
	}

	// Update if add image
	if mapUpdateRequest.Image != (helpers.FileUpload{}) && mapUpdateRequest.ImageName != "" {
		// Handle file-upload
		link, imageErr := helpers.SaveImageToCampaign(mapUpdateRequest.Image, pool.GetId(), mapUpdateRequest.ImageName)
		if imageErr != nil {
			return ecs.BaseEntity{}, SendManagementError("Error", imageErr.Error(), pool)
		}
		imageComponent := ecs_components.NewImageComponent().(*ecs_components.ImageComponent)
		imageComponent.Name = mapUpdateRequest.ImageName
		imageComponent.Active = isNewEntry
		imageComponent.Url = link
		imageComponent.Version = 1
		if addErr := mapEntity.AddComponent(imageComponent); addErr != nil {
			return ecs.BaseEntity{}, SendManagementError("Error", addErr.Error(), pool)
		}
	} else if isNewEntry {
		return ecs.BaseEntity{}, SendManagementError("Error", "Missing image on new Map", pool)
	}

	// Update the needed removal of images
	if !isNewEntry && len(mapUpdateRequest.RemoveImages) > 0 {
		var imageComponents = mapEntity.GetAllComponentsOfType(ecs.ImageComponentType)
		if len(imageComponents) <= len(mapUpdateRequest.RemoveImages) {
			return ecs.BaseEntity{}, SendManagementError("Error", "Can not delete all images; needs at least one", pool)
		}

		for index := range imageComponents {
			imageComponent := imageComponents[index].(*ecs_components.ImageComponent)
			if slices.Contains(mapUpdateRequest.RemoveImages, imageComponent.GetId().String()) {
				mapEntity.RemoveComponentByUuid(imageComponent.GetId())
			}
		}

		// Reload minus removed; set at least one of the leftover backgrounds to active
		imageComponents = mapEntity.GetAllComponentsOfType(ecs.ImageComponentType)
		var firstImageComponent *ecs_components.ImageComponent
		var hasActive = false
		for index := range imageComponents {
			imageComponent := imageComponents[index].(*ecs_components.ImageComponent)
			if index == 0 {
				firstImageComponent = imageComponent
			}
			hasActive = hasActive || imageComponent.Active
		}
		if !hasActive {
			firstImageComponent.Active = true
		}
	}

	// Update Basis Entity Properties
	if mapUpdateRequest.Name != "" && mapUpdateRequest.Name != mapEntity.Name {
		log.Printf("Update Name")
		mapEntity.SetName(mapUpdateRequest.Name)
		if !isNewEntry && rawMapEntity != nil {
			rawMapEntity.SetName(mapUpdateRequest.Name)
		}
	} else if isNewEntry {
		return ecs.BaseEntity{}, SendManagementError("Error", "Map should have a name", pool)
	}
	if mapUpdateRequest.Description != "" && mapUpdateRequest.Description != mapEntity.Description {
		log.Printf("Update Description")
		mapEntity.SetDescription(mapUpdateRequest.Description)
		if !isNewEntry && rawMapEntity != nil {
			rawMapEntity.SetDescription(mapUpdateRequest.Description)
		}
	}

	// Update Area
	if mapAreaComponent.Length >= 2 || mapAreaComponent.Width >= 2 {
		areaComponent := mapEntity.GetAllComponentsOfType(ecs.AreaComponentType)[0].(*ecs_components.AreaComponent)
		if mapAreaComponent.Length != areaComponent.Length && mapAreaComponent.Length >= 2 {
			areaComponent.Length = mapAreaComponent.Length
		}
		if mapAreaComponent.Width != areaComponent.Width && mapAreaComponent.Width >= 2 {
			areaComponent.Width = mapAreaComponent.Width
		}
	}

	// Add new Map Entity
	if isNewEntry {
		if addErr := pool.GetEngine().GetWorld().AddEntity(&mapEntity); addErr != nil {
			return ecs.BaseEntity{}, SendManagementError("Error", addErr.Error(), pool)
		}
	}

	return mapEntity, nil
}

func upsertItem(itemUpsertRequest itemUpsertRequest, pool CampaignPool) (ecs.BaseEntity, error) {
	var isNewEntity = itemUpsertRequest.Id == ""
	var itemEntity ecs.BaseEntity

	// Required Fields
	if itemUpsertRequest.Name == "" {
		return ecs.BaseEntity{}, SendManagementError("Error", "name is required for Item", pool)
	}

	if !isNewEntity {
		mapUuid, err := helpers.ParseStringToUuid(itemUpsertRequest.Id)
		if err != nil {
			return ecs.BaseEntity{}, err
		}
		itemEntityFromWorld, match := pool.GetEngine().GetWorld().GetItemEntityByUuid(mapUuid)
		if !match || itemEntityFromWorld == nil {
			return ecs.BaseEntity{}, SendManagementError("Error", "failure of loading Item by UUID", pool)
		}
		itemEntity = *itemEntityFromWorld.(*ecs.BaseEntity)
	} else {
		itemEntity = ecs.NewEntity()
		if addErr := itemEntity.AddComponent(ecs_components.NewItemComponent()); addErr != nil {
			return ecs.BaseEntity{}, SendManagementError("Error", addErr.Error(), pool)
		}
	}

	// Clear fields
	removeComponents := make([]ecs.Component, 0)
	if components := itemEntity.GetAllComponentsOfType(ecs.DamageComponentType); itemUpsertRequest.Damage == "" && len(components) > 0 {
		removeComponents = append(removeComponents, components...)
	}
	if components := itemEntity.GetAllComponentsOfType(ecs.RestoreComponentType); itemUpsertRequest.Restore == "" && len(components) > 0 {
		removeComponents = append(removeComponents, components...)
	}
	if components := itemEntity.GetAllComponentsOfType(ecs.RangeComponentType); itemUpsertRequest.RangeMin == "" && itemUpsertRequest.RangeMax == "" && len(components) > 0 {
		removeComponents = append(removeComponents, components...)
	}
	if components := itemEntity.GetAllComponentsOfType(ecs.WeightComponentType); itemUpsertRequest.Weight == "" && len(components) > 0 {
		removeComponents = append(removeComponents, components...)
	}
	// Clean
	for _, component := range removeComponents {
		itemEntity.RemoveComponentByUuid(component.GetId())
	}

	// Upsert Fields
	if components := itemEntity.GetAllComponentsOfType(ecs.ItemComponentType); len(components) > 0 {
		itemComponent := components[0].(*ecs_components.ItemComponent)
		itemComponent.Name = itemUpsertRequest.Name
		itemComponent.Description = itemUpsertRequest.Description
	}
	if itemUpsertRequest.Damage != "" {
		components := itemEntity.GetAllComponentsOfType(ecs.DamageComponentType)
		var damageComponent *ecs_components.DamageComponent
		if len(components) == 0 {
			damageComponent = ecs_components.NewDamageComponent().(*ecs_components.DamageComponent)
			if addErr := itemEntity.AddComponent(damageComponent); addErr != nil {
				return ecs.BaseEntity{}, SendManagementError("Error", addErr.Error(), pool)
			}
		} else {
			damageComponent = components[0].(*ecs_components.DamageComponent)
		}
		damageComponent.Amount = itemUpsertRequest.Damage
	}
	if itemUpsertRequest.Restore != "" {
		components := itemEntity.GetAllComponentsOfType(ecs.RestoreComponentType)
		var restoreComponent *ecs_components.RestoreComponent
		if len(components) == 0 {
			restoreComponent = ecs_components.NewRestoreComponent().(*ecs_components.RestoreComponent)
			if addErr := itemEntity.AddComponent(restoreComponent); addErr != nil {
				return ecs.BaseEntity{}, SendManagementError("Error", addErr.Error(), pool)
			}
		} else {
			restoreComponent = components[0].(*ecs_components.RestoreComponent)
		}
		restoreComponent.Amount = itemUpsertRequest.Damage
	}
	if itemUpsertRequest.RangeMin != "" || itemUpsertRequest.RangeMax != "" {
		components := itemEntity.GetAllComponentsOfType(ecs.RangeComponentType)
		var restoreComponent *ecs_components.RangeComponent
		if len(components) == 0 {
			restoreComponent = ecs_components.NewRangeComponent().(*ecs_components.RangeComponent)
			if addErr := itemEntity.AddComponent(restoreComponent); addErr != nil {
				return ecs.BaseEntity{}, SendManagementError("Error", addErr.Error(), pool)
			}
		} else {
			restoreComponent = components[0].(*ecs_components.RangeComponent)
		}
		_ = restoreComponent.MinFromString(itemUpsertRequest.RangeMin)
		_ = restoreComponent.MaxFromString(itemUpsertRequest.RangeMax)
	}
	if itemUpsertRequest.Weight != "" {
		components := itemEntity.GetAllComponentsOfType(ecs.WeightComponentType)
		var weightComponent *ecs_components.WeightComponent
		if len(components) == 0 {
			weightComponent = ecs_components.NewWeightComponent().(*ecs_components.WeightComponent)
			if addErr := itemEntity.AddComponent(weightComponent); addErr != nil {
				return ecs.BaseEntity{}, SendManagementError("Error", addErr.Error(), pool)
			}
		} else {
			weightComponent = components[0].(*ecs_components.WeightComponent)
		}
		weightComponent.Amount = itemUpsertRequest.Weight
	}

	// Add new Item Entity
	if isNewEntity {
		if addErr := pool.GetEngine().GetWorld().AddEntity(&itemEntity); addErr != nil {
			return ecs.BaseEntity{}, SendManagementError("Error", addErr.Error(), pool)
		}
	}

	return itemEntity, nil
}
