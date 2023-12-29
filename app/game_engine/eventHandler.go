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
		if message.Type == TypeLoadCharacters || message.Type == TypeLoadGame {
			log.Printf("Building Message: %+v\n", message)
			var transmitMessage = EventMessage{}
			transmitMessage.Type = TypeLoadCharacters
			chars := []models.Character{
				{Name: "Kaas - 1"}, {Name: "Kaas - 2"},
			}
			data := make(map[string]any)
			data["chars"] = chars

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

	if message.Type >= TypeChatBroadcast && message.Type <= TypeChatWhisper {
		// Just pass message trough
		pool.TransmitEventMessage(message)
	}

	return nil
}
