package ecs_components

import (
	"github.com/google/uuid"
	"github.com/hielkefellinger/go-dnd/app/ecs"
)

type SlotsComponent struct {
	ecs.BaseComponent
	Count uint `yaml:"count"`
}

func NewSlotsComponent() SlotsComponent {
	return SlotsComponent{
		BaseComponent: ecs.BaseComponent{Id: uuid.New()},
	}
}

func (c *SlotsComponent) WithCounter(count uint) *SlotsComponent {
	c.Count = count
	return c
}

func (c *SlotsComponent) ComponentType() uint64 {
	return ecs.SlotsComponentType
}
