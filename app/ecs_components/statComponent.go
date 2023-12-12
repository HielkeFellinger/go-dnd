package ecs_components

import (
	"github.com/google/uuid"
	"github.com/hielkefellinger/go-dnd/app/ecs"
)

type StatComponent struct {
	ecs.BaseComponent
	Name  string `yaml:"name"`
	Value uint   `yaml:"value"`
}

func NewStatComponent() StatComponent {
	return StatComponent{
		BaseComponent: ecs.BaseComponent{Id: uuid.New()},
	}
}

func (c *StatComponent) WidthName(name string) *StatComponent {
	c.Name = name
	return c
}

func (c *StatComponent) WidthValue(value uint) *StatComponent {
	c.Value = value
	return c
}

func (c *StatComponent) ComponentType() uint64 {
	return ecs.StatComponentType
}
