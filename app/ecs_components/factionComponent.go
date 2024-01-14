package ecs_components

import (
	"github.com/google/uuid"
	"github.com/hielkefellinger/go-dnd/app/ecs"
)

type FactionComponent struct {
	ecs.BaseComponent
	Name        string `yaml:"name"`
	ColourHex   string `yaml:"colour_hex"`
	Description string `yaml:"description"`
}

func NewFactionComponent() ecs.Component {
	return &FactionComponent{
		BaseComponent: ecs.BaseComponent{Id: uuid.New()},
	}
}

func (c *FactionComponent) LoadFromRawComponent(raw ecs.RawComponent) error {
	loadedValues := 0
	if value, ok := raw.Params["name"]; ok {
		c.Name = value
		loadedValues++
	}
	if value, ok := raw.Params["colour_hex"]; ok {
		// @todo Add HEX check
		c.ColourHex = value
		loadedValues++
	}
	if value, ok := raw.Params["description"]; ok {
		c.Description = value
		loadedValues++
	}

	return c.CheckValuesParsedFromRaw(loadedValues, raw)
}

func (c *FactionComponent) ComponentType() uint64 {
	return ecs.FactionComponentType
}
