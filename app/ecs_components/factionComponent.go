package ecs_components

import (
	"github.com/google/uuid"
	"github.com/hielkefellinger/go-dnd/app/ecs"
)

type FactionComponent struct {
	ecs.BaseComponent
	Name        string
	Description string
}

func NewFactionComponent() FactionComponent {
	return FactionComponent{
		BaseComponent: ecs.BaseComponent{Id: uuid.New()},
	}
}

func (c *FactionComponent) WidthName(name string) *FactionComponent {
	c.Name = name
	return c
}

func (c *FactionComponent) WithDescription(description string) *FactionComponent {
	c.Description = description
	return c
}

func (c *FactionComponent) ComponentType() uint64 {
	return ecs.FactionComponentType
}
