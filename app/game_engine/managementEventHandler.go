package game_engine

import (
	"encoding/json"
	"errors"
	"github.com/hielkefellinger/go-dnd/app/ecs_model_translation"
	"github.com/hielkefellinger/go-dnd/app/models"
	"log"
)

func (e *baseEventMessageHandler) handleManagementEvents(message EventMessage, pool CampaignPool) error {
	log.Printf("- Game Management Events Type: '%d' Message: '%s'", message.Type, message.Id)

	if message.Type == TypeManageMaps {
		return e.typeManageMaps(message, pool)
	}

	return nil
}

func (e *baseEventMessageHandler) typeManageMaps(message EventMessage, pool CampaignPool) error {

	if message.Source != pool.GetLeadId() {
		return errors.New("managing game Map elements is not allowed as non-lead")
	}

	// Update possible Map Entities
	mapEntries := make([]models.CampaignMap, 0)
	mapEntities := pool.GetEngine().GetWorld().GetMapEntities()
	for _, mapEntity := range mapEntities {
		mapEntries = append(mapEntries, ecs_model_translation.MapEntityToCampaignMapModel(mapEntity))
	}

	data := make(map[string]any)
	data["Maps"] = mapEntries
	rawJsonBytes, err := json.Marshal(
		e.handleLoadHtmlBodyMultipleTemplateFiles(
			[]string{"campaignManageMaps.html", "campaignManageMapSelectionBox.html"},
			"campaignManageMaps", data))
	if err != nil {
		return err
	}

	manageMaps := NewEventMessage()
	manageMaps.Source = message.Source
	manageMaps.Type = TypeManageMaps
	manageMaps.Body = string(rawJsonBytes)
	manageMaps.Destinations = append(manageMaps.Destinations, pool.GetLeadId())
	pool.TransmitEventMessage(manageMaps)

	return nil
}
