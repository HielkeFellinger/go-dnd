package game_engine

import (
	"bytes"
	"encoding/json"
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

func (e *baseEventMessageHandler) handleMapEvents(message EventMessage, pool CampaignPool) error {
	log.Printf("- Map. Event: '%s'", message.Id)
	if message.Type == TypeLoadMap || message.Type == TypeLoadFullGame {
		var transmitMessage = NewEventMessage()
		transmitMessage.Type = TypeLoadMap

		// Load Focus Map Related Details
		// - Gray out non-present players;

		var campaignScreenContent = models.NewCampaignScreenContent()
		mapEntities := pool.GetEngine().GetWorld().GetMapEntities()

		log.Printf("----- Maps: '%v'", mapEntities)
		for _, mapEntity := range mapEntities {
			var tab = models.CampaignTabItem{}
			var content = models.CampaignContentItem{}

			tab.Id = mapEntity.GetId().String()
			content.Id = mapEntity.GetId().String()

			data := make(map[string]any)
			data["id"] = tab.Id
			data["name"] = mapEntity.GetName()

			tab.Html = e.handleLoadHtmlBody("campaignSelector.html", "campaignSelector", data)
			content.Html = e.handleLoadHtmlBody("campaignContent.html", "campaignContent", data)

			campaignScreenContent.Tabs = append(campaignScreenContent.Tabs, tab)
			campaignScreenContent.Content = append(campaignScreenContent.Content, content)
		}

		rawJsonBytes, err := json.Marshal(campaignScreenContent)
		if err != nil {
			log.Printf("Error parsing Loading Map content `%s`", err.Error())
		}

		transmitMessage.Body = string(rawJsonBytes)
		pool.TransmitEventMessage(transmitMessage)
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
