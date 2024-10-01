package game_engine

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hielkefellinger/go-dnd/app/ecs"
	"github.com/hielkefellinger/go-dnd/app/ecs_components"
	"github.com/hielkefellinger/go-dnd/app/ecs_model_translation"
	"github.com/hielkefellinger/go-dnd/app/helpers"
	"golang.org/x/net/html"
	"log"
)

func (e *baseEventMessageHandler) handleManagementCrudEvents(message EventMessage, pool CampaignPool) error {
	log.Printf("- Game Management CRUD Events Type: '%d' Message: '%s'", message.Type, message.Id)
	var handled = false

	// Maps
	if message.Type == TypeLoadUpsertMap {
		err := e.typeLoadUpsertMap(message, pool)
		if err != nil {
			return err
		}
		handled = true
	} else if message.Type == TypeUpsertMap {
		err := e.typeUpsertMap(message, pool)
		if err != nil {
			return err
		}
		handled = true
	}

	// Items
	if message.Type == TypeLoadUpsertItem {
		err := e.typeLoadUpsertItem(message, pool)
		if err != nil {
			return err
		}
		handled = true
	} else if message.Type == TypeUpsertItem {
		err := e.typeUpsertItem(message, pool)
		if err != nil {
			return err
		}
		handled = true
	}

	// Characters
	if message.Type == TypeLoadUpsertCharacter {
		err := e.typeLoadUpsertCharacter(message, pool)
		if err != nil {
			return err
		}
		handled = true
	} else if message.Type == TypeUpsertCharacter {
		err := e.typeUpsertCharacter(message, pool)
		if err != nil {
			return err
		}
		handled = true
	}

	if !handled {
		return errors.New(fmt.Sprintf("message of type '%d' is not recognised by 'handleManagementCrudEvents()'", message.Type))
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

	// Check if there is an existing map with the supplied uuid
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
		e.handleLoadHtmlBodyMultipleTemplateFiles([]string{"campaignUpsertCharacter.html", "diceSpinnerSvg.html"},
			"campaignUpsertCharacter", data))
	if err != nil {
		return err
	}

	loadItemMessage := NewEventMessage()
	loadItemMessage.Source = message.Source
	loadItemMessage.Type = TypeLoadUpsertCharacter
	loadItemMessage.Body = string(rawJsonBytes)
	loadItemMessage.Destinations = append(loadItemMessage.Destinations, pool.GetLeadId())
	pool.TransmitEventMessage(loadItemMessage)

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
	err := json.Unmarshal([]byte(clearedBody), &charUpsertRequest)
	if err != nil {
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
		if addErr := newInventory.AddComponent(ecs_components.NewSlotsComponent()); addErr != nil {
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
		return e.sendManagementError("Error", typeLoadUpsertMapErr.Error(), pool)
	}
	if typeManageMapsErr := e.typeManageCharacters(loadUpsertCharMessage, pool); typeManageMapsErr != nil {
		return e.sendManagementError("Error", typeManageMapsErr.Error(), pool)
	}

	// Update Char info
	var reloadCharRibbon = NewEventMessage()
	reloadCharRibbon.Source = ServerUser
	reloadCharRibbon.Type = TypeLoadCharacters
	e.loadCharacters(reloadCharRibbon, pool)

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
		e.handleLoadHtmlBodyMultipleTemplateFiles([]string{"campaignUpsertMap.html", "diceSpinnerSvg.html"},
			"campaignUpsertMap", data))
	if err != nil {
		return err
	}

	loadItemMessage := NewEventMessage()
	loadItemMessage.Source = message.Source
	loadItemMessage.Type = TypeLoadUpsertMap
	loadItemMessage.Body = string(rawJsonBytes)
	loadItemMessage.Destinations = append(loadItemMessage.Destinations, pool.GetLeadId())
	pool.TransmitEventMessage(loadItemMessage)

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
	err := json.Unmarshal([]byte(clearedBody), &mapUpdateRequest)
	if err != nil {
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
		return e.sendManagementError("Error", typeLoadUpsertMapErr.Error(), pool)
	}
	if typeManageMapsErr := e.typeManageMaps(loadUpsertMapMessage, pool); typeManageMapsErr != nil {
		return e.sendManagementError("Error", typeManageMapsErr.Error(), pool)
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
		e.handleLoadHtmlBodyMultipleTemplateFiles([]string{"campaignUpsertItem.html", "diceSpinnerSvg.html"},
			"campaignUpsertItem", data))
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
	err := json.Unmarshal([]byte(clearedBody), &itemUpsertRequest)
	if err != nil {
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
		return e.sendManagementError("Error", typeLoadUpsertMapErr.Error(), pool)
	}
	if typeManageMapsErr := e.typeManageItems(loadUpsertItemMessage, pool); typeManageMapsErr != nil {
		return e.sendManagementError("Error", typeManageMapsErr.Error(), pool)
	}
	return nil
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
