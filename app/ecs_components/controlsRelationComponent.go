package ecs_components

import (
	"github.com/google/uuid"
	"github.com/hielkefellinger/go-dnd/app/ecs"
)

type ControlsRelationComponent struct {
	ecs.BaseComponent
	Entity ecs.BaseEntity
}

func NewControlsRelationComponent() ecs.Component {
	return &ControlsRelationComponent{
		BaseComponent: ecs.BaseComponent{Id: uuid.New()},
	}
}

func (c *ControlsRelationComponent) WithEntity(entity ecs.BaseEntity) *ControlsRelationComponent {
	c.Entity = entity
	return c
}

func (c *ControlsRelationComponent) ComponentType() uint64 {
	return ecs.ControlsRelationComponentType
}
