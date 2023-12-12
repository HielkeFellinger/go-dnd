package ecs_components

import (
	"github.com/google/uuid"
	"github.com/hielkefellinger/go-dnd/app/ecs"
)

type HealthComponent struct {
	ecs.BaseComponent
	Damage    uint `yaml:"damage"`
	Temporary uint `yaml:"temporary"`
	Maximum   uint `yaml:"maximum"`
}

func NewHealthComponent() HealthComponent {
	return HealthComponent{
		BaseComponent: ecs.BaseComponent{Id: uuid.New()},
	}
}

func (c *HealthComponent) WidthDamage(damage uint) *HealthComponent {
	c.Damage = damage
	return c
}

func (c *HealthComponent) WidthTemporary(temporary uint) *HealthComponent {
	c.Temporary = temporary
	return c
}

func (c *HealthComponent) WidthMaximum(maximum uint) *HealthComponent {
	c.Maximum = maximum
	return c
}

func (c *HealthComponent) ComponentType() uint64 {
	return ecs.HealthComponentType
}
