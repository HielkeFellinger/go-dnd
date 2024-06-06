package game_engine

import (
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
				if err := chatCom.handleRollChatCommand(message, pool, clearedBody); err != nil {
					return err
				}
			}
		}

		if justPassTroughMessage {
			// Just pass message trough
			pool.TransmitEventMessage(message)
		}
	}

	return nil
}

func getChatMessage(message string) EventMessage {
	chatMessage := NewEventMessage()
	chatMessage.Source = ServerUser
	chatMessage.Body = message
	chatMessage.Type = TypeChatBroadcast
	return chatMessage
}
