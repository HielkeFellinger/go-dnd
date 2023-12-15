package ecs_components

import (
	"github.com/google/uuid"
	"github.com/hielkefellinger/go-dnd/app/ecs"
)

type RequiresRelationComponent struct {
	ecs.BaseComponent
	Count  uint
	Entity ecs.BaseEntity
}

func NewRequiresRelationComponent() ecs.Component {
	return &RequiresRelationComponent{
		BaseComponent: ecs.BaseComponent{Id: uuid.New()},
		Count:         1,
	}
}

func (c *RequiresRelationComponent) WithCount(count uint) *RequiresRelationComponent {
	c.Count = count
	return c
}

func (c *RequiresRelationComponent) WithEntity(entity ecs.BaseEntity) *RequiresRelationComponent {
	c.Entity = entity
	return c
}

func (c *RequiresRelationComponent) ComponentType() uint64 {
	return ecs.RequiresRelationComponentType
}
