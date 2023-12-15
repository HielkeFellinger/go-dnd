package ecs_components

import (
	"github.com/google/uuid"
	"github.com/hielkefellinger/go-dnd/app/ecs"
	"strconv"
)

type RangeComponent struct {
	ecs.BaseComponent
	Min int `yaml:"min"`
	Max int `yaml:"max"`
}

func NewRangeComponent() ecs.Component {
	return &RangeComponent{
		BaseComponent: ecs.BaseComponent{Id: uuid.New()},
	}
}

func (c *RangeComponent) LoadFromRawComponent(raw ecs.RawComponent) error {
	loadedValues := 0
	if value, ok := raw.Params["min"]; ok {
		if err := c.MinFromString(value); err != nil {
			return err
		}
		loadedValues++
	}
	if value, ok := raw.Params["max"]; ok {
		if err := c.MaxFromString(value); err != nil {
			return err
		}
		loadedValues++
	}

	return c.CheckValuesParsedFromRaw(loadedValues, raw)
}

func (c *RangeComponent) MinFromString(min string) error {
	n, err := strconv.Atoi(min)
	c.Min = n
	return err
}

func (c *RangeComponent) MaxFromString(max string) error {
	n, err := strconv.Atoi(max)
	c.Max = n
	return err
}

func (c *RangeComponent) ComponentType() uint64 {
	return ecs.RangeComponentType
}
