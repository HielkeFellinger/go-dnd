package game_engine

import (
	"github.com/hielkefellinger/go-dnd/app/ecs"
	"github.com/hielkefellinger/go-dnd/app/models"
	"os"
	"strconv"
)

type Engine interface {
	GetWorld() ecs.World
	SaveWorld(world ecs.World, campaignId uint) error
	GetEventMessageHandler() EventMessageHandler
}

type baseEngine struct {
	World        ecs.World
	EventHandler EventMessageHandler
}

func (e *baseEngine) GetWorld() ecs.World {
	return e.World
}

func (e *baseEngine) GetEventMessageHandler() EventMessageHandler {
	return e.EventHandler
}

func (e *baseEngine) SaveWorld(world ecs.World, campaignId uint) error {
	baseLocation := os.Getenv("CAMPAIGN_DATA_DIR") + "/" + strconv.Itoa(int(campaignId))

	if _, err := os.Stat(baseLocation + "/save"); os.IsNotExist(err) {
		if mkSaveDirErr := os.MkdirAll(baseLocation+"/save", os.ModePerm); mkSaveDirErr != nil {
			return mkSaveDirErr
		}
	}
	if _, err := os.Stat(baseLocation + "/images"); os.IsNotExist(err) {
		if mkImageDirErr := os.MkdirAll(baseLocation+"/images", os.ModePerm); mkImageDirErr != nil {
			return mkImageDirErr
		}
	}

	gameFile := baseLocation + "/save/" + "campaign.yml"

	return saveGame(world, gameFile)
}

func InitGameEngine(campaign models.Campaign) Engine {
	var baseEngineInstance = baseEngine{}
	baseLocation := os.Getenv("CAMPAIGN_DATA_DIR") + "/" + strconv.Itoa(int(campaign.ID))
	gameFile := baseLocation + "/save/" + "campaign.yml"

	if _, err := os.Stat(gameFile); err == nil {
		campaign.GameFile = gameFile
	}

	if campaign.GameFile != "" {
		baseEngineInstance.World = loadGame(campaign.GameFile)
	} else {
		baseEngineInstance.World = loadGame(SpaceGame)
	}
	baseEngineInstance.EventHandler = &baseEventMessageHandler{}

	return &baseEngineInstance
}
