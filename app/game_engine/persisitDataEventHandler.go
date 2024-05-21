package game_engine

import (
	"errors"
	"log"
)

func (e *baseEventMessageHandler) handlePersistDataEvents(message EventMessage, pool CampaignPool) error {
	log.Printf("- Persist Data Events Type: '%d' Message: '%s'", message.Type, message.Id)

	if message.Type == TypeGameSave {
		return e.typeGameSave(message, pool)
	}

	return nil
}

func (e *baseEventMessageHandler) typeGameSave(message EventMessage, pool CampaignPool) error {

	if message.Source != pool.GetLeadId() {
		return errors.New("saving game is not allowed as non-lead")
	}

	if err := pool.GetEngine().SaveWorld(pool.GetId()); err != nil {
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
