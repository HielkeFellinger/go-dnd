package ecs_components

import (
	"github.com/google/uuid"
	"github.com/hielkefellinger/go-dnd/app/ecs"
)

type PositionComponent struct {
	ecs.BaseComponent
	X int `yaml:"x"` // Column
	Y int `yaml:"y"` // Row
}

func NewPositionComponent() PositionComponent {
	return PositionComponent{
		BaseComponent: ecs.BaseComponent{Id: uuid.New()},
	}
}

func (c *PositionComponent) WithX(x int) *PositionComponent {
	c.X = x
	return c
}

func (c *PositionComponent) WithY(y int) *PositionComponent {
	c.Y = y
	return c
}

func (c *PositionComponent) ComponentType() uint64 {
	return ecs.PositionComponentType
}
