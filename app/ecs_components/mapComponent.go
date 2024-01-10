package ecs_components

import (
	"github.com/google/uuid"
	"github.com/hielkefellinger/go-dnd/app/ecs"
	"strconv"
)

type MapComponent struct {
	ecs.BaseComponent
	Active bool `yaml:"active"`
}

func NewMapComponent() ecs.Component {
	return &MapComponent{
		BaseComponent: ecs.BaseComponent{Id: uuid.New()},
		Active:        false,
	}
}

func (c *MapComponent) LoadFromRawComponent(raw ecs.RawComponent) error {
	loadedValues := 0
	if value, ok := raw.Params["active"]; ok {
		if err := c.ActiveFromString(value); err != nil {
			return err
		}
		loadedValues++
	}

	return c.CheckValuesParsedFromRaw(loadedValues, raw)
}

func (c *MapComponent) ActiveFromString(bool string) error {
	b, err := strconv.ParseBool(bool)
	c.Active = b
	return err
}

func (c *MapComponent) ComponentType() uint64 {
	return ecs.MapComponentType
}
