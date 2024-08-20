package game_engine

import (
	"errors"
	"github.com/hielkefellinger/go-dnd/app/ecs"
	"github.com/hielkefellinger/go-dnd/app/ecs_components"
	"github.com/hielkefellinger/go-dnd/app/ecs_model_translation"
	"github.com/hielkefellinger/go-dnd/app/helpers"
	"github.com/hielkefellinger/go-dnd/app/models"
	"log"
	"slices"
	"sort"
)

func (e *baseEventMessageHandler) handleLoadCharacterEvents(message EventMessage, pool CampaignPool) error {
	log.Printf("- Char. Load Event Type: '%d' Message: '%s'", message.Type, message.Id)
	if message.Type == TypeLoadCharacters || message.Type == TypeLoadFullGame {
		e.loadCharacters(message, pool)
	}
	if message.Type == TypeLoadCharactersDetails {
		return e.loadCharactersDetails(message, pool)
	}

	return nil
}

func (e *baseEventMessageHandler) loadCharactersDetails(message EventMessage, pool CampaignPool) error {
	var transmitMessage = NewEventMessage()
	transmitMessage.Type = TypeLoadCharactersDetails
	transmitMessage.Source = message.Source

	isLead := message.Source == pool.GetLeadId()

	// Validate UUID Filter form message
	uuidCharFilter, err := helpers.ParseStringToUuid(message.Body)
	if err != nil {
		return err
	}

	// Test if Character exists
	var charEntity ecs.Entity
	charEntity, ok := pool.GetEngine().GetWorld().GetCharacterEntityByUuid(uuidCharFilter)
	if !ok || charEntity == nil {
		return errors.New("filter UUID has no match")
	}

	// Parse and check if the character could be parsed
	campaignCharacter := ecs_model_translation.CharacterEntityToCampaignCharacterModel(charEntity)

	// Send only to people allowed to view this character
	if message.Source == ServerUser {
		transmitMessage.Destinations = append(transmitMessage.Destinations, pool.GetLeadId())
		transmitMessage.Destinations = append(transmitMessage.Destinations, campaignCharacter.Controllers...)
	} else if isLead || slices.Contains(campaignCharacter.Controllers, message.Source) {
		transmitMessage.Destinations = append(transmitMessage.Destinations, message.Source)
	}

	if len(transmitMessage.Destinations) > 0 && len(campaignCharacter.Id) > 0 {
		data := make(map[string]any)
		data["character"] = campaignCharacter

		messageIdBody := EventMessageIdBody{
			Id: uuidCharFilter.String(),
			Html: e.handleLoadHtmlBodyMultipleTemplateFiles(
				[]string{"characterDetails.html", "inventory.html"},
				"characterDetails", data),
		}

		transmitMessage.Body = messageIdBody.ToBodyString()
	} else {
		transmitMessage.Destinations = append(transmitMessage.Destinations, message.Source)
	}

	pool.TransmitEventMessage(transmitMessage)
	return nil
}

func (e *baseEventMessageHandler) loadCharacters(message EventMessage, pool CampaignPool) {
	var transmitMessage = NewEventMessage()
	transmitMessage.Type = TypeLoadCharacters
	transmitMessage.Source = message.Source

	// On TypeLoadFullGame only load for relevant player
	if message.Type == TypeLoadFullGame {
		transmitMessage.Destinations = append(transmitMessage.Destinations, message.Source)
	} else {
		transmitMessage.Destinations = message.Destinations
	}

	charEntities := pool.GetEngine().GetWorld().GetCharacterEntities()
	allPlayers := pool.GetAllClientIds()

	var characters models.Characters
	for _, charEntity := range charEntities {
		// Only show Player Characters
		if !charEntity.HasComponentType(ecs.PlayerComponentType) {
			continue
		}

		// Check if (one of n of controlling) player(s) is online
		playerOnline := false
		for _, rawPlayerComponent := range charEntity.GetAllComponentsOfType(ecs.PlayerComponentType) {
			playerComponent := rawPlayerComponent.(*ecs_components.PlayerComponent)
			playerOnline = playerOnline || slices.Contains(allPlayers, playerComponent.Name)
		}

		var image *ecs_components.ImageComponent
		var imageDetails = charEntity.GetAllComponentsOfType(ecs.ImageComponentType)
		if imageDetails != nil && len(imageDetails) == 1 {
			image = imageDetails[0].(*ecs_components.ImageComponent)
		} else {
			// Set default
			image = ecs_components.NewMissingImageComponent()
		}

		characters = append(characters, models.Character{
			Id:     charEntity.GetId().String(),
			Name:   charEntity.GetName(),
			Online: playerOnline,
			Image: models.CampaignImage{
				Name: image.Name,
				Url:  image.Url,
			},
		})
	}
	// Sort the list to always show the same order
	sort.Sort(characters)

	data := make(map[string]any)
	data["chars"] = characters

	transmitMessage.Body = e.handleLoadHtmlBody("characterRibbon.html", "characterRibbon", data)
	pool.TransmitEventMessage(transmitMessage)
}
