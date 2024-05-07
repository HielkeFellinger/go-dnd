package ecs_components

import (
	"github.com/google/uuid"
	"github.com/hielkefellinger/go-dnd/app/ecs"
)

type MapItemRelationComponent struct {
	ecs.BaseComponent
	Position *PositionComponent
	Entity   ecs.Entity
}

func NewMapItemRelationComponent() ecs.Component {
	return &MapItemRelationComponent{
		BaseComponent: ecs.BaseComponent{Id: uuid.New()},
	}
}

func (c *MapItemRelationComponent) LoadFromRawComponentRelation(raw ecs.RawComponent, entity ecs.Entity) error {
	loadedValues := 1
	c.Entity = entity

	// Load the position (empty or missing X or Y is on the standby "bench")
	var position = NewPositionComponent().(*PositionComponent)
	if value, ok := raw.Params["x"]; ok {
		if err := position.XFromString(value); err != nil {
			return err
		}
		loadedValues++
	}
	if value, ok := raw.Params["y"]; ok {
		if err := position.YFromString(value); err != nil {
			return err
		}
		loadedValues++
	}
	if loadedValues == 3 {
		c.Position = position
	}

	return c.CheckValuesParsedFromRaw(loadedValues, raw)
}

func (c *MapItemRelationComponent) ComponentType() uint64 {
	return ecs.MapItemRelationComponentType
}

func (c *MapItemRelationComponent) IsRelationalComponent() bool {
	return true
}

func (c *MapItemRelationComponent) AllowMultipleOfType() bool {
	return true
}

func (c *MapItemRelationComponent) GetEntity() ecs.Entity {
	return c.Entity
}
