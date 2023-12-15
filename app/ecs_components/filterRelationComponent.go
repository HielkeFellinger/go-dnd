package ecs_components

import (
	"github.com/google/uuid"
	"github.com/hielkefellinger/go-dnd/app/ecs"
)

type FilterRelationComponent struct {
	ecs.BaseComponent
	Mode   ecs.FilterMode
	Entity ecs.BaseEntity
}

func NewFilterRelationComponent() ecs.Component {
	return &FilterRelationComponent{
		BaseComponent: ecs.BaseComponent{Id: uuid.New()},
		Mode:          ecs.NoFilterMode,
	}
}

func (c *FilterRelationComponent) WithMode(mode ecs.FilterMode) *FilterRelationComponent {
	c.Mode = mode
	return c
}

func (c *FilterRelationComponent) WidthEntity(entity ecs.BaseEntity) *FilterRelationComponent {
	c.Entity = entity
	return c
}

func (c *FilterRelationComponent) ComponentType() uint64 {
	return ecs.FilterRelationComponentType
}
