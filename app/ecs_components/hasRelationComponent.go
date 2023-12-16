package ecs_components

import (
	"github.com/google/uuid"
	"github.com/hielkefellinger/go-dnd/app/ecs"
	"strconv"
)

type HasRelationComponent struct {
	ecs.BaseComponent
	Count  uint
	Entity ecs.Entity
}

func NewHasRelationComponent() ecs.Component {
	return &HasRelationComponent{
		BaseComponent: ecs.BaseComponent{Id: uuid.New()},
		Count:         1,
	}
}

func (c *HasRelationComponent) LoadFromRawComponentRelation(raw ecs.RawComponent, entity ecs.Entity) error {
	loadedValues := 0
	if value, ok := raw.Params["count"]; ok {
		if err := c.CountFromString(value); err != nil {
			return err
		}
		loadedValues++
	}
	c.Entity = entity
	loadedValues++

	return c.CheckValuesParsedFromRaw(loadedValues, raw)
}

func (c *HasRelationComponent) CountFromString(count string) error {
	n, err := strconv.Atoi(count)
	c.Count = uint(n)
	return err
}

func (c *HasRelationComponent) ComponentType() uint64 {
	return ecs.HasRelationComponentType
}
