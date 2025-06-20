package ecs_components

import (
	"github.com/google/uuid"
	"github.com/hielkefellinger/go-dnd/app/ecs"
	"golang.org/x/net/html"
	"strconv"
)

type StatComponent struct {
	ecs.BaseComponent
	Name       string `yaml:"name"`
	Base       int    `yaml:"base"`
	Calculated int    `yaml:"calculated"`
}

func NewStatComponent() ecs.Component {
	return &StatComponent{
		BaseComponent: ecs.BaseComponent{Id: uuid.New()},
	}
}

func (c *StatComponent) LoadFromRawComponent(raw ecs.RawComponent) error {
	loadedValues := 0
	if value, ok := raw.Params["name"]; ok {
		c.Name = value
		loadedValues++
	}
	if value, ok := raw.Params["base"]; ok {
		if err := c.BaseFromString(value); err != nil {
			return err
		}
		loadedValues++
	}
	if value, ok := raw.Params["calculated"]; ok {
		if err := c.CalculatedFromString(value); err != nil {
			return err
		}
		loadedValues++
	}

	return c.CheckValuesParsedFromRaw(loadedValues, raw)
}

func (c *StatComponent) ParseToRawComponent() (ecs.RawComponent, error) {
	rawComponent := ecs.RawComponent{
		ComponentType: ecs.TypeNameToNthBit[c.ComponentType()].Name,
		Params: map[string]string{
			"name":       html.EscapeString(html.UnescapeString(c.Name)),
			"base":       strconv.Itoa(c.Base),
			"calculated": strconv.Itoa(c.Calculated),
		},
	}
	return rawComponent, nil
}

func (c *StatComponent) BaseFromString(value string) error {
	n, err := strconv.Atoi(value)
	c.Base = n
	return err
}

func (c *StatComponent) CalculatedFromString(value string) error {
	n, err := strconv.Atoi(value)
	c.Calculated = n
	return err
}

func (c *StatComponent) ComponentType() uint64 {
	return ecs.StatComponentType
}

func (c *StatComponent) AllowMultipleOfType() bool {
	return true
}
