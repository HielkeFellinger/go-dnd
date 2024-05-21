package ecs_components

import (
	"github.com/google/uuid"
	"github.com/hielkefellinger/go-dnd/app/ecs"
	"strconv"
)

type TurnDistanceComponent struct {
	ecs.BaseComponent
	Distance uint `yaml:"distance"`
}

func NewTurnDistanceComponent() ecs.Component {
	return &TurnDistanceComponent{
		BaseComponent: ecs.BaseComponent{Id: uuid.New()},
	}
}

func (c *TurnDistanceComponent) LoadFromRawComponent(raw ecs.RawComponent) error {
	loadedValues := 0
	if value, ok := raw.Params["distance"]; ok {
		if err := c.DistanceFromString(value); err != nil {
			return err
		}
		loadedValues++
	}

	return c.CheckValuesParsedFromRaw(loadedValues, raw)
}

func (c *TurnDistanceComponent) DistanceFromString(amount string) error {
	n, err := strconv.Atoi(amount)
	c.Distance = uint(n)
	return err
}

func (c *TurnDistanceComponent) ParseToRawComponent() (ecs.RawComponent, error) {
	rawComponent := ecs.RawComponent{
		ComponentType: ecs.TypeNameToNthBit[c.ComponentType()].Name,
		Params: map[string]string{
			"distance": strconv.Itoa(int(c.Distance)),
		},
	}
	return rawComponent, nil
}

func (c *TurnDistanceComponent) ComponentType() uint64 {
	return ecs.TurnDistanceComponentType
}
