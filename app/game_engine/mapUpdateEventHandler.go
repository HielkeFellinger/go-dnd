package game_engine

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hielkefellinger/go-dnd/app/ecs"
	"github.com/hielkefellinger/go-dnd/app/ecs_components"
	"github.com/hielkefellinger/go-dnd/app/models"
	"golang.org/x/net/html"
	"log"
	"slices"
)

func (e *baseEventMessageHandler) handleMapUpdateEvents(message EventMessage, pool CampaignPool) error {
	log.Printf("- Map Update Event Type: '%d' Message: '%s'", message.Type, message.Id)

	if message.Type == TypeUpdateMapEntity {
		err := e.typeUpdateMapEntity(message, pool)
		if err != nil {
			return err
		}
	}
	if message.Type == TypeUpdateMapVisibility {
		err := e.typeUpdateMapVisibility(message, pool)
		if err != nil {
			return err
		}
	}
	if message.Type == TypeAddMapItem {
		err := e.typeAddMapItem(message, pool)
		if err != nil {
			return err
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

	var newMapEntity models.AddMapItem
	err := json.Unmarshal([]byte(clearedBody), &newMapEntity)
	if err != nil {
		return err
	}

	entityUuid, err := parseStingToUuid(newMapEntity.EntityId)
	if err != nil {
		return err
	}
	mapUuid, err := parseStingToUuid(newMapEntity.MapId)
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
	posissions := make([]string, len(mapItems))
	for _, mapItem := range mapItems {
		mapItemRelation := mapItem.(*ecs_components.MapItemRelationComponent)

		// Test if char is not already linked
		if mapItemRelation.Entity != nil && characterEntity.GetId() == mapItemRelation.Entity.GetId() {
			return errors.New("character already added to that map")
		}

		// Reserve position
		if mapItemRelation.Position != nil {
			posissions = append(posissions, fmt.Sprintf("%d-%d", mapItemRelation.Position.X, mapItemRelation.Position.Y))
		}
	}

	// Build the new NewMapItemRelationComponent
	newMapItemRelation := ecs_components.NewMapItemRelationComponent().(*ecs_components.MapItemRelationComponent)
	newMapItemRelation.Entity = characterEntity

	// Get an empty space on grid
	for x := 0; x < int(mapArea.Width); x++ {
		for y := 0; y < int(mapArea.Length); y++ {
			if !slices.Contains(posissions, fmt.Sprintf("%d-%d", x, y)) {
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

	var mapActivity models.SetActivity
	err := json.Unmarshal([]byte(clearedBody), &mapActivity)
	if err != nil {
		return err
	}

	mapUuid, err := parseStingToUuid(mapActivity.Id)
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

	mapUuid, err := parseStingToUuid(messageMapItem.MapId)
	if err != nil {
		return err
	}
	mapItemUuid, err := parseStingToUuid(messageMapItem.Id)
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

		err = e.typeLoadMapEntity(updateMessage, pool)
		if err != nil {
			return err
		}
	}

	// - @todo Change visibility

	return nil
}
