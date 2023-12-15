package ecs_components

import (
	"github.com/google/uuid"
	"github.com/hielkefellinger/go-dnd/app/ecs"
)

type WeightComponent struct {
	ecs.BaseComponent
	Amount string `yaml:"amount"`
}

func NewWeightComponent() ecs.Component {
	return &WeightComponent{
		BaseComponent: ecs.BaseComponent{Id: uuid.New()},
	}
}

func (c *WeightComponent) LoadFromRawComponent(raw ecs.RawComponent) error {
	loadedValues := 0
	if value, ok := raw.Params["amount"]; ok {
		c.Amount = value
		loadedValues++
	}

	return c.CheckValuesParsedFromRaw(loadedValues, raw)
}

func (c *WeightComponent) ComponentType() uint64 {
	return ecs.WeightComponentType
}
