package game_engine

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hielkefellinger/go-dnd/app/ecs"
	"github.com/hielkefellinger/go-dnd/app/ecs_components"
	"github.com/hielkefellinger/go-dnd/app/helpers"
	"github.com/hielkefellinger/go-dnd/app/models"
	"golang.org/x/net/html"
	"log"
	"slices"
)

func (e *baseEventMessageHandler) handleMapUpdateEvents(message EventMessage, pool CampaignPool) error {
	log.Printf("- Map Update Event Type: '%d' Message: '%s'", message.Type, message.Id)
	var handled = false

	if message.Type == TypeUpdateMapEntity {
		err := e.typeUpdateMapEntity(message, pool)
		if err != nil {
			return err
		}
		handled = true
	}
	if message.Type == TypeUpdateMapVisibility {
		err := e.typeUpdateMapVisibility(message, pool)
		if err != nil {
			return err
		}
		handled = true
	}
	if message.Type == TypeAddMapItem {
		err := e.typeAddMapItem(message, pool)
		if err != nil {
			return err
		}
		handled = true
	}
	if message.Type == TypeRemoveMapItem {
		err := e.typeRemoveMapItem(message, pool)
		if err != nil {
			return err
		}
		handled = true
	}
	if message.Type == TypeSignalMapItem {
		err := e.typeSignalMapItem(message, pool)
		if err != nil {
			return err
		}
		handled = true
	}
	if message.Type == TypeChangeMapBackgroundImage {
		err := e.typeChangeMapBackgroundImage(message, pool)
		if err != nil {
			return err
		}
		handled = true
	}

	if !handled {
		return errors.New(fmt.Sprintf("message of type '%d' is not recognised by 'handleMapUpdateEvents()'", message.Type))
	}
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
			if ok := mapEntity.RemoveComponentByUuid(mapItemRelComponent.Id); ok {
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

	// Get an empty space on grid
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
	if err := mapEntity.AddComponent(newMapItemRelation); err != nil {
		return err
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

	// Update Clients, if any are found
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

		err := e.typeLoadMapEntity(updateMessage, pool)
		if err != nil {
			return err
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
			err := mapItemComponent.Entity.AddComponent(visibilityComponent)
			if err != nil {
				return err
			}
		}

		// Update possible Map Entities
		mapEntities := pool.GetEngine().GetWorld().GetMapEntities()
		for _, mapEntity := range mapEntities {

			// Only get the map with the relevant relation to entity
			if !mapEntity.HasRelationWithEntityByUuid(mapItemComponent.Entity.GetId()) {
				continue
			}

			for _, mapItem := range mapEntity.GetAllComponentsOfType(ecs.MapItemRelationComponentType) {
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
						// Get list of controlling users to use as exclude filter
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
