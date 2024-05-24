package game_engine

import (
	"github.com/hielkefellinger/go-dnd/app/ecs"
	"github.com/hielkefellinger/go-dnd/app/models"
	"os"
	"strconv"
)

type Engine interface {
	GetWorld() ecs.World
	SaveWorld(campaignId uint) error
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

func (e *baseEngine) SaveWorld(campaignId uint) error {
	baseLocation := os.Getenv("CAMPAIGN_DATA_DIR") + "/" + strconv.Itoa(int(campaignId))

	if _, err := os.Stat(baseLocation + "/save"); os.IsNotExist(err) {
		if err := os.MkdirAll(baseLocation+"/save", os.ModePerm); err != nil {
			return err
		}
	}
	if _, err := os.Stat(baseLocation + "/images"); os.IsNotExist(err) {
		if err := os.MkdirAll(baseLocation+"/images", os.ModePerm); err != nil {
			return err
		}
	}

	gameFile := baseLocation + "/save/" + "campaign.yml"

	return saveGame(e.World, gameFile)
}

func InitGameEngine(campaign models.Campaign) Engine {
	var baseEngine = baseEngine{}
	baseLocation := os.Getenv("CAMPAIGN_DATA_DIR") + "/" + strconv.Itoa(int(campaign.ID))
	gameFile := baseLocation + "/save/" + "campaign.yml"

	if _, err := os.Stat(gameFile); err == nil {
		campaign.GameFile = gameFile
	}

	if campaign.GameFile != "" {
		baseEngine.World = loadGame(campaign.GameFile)
	} else {
		baseEngine.World = loadGame(SpaceGame)
	}
	baseEngine.EventHandler = &baseEventMessageHandler{}

	return &baseEngine
}
