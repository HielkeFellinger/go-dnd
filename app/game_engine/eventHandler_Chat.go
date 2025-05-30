package game_engine

import (
	"errors"
	"fmt"
	"golang.org/x/net/html"
	"log"
	"strings"
)

func (e *baseEventMessageHandler) handleChatEventMessage(message EventMessage, pool CampaignPool) error {
	log.Printf("- Chat Data Events Type: '%d' Message: '%s'", message.Type, message.Id)

	if message.Type >= TypeChatBroadcast && message.Type <= TypeChatWhisper {
		justPassTroughMessage := true
		if message.Type == TypeChatBroadcast {
			// Undo escaping
			chatCom := chatCommands{}
			clearedBody := html.UnescapeString(message.Body)

			if strings.HasPrefix(clearedBody, "/roll") {
				justPassTroughMessage = false
				return chatCom.handleRollChatCommand(message, pool, clearedBody)
			}
		}

		if justPassTroughMessage {
			// Just pass message trough
			pool.TransmitEventMessage(message)
			return nil
		}
	}

	return errors.New(fmt.Sprintf("message of type '%d' is not recognised by 'handleChatEventMessage()'", message.Type))
}

func getChatMessage(message string) EventMessage {
	chatMessage := NewEventMessage()
	chatMessage.Source = ServerUser
	chatMessage.Body = message
	chatMessage.Type = TypeChatBroadcast
	return chatMessage
}
