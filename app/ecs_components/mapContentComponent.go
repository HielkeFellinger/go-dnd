package ecs_components

import (
	"github.com/google/uuid"
	"github.com/hielkefellinger/go-dnd/app/ecs"
	"strconv"
)

type MapContentComponent struct {
	ecs.BaseComponent
	Active bool `yaml:"active"`
}

func NewMapContentComponent() ecs.Component {
	return &BlockerComponent{
		BaseComponent: ecs.BaseComponent{Id: uuid.New()},
		Active:        true,
	}
}

func (c *MapContentComponent) LoadFromRawComponent(raw ecs.RawComponent) error {
	loadedValues := 0
	if value, ok := raw.Params["active"]; ok {
		if err := c.ActiveFromString(value); err != nil {
			return err
		}
		loadedValues++
	}

	return c.CheckValuesParsedFromRaw(loadedValues, raw)
}

func (c *MapContentComponent) ParseToRawComponent() (ecs.RawComponent, error) {
	rawComponent := ecs.RawComponent{
		ComponentType: ecs.TypeNameToNthBit[c.ComponentType()].Name,
		Params: map[string]string{
			"active": strconv.FormatBool(c.Active),
		},
	}
	return rawComponent, nil
}

func (c *MapContentComponent) ComponentType() uint64 {
	return ecs.MapContentComponentType
}

func (c *MapContentComponent) ActiveFromString(bool string) error {
	b, err := strconv.ParseBool(bool)
	c.Active = b
	return err
}
