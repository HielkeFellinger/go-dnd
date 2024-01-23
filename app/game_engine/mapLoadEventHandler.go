package game_engine

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/hielkefellinger/go-dnd/app/ecs"
	"github.com/hielkefellinger/go-dnd/app/ecs_model_translation"
	"github.com/hielkefellinger/go-dnd/app/models"
	"log"
	"slices"
	"strconv"
)

func (e *baseEventMessageHandler) handleMapLoadEvents(message EventMessage, pool CampaignPool) error {
	log.Printf("- Map. Event: '%s'", message.Id)
	if message.Type == TypeLoadMap || message.Type == TypeLoadFullGame {
		var transmitMessage = NewEventMessage()
		transmitMessage.Type = TypeLoadMap
		transmitMessage.Destinations = append(transmitMessage.Destinations, message.Source)

		// Check if is GM:
		isLead := message.Source == pool.GetLeadId()

		// @todo Load Focus Map Related Details
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
	if message.Type == TypeLoadMapEntities || message.Type == TypeLoadFullGame {
		var transmitMessage = NewEventMessage()
		transmitMessage.Type = TypeLoadMapEntities
		transmitMessage.Destinations = append(transmitMessage.Destinations, message.Source)

		// Check if is GM:
		isLead := message.Source == pool.GetLeadId()

		mapEntities := pool.GetEngine().GetWorld().GetMapEntities()
		for _, mapEntity := range mapEntities {

			componentMap := ecs_model_translation.MapEntityToCampaignMapModel(mapEntity)

			// Only show enabled maps for player
			if !componentMap.Enabled && !isLead {
				continue
			}

			// Only show filtered form body
			if len(message.Body) > 0 && componentMap.Id != message.Body {
				continue
			}

			// load models
			mapItems := mapEntity.GetAllComponentsOfType(ecs.MapItemRelationComponentType)
			mapItemsModel := models.CampaignScreenMapItems{
				MapId:    componentMap.Id,
				Elements: make(map[string]models.CampaignScreenMapItemElement, len(mapItems)),
			}

			// Translate all items
			for _, mapItem := range mapItems {
				var mapItemModel = ecs_model_translation.MapItemEntityToCampaignMapItemElement(mapItem, mapItemsModel.MapId)
				mapItemModel = e.buildMapItem(mapItemModel, isLead || slices.Contains(mapItemModel.Controllers, message.Source))
				mapItemsModel.Elements[mapItemModel.Id] = mapItemModel
			}

			rawJsonBytes, err := json.Marshal(mapItemsModel)
			if err != nil {
				log.Printf("Error parsing Loading Map Item content `%s`", err.Error())
			}

			transmitMessage.Body = string(rawJsonBytes)
			pool.TransmitEventMessage(transmitMessage)
		}
	}
	if message.Type == TypeLoadMapEntity {
		// No filter in body equals no map entity to load
		if len(message.Body) == 0 {
			return nil
		}

		// Validate Filter
		var uuidMapItemFilter uuid.UUID
		if savedUuid, err := uuid.Parse(message.Body); err == nil {
			uuidMapItemFilter = savedUuid
		} else {
			return nil
		}

		var transmitMessage = NewEventMessage()
		transmitMessage.Type = TypeLoadMapEntity

		mapEntities := pool.GetEngine().GetWorld().GetMapEntities()
		for _, mapEntity := range mapEntities {

			// Only get the map with the relevant entity
			if mapEntity.HasComponentByUuid(uuidMapItemFilter) {
				continue
			}
			componentMap := ecs_model_translation.MapEntityToCampaignMapModel(mapEntity)
			mapItem := mapEntity.GetComponentByUuid(uuidMapItemFilter)

			// Translate Entity to controlling
			var mapItemModel = ecs_model_translation.MapItemEntityToCampaignMapItemElement(mapItem, componentMap.Id)

			// Get list of controlling (default the GM + controlling player if map is enabled)
			controllingPlayers := make([]string, 1)
			controllingPlayers[0] = pool.GetLeadId()
			if componentMap.Enabled {
				controllingPlayers = append(controllingPlayers, mapItemModel.Controllers...)
			}

			// Controlling users
			mapItemModel = e.buildMapItem(mapItemModel, true)
			transmitMessage.Body = string(e.parseObjectToJson(mapItemModel))
			transmitMessage.Destinations = controllingPlayers
			pool.TransmitEventMessage(transmitMessage)

			// Non-controlling users (only if available)
			if componentMap.Enabled {
				mapItemModel = e.buildMapItem(mapItemModel, false)
				transmitMessage.Body = string(e.parseObjectToJson(mapItemModel))
				transmitMessage.Destinations = pool.GetAllClientIds(controllingPlayers...)
				pool.TransmitEventMessage(transmitMessage)
			}
		}
	}

	return nil
}

func (e *baseEventMessageHandler) parseObjectToJson(object any) []byte {
	rawJsonBytes, err := json.Marshal(object)
	if err != nil {
		log.Printf("Error parsing content to Json `%s`", err.Error())
	}
	return rawJsonBytes
}

func (e *baseEventMessageHandler) buildMapItem(mapItemModel models.CampaignScreenMapItemElement, hasControl bool) models.CampaignScreenMapItemElement {
	data := make(map[string]any)
	data["id"] = mapItemModel.Id
	data["mapId"] = mapItemModel.MapId
	data["hasControl"] = hasControl
	data["entityName"] = mapItemModel.EntityName
	data["backgroundImage"] = mapItemModel.Image.Url
	data["healthPercentage"] = 70

	mapItemModel.Html = e.handleLoadHtmlBody("campaignContentMapCell.html", "campaignContentMapCell", data)
	return mapItemModel
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
