package game_engine

import (
	"github.com/hielkefellinger/go-dnd/app/ecs"
	"github.com/hielkefellinger/go-dnd/app/ecs_components"
	"github.com/hielkefellinger/go-dnd/app/helpers"
	"golang.org/x/net/html"
	"slices"
	"strconv"
)

func upsertInventory(inventUpdateRequest inventoryUpsertRequest, pool CampaignPool) (ecs.BaseEntity, error) {
	var isNewEntry = inventUpdateRequest.Id == ""
	var inventoryEntity ecs.BaseEntity
	var rawInventoryEntity ecs.Entity = nil

	if inventUpdateRequest.Name == "" {
		return ecs.BaseEntity{}, SendManagementError("Error", "missing Required Inventory Name", pool)
	}

	if !isNewEntry {
		inventoryUuid, err := helpers.ParseStringToUuid(inventUpdateRequest.Id)
		if err != nil {
			return ecs.BaseEntity{}, err
		}
		inventoryEntityFromWorld, match := pool.GetEngine().GetWorld().GetInventoryEntityByUuid(inventoryUuid)
		if !match || inventoryEntityFromWorld == nil {
			return ecs.BaseEntity{}, SendManagementError("Error", "failure of loading Char by UUID", pool)
		}
		rawInventoryEntity = inventoryEntityFromWorld
		inventoryEntity = *inventoryEntityFromWorld.(*ecs.BaseEntity)
	} else {
		inventoryEntity = ecs.NewEntity()
		if addErr := inventoryEntity.AddComponent(ecs_components.NewInventoryComponent()); addErr != nil {
			return ecs.BaseEntity{}, SendManagementError("Error", addErr.Error(), pool)
		}
	}

	// Slots
	inventoryComponents := inventoryEntity.GetAllComponentsOfType(ecs.InventoryComponentType)
	for _, rawInventoryComponent := range inventoryComponents {
		inventoryComponent := rawInventoryComponent.(*ecs_components.InventoryComponent)
		if strconv.Itoa(int(inventoryComponent.Slots)) != inventUpdateRequest.Slots {
			if slotsErr := inventoryComponent.SlotsFromString(inventUpdateRequest.Slots); slotsErr != nil {
				return ecs.BaseEntity{}, SendManagementError("Error", slotsErr.Error(), pool)
			}
		}
	}

	// Update Basis Entity Properties
	if !isNewEntry && rawInventoryEntity != nil {
		rawInventoryEntity.SetName(html.EscapeString(inventUpdateRequest.Name))
		rawInventoryEntity.SetDescription(html.EscapeString(inventUpdateRequest.Description))
	} else {
		inventoryEntity.Name = html.EscapeString(inventUpdateRequest.Name)
		inventoryEntity.Description = html.EscapeString(inventUpdateRequest.Description)

		// Add to world
		if addErr := pool.GetEngine().GetWorld().AddEntity(&inventoryEntity); addErr != nil {
			return ecs.BaseEntity{}, SendManagementError("Error", addErr.Error(), pool)
		}
	}

	// Character Cleanup: add or remove
	for _, characterEntity := range pool.GetEngine().GetWorld().GetCharacterEntities() {
		if slices.Contains(inventUpdateRequest.Characters, characterEntity.GetId().String()) {
			// Only add if character has no relation to inventory
			if !characterEntity.HasRelationWithEntityByUuid(inventoryEntity.GetId()) {
				hasRelation := ecs_components.NewHasRelationComponent().(*ecs_components.HasRelationComponent)
				hasRelation.Entity = &inventoryEntity
				if err := characterEntity.AddComponent(hasRelation); err != nil {
					return ecs.BaseEntity{}, SendManagementError("Error", err.Error(), pool)
				}
			}
		} else if !isNewEntry && characterEntity.HasRelationWithEntityByUuid(inventoryEntity.GetId()) {
			// Remove if needed
			for _, hasRelation := range characterEntity.GetAllComponentsOfType(ecs.HasRelationComponentType) {
				hasRelationComponent := hasRelation.(*ecs_components.HasRelationComponent)
				if hasRelationComponent.Entity.GetId() == inventoryEntity.GetId() {
					characterEntity.RemoveComponentByUuid(hasRelationComponent.GetId())
				}
			}
		}
	}

	return inventoryEntity, nil
}

func upsertCharacter(charUpdateRequest characterUpsertRequest, pool CampaignPool) (ecs.BaseEntity, error) {
	var isNewEntry = charUpdateRequest.Id == ""
	var charEntity ecs.BaseEntity

	if charUpdateRequest.Name == "" {
		return ecs.BaseEntity{}, SendManagementError("Error", "missing Required Character Name", pool)
	}

	if !isNewEntry {
		charUuid, err := helpers.ParseStringToUuid(charUpdateRequest.Id)
		if err != nil {
			return ecs.BaseEntity{}, err
		}
		charEntityFromWorld, match := pool.GetEngine().GetWorld().GetCharacterEntityByUuid(charUuid)
		if !match || charEntityFromWorld == nil {
			return ecs.BaseEntity{}, SendManagementError("Error", "failure of loading Char by UUID", pool)
		}
		charEntityFromWorld.SetName(html.EscapeString(charUpdateRequest.Name))
		charEntityFromWorld.SetDescription(html.EscapeString(charUpdateRequest.Description))
		charEntity = *charEntityFromWorld.(*ecs.BaseEntity)
	} else {
		charEntity = ecs.NewEntity()
		charEntity.Name = html.EscapeString(charUpdateRequest.Name)
		charEntity.Description = html.EscapeString(charUpdateRequest.Description)
		if addErr := charEntity.AddComponent(ecs_components.NewCharacterComponent()); addErr != nil {
			return ecs.BaseEntity{}, SendManagementError("Error", addErr.Error(), pool)
		}
	}

	// Update Character Details
	if charDetails := charEntity.GetAllComponentsOfType(ecs.CharacterComponentType); len(charDetails) > 0 {
		charDetail := charDetails[0].(*ecs_components.CharacterComponent)
		charDetail.Name = html.EscapeString(charUpdateRequest.Name)
		charDetail.Description = html.EscapeString(charUpdateRequest.Description)
	}

	// Add Health
	if charUpdateRequest.HealthDamage != "" && charUpdateRequest.HealthTmp != "" && charUpdateRequest.HealthMax != "" {
		health := charEntity.GetAllComponentsOfType(ecs.HealthComponentType)
		if isNewEntry || len(health) == 0 {
			if addErr := charEntity.AddComponent(ecs_components.NewHealthComponent()); addErr != nil {
				return ecs.BaseEntity{}, SendManagementError("Error", addErr.Error(), pool)
			}
			health = charEntity.GetAllComponentsOfType(ecs.HealthComponentType)
		}
		if len(health) > 0 {
			if charUpdateRequest.HealthDamage == "" {
				charUpdateRequest.HealthDamage = "0"
			}
			if charUpdateRequest.HealthTmp == "" {
				charUpdateRequest.HealthTmp = "0"
			}
			if charUpdateRequest.HealthMax == "" {
				charUpdateRequest.HealthMax = "0"
			}
			healthComponent := health[0].(*ecs_components.HealthComponent)
			if updateErr := healthComponent.MaximumFromString(charUpdateRequest.HealthDamage); updateErr != nil {
				return ecs.BaseEntity{}, SendManagementError("Error", updateErr.Error(), pool)
			}
			if updateErr := healthComponent.MaximumFromString(charUpdateRequest.HealthTmp); updateErr != nil {
				return ecs.BaseEntity{}, SendManagementError("Error", updateErr.Error(), pool)
			}
			if updateErr := healthComponent.MaximumFromString(charUpdateRequest.HealthMax); updateErr != nil {
				return ecs.BaseEntity{}, SendManagementError("Error", updateErr.Error(), pool)
			}
		}
	} else {
		if health := charEntity.GetAllComponentsOfType(ecs.HealthComponentType); len(health) > 0 {
			for _, healthComponent := range health {
				charEntity.RemoveComponentByUuid(healthComponent.GetId())
			}
		}
	}

	// Update Level
	if charUpdateRequest.Level != "" {
		levels := charEntity.GetAllComponentsOfType(ecs.LevelComponentType)
		if isNewEntry || len(levels) == 0 {
			if addErr := charEntity.AddComponent(ecs_components.NewLevelComponent()); addErr != nil {
				return ecs.BaseEntity{}, SendManagementError("Error", addErr.Error(), pool)
			}
			levels = charEntity.GetAllComponentsOfType(ecs.LevelComponentType)
		}
		if len(levels) > 0 {
			level := levels[0].(*ecs_components.LevelComponent)
			if updateErr := level.LevelFromString(charUpdateRequest.Level); updateErr != nil {
				return ecs.BaseEntity{}, SendManagementError("Error", updateErr.Error(), pool)
			}
		}
	} else {
		if levels := charEntity.GetAllComponentsOfType(ecs.LevelComponentType); len(levels) > 0 {
			for _, level := range levels {
				charEntity.RemoveComponentByUuid(level.GetId())
			}
		}
	}

	// Player Playable
	if charUpdateRequest.PlayerPlayable {
		if playerRelations := charEntity.GetAllComponentsOfType(ecs.PlayerComponentType); isNewEntry || len(playerRelations) == 0 {
			if addErr := charEntity.AddComponent(ecs_components.NewPlayerComponent()); addErr != nil {
				return ecs.BaseEntity{}, SendManagementError("Error", addErr.Error(), pool)
			}
		}
	} else {
		if playerRelations := charEntity.GetAllComponentsOfType(ecs.PlayerComponentType); len(playerRelations) > 0 {
			for _, playerRelation := range playerRelations {
				charEntity.RemoveComponentByUuid(playerRelation.GetId())
			}
		}
	}

	// Player Hidden
	hidden := charEntity.GetAllComponentsOfType(ecs.VisibilityComponentType)
	if charUpdateRequest.Hidden {
		if len(hidden) > 0 {
			visibility := hidden[0].(*ecs_components.VisibilityComponent)
			visibility.Hidden = true
		} else {
			visibility := ecs_components.NewVisibilityComponent().(*ecs_components.VisibilityComponent)
			visibility.Hidden = true
			if addErr := charEntity.AddComponent(visibility); addErr != nil {
				return ecs.BaseEntity{}, SendManagementError("Error", addErr.Error(), pool)
			}
		}
	} else {
		for _, hiddenComponent := range hidden {
			charEntity.RemoveComponentByUuid(hiddenComponent.GetId())
		}
	}

	// Update Image if needed
	if charUpdateRequest.Image != (helpers.FileUpload{}) && charUpdateRequest.ImageName != "" {
		// Handle file-upload
		link, imageErr := helpers.SaveImageToCampaign(charUpdateRequest.Image, pool.GetId(), charUpdateRequest.ImageName)
		if imageErr != nil {
			return ecs.BaseEntity{}, SendManagementError("Error", imageErr.Error(), pool)
		}
		imageComponent := ecs_components.NewImageComponent().(*ecs_components.ImageComponent)
		imageComponent.Name = html.EscapeString(charUpdateRequest.ImageName)
		imageComponent.Active = isNewEntry
		imageComponent.Url = link
		imageComponent.Version = 1
		if addErr := charEntity.AddComponent(imageComponent); addErr != nil {
			return ecs.BaseEntity{}, SendManagementError("Error", addErr.Error(), pool)
		}
	}

	// Add Char Entity
	if isNewEntry {
		if addErr := pool.GetEngine().GetWorld().AddEntity(&charEntity); addErr != nil {
			return ecs.BaseEntity{}, SendManagementError("Error", addErr.Error(), pool)
		}
	}

	return charEntity, nil
}

func upsertMap(mapUpdateRequest mapUpsertRequest, pool CampaignPool) (ecs.BaseEntity, error) {
	var isNewEntry = mapUpdateRequest.Id == ""
	var mapEntity ecs.BaseEntity
	var rawMapEntity ecs.Entity = nil

	if mapUpdateRequest.Name == "" {
		return ecs.BaseEntity{}, SendManagementError("Error", "missing Required Map Name", pool)
	}

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
		imageComponent.Name = html.EscapeString(mapUpdateRequest.ImageName)
		imageComponent.Active = isNewEntry
		imageComponent.Url = link
		imageComponent.Version = 1
		if addErr := mapEntity.AddComponent(imageComponent); addErr != nil {
			return ecs.BaseEntity{}, SendManagementError("Error", addErr.Error(), pool)
		}
	} else if isNewEntry {
		return ecs.BaseEntity{}, SendManagementError("Error", "missing image on new Map", pool)
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
	if !isNewEntry {
		rawMapEntity.SetName(html.EscapeString(mapUpdateRequest.Name))
		rawMapEntity.SetDescription(html.EscapeString(mapUpdateRequest.Description))
	} else {
		mapEntity.SetName(html.EscapeString(mapUpdateRequest.Name))
		mapEntity.SetDescription(html.EscapeString(mapUpdateRequest.Description))
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
		return ecs.BaseEntity{}, SendManagementError("Error", "missing Required Item Name", pool)
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
		itemComponent.Name = html.EscapeString(itemUpsertRequest.Name)
		itemComponent.Description = html.EscapeString(itemUpsertRequest.Description)
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
		weightComponent.Amount = html.EscapeString(itemUpsertRequest.Weight)
	}

	// Add new Item Entity
	if isNewEntity {
		if addErr := pool.GetEngine().GetWorld().AddEntity(&itemEntity); addErr != nil {
			return ecs.BaseEntity{}, SendManagementError("Error", addErr.Error(), pool)
		}
	}

	return itemEntity, nil
}
