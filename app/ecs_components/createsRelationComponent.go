package ecs_components

import (
	"github.com/google/uuid"
	"github.com/hielkefellinger/go-dnd/app/ecs"
)

type CreatesRelationComponent struct {
	ecs.BaseComponent
	Count  uint `yaml:"count"`
	Entity ecs.BaseEntity
}

func NewCreatesRelationComponent() CreatesRelationComponent {
	return CreatesRelationComponent{
		BaseComponent: ecs.BaseComponent{Id: uuid.New()},
		Count:         1,
	}
}

func (c *CreatesRelationComponent) WithCount(count uint) *CreatesRelationComponent {
	c.Count = count
	return c
}

func (c *CreatesRelationComponent) WithEntity(entity ecs.BaseEntity) *CreatesRelationComponent {
	c.Entity = entity
	return c
}

func (c *CreatesRelationComponent) ComponentType() uint64 {
	return ecs.CreatesRelationComponentType
}
