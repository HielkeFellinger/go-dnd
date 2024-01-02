package game_engine

import (
	"bytes"
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

	if message.Type >= TypeLoadGame && message.Type <= TypeRemoveCharacter {
		return e.handleGameLoadEvents(message, pool)
	}

	if message.Type >= TypeChatBroadcast && message.Type <= TypeChatWhisper {
		// Just pass message trough
		pool.TransmitEventMessage(message)
	}

	return nil
}

func (e *baseEventMessageHandler) handleGameLoadEvents(message EventMessage, pool CampaignPool) error {
	if message.Type == TypeLoadCharacters || message.Type == TypeLoadGame {
		log.Printf("Building Message: %+v\n", message)
		var transmitMessage = EventMessage{}
		transmitMessage.Type = TypeLoadCharacters

		charEntities := pool.GetEngine().GetWorld().GetCharacterEntities()
		characters := make([]models.Character, len(charEntities))

		for _, charEntity := range charEntities {
			characters = append(characters, models.Character{Name: charEntity.GetName()})
		}

		data := make(map[string]any)
		data["chars"] = characters

		var buf bytes.Buffer
		tmpl := template.Must(template.ParseFiles("web/templates/test.html"))
		err := tmpl.ExecuteTemplate(&buf, "chars", data)
		if err != nil {
			log.Printf("Error parsing test.html `%s`", err.Error())
		}
		transmitMessage.Body = string(buf.Bytes())
		log.Printf("Build Message: %+v\n", transmitMessage)
		pool.TransmitEventMessage(transmitMessage)
	}
	return nil
}
