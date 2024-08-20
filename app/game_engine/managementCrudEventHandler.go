package game_engine

import (
	"encoding/json"
	"errors"
	"github.com/hielkefellinger/go-dnd/app/ecs"
	"github.com/hielkefellinger/go-dnd/app/ecs_model_translation"
	"github.com/hielkefellinger/go-dnd/app/helpers"
	"golang.org/x/net/html"
	"log"
)

func (e *baseEventMessageHandler) handleManagementCrudEvents(message EventMessage, pool CampaignPool) error {
	log.Printf("- Game Management CRUD Events Type: '%d' Message: '%s'", message.Type, message.Id)

	// Maps
	if message.Type == TypeLoadUpsertMap {
		err := e.typeLoadUpsertMap(message, pool)
		if err != nil {
			return err
		}
	} else if message.Type == TypeUpsertMap {
		err := e.typeUpsertMap(message, pool)
		if err != nil {
			return err
		}
	}

	// Items
	if message.Type == TypeLoadUpsertItem {
		err := e.typeLoadUpsertItem(message, pool)
		if err != nil {
			return err
		}
	} else if message.Type == TypeUpsertItem {
		err := e.typeUpsertItem(message, pool)
		if err != nil {
			return err
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
	// Undo escaping
	clearedBody := html.UnescapeString(message.Body)

	// Check if user is lead
	if message.Source != pool.GetLeadId() {
		return errors.New("modifying items is not allowed as non-lead")
	}

	var mapUpdateRequest MapUpsertRequest
	err := json.Unmarshal([]byte(clearedBody), &mapUpdateRequest)
	if err != nil {
		return err
	}

	// @TODO: Split after this part
	// Check if it needs to be updated; or inserted
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
	// Undo escaping
	clearedBody := html.UnescapeString(message.Body)

	// Check if user is lead
	if message.Source != pool.GetLeadId() {
		return errors.New("modifying items is not allowed as non-lead")
	}

	data := make(map[string]any)

	// Check if there is an existing item with the supplied uuid
	uuidItemFilter, err := helpers.ParseStringToUuid(clearedBody)
	if err == nil {
		itemCandidate, match := pool.GetEngine().GetWorld().GetEntityByUuid(uuidItemFilter)
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
	// Undo escaping
	clearedBody := html.UnescapeString(message.Body)

	// Check if user is lead
	if message.Source != pool.GetLeadId() {
		return errors.New("modifying items is not allowed as non-lead")
	}

	// @TODO translate to object and update

	log.Printf("- Item ID: '%v'", clearedBody)

	// Check if it needs to be updated; or inserted

	loadItemMessage := NewEventMessage()
	loadItemMessage.Source = pool.GetLeadId()
	loadItemMessage.Body = ""
	return e.typeLoadUpsertItem(loadItemMessage, pool)
}
