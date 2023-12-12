package ecs_components

import (
	"github.com/google/uuid"
	"github.com/hielkefellinger/go-dnd/app/ecs"
)

type RestoreComponent struct {
	ecs.BaseComponent
	Amount string `yaml:"amount"`
}

func NewRestoreComponent() RestoreComponent {
	return RestoreComponent{
		BaseComponent: ecs.BaseComponent{Id: uuid.New()},
	}
}

func (c *RestoreComponent) WithAmount(amount string) *RestoreComponent {
	c.Amount = amount
	return c
}

func (c *RestoreComponent) ComponentType() uint64 {
	return ecs.RestoreComponentType
}
