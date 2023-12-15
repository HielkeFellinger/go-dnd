package ecs_components

import (
	"github.com/google/uuid"
	"github.com/hielkefellinger/go-dnd/app/ecs"
)

type TransportComponent struct {
	ecs.BaseComponent
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
}

func NewTransportComponent() ecs.Component {
	return &TransportComponent{
		BaseComponent: ecs.BaseComponent{Id: uuid.New()},
	}
}

func (c *TransportComponent) LoadFromRawComponent(raw ecs.RawComponent) error {
	loadedValues := 0
	if value, ok := raw.Params["name"]; ok {
		c.Name = value
		loadedValues++
	}
	if value, ok := raw.Params["description"]; ok {
		c.Description = value
		loadedValues++
	}

	return c.CheckValuesParsedFromRaw(loadedValues, raw)
}
func (c *TransportComponent) ComponentType() uint64 {
	return ecs.TransportComponentType
}
