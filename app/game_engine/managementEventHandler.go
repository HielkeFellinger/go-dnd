package game_engine

import (
	"encoding/json"
	"errors"
	"github.com/hielkefellinger/go-dnd/app/ecs"
	"github.com/hielkefellinger/go-dnd/app/ecs_components"
	"github.com/hielkefellinger/go-dnd/app/ecs_model_translation"
	"github.com/hielkefellinger/go-dnd/app/models"
	"log"
	"slices"
	"sort"
)

func (e *baseEventMessageHandler) handleManagementEvents(message EventMessage, pool CampaignPool) error {
	log.Printf("- Game Management Events Type: '%d' Message: '%s'", message.Type, message.Id)

	if message.Type == TypeManageMaps {
		return e.typeManageMaps(message, pool)
	} else if message.Type == TypeManageCharacters {
		return e.typeManageCharacters(message, pool)
	} else if message.Type == TypeManageInventory {
		return e.typeManageInventory(message, pool)
	} else if message.Type == TypeManageItems {
		return e.typeManageItems(message, pool)
	} else if message.Type == TypeManageCampaign {
		return e.typeManageCampaign(message, pool)
	}

	return nil
}

func (e *baseEventMessageHandler) typeManageCampaign(message EventMessage, pool CampaignPool) error {
	if message.Source != pool.GetLeadId() {
		return errors.New("managing the campaign is not allowed as non-lead")
	}

	return nil
}

func (e *baseEventMessageHandler) typeManageItems(message EventMessage, pool CampaignPool) error {
	if message.Source != pool.GetLeadId() {
		return errors.New("managing game Item elements is not allowed as non-lead")
	}

	return nil
}

func (e *baseEventMessageHandler) typeManageInventory(message EventMessage, pool CampaignPool) error {
	if message.Source != pool.GetLeadId() {
		return errors.New("managing game Inventory elements is not allowed as non-lead")
	}

	return nil
}

func (e *baseEventMessageHandler) typeManageCharacters(message EventMessage, pool CampaignPool) error {
	if message.Source != pool.GetLeadId() {
		return errors.New("managing game Characters elements is not allowed as non-lead")
	}

	charEntities := pool.GetEngine().GetWorld().GetCharacterEntities()
	allPlayers := pool.GetAllClientIds()

	var characters models.Characters
	for _, charEntity := range charEntities {

		charModel := models.Character{
			Id:          charEntity.GetId().String(),
			Name:        charEntity.GetName(),
			Description: charEntity.GetDescription(),
		}

		// Check if (one of n of controlling) player(s) is online
		for _, rawPlayerComponent := range charEntity.GetAllComponentsOfType(ecs.PlayerComponentType) {
			playerComponent := rawPlayerComponent.(*ecs_components.PlayerComponent)
			charModel.Online = charModel.Online || slices.Contains(allPlayers, playerComponent.Name)
		}

		var charDetails = charEntity.GetAllComponentsOfType(ecs.CharacterComponentType)
		if charDetails != nil && len(charDetails) > 0 {
			characterComponent := charDetails[0].(*ecs_components.CharacterComponent)
			charModel.Name = characterComponent.Name
			charModel.Description = characterComponent.Description
		}

		var image *ecs_components.ImageComponent
		var imageDetails = charEntity.GetAllComponentsOfType(ecs.ImageComponentType)
		if imageDetails != nil && len(imageDetails) == 1 {
			image = imageDetails[0].(*ecs_components.ImageComponent)
		} else {
			// Set default
			image = ecs_components.NewMissingImageComponent()
		}

		charModel.Image = models.CampaignImage{
			Name: image.Name,
			Url:  image.Url,
		}

		characters = append(characters, charModel)
	}

	// Sort the list to always show the same order
	sort.Sort(characters)

	data := make(map[string]any)
	data["chars"] = characters

	rawJsonBytes, err := json.Marshal(
		e.handleLoadHtmlBodyMultipleTemplateFiles(
			[]string{"campaignManageCharacters.html", "campaignManageCharacterSelectionBox.html"},
			"campaignManageCharacters", data))
	if err != nil {
		return err
	}

	manageChars := NewEventMessage()
	manageChars.Source = message.Source
	manageChars.Type = TypeManageCharacters
	manageChars.Body = string(rawJsonBytes)
	manageChars.Destinations = append(manageChars.Destinations, pool.GetLeadId())
	pool.TransmitEventMessage(manageChars)

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