package ecs_components

import (
	"github.com/google/uuid"
	"github.com/hielkefellinger/go-dnd/app/ecs"
)

type WeightComponent struct {
	ecs.BaseComponent
	Amount string `yaml:"amount"`
}

func NewWeightComponent() AmountComponent {
	return AmountComponent{
		BaseComponent: ecs.BaseComponent{Id: uuid.New()},
	}
}

func (c *WeightComponent) WithAmount(amount string) *WeightComponent {
	c.Amount = amount
	return c
}

func (c *WeightComponent) ComponentType() uint64 {
	return ecs.WeightComponentType
}
