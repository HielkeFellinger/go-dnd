package ecs_components

import (
	"github.com/google/uuid"
	"github.com/hielkefellinger/go-dnd/app/ecs"
)

type DamageComponent struct {
	ecs.BaseComponent
	Amount string
}

func NewDamageComponent() DamageComponent {
	return DamageComponent{
		BaseComponent: ecs.BaseComponent{Id: uuid.New()},
	}
}

func (c *DamageComponent) WithAmount(amount string) *DamageComponent {
	c.Amount = amount
	return c
}

func (c *DamageComponent) ComponentType() uint64 {
	return ecs.DamageComponentType
}
