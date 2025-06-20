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

func (c *HasRelationComponent) ParseToRawComponent() (ecs.RawComponent, error) {
	rawComponent := ecs.RawComponent{
		ComponentType: ecs.TypeNameToNthBit[c.ComponentType()].Name,
		Params: map[string]string{
			"entity": c.Entity.GetId().String(),
			"count":  strconv.Itoa(int(c.Count)),
		},
	}
	return rawComponent, nil
}

func (c *HasRelationComponent) AllowMultipleOfType() bool {
	return true
}

func (c *HasRelationComponent) CountFromString(count string) error {
	n, err := strconv.Atoi(count)
	c.Count = uint(n)
	return err
}

func (c *HasRelationComponent) ComponentType() uint64 {
	return ecs.HasRelationComponentType
}

func (c *HasRelationComponent) IsRelationalComponent() bool {
	return true
}

func (c *HasRelationComponent) GetEntity() ecs.Entity {
	return c.Entity
}

func (c *HasRelationComponent) IsLessThanValue(value int) bool {
	return int(c.Count) < value
}

func (c *HasRelationComponent) IsMoreThanValue(value int) bool {
	return int(c.Count) > value
}
