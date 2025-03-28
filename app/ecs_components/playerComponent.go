package ecs_components

import (
	"github.com/google/uuid"
	"github.com/hielkefellinger/go-dnd/app/ecs"
	"golang.org/x/net/html"
)

type PlayerComponent struct {
	ecs.BaseComponent
	Name string `yaml:"name"`
}

func NewPlayerComponent() ecs.Component {
	return &PlayerComponent{
		BaseComponent: ecs.BaseComponent{Id: uuid.New()},
	}
}

func (c *PlayerComponent) LoadFromRawComponent(raw ecs.RawComponent) error {
	loadedValues := 0
	if value, ok := raw.Params["name"]; ok {
		c.Name = value
		loadedValues++
	}

	return c.CheckValuesParsedFromRaw(loadedValues, raw)
}

func (c *PlayerComponent) ParseToRawComponent() (ecs.RawComponent, error) {
	rawComponent := ecs.RawComponent{
		ComponentType: ecs.TypeNameToNthBit[c.ComponentType()].Name,
		Params: map[string]string{
			"name": html.EscapeString(html.UnescapeString(c.Name)),
		},
	}
	return rawComponent, nil
}

func (c *PlayerComponent) ComponentType() uint64 {
	return ecs.PlayerComponentType
}

func (c *PlayerComponent) AllowMultipleOfType() bool {
	return true
}
