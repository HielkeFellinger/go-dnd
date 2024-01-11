package game_engine

import (
	"encoding/json"
	"github.com/hielkefellinger/go-dnd/app/ecs_model_translation"
	"github.com/hielkefellinger/go-dnd/app/models"
	"log"
	"strconv"
)

func (e *baseEventMessageHandler) handleMapEvents(message EventMessage, pool CampaignPool) error {
	log.Printf("- Map. Event: '%s'", message.Id)
	if message.Type == TypeLoadMap || message.Type == TypeLoadFullGame {
		var transmitMessage = NewEventMessage()
		transmitMessage.Type = TypeLoadMap
		transmitMessage.Destinations = append(transmitMessage.Destinations, message.Source)

		// Check if is GM:
		isLead := message.Source == pool.GetLeadId()

		// Load Focus Map Related Details
		// - Gray out non-present players;

		var campaignScreenContent = models.NewCampaignScreenContent()
		mapEntities := pool.GetEngine().GetWorld().GetMapEntities()

		for _, mapEntity := range mapEntities {
			// Translate
			componentMap := ecs_model_translation.MapEntityToCampaignMapModel(mapEntity)

			// Only show enabled maps for player
			if !componentMap.Enabled && !isLead {
				continue
			}

			// Only show filtered form body
			if len(message.Body) > 0 && componentMap.Id != message.Body {
				continue
			}

			var data = buildMapData(componentMap, isLead)
			var content = models.CampaignContentItem{}
			var tab = models.CampaignTabItem{}

			tab.Id = componentMap.Id
			content.Id = componentMap.Id
			tab.Html = e.handleLoadHtmlBody("campaignSelector.html", "campaignSelector", data)
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

func buildMapData(model models.CampaignMap, isLead bool) map[string]any {
	data := make(map[string]any)
	data["id"] = model.Id
	data["name"] = model.Name
	data["lead"] = isLead
	xVal := make([]string, model.X)
	yVal := make([]string, model.Y)
	for i := range xVal {
		xVal[i] = strconv.Itoa(i)
	}
	for i := range yVal {
		yVal[i] = strconv.Itoa(i)
	}

	data["x"] = xVal
	data["y"] = yVal
	data["backgroundImage"] = model.Image.Url

	return data
}
