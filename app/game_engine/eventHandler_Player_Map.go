package game_engine

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hielkefellinger/go-dnd/app/ecs"
	"github.com/hielkefellinger/go-dnd/app/ecs_builders"
	"github.com/hielkefellinger/go-dnd/app/ecs_components"
	"github.com/hielkefellinger/go-dnd/app/ecs_model_translation"
	"github.com/hielkefellinger/go-dnd/app/helpers"
	"github.com/hielkefellinger/go-dnd/app/models"
	"golang.org/x/net/html"
	"log"
	"slices"
	"strconv"
)

func (e *baseEventMessageHandler) typeLoadMapEntity(message EventMessage, pool CampaignPool) error {
	// No filter in the body equals no map entity to load
	if len(message.Body) == 0 {
		return nil
	}

	// Validate Filter
	uuidMapItemFilter, err := helpers.ParseStringToUuid(message.Body)
	if err != nil {
		return err
	}

	var transmitMessage = NewEventMessage()
	transmitMessage.Type = TypeLoadMapEntity
	transmitMessage.Source = message.Source

	isLead := message.Source == pool.GetLeadId()

	mapEntities := pool.GetEngine().GetWorld().GetMapEntities()
	for _, mapEntity := range mapEntities {
		// Only get the map with the relevant entity
		if !mapEntity.HasComponentByUuid(uuidMapItemFilter) {
			continue
		}

		componentMap := ecs_model_translation.MapEntityToCampaignMapModel(mapEntity)
		mapItem, ok := mapEntity.GetComponentByUuid(uuidMapItemFilter)
		if !ok || mapItem == nil {
			return errors.New("failure loading MapItem")
		}

		// Translate Entity to controlling
		var mapItemModel = ecs_model_translation.MapItemEntityToCampaignMapItemElement(mapItem, componentMap.Id)

		// Lead Message
		mapItemModel = e.buildMapItem(mapItemModel, true, true)
		transmitMessage.Body = string(e.parseObjectToJson(mapItemModel))
		transmitMessage.Destinations = make([]string, 1)
		transmitMessage.Destinations[0] = pool.GetLeadId()
		pool.TransmitEventMessage(transmitMessage)

		// Get list of controlling players
		controllingPlayers := make([]string, 0)
		if componentMap.Active {
			controllingPlayers = append(controllingPlayers, mapItemModel.Controllers...)
		}

		// Controlling users Message
		if len(controllingPlayers) > 0 {
			mapItemModel = e.buildMapItem(mapItemModel, false, true)
			transmitMessage.Body = string(e.parseObjectToJson(mapItemModel))
			transmitMessage.Destinations = controllingPlayers
			pool.TransmitEventMessage(transmitMessage)
		}

		// - @todo team visibility
		// Check hidden non owner
		if mapItemModel.Hidden && !isLead && !slices.Contains(mapItemModel.Controllers, message.Source) {
			continue
		}

		managingPlayersFilter := append(controllingPlayers, pool.GetLeadId())
		nonControllingPlayers := pool.GetAllClientIds(managingPlayersFilter...)
		if componentMap.Active && len(nonControllingPlayers) > 0 {
			mapItemModel = e.buildMapItem(mapItemModel, false, false)
			transmitMessage.Body = string(e.parseObjectToJson(mapItemModel))
			transmitMessage.Destinations = nonControllingPlayers

			log.Printf("    Non-controlling players: %v", transmitMessage.Destinations)
			pool.TransmitEventMessage(transmitMessage)
		}
	}
	return nil
}

func (e *baseEventMessageHandler) typeLoadMapEntities(message EventMessage, pool CampaignPool) error {
	var transmitMessage = NewEventMessage()
	transmitMessage.Type = TypeLoadMapEntities
	transmitMessage.Source = message.Source
	transmitMessage.Destinations = append(transmitMessage.Destinations, message.Source)

	// Check if is GM:
	isLead := message.Source == pool.GetLeadId()

	mapEntities := pool.GetEngine().GetWorld().GetMapEntities()
	for _, mapEntity := range mapEntities {
		// Only show filtered form body
		if len(message.Body) > 0 && mapEntity.GetId().String() != message.Body {
			continue
		}

		componentMap := ecs_model_translation.MapEntityToCampaignMapModel(mapEntity)

		// Only show enabled maps for player
		if !componentMap.Active && !isLead {
			continue
		}

		// load models
		mapItems := mapEntity.GetAllComponentsOfType(ecs.MapItemRelationComponentType)
		log.Printf("Map mapItems count: %v", len(mapItems))
		mapItemsModel := models.CampaignScreenMapItems{
			MapId:    componentMap.Id,
			Elements: make(map[string]models.CampaignScreenMapItemElement, len(mapItems)),
		}

		// Translate all items
		for _, mapItem := range mapItems {
			var mapItemModel = ecs_model_translation.MapItemEntityToCampaignMapItemElement(mapItem, mapItemsModel.MapId)

			// - @todo team visibility
			// Check hidden non owner
			if mapItemModel.Hidden && !isLead && !slices.Contains(mapItemModel.Controllers, message.Source) {
				continue
			}

			mapItemModel = e.buildMapItem(mapItemModel, isLead, isLead || slices.Contains(mapItemModel.Controllers, message.Source))
			mapItemsModel.Elements[mapItemModel.Id] = mapItemModel
		}

		rawJsonBytes, err := json.Marshal(mapItemsModel)
		if err != nil {
			return err
		}

		transmitMessage.Body = string(rawJsonBytes)
		pool.TransmitEventMessage(transmitMessage)
	}
	return nil
}

func (e *baseEventMessageHandler) typeLoadMap(message EventMessage, pool CampaignPool) error {
	var transmitMessage = NewEventMessage()
	transmitMessage.Type = TypeLoadMap
	transmitMessage.Source = message.Source
	transmitMessage.Destinations = append(transmitMessage.Destinations, message.Source)

	// Check if is GM:
	isLead := message.Source == pool.GetLeadId()

	// @todo Load Focus Map Related Details
	// - Gray out non-present players;

	var campaignScreenContent = models.NewCampaignScreenContent()
	mapEntities := pool.GetEngine().GetWorld().GetMapEntities()
	charEntities := pool.GetEngine().GetWorld().GetCharacterEntities()
	for _, mapEntity := range mapEntities {
		// Only handle if id is requested from body
		if len(message.Body) > 0 && mapEntity.GetId().String() != message.Body {
			continue
		}

		// Translate
		componentMap := ecs_model_translation.MapEntityToCampaignMapModel(mapEntity)

		// Only show enabled maps for player (if not directly requested by lead)
		if !componentMap.Active && (!isLead || (isLead && len(message.Body) == 0)) {
			continue
		}

		var data = e.buildMapData(componentMap, isLead, charEntities)
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
		return err
	}

	transmitMessage.Body = string(rawJsonBytes)
	pool.TransmitEventMessage(transmitMessage)
	return nil
}

func (e *baseEventMessageHandler) typeChangeMapBackgroundImage(message EventMessage, pool CampaignPool) error {
	// Undo escaping
	clearedBody := html.UnescapeString(message.Body)

	// Check if user is lead
	if message.Source != pool.GetLeadId() {
		return errors.New("adding items to map is not allowed as non-lead")
	}

	var activeMapBackground ActiveMapBackground
	if err := json.Unmarshal([]byte(clearedBody), &activeMapBackground); err != nil {
		return err
	}

	componentUuid, err := helpers.ParseStringToUuid(activeMapBackground.ImageId)
	if err != nil {
		return err
	}
	mapUuid, err := helpers.ParseStringToUuid(activeMapBackground.MapId)
	if err != nil {
		return err
	}

	mapEntity, match := pool.GetEngine().GetWorld().GetMapEntityByUuid(mapUuid)
	if !match || mapEntity == nil {
		return errors.New("failure of loading MAP by UUID")
	}

	if !mapEntity.HasComponentByUuid(componentUuid) {
		return errors.New("requested image component does not exist on map")
	}

	// Get map size
	hasUpdate := false
	var updateImage ecs_components.ImageComponent
	backgroundImages := mapEntity.GetAllComponentsOfType(ecs.ImageComponentType)
	for _, backgroundImage := range backgroundImages {
		image := backgroundImage.(*ecs_components.ImageComponent)
		if image.Id == componentUuid {
			if !image.Active {
				image.Active = true
				updateImage = *image
				hasUpdate = true
			}
		} else {
			image.Active = false
		}
	}

	if hasUpdate {
		rawJsonBytes, err := json.Marshal(SendNewMapBackgroundImage{
			MapId: activeMapBackground.MapId,
			Id:    activeMapBackground.ImageId,
			Url:   updateImage.Url,
		})
		if err != nil {
			return err
		}
		var updateMessage = NewEventMessage()
		updateMessage.Type = TypeChangeMapBackgroundImage
		updateMessage.Body = string(rawJsonBytes)
		updateMessage.Source = message.Source
		pool.TransmitEventMessage(updateMessage)
	}

	return nil
}

func (e *baseEventMessageHandler) typeSignalMapItem(message EventMessage, pool CampaignPool) error {
	// Undo escaping
	clearedBody := html.UnescapeString(message.Body)

	var sendSignal SendSignal
	if err := json.Unmarshal([]byte(clearedBody), &sendSignal); err != nil {
		return err
	}

	// Get the map and its MapItemRelationComponent and remove it
	mapUuid, err := helpers.ParseStringToUuid(sendSignal.Id)
	if err != nil {
		return err
	}
	mapEntity, match := pool.GetEngine().GetWorld().GetMapEntityByUuid(mapUuid)
	if !match || mapEntity == nil {
		return errors.New("failure of loading MAP by UUID")
	}

	data := make(map[string]any)
	data["id"] = sendSignal.Id
	data["x"] = sendSignal.X
	data["y"] = sendSignal.Y

	if sendSignal.Type == "danger" {
		sendSignal.Html = e.handleLoadHtmlBody("signalDangerContent.html", "signalDangerContent", data)
	} else if sendSignal.Type == "warn" {
		sendSignal.Html = e.handleLoadHtmlBody("signalWarnContent.html", "signalWarnContent", data)
	} else {
		sendSignal.Html = e.handleLoadHtmlBody("signalInfoContent.html", "signalInfoContent", data)
	}

	rawJsonBytes, err := json.Marshal(sendSignal)
	if err != nil {
		return err
	}

	signalMapItemMessage := NewEventMessage()
	signalMapItemMessage.Source = message.Source
	signalMapItemMessage.Destinations = pool.GetAllClientIds()
	signalMapItemMessage.Type = TypeSignalMapItem
	signalMapItemMessage.Body = string(rawJsonBytes)
	pool.TransmitEventMessage(signalMapItemMessage)

	return nil
}

func (e *baseEventMessageHandler) typeMapInteraction(message EventMessage, pool CampaignPool) error {
	// Undo escaping @ TODO: Interaction!
	clearedBody := html.UnescapeString(message.Body)

	var mapInteraction MapInteraction
	if err := json.Unmarshal([]byte(clearedBody), &mapInteraction); err != nil {
		return err
	}

	// Check if user is lead
	if message.Source != pool.GetLeadId() {
		return errors.New("removing items to map is not allowed as non-lead")
	}

	// Get Map
	mapUuid, err := helpers.ParseStringToUuid(mapInteraction.MapId)
	if err != nil {
		return err
	}
	mapEntity, match := pool.GetEngine().GetWorld().GetMapEntityByUuid(mapUuid)
	if !match || mapEntity == nil {
		return errors.New("failure of loading MAP by UUID")
	}

	// CheckInteraction
	reqX, err := strconv.ParseUint(mapInteraction.X, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid X coordinate: %v", err)
	}
	reqY, err := strconv.ParseUint(mapInteraction.Y, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid Y coordinate: %v", err)
	}

	// Get map size
	var maxX uint = 0
	var maxY uint = 0
	if area := mapEntity.GetFirstComponentsOfType(ecs.AreaComponentType); area != nil {
		areaComponent := area.(*ecs_components.AreaComponent)
		maxX = areaComponent.Width
		maxY = areaComponent.Length
	}
	if maxY == 0 || maxX == 0 {
		return errors.New("map has no Width or Length; do not change map")
	}

	// Check if interaction is inside the area
	interactionX := uint(reqX)
	interactionY := uint(reqY)
	if interactionX >= maxX || interactionY >= maxY {
		return errors.New("Request is outside of ")
	}

	// Collect all Elements on the map; clean possible old missing links / off map items
	mapGrid := make([][]ecs.RelationalComponent, maxX)
	for i := range mapGrid {
		mapGrid[i] = make([]ecs.RelationalComponent, maxY)
	}
	types := []uint64{ecs.MapItemRelationComponentType, ecs.MapItemRelationComponentType}
	mapContent := mapEntity.GetAllComponentsOfTypes(types)
	for _, item := range mapContent {
		itemRel := item.(ecs.RelationalComponent)
		cleanup := false
		if itemRel.GetEntity() == nil {
			cleanup = true
		} else if itemRel.ComponentType() == ecs.MapItemRelationComponentType {
			mapItem := item.(*ecs_components.MapItemRelationComponent)
			// Check position
			if mapItem.Position.CheckIfOnArea(maxX, maxY) {
				mapGrid[mapItem.Position.X][mapItem.Position.Y] = mapItem
			} else {
				cleanup = true
			}
		} else if itemRel.ComponentType() == ecs.MapLinkRelationComponentType {
			mapLink := item.(*ecs_components.MapLinkRelationComponent)
			if mapLink.SourcePosition.CheckIfOnArea(maxX, maxY) {
				mapGrid[mapLink.SourcePosition.X][mapLink.SourcePosition.Y] = mapLink
			} else {
				cleanup = true
			}
		}

		// Item is no longer valid or is missing a position; clear from map
		if cleanup {
			mapEntity.RemoveComponentByUuid(item.GetId())
		}
	}

	addIds := make([]string, 0)
	removeIds := make([]string, 0)

	// Add/Remove Blocker
	if mapInteraction.Type == Blocker {
		// If the position has blocker Remove; else Add if space
		gridRelComponent := mapGrid[interactionX][interactionY]
		if gridRelComponent == nil {
			if newRelId, errBlockAdd := ecs_builders.AddBasicBlockerToMap(
				pool.GetEngine().GetWorld(), mapEntity, interactionX, interactionY); errBlockAdd != nil {
				return errBlockAdd
			} else {
				addIds = append(addIds, newRelId)
			}
		} else if gridRelComponent.GetEntity().HasComponentType(ecs.BlockerComponentType) {
			removeIds = append(removeIds, gridRelComponent.GetId().String())
		} else {
			// Ignore
		}

	} else if mapInteraction.Type == AddBlocker {
		// @TODO
	} else if mapInteraction.Type == RemoveBlocker {
		// @TODO
	}

	// Add/Remove Portal (Map link to another map)
	if mapInteraction.Type == AddPortal {
		// @TODO
	} else if mapInteraction.Type == RemovePortal {
		// @TODO
	}

	// Interact with stuff
	log.Printf("Map interaction: %v", mapInteraction)

	for _, addId := range addIds {
		// Trigger update of map
		updateMessage := NewEventMessage()
		updateMessage.Source = ServerUser
		updateMessage.Body = addId
		updateMessage.Destinations = pool.GetAllClientIds()
		updateMessage.Type = TypeLoadMapEntity
		_ = e.typeLoadMapEntity(updateMessage, pool)
	}
	for _, removeId := range removeIds {
		removeMapItemJson, _ := json.Marshal(RemoveMapItem{
			MapId:     mapEntity.GetId().String(),
			MapItemId: removeId,
		})
		removeMapItemMessage := NewEventMessage()
		removeMapItemMessage.Source = pool.GetLeadId()
		removeMapItemMessage.Type = TypeRemoveMapItem
		removeMapItemMessage.Destinations = pool.GetAllClientIds()
		removeMapItemMessage.Body = string(removeMapItemJson)
		if removeErr := e.typeRemoveMapItem(removeMapItemMessage, pool); removeErr != nil {
			log.Printf("Err removing map item: %v", removeErr.Error())
		}
	}

	return nil
}

func (e *baseEventMessageHandler) typeRemoveMapItem(message EventMessage, pool CampaignPool) error {
	// Undo escaping
	clearedBody := html.UnescapeString(message.Body)

	// Check if user is lead
	if message.Source != pool.GetLeadId() {
		return errors.New("removing items to map is not allowed as non-lead")
	}

	var removeMapItem RemoveMapItem
	if err := json.Unmarshal([]byte(clearedBody), &removeMapItem); err != nil {
		return err
	}
	mapUuid, err := helpers.ParseStringToUuid(removeMapItem.MapId)
	if err != nil {
		return err
	}
	mapItemUuid, err := helpers.ParseStringToUuid(removeMapItem.MapItemId)
	if err != nil {
		return err
	}

	// Get the map and its MapItemRelationComponent and remove it
	mapEntity, match := pool.GetEngine().GetWorld().GetMapEntityByUuid(mapUuid)
	if !match || mapEntity == nil {
		return errors.New("failure of loading MAP by UUID")
	}

	// Check on the component
	if component, ok := mapEntity.GetComponentByUuid(mapItemUuid); ok {
		if component.ComponentType() == ecs.MapItemRelationComponentType {
			mapItemRelComponent := component.(*ecs_components.MapItemRelationComponent)
			if oke := mapEntity.RemoveComponentByUuid(mapItemRelComponent.Id); oke {
				removeMapItemMessage := NewEventMessage()
				removeMapItemMessage.Source = ServerUser
				removeMapItemMessage.Type = TypeRemoveMapItem
				removeMapItemMessage.Destinations = make([]string, 0)
				removeMapItemMessage.Body = mapItemRelComponent.Id.String()
				pool.TransmitEventMessage(removeMapItemMessage)
				return nil
			}
		}
	}
	return nil
}

func (e *baseEventMessageHandler) typeAddMapItem(message EventMessage, pool CampaignPool) error {
	// Undo escaping
	clearedBody := html.UnescapeString(message.Body)

	// Check if user is lead
	if message.Source != pool.GetLeadId() {
		return errors.New("adding items to map is not allowed as non-lead")
	}

	var newMapEntity AddMapItem
	if err := json.Unmarshal([]byte(clearedBody), &newMapEntity); err != nil {
		return err
	}

	entityUuid, err := helpers.ParseStringToUuid(newMapEntity.EntityId)
	if err != nil {
		return err
	}
	mapUuid, err := helpers.ParseStringToUuid(newMapEntity.MapId)
	if err != nil {
		return err
	}

	// @todo allow adding additional (non-char) items
	characterEntity, match := pool.GetEngine().GetWorld().GetCharacterEntityByUuid(entityUuid)
	if !match || characterEntity == nil {
		return errors.New("failure of loading CHARACTER by UUID")
	}

	mapEntity, match := pool.GetEngine().GetWorld().GetMapEntityByUuid(mapUuid)
	if !match || mapEntity == nil {
		return errors.New("failure of loading MAP by UUID")
	}

	// Get map size
	var mapArea *ecs_components.AreaComponent
	mapAreas := mapEntity.GetAllComponentsOfType(ecs.AreaComponentType)
	if len(mapAreas) > 0 {
		mapArea = mapAreas[0].(*ecs_components.AreaComponent)
	} else {
		return errors.New("cloud not find place for player on map")
	}

	// Get all
	mapItems := mapEntity.GetAllComponentsOfType(ecs.MapItemRelationComponentType)
	positions := make([]string, len(mapItems))
	for _, mapItem := range mapItems {
		mapItemRelation := mapItem.(*ecs_components.MapItemRelationComponent)

		// Test if char is not already linked
		if mapItemRelation.Entity != nil && characterEntity.GetId() == mapItemRelation.Entity.GetId() {
			return errors.New("character already added to that map")
		}

		// Reserve position
		if mapItemRelation.Position != nil {
			positions = append(positions, fmt.Sprintf("%d-%d", mapItemRelation.Position.X, mapItemRelation.Position.Y))
		}
	}

	// Build the new NewMapItemRelationComponent
	newMapItemRelation := ecs_components.NewMapItemRelationComponent().(*ecs_components.MapItemRelationComponent)
	newMapItemRelation.Entity = characterEntity

	// Get an empty default on grid
	for x := 0; x < int(mapArea.Width); x++ {
		for y := 0; y < int(mapArea.Length); y++ {
			if !slices.Contains(positions, fmt.Sprintf("%d-%d", x, y)) {
				var position = ecs_components.NewPositionComponent().(*ecs_components.PositionComponent)
				position.X = uint(x)
				position.Y = uint(y)
				newMapItemRelation.Position = position
				break
			}
		}
		if newMapItemRelation.Position != nil {
			break
		}
	}

	// Add if possible
	if errAddComponent := mapEntity.AddComponent(newMapItemRelation); errAddComponent != nil {
		return errAddComponent
	}

	// Trigger update of map
	updateMessage := NewEventMessage()
	updateMessage.Source = ServerUser
	updateMessage.Body = newMapItemRelation.Id.String()
	updateMessage.Type = TypeLoadMapEntity
	return e.typeLoadMapEntity(updateMessage, pool)
}

func (e *baseEventMessageHandler) typeUpdateMapVisibility(message EventMessage, pool CampaignPool) error {
	// Undo escaping
	clearedBody := html.UnescapeString(message.Body)

	var mapActivity SetActivity
	err := json.Unmarshal([]byte(clearedBody), &mapActivity)
	if err != nil {
		return err
	}

	mapUuid, err := helpers.ParseStringToUuid(mapActivity.Id)
	if err != nil {
		return err
	}

	mapEntity, match := pool.GetEngine().GetWorld().GetMapEntityByUuid(mapUuid)
	if !match || mapEntity == nil {
		return errors.New("failure of loading MAP by UUID")
	}

	// Update Entity
	mapComponents := mapEntity.GetAllComponentsOfType(ecs.MapComponentType)
	if len(mapComponents) >= 1 {
		mapComponent := mapComponents[0].(*ecs_components.MapComponent)
		mapComponent.Active = mapActivity.Active
	}

	// Update Clients if any are found
	rawJsonBytes, err := json.Marshal(mapActivity)
	if err != nil {
		return err
	}
	var updateMessage = NewEventMessage()
	updateMessage.Type = TypeUpdateMapVisibility
	updateMessage.Body = string(rawJsonBytes)
	updateMessage.Destinations = pool.GetAllClientIds(pool.GetLeadId())
	if len(updateMessage.Destinations) > 0 {
		pool.TransmitEventMessage(updateMessage)
	}

	return nil
}

func (e *baseEventMessageHandler) typeUpdateMapEntity(message EventMessage, pool CampaignPool) error {
	// Undo escaping
	clearedBody := html.UnescapeString(message.Body)

	// Attempt to parse the campaign screen map messageMapItem
	var messageMapItem models.CampaignScreenMapItemElement
	err := json.Unmarshal([]byte(clearedBody), &messageMapItem)
	if err != nil {
		return err
	}

	mapUuid, err := helpers.ParseStringToUuid(messageMapItem.MapId)
	if err != nil {
		return err
	}
	mapItemUuid, err := helpers.ParseStringToUuid(messageMapItem.Id)
	if err != nil {
		return err
	}

	mapEntity, match := pool.GetEngine().GetWorld().GetMapEntityByUuid(mapUuid)
	if !match || mapEntity == nil {
		return errors.New("failure of loading MAP by UUID")
	}
	rawMapItemComponent, ok := mapEntity.GetComponentByUuid(mapItemUuid)
	if !ok || rawMapItemComponent == nil {
		return errors.New("failure loading MapItemComponent")
	}
	mapItemComponent := rawMapItemComponent.(*ecs_components.MapItemRelationComponent)

	messagePosition := ecs_components.NewPositionComponent().(*ecs_components.PositionComponent)
	err = messagePosition.XFromString(messageMapItem.Position.X)
	if err != nil {
		return err
	}
	err = messagePosition.YFromString(messageMapItem.Position.Y)
	if err != nil {
		return err
	}

	// Test and update the properties
	// - Get the new X and Y and update the
	if mapItemComponent.Position.Y != messagePosition.Y || mapItemComponent.Position.X != messagePosition.X {
		// Update Position
		mapItemComponent.Position.Y = messagePosition.Y
		mapItemComponent.Position.X = messagePosition.X

		// Trigger sending update
		var updateMessage = NewEventMessage()
		updateMessage.Type = TypeLoadMapEntity
		updateMessage.Body = mapItemComponent.Id.String()
		updateMessage.Source = ServerUser

		errLoadMap := e.typeLoadMapEntity(updateMessage, pool)
		if errLoadMap != nil {
			return errLoadMap
		}
	}

	// Update visibility on Entity
	if mapItemComponent.Entity != nil {
		updatedVisibility := false
		visibilities := mapItemComponent.Entity.GetAllComponentsOfType(ecs.VisibilityComponentType)
		if len(visibilities) > 0 {
			visibilityComponent := visibilities[0].(*ecs_components.VisibilityComponent)
			updatedVisibility = visibilityComponent.Hidden != messageMapItem.Hidden
			visibilityComponent.Hidden = messageMapItem.Hidden
		} else {
			visibilityComponent := ecs_components.NewVisibilityComponent().(*ecs_components.VisibilityComponent)
			visibilityComponent.Hidden = messageMapItem.Hidden
			updatedVisibility = true
			errAddComponent := mapItemComponent.Entity.AddComponent(visibilityComponent)
			if errAddComponent != nil {
				return errAddComponent
			}
		}

		// Update possible Map Entities
		mapEntities := pool.GetEngine().GetWorld().GetMapEntities()
		for _, mapUpdateEntity := range mapEntities {

			// Only get the map with the relevant relation to entity
			if !mapUpdateEntity.HasRelationWithEntityByUuid(mapItemComponent.Entity.GetId()) {
				continue
			}

			for _, mapItem := range mapUpdateEntity.GetAllComponentsOfType(ecs.MapItemRelationComponentType) {
				mapItemRelComponent := mapItem.(*ecs_components.MapItemRelationComponent)

				if mapItemRelComponent.Entity.GetId() == mapItemComponent.Entity.GetId() {
					reloadMapItemMessage := NewEventMessage()
					reloadMapItemMessage.Source = ServerUser
					reloadMapItemMessage.Body = mapItemRelComponent.Id.String()
					reloadMapItemErr := e.typeLoadMapEntity(reloadMapItemMessage, pool)
					if reloadMapItemErr != nil {
						return reloadMapItemErr
					}

					// Remove old visible ghost mapItem of non-lead / non-controlling users
					if messageMapItem.Hidden && updatedVisibility {
						// Get list of controlling users to use as an exclude filter
						controllers := make([]string, 0)
						controllers = append(controllers, pool.GetLeadId())
						players := mapItemRelComponent.Entity.GetAllComponentsOfType(ecs.PlayerComponentType)
						for _, player := range players {
							controllers = append(controllers, player.(*ecs_components.PlayerComponent).Name)
						}

						removeMapItemMessage := NewEventMessage()
						removeMapItemMessage.Source = ServerUser
						removeMapItemMessage.Type = TypeRemoveMapItem
						removeMapItemMessage.Destinations = pool.GetAllClientIds(controllers...)
						removeMapItemMessage.Body = mapItemRelComponent.Id.String()
						if len(removeMapItemMessage.Destinations) > 0 {
							pool.TransmitEventMessage(removeMapItemMessage)
						}
					}
				}
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

// @TODO allow rendering of other non-player entities!
func (e *baseEventMessageHandler) buildMapItem(mapItemModel models.CampaignScreenMapItemElement, isLead bool, hasControl bool) models.CampaignScreenMapItemElement {
	data := make(map[string]any)
	data["id"] = mapItemModel.Id
	data["type"] = mapItemModel.Type
	data["mapId"] = mapItemModel.MapId
	data["hidden"] = mapItemModel.Hidden
	data["lead"] = isLead
	data["hasControl"] = hasControl
	data["entityName"] = mapItemModel.EntityName
	data["backgroundImage"] = mapItemModel.Image.Url

	if mapItemModel.HasHealth() {
		data["healthBar"] = true
		data["healthTotal"] = mapItemModel.Health.Total
		data["healthCurrent"] = mapItemModel.Health.Current

		data["healthColour"] = "green"
		if mapItemModel.Health.Current < (int(mapItemModel.Health.Total) / 2) {
			data["healthColour"] = "yellow"
		}
		if mapItemModel.Health.Current < int(float64(mapItemModel.Health.Total)*0.334) {
			data["healthColour"] = "red"
		}
	}

	mapItemModel.Html = e.handleLoadHtmlBody("campaignContentMapCell.html", "campaignContentMapCell", data)
	return mapItemModel
}

func (e *baseEventMessageHandler) buildMapData(model models.CampaignMap, isLead bool, characters []ecs.Entity) map[string]any {
	data := make(map[string]any)
	data["type"] = "Map"
	data["id"] = model.Id
	data["name"] = model.Name
	data["lead"] = isLead
	data["active"] = model.Active
	data["backgroundImage"] = model.ActiveImage
	data["altImage"] = model.Images
	data["characters"] = characters

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
	return data
}

type MapInteractionType string

const (
	Blocker       MapInteractionType = "Blocker"
	AddBlocker    MapInteractionType = "AddBlocker"
	RemoveBlocker MapInteractionType = "RemoveBlocker"
	AddPortal     MapInteractionType = "AddPortal"
	RemovePortal  MapInteractionType = "RemovePortal"
)

type MapInteraction struct {
	MapId string             `json:"mapId"`
	Type  MapInteractionType `json:"type"`
	X     string             `json:"cellX"`
	Y     string             `json:"cellY"`
}

type SendSignal struct {
	Id   string `json:"Id"`
	X    string `json:"X"`
	Y    string `json:"Y"`
	Type string `json:"Type"`
	Html string `json:"Html"`
}

type SendNewMapBackgroundImage struct {
	Id    string `json:"Id"`
	MapId string `json:"MapId"`
	Url   string `json:"Url"`
}

type SetActivity struct {
	Id     string `json:"Id"`
	Active bool   `json:"Active"`
}

type AddMapItem struct {
	EntityId string `json:"EntityId"`
	MapId    string `json:"MapId"`
}

type ActiveMapBackground struct {
	ImageId string `json:"ImageId"`
	MapId   string `json:"MapId"`
}

type RemoveMapItem struct {
	MapId     string `json:"MapId"`
	MapItemId string `json:"Id"`
}
