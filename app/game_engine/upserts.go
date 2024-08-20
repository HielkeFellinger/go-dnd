package game_engine

import (
	"github.com/hielkefellinger/go-dnd/app/ecs"
	"github.com/hielkefellinger/go-dnd/app/ecs_components"
	"github.com/hielkefellinger/go-dnd/app/helpers"
	"log"
	"slices"
)

func upsertMap(mapUpdateRequest MapUpsertRequest, pool CampaignPool) (ecs.BaseEntity, error) {
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
