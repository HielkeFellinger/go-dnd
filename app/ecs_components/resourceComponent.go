package ecs_components

import (
	"github.com/google/uuid"
	"github.com/hielkefellinger/go-dnd/app/ecs"
)

type ResourceComponent struct {
	ecs.BaseComponent
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
}

func NewResourceComponent() ResourceComponent {
	return ResourceComponent{
		BaseComponent: ecs.BaseComponent{Id: uuid.New()},
	}
}

func (c *ResourceComponent) WidthName(name string) *ResourceComponent {
	c.Name = name
	return c
}

func (c *ResourceComponent) WidthDescription(description string) *ResourceComponent {
	c.Description = description
	return c
}

func (c *ResourceComponent) ComponentType() uint64 {
	return ecs.ResourceComponentType
}
