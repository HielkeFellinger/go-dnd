package ecs_components

import (
	"github.com/google/uuid"
	"github.com/hielkefellinger/go-dnd/app/ecs"
	"strconv"
)

type VisibilityComponent struct {
	ecs.BaseComponent
	Hidden bool `yaml:"hidden"`
}

func NewVisibilityComponent() ecs.Component {
	return &VisibilityComponent{
		BaseComponent: ecs.BaseComponent{Id: uuid.New()},
	}
}

func (c *VisibilityComponent) LoadFromRawComponent(raw ecs.RawComponent) error {
	loadedValues := 0
	if value, ok := raw.Params["hidden"]; ok {
		if err := c.HiddenFromString(value); err != nil {
			return err
		}
		loadedValues++
	}

	return c.CheckValuesParsedFromRaw(loadedValues, raw)
}

func (c *VisibilityComponent) HiddenFromString(bool string) error {
	b, err := strconv.ParseBool(bool)
	c.Hidden = b
	return err
}

func (c *VisibilityComponent) ComponentType() uint64 {
	return ecs.VisibilityComponentType
}
