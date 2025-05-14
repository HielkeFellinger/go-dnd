package ecs_components

import (
	"github.com/google/uuid"
	"github.com/hielkefellinger/go-dnd/app/ecs"
	"strconv"
)

type BlockerComponent struct {
	ecs.BaseComponent
	Active bool `yaml:"active"`
}

func NewBlockerComponent() ecs.Component {
	return &BlockerComponent{
		BaseComponent: ecs.BaseComponent{Id: uuid.New()},
		Active:        true,
	}
}

func (c *BlockerComponent) LoadFromRawComponent(raw ecs.RawComponent) error {
	loadedValues := 0
	if value, ok := raw.Params["active"]; ok {
		if err := c.ActiveFromString(value); err != nil {
			return err
		}
		loadedValues++
	}

	return c.CheckValuesParsedFromRaw(loadedValues, raw)
}

func (c *BlockerComponent) ParseToRawComponent() (ecs.RawComponent, error) {
	rawComponent := ecs.RawComponent{
		ComponentType: ecs.TypeNameToNthBit[c.ComponentType()].Name,
		Params: map[string]string{
			"active": strconv.FormatBool(c.Active),
		},
	}
	return rawComponent, nil
}

func (c *BlockerComponent) ComponentType() uint64 {
	return ecs.BlockerComponentType
}

func (c *BlockerComponent) ActiveFromString(bool string) error {
	b, err := strconv.ParseBool(bool)
	c.Active = b
	return err
}
