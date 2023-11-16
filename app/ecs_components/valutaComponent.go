package ecs_components

import (
	"github.com/google/uuid"
	"github.com/hielkefellinger/go-dnd/app/ecs"
)

type ValutaComponent struct {
	ecs.BaseComponent
	Name        string
	Description string
}

func NewValutaComponent() ValutaComponent {
	return ValutaComponent{
		BaseComponent: ecs.BaseComponent{Id: uuid.New()},
	}
}

func (c *ValutaComponent) WidthName(name string) *ValutaComponent {
	c.Name = name
	return c
}

func (c *ValutaComponent) WidthDescription(description string) *ValutaComponent {
	c.Description = description
	return c
}

func (c *ValutaComponent) ComponentType() uint64 {
	return ecs.ValutaComponentType
}
