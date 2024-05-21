package ecs_components

import (
	"github.com/google/uuid"
	"github.com/hielkefellinger/go-dnd/app/ecs"
)

type ControlsRelationComponent struct {
	ecs.BaseComponent
	Entity ecs.Entity
}

func NewControlsRelationComponent() ecs.Component {
	return &ControlsRelationComponent{
		BaseComponent: ecs.BaseComponent{Id: uuid.New()},
	}
}

func (c *ControlsRelationComponent) LoadFromRawComponentRelation(raw ecs.RawComponent, entity ecs.Entity) error {
	loadedValues := 1
	c.Entity = entity

	return c.CheckValuesParsedFromRaw(loadedValues, raw)
}

func (c *ControlsRelationComponent) ParseToRawComponent() (ecs.RawComponent, error) {
	rawComponent := ecs.RawComponent{
		ComponentType: ecs.TypeNameToNthBit[c.ComponentType()].Name,
		Params: map[string]string{
			"entity": c.Entity.GetId().String(),
		},
	}
	return rawComponent, nil
}

func (c *ControlsRelationComponent) AllowMultipleOfType() bool {
	return true
}

func (c *ControlsRelationComponent) ComponentType() uint64 {
	return ecs.ControlsRelationComponentType
}

func (c *ControlsRelationComponent) IsRelationalComponent() bool {
	return true
}

func (c *ControlsRelationComponent) GetEntity() ecs.Entity {
	return c.Entity
}
