package game_engine

import (
	"encoding/json"
	"errors"
	"github.com/hielkefellinger/go-dnd/app/ecs"
	"github.com/hielkefellinger/go-dnd/app/ecs_model_translation"
	"golang.org/x/net/html"
	"log"
)

func (e *baseEventMessageHandler) handleItemEvents(message EventMessage, pool CampaignPool) error {
	log.Printf("- Item Update Event Type: '%d' Message: '%s'", message.Type, message.Id)

	if message.Type == TypeLoadItem {
		err := e.typeLoadItem(message, pool)
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

func (e *baseEventMessageHandler) typeLoadItem(message EventMessage, pool CampaignPool) error {
	// Undo escaping
	clearedBody := html.UnescapeString(message.Body)

	// Check if user is lead
	if message.Source != pool.GetLeadId() {
		return errors.New("modifying items is not allowed as non-lead")
	}

	data := make(map[string]any)

	// Check if there is an existing item with the supplied uuid
	uuidItemFilter, err := parseStingToUuid(clearedBody)
	if err == nil {
		itemCandidate, match := pool.GetEngine().GetWorld().GetEntityByUuid(uuidItemFilter)
		if match && itemCandidate.HasComponentType(ecs.ItemComponentType) {
			data["Item"] = ecs_model_translation.ItemEntityToCampaignInventoryItem(itemCandidate, 0)
		}
	}

	log.Printf("- Item ID: '%v' Obj: '%v'", clearedBody, data)

	rawJsonBytes, err := json.Marshal(
		e.handleLoadHtmlBodyMultipleTemplateFiles([]string{"campaignUpsertItem.html", "diceSpinnerSvg.html"},
			"campaignUpsertItem", data))
	if err != nil {
		return err
	}

	loadItemMessage := NewEventMessage()
	loadItemMessage.Source = message.Source
	loadItemMessage.Type = TypeLoadItem
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

	return nil
}
