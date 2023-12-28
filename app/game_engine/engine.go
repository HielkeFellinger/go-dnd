package game_engine

import (
	"github.com/hielkefellinger/go-dnd/app/ecs"
	"github.com/hielkefellinger/go-dnd/app/models"
)

func InitCampaignWorld(campaign models.Campaign) ecs.World {
	if campaign.GameFile != "" {
		return loadGame(campaign.GameFile)
	}
	return loadGame(SpaceGame)
}
