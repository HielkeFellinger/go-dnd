package ecs_components

import (
	"github.com/google/uuid"
	"github.com/hielkefellinger/go-dnd/app/ecs"
)

type RangeComponent struct {
	ecs.BaseComponent
	Min int `yaml:"min"`
	Max int `yaml:"max"`
}

func NewRangeComponent() RangeComponent {
	return RangeComponent{
		BaseComponent: ecs.BaseComponent{Id: uuid.New()},
	}
}

func (c *RangeComponent) WithMin(min int) *RangeComponent {
	c.Min = min
	return c
}

func (c *RangeComponent) WithMax(max int) *RangeComponent {
	c.Max = max
	return c
}

func (c *RangeComponent) ComponentType() uint64 {
	return ecs.RangeComponentType
}
