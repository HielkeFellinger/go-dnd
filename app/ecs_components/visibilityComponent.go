package ecs_components

import (
	"github.com/google/uuid"
	"github.com/hielkefellinger/go-dnd/app/ecs"
)

type VisibilityComponent struct {
	ecs.BaseComponent
	Hidden bool `yaml:"hidden"`
}

func NewVisibilityComponent() VisibilityComponent {
	return VisibilityComponent{
		BaseComponent: ecs.BaseComponent{Id: uuid.New()},
	}
}

func (c *VisibilityComponent) WidthHidden(hidden bool) *VisibilityComponent {
	c.Hidden = hidden
	return c
}

func (c *VisibilityComponent) ComponentType() uint64 {
	return ecs.VisibilityComponentType
}
