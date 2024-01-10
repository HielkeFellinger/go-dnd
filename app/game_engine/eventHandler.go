package game_engine

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hielkefellinger/go-dnd/app/ecs_model_translation"
	"github.com/hielkefellinger/go-dnd/app/models"
	"log"
	"strconv"
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
			var content = models.CampaignContentItem{}
			var tab = models.CampaignTabItem{}

			componentMap := ecs_model_translation.MapEntityToCampaignMapModel(mapEntity)
			tab.Id = componentMap.Id
			content.Id = componentMap.Id

			data := make(map[string]any)
			data["id"] = tab.Id
			data["name"] = componentMap.Name

			tab.Html = e.handleLoadHtmlBody("campaignSelector.html", "campaignSelector", data)

			// Add extra data, like chars
			x := componentMap.X
			y := componentMap.Y

			xVal := make([]string, x)
			yVal := make([]string, y)
			for i := range xVal {
				xVal[i] = strconv.Itoa(i)
			}
			for i := range yVal {
				yVal[i] = strconv.Itoa(i)
			}

			data["x"] = xVal
			data["y"] = yVal
			data["backgroundImage"] = componentMap.Image.Url
			content.Html = e.handleLoadHtmlBody("campaignContentMap.html", "campaignContentMap", data)

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
