package game_engine

import (
	"errors"
	"github.com/hielkefellinger/go-dnd/app/ecs"
	"github.com/hielkefellinger/go-dnd/app/ecs_model_translation"
	"github.com/hielkefellinger/go-dnd/app/helpers"
)

func (e *baseEventMessageHandler) typeLoadItemDetails(message EventMessage, pool CampaignPool) error {
	var transmitMessage = NewEventMessage()
	transmitMessage.Type = TypeLoadItemDetails
	transmitMessage.Source = message.Source
	transmitMessage.Destinations = append(transmitMessage.Destinations, message.Source)

	// Validate UUID Filter form message
	uuidCharFilter, err := helpers.ParseStringToUuid(message.Body)
	if err != nil {
		return err
	}

	// Test if item exists
	var itemEntity ecs.Entity
	itemEntity, ok := pool.GetEngine().GetWorld().GetItemEntityByUuid(uuidCharFilter)
	if !ok || itemEntity == nil {
		return errors.New("filter UUID has no match")
	}

	// Parse and check if the item could be parsed
	inventoryItem := ecs_model_translation.ItemEntityToCampaignInventoryItem(itemEntity, 1)

	if len(inventoryItem.Id) > 0 {
		data := make(map[string]any)
		data["Item"] = inventoryItem
		messageIdBody := EventMessageIdBody{
			Id:   inventoryItem.Id,
			Html: e.handleLoadHtmlBody("itemDetails.html", "itemDetails", data),
		}

		transmitMessage.Body = messageIdBody.ToBodyString()
		pool.TransmitEventMessage(transmitMessage)
	}
	return nil
}
