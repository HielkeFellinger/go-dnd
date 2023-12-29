package game_engine

import (
	"github.com/hielkefellinger/go-dnd/app/ecs"
	"github.com/hielkefellinger/go-dnd/app/models"
)

type Engine interface {
	GetWorld() ecs.World
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

func InitGameEngine(campaign models.Campaign) Engine {
	var baseEngine = baseEngine{}

	if campaign.GameFile != "" {
		baseEngine.World = loadGame(campaign.GameFile)
	} else {
		baseEngine.World = loadGame(SpaceGame)
	}
	baseEngine.EventHandler = &baseEventMessageHandler{}

	return &baseEngine
}
