package ecs_components

import (
	"github.com/google/uuid"
	"github.com/hielkefellinger/go-dnd/app/ecs"
	"strconv"
)

type SlotsComponent struct {
	ecs.BaseComponent
	Count uint `yaml:"count"`
}

func NewSlotsComponent() ecs.Component {
	return &SlotsComponent{
		BaseComponent: ecs.BaseComponent{Id: uuid.New()},
	}
}

func (c *SlotsComponent) LoadFromRawComponent(raw ecs.RawComponent) error {
	loadedValues := 0
	if value, ok := raw.Params["count"]; ok {
		if err := c.CounterFromString(value); err != nil {
			return err
		}
		loadedValues++
	}

	return c.CheckValuesParsedFromRaw(loadedValues, raw)
}

func (c *SlotsComponent) CounterFromString(counter string) error {
	n, err := strconv.Atoi(counter)
	c.Count = uint(n)
	return err
}

func (c *SlotsComponent) ComponentType() uint64 {
	return ecs.SlotsComponentType
}
