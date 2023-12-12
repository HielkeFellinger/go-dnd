package ecs_components

import (
	"github.com/google/uuid"
	"github.com/hielkefellinger/go-dnd/app/ecs"
)

type TransportComponent struct {
	ecs.BaseComponent
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
}

func NewTransportComponent() TransportComponent {
	return TransportComponent{
		BaseComponent: ecs.BaseComponent{Id: uuid.New()},
	}
}

func (c *TransportComponent) WidthName(name string) *TransportComponent {
	c.Name = name
	return c
}

func (c *TransportComponent) WidthDescription(description string) *TransportComponent {
	c.Description = description
	return c
}

func (c *TransportComponent) ComponentType() uint64 {
	return ecs.TransportComponentType
}
