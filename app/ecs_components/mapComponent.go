package ecs_components

import (
	"github.com/google/uuid"
	"github.com/hielkefellinger/go-dnd/app/ecs"
)

type MapComponent struct {
	ecs.BaseComponent
	Active bool
}

func NewMapComponent() MapComponent {
	return MapComponent{
		BaseComponent: ecs.BaseComponent{Id: uuid.New()},
		Active:        false,
	}
}

func (c *MapComponent) WidthActive(active bool) *MapComponent {
	c.Active = active
	return c
}

func (c *MapComponent) ComponentType() uint64 {
	return ecs.MapComponentType
}
