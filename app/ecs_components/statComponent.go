package ecs_components

import (
	"github.com/google/uuid"
	"github.com/hielkefellinger/go-dnd/app/ecs"
	"strconv"
)

type StatComponent struct {
	ecs.BaseComponent
	Name  string `yaml:"name"`
	Value uint   `yaml:"value"`
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
	if value, ok := raw.Params["count"]; ok {
		if err := c.ValueFromString(value); err != nil {
			return err
		}
		loadedValues++
	}

	return c.CheckValuesParsedFromRaw(loadedValues, raw)
}

func (c *StatComponent) ValueFromString(value string) error {
	n, err := strconv.Atoi(value)
	c.Value = uint(n)
	return err
}

func (c *StatComponent) ComponentType() uint64 {
	return ecs.StatComponentType
}
