package ecs_components

import (
	"github.com/google/uuid"
	"github.com/hielkefellinger/go-dnd/app/ecs"
)

type ImageComponent struct {
	ecs.BaseComponent
	Name string `yaml:"name"`
	Url  string `yaml:"url"`
}

func NewImageComponent() ecs.Component {
	return &ImageComponent{
		BaseComponent: ecs.BaseComponent{Id: uuid.New()},
	}
}

func (c *ImageComponent) LoadFromRawComponent(raw ecs.RawComponent) error {
	loadedValues := 0
	if value, ok := raw.Params["name"]; ok {
		c.Name = value
		loadedValues++
	}
	if value, ok := raw.Params["url"]; ok {
		c.Url = value
		loadedValues++
	}

	return c.CheckValuesParsedFromRaw(loadedValues, raw)
}

func (c *ImageComponent) ComponentType() uint64 {
	return ecs.ImageComponentType
}
