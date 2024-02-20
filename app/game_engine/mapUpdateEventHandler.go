package game_engine

import (
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/hielkefellinger/go-dnd/app/ecs_components"
	"github.com/hielkefellinger/go-dnd/app/models"
	"golang.org/x/net/html"
	"log"
)

func (e *baseEventMessageHandler) handleMapUpdateEvents(message EventMessage, pool CampaignPool) error {
	log.Printf("- Map Update Event Type: '%d' Message: '%s'", message.Type, message.Id)

	if message.Type == TypeUpdateMapEntity {
		err := e.updateMapEntity(message, pool)
		if err != nil {
			return nil
		}
	}
	if message.Type == TypeUpdateMapVisibility {
		return nil
	}

	return nil
}

func (e *baseEventMessageHandler) updateMapEntity(message EventMessage, pool CampaignPool) error {
	// Undo escaping
	clearedBody := html.UnescapeString(message.Body)

	// Attempt to parse the campaign screen map messageMapItem
	var messageMapItem models.CampaignScreenMapItemElement
	err := json.Unmarshal([]byte(clearedBody), &messageMapItem)
	if err != nil {
		return err
	}

	var mapUuid uuid.UUID
	if parsedUuid, err := uuid.Parse(messageMapItem.MapId); err == nil {
		mapUuid = parsedUuid
	} else {
		return err
	}

	var mapItemUuid uuid.UUID
	if parsedUuid, err := uuid.Parse(messageMapItem.Id); err == nil {
		mapItemUuid = parsedUuid
	} else {
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
		updateMessage.Source = "system"

		err = e.typeLoadMapEntity(updateMessage, pool)
		if err != nil {
			return err
		}
	}

	// - @todo Change visibility
	return nil
}
