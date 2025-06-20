package ecs_components

import (
	"github.com/google/uuid"
	"github.com/hielkefellinger/go-dnd/app/ecs"
	"strconv"
)

type AreaComponent struct {
	ecs.BaseComponent
	Length uint `yaml:"length"`
	Width  uint `yaml:"width"`
}

func NewAreaComponent() ecs.Component {
	return &AreaComponent{
		BaseComponent: ecs.BaseComponent{Id: uuid.New()},
	}
}

func (c *AreaComponent) LoadFromRawComponent(raw ecs.RawComponent) error {
	loadedValues := 0
	if value, ok := raw.Params["length"]; ok {
		if err := c.LengthFromString(value); err != nil {
			return err
		}
		loadedValues++
	}
	if value, ok := raw.Params["width"]; ok {
		if err := c.WidthFromString(value); err != nil {
			return err
		}
		loadedValues++
	}

	return c.CheckValuesParsedFromRaw(loadedValues, raw)
}

func (c *AreaComponent) ParseToRawComponent() (ecs.RawComponent, error) {
	rawComponent := ecs.RawComponent{
		ComponentType: ecs.TypeNameToNthBit[c.ComponentType()].Name,
		Params: map[string]string{
			"length": strconv.Itoa(int(c.Length)),
			"width":  strconv.Itoa(int(c.Width)),
		},
	}
	return rawComponent, nil
}

func (c *AreaComponent) LengthFromString(length string) error {
	n, err := strconv.Atoi(length)
	c.Length = uint(n)
	return err
}

func (c *AreaComponent) WidthFromString(width string) error {
	n, err := strconv.Atoi(width)
	c.Width = uint(n)
	return err
}

func (c *AreaComponent) ComponentType() uint64 {
	return ecs.AreaComponentType
}

func (c *AreaComponent) IsLessThanValue(value int) bool {
	return int(c.Length*c.Width) < value
}

func (c *AreaComponent) IsMoreThanValue(value int) bool {
	return int(c.Length*c.Width) > value
}
