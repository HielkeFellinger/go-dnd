package ecs_components

import (
	"github.com/google/uuid"
	"github.com/hielkefellinger/go-dnd/app/ecs"
)

type AmountComponent struct {
	ecs.BaseComponent
	Amount int
}

func NewValueComponent() AmountComponent {
	return AmountComponent{
		BaseComponent: ecs.BaseComponent{Id: uuid.New()},
	}
}

func (c *AmountComponent) WithAmount(amount int) *AmountComponent {
	c.Amount = amount
	return c
}

func (c *AmountComponent) ComponentType() uint64 {
	return ecs.AmountComponentType
}
