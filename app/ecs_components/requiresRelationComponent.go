package ecs_components

import (
	"github.com/google/uuid"
	"github.com/hielkefellinger/go-dnd/app/ecs"
	"strconv"
)

type RequiresRelationComponent struct {
	ecs.BaseComponent
	Count  uint
	Entity ecs.Entity
}

func NewRequiresRelationComponent() ecs.Component {
	return &RequiresRelationComponent{
		BaseComponent: ecs.BaseComponent{Id: uuid.New()},
		Count:         1,
	}
}

func (c *RequiresRelationComponent) LoadFromRawComponentRelation(raw ecs.RawComponent, entity ecs.Entity) error {
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

func (c *RequiresRelationComponent) ParseToRawComponent() (ecs.RawComponent, error) {
	rawComponent := ecs.RawComponent{
		ComponentType: ecs.TypeNameToNthBit[c.ComponentType()].Name,
		Params: map[string]string{
			"entity": c.Entity.GetId().String(),
			"count":  strconv.Itoa(int(c.Count)),
		},
	}
	return rawComponent, nil
}

func (c *RequiresRelationComponent) AllowMultipleOfType() bool {
	return true
}

func (c *RequiresRelationComponent) CountFromString(count string) error {
	n, err := strconv.Atoi(count)
	c.Count = uint(n)
	return err
}

func (c *RequiresRelationComponent) ComponentType() uint64 {
	return ecs.RequiresRelationComponentType
}

func (c *RequiresRelationComponent) IsRelationalComponent() bool {
	return true
}

func (c *RequiresRelationComponent) GetEntity() ecs.Entity {
	return c.Entity
}
