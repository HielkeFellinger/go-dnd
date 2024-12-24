package game_engine

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hielkefellinger/go-dnd/app/ecs"
	"github.com/hielkefellinger/go-dnd/app/ecs_components"
	"github.com/hielkefellinger/go-dnd/app/ecs_model_translation"
	"github.com/hielkefellinger/go-dnd/app/helpers"
	"github.com/hielkefellinger/go-dnd/app/models"
	"golang.org/x/net/html"
	"sort"
	"strconv"
)

func (e *baseEventMessageHandler) typeLoadUpsertInventory(message EventMessage, pool CampaignPool) error {
	// Undo escaping
	clearedBody := html.UnescapeString(message.Body)

	// Check if user is lead
	if message.Source != pool.GetLeadId() {
		return errors.New("modifying inventories is not allowed as non-lead")
	}

	data := make(map[string]any)

	allSelectableChars := make([]models.CampaignDropdownCharacter, 0)
	for _, characterEntity := range pool.GetEngine().GetWorld().GetCharacterEntities() {
		allSelectableChars = append(allSelectableChars, models.CampaignDropdownCharacter{
			Id:       characterEntity.GetId().String(),
			Name:     characterEntity.GetName(),
			Selected: false,
			Source:   characterEntity,
		})
	}

	// Check if there is an existing character with the supplied uuid
	uuidInventoryFilter, err := helpers.ParseStringToUuid(clearedBody)
	if err == nil {
		inventoryEntity, match := pool.GetEngine().GetWorld().GetEntityByUuid(uuidInventoryFilter)
		if match && inventoryEntity.HasComponentType(ecs.InventoryComponentType) {
			inventoryModel := ecs_model_translation.InventoryEntityToCampaignInventoryModel(inventoryEntity)

			// Link to characters
			for i, characterEntity := range allSelectableChars {
				allSelectableChars[i].Selected = characterEntity.Source.HasRelationWithEntityByUuid(inventoryEntity.GetId())
			}

			data["Inventory"] = inventoryModel
		} else {
			return errors.New("no inventories found with matching identifier")
		}
	}
	parsedItems := make([]*models.CampaignInventoryItem, 0)
	allItemEntities := pool.GetEngine().GetWorld().GetItemEntities()
	for _, itemEntity := range allItemEntities {
		parsedItems = append(parsedItems, ecs_model_translation.ItemEntityToCampaignInventoryItem(itemEntity, 0))
	}
	sort.Slice(parsedItems, func(i, j int) bool {
		return parsedItems[i].Name < parsedItems[j].Name
	})
	data["Items"] = parsedItems

	// Sort
	sort.Slice(allSelectableChars, func(i, j int) bool {
		return allSelectableChars[i].Selected || (allSelectableChars[i].Name < allSelectableChars[j].Name && !allSelectableChars[j].Selected)
	})
	data["Characters"] = allSelectableChars

	rawJsonBytes, err := json.Marshal(
		e.handleLoadHtmlBodyMultipleTemplateFiles([]string{"manageInventoryCrud.html", "diceSpinnerSvg.html", "inventory.html"},
			"manageInventoryCrud", data))
	if err != nil {
		return err
	}

	loadItemMessage := NewEventMessage()
	loadItemMessage.Source = message.Source
	loadItemMessage.Type = TypeLoadUpsertInventory
	loadItemMessage.Body = string(rawJsonBytes)
	loadItemMessage.Destinations = append(loadItemMessage.Destinations, pool.GetLeadId())
	pool.TransmitEventMessage(loadItemMessage)

	return nil
}

func (e *baseEventMessageHandler) typeUpsertInventory(message EventMessage, pool CampaignPool) error {
	// Check if user is lead
	if message.Source != pool.GetLeadId() {
		return errors.New("modifying Inventory details and ownership is not allowed as non-lead")
	}

	// Undo escaping
	clearedBody := html.UnescapeString(message.Body)

	var inventUpsertRequest inventoryUpsertRequest
	if err := json.Unmarshal([]byte(clearedBody), &inventUpsertRequest); err != nil {
		return err
	}

	// Upsert
	inventoryEntity, upsertError := upsertInventory(inventUpsertRequest, pool)
	if upsertError != nil {
		return upsertError
	}

	// Update the CRUD box AND the "Manage Inventory screen"
	loadUpsertInventoryMessage := NewEventMessage()
	loadUpsertInventoryMessage.Source = pool.GetLeadId()
	loadUpsertInventoryMessage.Body = inventoryEntity.GetId().String()
	if typeLoadUpsertInventoryErr := e.typeLoadUpsertInventory(loadUpsertInventoryMessage, pool); typeLoadUpsertInventoryErr != nil {
		return SendManagementError("Error", typeLoadUpsertInventoryErr.Error(), pool)
	}
	if typeManageInventoryErr := e.typeManageInventory(loadUpsertInventoryMessage, pool); typeManageInventoryErr != nil {
		return SendManagementError("Error", typeManageInventoryErr.Error(), pool)
	}
	return nil
}

func (e *baseEventMessageHandler) typeCloneInventory(message EventMessage, pool CampaignPool) error {
	// Check if user is lead
	if message.Source != pool.GetLeadId() {
		return errors.New("modifying Inventory details and ownership is not allowed as non-lead")
	}

	// Undo escaping
	clearedBody := html.UnescapeString(message.Body)

	uuidInventoryFilter, err := helpers.ParseStringToUuid(clearedBody)
	if err == nil {
		inventoryEntity, match := pool.GetEngine().GetWorld().GetEntityByUuid(uuidInventoryFilter)
		if match && inventoryEntity.HasComponentType(ecs.InventoryComponentType) {

			inventUpsertRequest := inventoryUpsertRequest{
				Name:        inventoryEntity.GetName() + " (clone)",
				Description: inventoryEntity.GetDescription(),
				Characters:  make([]string, 0),
			}

			for _, component := range inventoryEntity.GetAllComponentsOfType(ecs.InventoryComponentType) {
				inventoryComponent := component.(*ecs_components.InventoryComponent)
				inventUpsertRequest.Slots = strconv.Itoa(int(inventoryComponent.Slots))
			}

			clonedInventoryEntity, upsertError := upsertInventory(inventUpsertRequest, pool)
			if upsertError != nil {
				return upsertError
			}

			// Add the items
			// Loop over all the hasRelations and get the items
			rawHasRelations := inventoryEntity.GetAllComponentsOfType(ecs.HasRelationComponentType)
			for _, rawHasRelation := range rawHasRelations {
				hasRelation := rawHasRelation.(*ecs_components.HasRelationComponent)

				// Check if containing entity is an Item
				if hasRelation.Entity != nil && hasRelation.Entity.HasComponentType(ecs.ItemComponentType) {
					rawClonedHasRelation := ecs_components.NewHasRelationComponent()
					clonedHasRelation := rawClonedHasRelation.(*ecs_components.HasRelationComponent)
					clonedHasRelation.Entity = hasRelation.Entity
					clonedHasRelation.Count = hasRelation.Count

					if addErr := clonedInventoryEntity.AddComponent(clonedHasRelation); addErr != nil {
						return SendManagementError("Error", addErr.Error(), pool)
					}
				}
			}

			// Update the CRUD box AND the "Manage Inventory screen"
			loadUpsertInventoryMessage := NewEventMessage()
			loadUpsertInventoryMessage.Source = pool.GetLeadId()
			loadUpsertInventoryMessage.Body = clonedInventoryEntity.GetId().String()
			if typeLoadUpsertInventoryErr := e.typeLoadUpsertInventory(loadUpsertInventoryMessage, pool); typeLoadUpsertInventoryErr != nil {
				return SendManagementError("Error", typeLoadUpsertInventoryErr.Error(), pool)
			}
			if typeManageInventoryErr := e.typeManageInventory(loadUpsertInventoryMessage, pool); typeManageInventoryErr != nil {
				return SendManagementError("Error", typeManageInventoryErr.Error(), pool)
			}
		}
	}
	return nil
}

func (e *baseEventMessageHandler) typeRemoveInventory(message EventMessage, pool CampaignPool) error {
	// Check if user is lead
	if message.Source != pool.GetLeadId() {
		return errors.New("modifying Inventory details and ownership is not allowed as non-lead")
	}

	// Undo escaping
	clearedBody := html.UnescapeString(message.Body)

	uuidInventoryFilter, err := helpers.ParseStringToUuid(clearedBody)
	if err == nil {
		inventoryEntity, match := pool.GetEngine().GetWorld().GetEntityByUuid(uuidInventoryFilter)
		if match && inventoryEntity.HasComponentType(ecs.InventoryComponentType) {

			// Check if linked to characters
			for _, characterEntity := range pool.GetEngine().GetWorld().GetCharacterEntities() {
				if characterEntity.HasRelationWithEntityByUuid(inventoryEntity.GetId()) {
					return SendManagementError("Error", "This Specific Inventory is still linked to one or multiple Characters", pool)
				}
			}

			// Remove without clearing items
			if removeErr := pool.GetEngine().GetWorld().RemoveEntity(inventoryEntity); removeErr != nil {
				return removeErr
			}

			removeInventoryMessage := NewEventMessage()
			removeInventoryMessage.Type = TypeRemoveInventory
			removeInventoryMessage.Source = pool.GetLeadId()
			removeInventoryMessage.Body = ""
			removeInventoryMessage.Destinations = append(message.Destinations, pool.GetLeadId())
			if typeManageInventoryErr := e.typeManageInventory(removeInventoryMessage, pool); typeManageInventoryErr != nil {
				return SendManagementError("Error", typeManageInventoryErr.Error(), pool)
			}
			pool.TransmitEventMessage(removeInventoryMessage)

		} else {
			return errors.New("no inventories found with matching identifier")
		}
	}

	return err
}

func (e *baseEventMessageHandler) typeAddItemToInventory(message EventMessage, pool CampaignPool) error {
	// Check if user is lead
	if message.Source != pool.GetLeadId() {
		return errors.New("modifying Inventory details and ownership is not allowed as non-lead")
	}

	// Undo escaping
	clearedBody := html.UnescapeString(message.Body)

	var inventAddItemRequest inventoryAddRemoveItemRequest
	if err := json.Unmarshal([]byte(clearedBody), &inventAddItemRequest); err != nil {
		return err
	}

	// Check if there is an item and an inventory that match the request
	var item ecs.Entity = nil
	if uuidItemFilter, errItem := helpers.ParseStringToUuid(inventAddItemRequest.ItemId); errItem == nil {
		if entityFound, ok := pool.GetEngine().GetWorld().GetEntityByUuid(uuidItemFilter); ok {
			if entityFound.HasComponentType(ecs.ItemComponentType) {
				item = entityFound
			}
		}
	}
	var inventory ecs.Entity = nil
	if uuidInventoryFilter, errInv := helpers.ParseStringToUuid(inventAddItemRequest.InventoryId); errInv == nil {
		if entityFound, ok := pool.GetEngine().GetWorld().GetEntityByUuid(uuidInventoryFilter); ok {
			if entityFound.HasComponentType(ecs.InventoryComponentType) {
				inventory = entityFound
			}
		}
	}

	// Check integrity
	if item == nil || inventory == nil {
		return errors.New("no item and/or inventory found with matching identifier")
	}
	if inventory.HasRelationWithEntityByUuid(item.GetId()) {
		return SendManagementError("Warning", "This Specific Item is already present in inventory", pool)
	}

	// Find owning characters
	owningPcCharIds := make([]string, 0)
	for _, characterEntity := range pool.GetEngine().GetWorld().GetCharacterEntities() {
		if characterEntity.HasRelationWithEntityByUuid(inventory.GetId()) && characterEntity.HasComponentType(ecs.PlayerComponentType) {
			owningPcCharIds = append(owningPcCharIds, characterEntity.GetId().String())
		}
	}

	// Add Item
	hasRelationComponent := ecs_components.NewHasRelationComponent().(*ecs_components.HasRelationComponent)
	hasRelationComponent.Entity = item
	if err := inventory.AddComponent(hasRelationComponent); err != nil {
		return SendManagementError("Warning", err.Error(), pool)
	}

	// Trigger Visual Updates on chars (details)
	for _, charId := range owningPcCharIds {
		var reloadCharDetailMessage = NewEventMessage()
		reloadCharDetailMessage.Source = ServerUser
		reloadCharDetailMessage.Body = charId
		loadCharErr := e.loadCharactersDetails(reloadCharDetailMessage, pool)
		if loadCharErr != nil {
			return loadCharErr
		}
	}
	loadUpsertInventoryMessage := NewEventMessage()
	loadUpsertInventoryMessage.Source = pool.GetLeadId()
	loadUpsertInventoryMessage.Body = inventory.GetId().String()
	if typeLoadUpsertInventoryErr := e.typeLoadUpsertInventory(loadUpsertInventoryMessage, pool); typeLoadUpsertInventoryErr != nil {
		return SendManagementError("Error", typeLoadUpsertInventoryErr.Error(), pool)
	}

	return nil
}

func (e *baseEventMessageHandler) typeRemoveItemFromInventory(message EventMessage, pool CampaignPool) error {
	// Check if user is lead
	if message.Source != pool.GetLeadId() {
		return errors.New("modifying Inventory details and ownership is not allowed as non-lead")
	}

	// Undo escaping
	clearedBody := html.UnescapeString(message.Body)

	var inventRemoveItemRequest inventoryAddRemoveItemRequest
	if err := json.Unmarshal([]byte(clearedBody), &inventRemoveItemRequest); err != nil {
		return err
	}

	// Check if there is an item and an inventory that match the request
	var item ecs.Entity = nil
	if uuidItemFilter, errItem := helpers.ParseStringToUuid(inventRemoveItemRequest.ItemId); errItem == nil {
		if entityFound, ok := pool.GetEngine().GetWorld().GetEntityByUuid(uuidItemFilter); ok {
			if entityFound.HasComponentType(ecs.ItemComponentType) {
				item = entityFound
			}
		}
	}
	var inventory ecs.Entity = nil
	if uuidInventoryFilter, errInv := helpers.ParseStringToUuid(inventRemoveItemRequest.InventoryId); errInv == nil {
		if entityFound, ok := pool.GetEngine().GetWorld().GetEntityByUuid(uuidInventoryFilter); ok {
			if entityFound.HasComponentType(ecs.InventoryComponentType) {
				inventory = entityFound
			}
		}
	}

	// Check integrity
	if item == nil || inventory == nil {
		return errors.New("no item and/or inventory found with matching identifier")
	}
	if !inventory.HasRelationWithEntityByUuid(item.GetId()) {
		return SendManagementError("Warning", "item is already removed from inventory", pool)
	}

	// Find owning characters
	owningPcCharIds := make([]string, 0)
	for _, characterEntity := range pool.GetEngine().GetWorld().GetCharacterEntities() {
		if characterEntity.HasRelationWithEntityByUuid(inventory.GetId()) && characterEntity.HasComponentType(ecs.PlayerComponentType) {
			owningPcCharIds = append(owningPcCharIds, characterEntity.GetId().String())
		}
	}

	// Remove Item
	for _, hasRelations := range inventory.GetAllComponentsOfType(ecs.HasRelationComponentType) {
		hasRelationComponent := hasRelations.(*ecs_components.HasRelationComponent)
		if hasRelationComponent.Entity.GetId() == item.GetId() {
			inventory.RemoveComponentByUuid(hasRelationComponent.Id)
		}
	}

	if uuidInventoryFilter, errInv := helpers.ParseStringToUuid(inventRemoveItemRequest.InventoryId); errInv == nil {
		if entityFound, ok := pool.GetEngine().GetWorld().GetEntityByUuid(uuidInventoryFilter); ok {
			if entityFound.HasComponentType(ecs.InventoryComponentType) {
				inventory = entityFound
			}
		}
	}

	// Trigger Visual Updates on chars (details)
	for _, charId := range owningPcCharIds {
		var reloadCharDetailMessage = NewEventMessage()
		reloadCharDetailMessage.Source = ServerUser
		reloadCharDetailMessage.Body = charId
		loadCharErr := e.loadCharactersDetails(reloadCharDetailMessage, pool)
		if loadCharErr != nil {
			return loadCharErr
		}
	}
	loadUpsertInventoryMessage := NewEventMessage()
	loadUpsertInventoryMessage.Source = pool.GetLeadId()
	loadUpsertInventoryMessage.Body = inventory.GetId().String()
	if typeLoadUpsertInventoryErr := e.typeLoadUpsertInventory(loadUpsertInventoryMessage, pool); typeLoadUpsertInventoryErr != nil {
		return SendManagementError("Error", typeLoadUpsertInventoryErr.Error(), pool)
	}

	return nil
}

func (e *baseEventMessageHandler) typeLoadUpsertCharacter(message EventMessage, pool CampaignPool) error {
	// Undo escaping
	clearedBody := html.UnescapeString(message.Body)

	// Check if user is lead
	if message.Source != pool.GetLeadId() {
		return errors.New("modifying characters is not allowed as non-lead")
	}

	data := make(map[string]any)

	// Check if there is an existing character with the supplied uuid
	uuidItemFilter, err := helpers.ParseStringToUuid(clearedBody)
	if err == nil {
		charCandidate, match := pool.GetEngine().GetWorld().GetEntityByUuid(uuidItemFilter)
		if match && charCandidate.HasComponentType(ecs.CharacterComponentType) {
			data["Character"] = ecs_model_translation.CharacterEntityToCampaignCharacterModel(charCandidate)
		} else {
			return errors.New("no characters found with matching identifier")
		}
	}

	rawJsonBytes, err := json.Marshal(
		e.handleLoadHtmlBodyMultipleTemplateFiles([]string{"manageCharacterCrud.html", "diceSpinnerSvg.html"},
			"manageCharacterCrud", data))
	if err != nil {
		return err
	}

	loadItemCharacter := NewEventMessage()
	loadItemCharacter.Source = message.Source
	loadItemCharacter.Type = TypeLoadUpsertCharacter
	loadItemCharacter.Body = string(rawJsonBytes)
	loadItemCharacter.Destinations = append(loadItemCharacter.Destinations, pool.GetLeadId())
	pool.TransmitEventMessage(loadItemCharacter)

	return nil
}

func (e *baseEventMessageHandler) typeUpsertCharacter(message EventMessage, pool CampaignPool) error {
	// Check if user is lead
	if message.Source != pool.GetLeadId() {
		return errors.New("modifying characters is not allowed as non-lead")
	}

	// Undo escaping
	clearedBody := html.UnescapeString(message.Body)

	var charUpsertRequest characterUpsertRequest
	if err := json.Unmarshal([]byte(clearedBody), &charUpsertRequest); err != nil {
		return err
	}

	// Upsert
	charEntity, upsertError := upsertCharacter(charUpsertRequest, pool)
	if upsertError != nil {
		return upsertError
	}

	// Add an Inventory
	if charUpsertRequest.AddInventory {
		newInventory := ecs.NewEntity()
		// Set Entity to SlotType
		if addErr := newInventory.AddComponent(ecs_components.NewInventoryComponent()); addErr != nil {
			return SendManagementError("Error", addErr.Error(), pool)
		}
		// Add to world
		if addErr := pool.GetEngine().GetWorld().AddEntity(&newInventory); addErr != nil {
			return SendManagementError("Error", addErr.Error(), pool)
		}

		// Build relation and add
		hasRelation := ecs_components.NewHasRelationComponent().(*ecs_components.HasRelationComponent)
		hasRelation.Count = 1
		hasRelation.Entity = &newInventory
		if addErr := charEntity.AddComponent(hasRelation); addErr != nil {
			return SendManagementError("Error", addErr.Error(), pool)
		}
	}

	// Refresh Character Ribbon
	loadUpsertCharMessage := NewEventMessage()
	loadUpsertCharMessage.Source = pool.GetLeadId()
	loadUpsertCharMessage.Body = charEntity.GetId().String()
	if typeLoadUpsertMapErr := e.typeLoadUpsertCharacter(loadUpsertCharMessage, pool); typeLoadUpsertMapErr != nil {
		return SendManagementError("Error", typeLoadUpsertMapErr.Error(), pool)
	}
	if typeManageMapsErr := e.typeManageCharacters(loadUpsertCharMessage, pool); typeManageMapsErr != nil {
		return SendManagementError("Error", typeManageMapsErr.Error(), pool)
	}

	// Update Char info
	var reloadCharRibbon = NewEventMessage()
	reloadCharRibbon.Source = ServerUser
	reloadCharRibbon.Type = TypeLoadCharacters
	_ = e.loadCharacters(reloadCharRibbon, pool)

	var closeCharDetailMessage = NewEventMessage()
	closeCharDetailMessage.Source = ServerUser
	closeCharDetailMessage.Type = TypeLoadCharactersDetails
	closeCharDetailMessage.Body = charEntity.GetId().String()
	if err := e.loadCharactersDetails(closeCharDetailMessage, pool); err != nil {
		return err
	}

	// Reload Map info
	// Update possible Map Entities
	mapEntities := pool.GetEngine().GetWorld().GetMapEntities()
	for _, mapEntity := range mapEntities {

		// Only get the map with the relevant relation to entity
		if !mapEntity.HasRelationWithEntityByUuid(charEntity.GetId()) {
			continue
		}

		for _, mapItem := range mapEntity.GetAllComponentsOfType(ecs.MapItemRelationComponentType) {
			mapItemRelComponent := mapItem.(*ecs_components.MapItemRelationComponent)

			if mapItemRelComponent.Entity.GetId() == charEntity.GetId() {
				reloadMapItemMessage := NewEventMessage()
				reloadMapItemMessage.Source = ServerUser
				reloadMapItemMessage.Body = mapItemRelComponent.Id.String()
				reloadMapItemErr := e.typeLoadMapEntity(reloadMapItemMessage, pool)
				if reloadMapItemErr != nil {
					return reloadMapItemErr
				}
			}
		}
	}
	return nil
}

func (e *baseEventMessageHandler) typeRemoveCharacter(message EventMessage, pool CampaignPool) error {
	// Undo escaping
	clearedBody := html.UnescapeString(message.Body)

	// Check if user is lead
	if message.Source != pool.GetLeadId() {
		return errors.New("modifying maps is not allowed as non-lead")
	}

	uuidCharacterFilter, err := helpers.ParseStringToUuid(clearedBody)
	if err == nil {
		characterEntity, match := pool.GetEngine().GetWorld().GetEntityByUuid(uuidCharacterFilter)
		if match && characterEntity.HasComponentType(ecs.CharacterComponentType) {
			// Test if on a map
			isOnMapNames := make([]string, 0)
			for _, rawMapEntity := range pool.GetEngine().GetWorld().GetMapEntities() {
				if rawMapEntity.HasRelationWithEntityByUuid(uuidCharacterFilter) {
					isOnMapNames = append(isOnMapNames, rawMapEntity.GetName())
				}
			}

			// Test if linked to player
			playersAssigned := make([]string, 0)
			isAPlayerChar := false
			for _, rawPlayerComponent := range characterEntity.GetAllComponentsOfType(ecs.PlayerComponentType) {
				isAPlayerChar = true
				playerComponent := rawPlayerComponent.(*ecs_components.PlayerComponent)
				if playerComponent.Name != "" {
					playersAssigned = append(playersAssigned, playerComponent.Name)
				}
			}

			// Handle player linked to game.
			if len(isOnMapNames) > 0 || len(playersAssigned) > 0 {
				errorMessage := "Deletion of Character '' is not allowed!"
				if len(isOnMapNames) > 0 {
					errorMessage += fmt.Sprintf("\nCharacter is still present on the following maps: '%v'", isOnMapNames)
				}
				if len(playersAssigned) > 0 {
					errorMessage += fmt.Sprintf("\nCharacter is still linked to the following players: '%v'", playersAssigned)
				}
				return SendManagementError("Error", errorMessage, pool)
			}

			// Remove without clearing items
			if removeErr := pool.GetEngine().GetWorld().RemoveEntity(characterEntity); removeErr != nil {
				return removeErr
			}

			removeCharacterMessage := NewEventMessage()
			removeCharacterMessage.Type = TypeRemoveCharacter
			removeCharacterMessage.Source = pool.GetLeadId()
			removeCharacterMessage.Body = ""
			removeCharacterMessage.Destinations = append(message.Destinations, pool.GetLeadId())
			if typeManageCharactersErr := e.typeManageCharacters(removeCharacterMessage, pool); typeManageCharactersErr != nil {
				return SendManagementError("Error", typeManageCharactersErr.Error(), pool)
			}
			pool.TransmitEventMessage(removeCharacterMessage)

			// If a player Char; update ribbon
			if isAPlayerChar {
				var reloadCharRibbon = NewEventMessage()
				reloadCharRibbon.Source = ServerUser
				reloadCharRibbon.Type = TypeLoadCharacters
				_ = e.loadCharacters(reloadCharRibbon, pool)
			}

		} else {
			return errors.New("no Characters found with matching identifier")
		}
	}
	return err
}

func (e *baseEventMessageHandler) typeLoadUpsertMap(message EventMessage, pool CampaignPool) error {
	// Undo escaping
	clearedBody := html.UnescapeString(message.Body)

	// Check if user is lead
	if message.Source != pool.GetLeadId() {
		return errors.New("modifying maps is not allowed as non-lead")
	}

	data := make(map[string]any)

	// Check if there is an existing map with the supplied uuid
	uuidItemFilter, err := helpers.ParseStringToUuid(clearedBody)
	if err == nil {
		mapCandidate, match := pool.GetEngine().GetWorld().GetEntityByUuid(uuidItemFilter)
		if match && mapCandidate.HasComponentType(ecs.MapComponentType) {
			data["Map"] = ecs_model_translation.MapEntityToCampaignMapModel(mapCandidate)
		}
	}

	rawJsonBytes, err := json.Marshal(
		e.handleLoadHtmlBodyMultipleTemplateFiles([]string{"manageMapCrud.html", "diceSpinnerSvg.html"},
			"manageMapCrud", data))
	if err != nil {
		return err
	}

	loadMapMessage := NewEventMessage()
	loadMapMessage.Source = message.Source
	loadMapMessage.Type = TypeLoadUpsertMap
	loadMapMessage.Body = string(rawJsonBytes)
	loadMapMessage.Destinations = append(loadMapMessage.Destinations, pool.GetLeadId())
	pool.TransmitEventMessage(loadMapMessage)

	return nil
}

func (e *baseEventMessageHandler) typeUpsertMap(message EventMessage, pool CampaignPool) error {
	// Check if user is lead
	if message.Source != pool.GetLeadId() {
		return errors.New("modifying maps is not allowed as non-lead")
	}

	// Undo escaping
	clearedBody := html.UnescapeString(message.Body)

	var mapUpdateRequest mapUpsertRequest
	if err := json.Unmarshal([]byte(clearedBody), &mapUpdateRequest); err != nil {
		return err
	}

	// Upsert
	mapEntity, upsertError := upsertMap(mapUpdateRequest, pool)
	if upsertError != nil {
		return upsertError
	}

	// Update the CRUD box AND the "Manage Maps screen"
	loadUpsertMapMessage := NewEventMessage()
	loadUpsertMapMessage.Source = pool.GetLeadId()
	loadUpsertMapMessage.Body = mapEntity.GetId().String()
	if typeLoadUpsertMapErr := e.typeLoadUpsertMap(loadUpsertMapMessage, pool); typeLoadUpsertMapErr != nil {
		return SendManagementError("Error", typeLoadUpsertMapErr.Error(), pool)
	}
	if typeManageMapsErr := e.typeManageMaps(loadUpsertMapMessage, pool); typeManageMapsErr != nil {
		return SendManagementError("Error", typeManageMapsErr.Error(), pool)
	}
	return nil
}

func (e *baseEventMessageHandler) typeLoadUpsertItem(message EventMessage, pool CampaignPool) error {
	// Check if user is lead
	if message.Source != pool.GetLeadId() {
		return errors.New("modifying items is not allowed as non-lead")
	}

	// Undo escaping
	clearedBody := html.UnescapeString(message.Body)

	data := make(map[string]any)

	// Check if there is an existing item with the supplied uuid
	uuidItemFilter, err := helpers.ParseStringToUuid(clearedBody)
	if err == nil {
		itemCandidate, match := pool.GetEngine().GetWorld().GetItemEntityByUuid(uuidItemFilter)
		if match && itemCandidate.HasComponentType(ecs.ItemComponentType) {
			data["Item"] = ecs_model_translation.ItemEntityToCampaignInventoryItem(itemCandidate, 0)
		}
	}

	rawJsonBytes, err := json.Marshal(
		e.handleLoadHtmlBodyMultipleTemplateFiles([]string{"manageItemCrud.html", "diceSpinnerSvg.html"},
			"manageItemCrud", data))
	if err != nil {
		return err
	}

	loadItemMessage := NewEventMessage()
	loadItemMessage.Source = message.Source
	loadItemMessage.Type = TypeLoadUpsertItem
	loadItemMessage.Body = string(rawJsonBytes)
	loadItemMessage.Destinations = append(loadItemMessage.Destinations, pool.GetLeadId())
	pool.TransmitEventMessage(loadItemMessage)

	return nil
}

func (e *baseEventMessageHandler) typeUpsertItem(message EventMessage, pool CampaignPool) error {
	// Check if user is lead
	if message.Source != pool.GetLeadId() {
		return errors.New("modifying items is not allowed as non-lead")
	}

	// Undo escaping
	clearedBody := html.UnescapeString(message.Body)

	var itemUpsertRequest itemUpsertRequest
	if err := json.Unmarshal([]byte(clearedBody), &itemUpsertRequest); err != nil {
		return err
	}

	// Upsert
	itemEntity, upsertError := upsertItem(itemUpsertRequest, pool)
	if upsertError != nil {
		return upsertError
	}

	// Update the CRUD box AND the "Manage Maps screen"
	loadUpsertItemMessage := NewEventMessage()
	loadUpsertItemMessage.Source = pool.GetLeadId()
	loadUpsertItemMessage.Body = itemEntity.GetId().String()
	if typeLoadUpsertMapErr := e.typeLoadUpsertItem(loadUpsertItemMessage, pool); typeLoadUpsertMapErr != nil {
		return SendManagementError("Error", typeLoadUpsertMapErr.Error(), pool)
	}
	if typeManageMapsErr := e.typeManageItems(loadUpsertItemMessage, pool); typeManageMapsErr != nil {
		return SendManagementError("Error", typeManageMapsErr.Error(), pool)
	}
	return nil
}

type inventoryAddRemoveItemRequest struct {
	InventoryId string `json:"InventoryId"`
	ItemId      string `json:"ItemId"`
}

type inventoryUpsertRequest struct {
	Id          string   `json:"Id"`
	Name        string   `json:"Name"`
	Slots       string   `json:"Slots"`
	Description string   `json:"Description"`
	Characters  []string `json:"Characters"`
}

type characterUpsertRequest struct {
	Id             string             `json:"Id"`
	Name           string             `json:"Name"`
	Description    string             `json:"Description"`
	HealthDamage   string             `json:"HealthDamage"`
	HealthTmp      string             `json:"HealthTmp"`
	HealthMax      string             `json:"HealthMax"`
	Level          string             `json:"Level"`
	Image          helpers.FileUpload `json:"Image"`
	ImageName      string             `json:"ImageName"`
	PlayerPlayable bool               `json:"PlayerPlayable"`
	AddInventory   bool               `json:"AddInventory"`
	Hidden         bool               `json:"Hidden"`
}

type mapUpsertRequest struct {
	Id           string             `json:"Id"`
	Name         string             `json:"Name"`
	Description  string             `json:"Description"`
	X            string             `json:"X"`
	Y            string             `json:"Y"`
	ImageName    string             `json:"ImageName"`
	RemoveImages []string           `json:"RemoveImages"`
	Image        helpers.FileUpload `json:"Image"`
}

type itemUpsertRequest struct {
	Id          string `json:"Id"`
	Name        string `json:"Name"`
	Description string `json:"Description"`
	Damage      string `json:"Damage"`
	Restore     string `json:"Restore"`
	RangeMin    string `json:"RangeMin"`
	RangeMax    string `json:"RangeMax"`
	Weight      string `json:"Weight"`
}
