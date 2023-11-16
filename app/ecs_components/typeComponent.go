package ecs_components

import (
	"github.com/google/uuid"
	"github.com/hielkefellinger/go-dnd/app/ecs"
)

type TypeComponent struct {
	ecs.BaseComponent
	Name        string
	Description string
}

func NewTypeComponent() TypeComponent {
	return TypeComponent{
		BaseComponent: ecs.BaseComponent{Id: uuid.New()},
	}
}

func (c *TypeComponent) WithName(name string) *TypeComponent {
	c.Name = name
	return c
}

func (c *TypeComponent) WithDescription(Description string) *TypeComponent {
	c.Description = Description
	return c
}

func (c *TypeComponent) ComponentType() uint64 {
	return ecs.TypeComponentType
}
