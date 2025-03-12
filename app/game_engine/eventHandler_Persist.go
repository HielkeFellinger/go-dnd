package game_engine

import (
	"errors"
	"fmt"
	"log"
)

func (e *baseEventMessageHandler) handlePersistDataEvents(message EventMessage, pool CampaignPool) error {
	log.Printf("- Persist Data Events Type: '%d' Message: '%s'", message.Type, message.Id)

	if message.Type == TypeGameSave {
		return e.typeGameSave(message, pool)
	}

	return errors.New(fmt.Sprintf("message of type '%d' is not recognised by 'handlePersistDataEvents()'", message.Type))
}

func (e *baseEventMessageHandler) typeGameSave(message EventMessage, pool CampaignPool) error {

	if message.Source != pool.GetLeadId() {
		return errors.New("saving game is not allowed as non-lead")
	}

	if err := pool.GetEngine().SaveWorld(pool.GetEngine().GetWorld(), pool.GetId()); err != nil {
		return err
	}

	// Trigger send game saved message
	updateMessage := NewEventMessage()
	updateMessage.Source = ServerUser
	updateMessage.Body = "Game Has been saved"
	updateMessage.Type = TypeChatServerMsg
	pool.TransmitEventMessage(updateMessage)

	return nil
}
