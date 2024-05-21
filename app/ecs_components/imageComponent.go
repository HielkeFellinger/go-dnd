package ecs_components

import (
	"github.com/google/uuid"
	"github.com/hielkefellinger/go-dnd/app/ecs"
)

type ImageComponent struct {
	ecs.BaseComponent
	Name   string `yaml:"name"`
	Url    string `yaml:"url"`
	Base64 string `yaml:"base64"`
}

func NewImageComponent() ecs.Component {
	return &ImageComponent{
		BaseComponent: ecs.BaseComponent{Id: uuid.New()},
	}
}

func NewMissingImageComponent() *ImageComponent {
	return &ImageComponent{
		BaseComponent: ecs.BaseComponent{Id: uuid.New()},
		Name:          "MISSING IMAGE",
		Url:           "/images/unknown_item.png",
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
	if value, ok := raw.Params["base64"]; ok {
		c.Base64 = value
		loadedValues++
	}

	return c.CheckValuesParsedFromRaw(loadedValues, raw)
}

func (c *ImageComponent) ParseToRawComponent() (ecs.RawComponent, error) {
	rawComponent := ecs.RawComponent{
		ComponentType: ecs.TypeNameToNthBit[c.ComponentType()].Name,
		Params: map[string]string{
			"name":   c.Name,
			"url":    c.Url,
			"base64": c.Base64,
		},
	}
	return rawComponent, nil
}

func (c *ImageComponent) ComponentType() uint64 {
	return ecs.ImageComponentType
}
