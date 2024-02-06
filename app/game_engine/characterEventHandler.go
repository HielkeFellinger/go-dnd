package game_engine

import (
	"github.com/hielkefellinger/go-dnd/app/ecs"
	"github.com/hielkefellinger/go-dnd/app/ecs_components"
	"github.com/hielkefellinger/go-dnd/app/models"
	"log"
	"slices"
	"sort"
)

func (e *baseEventMessageHandler) handleCharacterEvents(message EventMessage, pool CampaignPool) error {
	log.Printf("- Char. Event Type: '%d' Message: '%s'", message.Type, message.Id)
	if message.Type == TypeLoadCharacters || message.Type == TypeLoadFullGame {
		e.loadCharacters(message, pool)
	}
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

	transmitMessage.Body = e.handleLoadHtmlBody("characterRibbon.html", "chars", data)
	pool.TransmitEventMessage(transmitMessage)
}
