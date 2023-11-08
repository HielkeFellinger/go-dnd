package ecs_components

import (
	"github.com/google/uuid"
	"github.com/hielkefellinger/go-dnd/app/ecs"
)

type AreaComponent struct {
	ecs.BaseComponent
	Length int
	Width  int
}

func NewAreaComponent() AreaComponent {
	return AreaComponent{
		BaseComponent: ecs.BaseComponent{Id: uuid.New()},
	}
}

func (c *AreaComponent) WithLength(length int) *AreaComponent {
	c.Length = length
	return c
}

func (c *AreaComponent) WithWidth(width int) *AreaComponent {
	c.Width = width
	return c
}

func (c *AreaComponent) ComponentType() uint64 {
	return ecs.AreaComponentType
}
