package ecs_components

import (
	"github.com/google/uuid"
	"github.com/hielkefellinger/go-dnd/app/ecs"
)

type ItemComponent struct {
	ecs.BaseComponent
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
}

func NewItemComponent() ItemComponent {
	return ItemComponent{
		BaseComponent: ecs.BaseComponent{Id: uuid.New()},
	}
}

func (c *ItemComponent) WidthName(name string) *ItemComponent {
	c.Name = name
	return c
}

func (c *ItemComponent) WidthDescription(description string) *ItemComponent {
	c.Description = description
	return c
}

func (c *ItemComponent) ComponentType() uint64 {
	return ecs.ItemComponentType
}
