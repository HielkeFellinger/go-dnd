package ecs_components

import (
	"github.com/google/uuid"
	"github.com/hielkefellinger/go-dnd/app/ecs"
)

type ValueComponent struct {
	ecs.BaseComponent
	Amount string
}

func NewValueComponent() ValueComponent {
	return ValueComponent{
		BaseComponent: ecs.BaseComponent{Id: uuid.New()},
	}
}

func (c *ValueComponent) WithAmount(amount string) *ValueComponent {
	c.Amount = amount
	return c
}

func (c *ValueComponent) ComponentType() uint64 {
	return ecs.ValueComponentType
}
