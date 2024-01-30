package game_engine

import (
	"bytes"
	"fmt"
	"github.com/hielkefellinger/go-dnd/app/ecs"
	"github.com/hielkefellinger/go-dnd/app/ecs_components"
	"github.com/hielkefellinger/go-dnd/app/models"
	"log"
	"text/template"
)

type EventMessageHandler interface {
	HandleEventMessage(message EventMessage, pool CampaignPool) error
}

type baseEventMessageHandler struct {
}

func (e *baseEventMessageHandler) HandleEventMessage(message EventMessage, pool CampaignPool) error {
	log.Printf("Message Handler Parsing: %+v\n", message)

	if message.Type == TypeLoadFullGame || (message.Type >= TypeLoadCharacters && message.Type <= TypeRemoveCharacter) {
		err := e.handleCharacterEvents(message, pool)
		if err != nil {
			return err
		}
	}

	if message.Type == TypeLoadFullGame || (message.Type >= TypeLoadMap && message.Type <= TypeRemoveMap) {
		err := e.handleMapLoadEvents(message, pool)
		if err != nil {
			return err
		}
	}

	if message.Type == TypeUpdateMapEntity {
		err := e.handleMapUpdateEvents(message, pool)
		if err != nil {
			return err
		}
	}

	if message.Type >= TypeChatBroadcast && message.Type <= TypeChatWhisper {
		// Just pass message trough
		pool.TransmitEventMessage(message)
	}

	return nil
}

func (e *baseEventMessageHandler) handleCharacterEvents(message EventMessage, pool CampaignPool) error {
	log.Printf("- Char. Event: '%s'", message.Id)
	if message.Type == TypeLoadCharacters || message.Type == TypeLoadFullGame {

		var transmitMessage = NewEventMessage()
		transmitMessage.Type = TypeLoadCharacters

		charEntities := pool.GetEngine().GetWorld().GetCharacterEntities()

		// Check if GM/DM if not filter non-player controlled characters

		// @todo Load Focus Map Related Details
		// - Gray out non-present players;

		var characters []models.Character
		for _, charEntity := range charEntities {

			var image *ecs_components.ImageComponent
			var imageDetails = charEntity.GetAllComponentsOfType(ecs.ImageComponentType)
			if imageDetails != nil && len(imageDetails) == 1 {
				image = imageDetails[0].(*ecs_components.ImageComponent)
			} else {
				// Set default
				image = ecs_components.NewImageComponent().(*ecs_components.ImageComponent)
				image.Name = "MISSING IMAGE"
				image.Url = "/images/unknown_item.png"
			}

			characters = append(characters, models.Character{
				Name: charEntity.GetName(),
				Image: models.CampaignImage{
					Name: image.Name,
					Url:  image.Url,
				},
			})
		}

		data := make(map[string]any)
		data["chars"] = characters

		transmitMessage.Body = e.handleLoadHtmlBody("characterRibbon.html", "chars", data)
		pool.TransmitEventMessage(transmitMessage)
	}
	return nil
}

func (e *baseEventMessageHandler) handleLoadHtmlBody(fileName string, templateName string, data map[string]any) string {
	var buf bytes.Buffer
	tmpl := template.Must(template.ParseFiles(fmt.Sprintf("web/templates/%s", fileName)))
	err := tmpl.ExecuteTemplate(&buf, templateName, data)
	if err != nil {
		log.Printf("Error parsing %s `%s`", fileName, err.Error())
	}
	return string(buf.Bytes())
}
