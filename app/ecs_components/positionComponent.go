package ecs_components

import (
	"github.com/google/uuid"
	"github.com/hielkefellinger/go-dnd/app/ecs"
	"strconv"
)

type PositionComponent struct {
	ecs.BaseComponent
	X uint `yaml:"x"` // Column
	Y uint `yaml:"y"` // Row
}

func NewPositionComponent() ecs.Component {
	return &PositionComponent{
		BaseComponent: ecs.BaseComponent{Id: uuid.New()},
	}
}

func (c *PositionComponent) LoadFromRawComponent(raw ecs.RawComponent) error {
	loadedValues := 0
	if value, ok := raw.Params["x"]; ok {
		if err := c.XFromString(value); err != nil {
			return err
		}
		loadedValues++
	}
	if value, ok := raw.Params["y"]; ok {
		if err := c.YFromString(value); err != nil {
			return err
		}
		loadedValues++
	}

	return c.CheckValuesParsedFromRaw(loadedValues, raw)
}

func (c *PositionComponent) XFromString(x string) error {
	n, err := strconv.Atoi(x)
	c.X = uint(n)
	return err
}

func (c *PositionComponent) YFromString(y string) error {
	n, err := strconv.Atoi(y)
	c.Y = uint(n)
	return err
}

func (c *PositionComponent) ComponentType() uint64 {
	return ecs.PositionComponentType
}
