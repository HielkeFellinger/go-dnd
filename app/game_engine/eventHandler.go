package game_engine

import (
	"bytes"
	"fmt"
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
		err := e.handleMapEvents(message, pool)
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

		// Load Focus Map Related Details
		// - Gray out non-present players;

		var characters []models.Character
		for _, charEntity := range charEntities {
			characters = append(characters, models.Character{Name: charEntity.GetName()})
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
