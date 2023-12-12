package ecs_components

import (
	"github.com/google/uuid"
	"github.com/hielkefellinger/go-dnd/app/ecs"
)

type CharacterComponent struct {
	ecs.BaseComponent
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
}

func NewCharacterComponent() CharacterComponent {
	return CharacterComponent{
		BaseComponent: ecs.BaseComponent{Id: uuid.New()},
	}
}

func (c *CharacterComponent) WidthName(name string) *CharacterComponent {
	c.Name = name
	return c
}

func (c *CharacterComponent) WidthDescription(description string) *CharacterComponent {
	c.Description = description
	return c
}

func (c *CharacterComponent) ComponentType() uint64 {
	return ecs.CharacterComponentType
}
