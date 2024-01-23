package game_engine

import (
	"log"
)

func (e *baseEventMessageHandler) handleMapUpdateEvents(message EventMessage, pool CampaignPool) error {
	log.Printf("- Map. Event: '%s'", message.Id)

	if message.Type == TypeUpdateMapEntity {
		// @todo ^ Fix / implement

		// ...

	}

	return nil
}
