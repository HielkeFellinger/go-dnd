package ecs_components

import (
	"github.com/google/uuid"
	"github.com/hielkefellinger/go-dnd/app/ecs"
	"strconv"
)

type LevelComponent struct {
	ecs.BaseComponent
	Level uint `yaml:"level"`
}

func NewLevelComponent() ecs.Component {
	return &LevelComponent{
		BaseComponent: ecs.BaseComponent{Id: uuid.New()},
	}
}

func (c *LevelComponent) LoadFromRawComponent(raw ecs.RawComponent) error {
	loadedValues := 0
	if value, ok := raw.Params["level"]; ok {
		if err := c.LevelFromString(value); err != nil {
			return err
		}
		loadedValues++
	}

	return c.CheckValuesParsedFromRaw(loadedValues, raw)
}

func (c *LevelComponent) LevelFromString(level string) error {
	n, err := strconv.Atoi(level)
	c.Level = uint(n)
	return err
}

func (c *LevelComponent) ComponentType() uint64 {
	return ecs.LevelComponentType
}
