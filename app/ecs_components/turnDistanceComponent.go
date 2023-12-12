package ecs_components

import (
	"github.com/google/uuid"
	"github.com/hielkefellinger/go-dnd/app/ecs"
)

type TurnDistanceComponent struct {
	ecs.BaseComponent
	Distance uint `yaml:"distance"`
}

func NewTurnDistanceComponent() TurnDistanceComponent {
	return TurnDistanceComponent{
		BaseComponent: ecs.BaseComponent{Id: uuid.New()},
	}
}

func (c *TurnDistanceComponent) WidthDistance(distance uint) *TurnDistanceComponent {
	c.Distance = distance
	return c
}

func (c *TurnDistanceComponent) ComponentType() uint64 {
	return ecs.TurnDistanceComponentType
}
