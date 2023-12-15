package ecs_components

import (
	"github.com/google/uuid"
	"github.com/hielkefellinger/go-dnd/app/ecs"
)

type HasRelationComponent struct {
	ecs.BaseComponent
	Count  uint
	Entity ecs.BaseEntity
}

func NewHasRelationComponent() ecs.Component {
	return &HasRelationComponent{
		BaseComponent: ecs.BaseComponent{Id: uuid.New()},
		Count:         1,
	}
}

func (c *HasRelationComponent) WithCount(count uint) *HasRelationComponent {
	c.Count = count
	return c
}

func (c *HasRelationComponent) WithEntity(entity ecs.BaseEntity) *HasRelationComponent {
	c.Entity = entity
	return c
}

func (c *HasRelationComponent) ComponentType() uint64 {
	return ecs.HasRelationComponentType
}
