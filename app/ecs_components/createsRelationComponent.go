package ecs_components

import (
	"github.com/google/uuid"
	"github.com/hielkefellinger/go-dnd/app/ecs"
	"strconv"
)

type CreatesRelationComponent struct {
	ecs.BaseComponent
	Count  uint `yaml:"count"`
	Entity ecs.Entity
}

func NewCreatesRelationComponent() ecs.Component {
	return &CreatesRelationComponent{
		BaseComponent: ecs.BaseComponent{Id: uuid.New()},
		Count:         1,
	}
}

func (c *CreatesRelationComponent) LoadFromRawComponentRelation(raw ecs.RawComponent, entity ecs.Entity) error {
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

func (c *CreatesRelationComponent) ParseToRawComponent() (ecs.RawComponent, error) {
	rawComponent := ecs.RawComponent{
		ComponentType: ecs.TypeNameToNthBit[c.ComponentType()].Name,
		Params: map[string]string{
			"entity": c.Entity.GetId().String(),
			"count":  strconv.Itoa(int(c.Count)),
		},
	}
	return rawComponent, nil
}

func (c *CreatesRelationComponent) AllowMultipleOfType() bool {
	return true
}

func (c *CreatesRelationComponent) CountFromString(count string) error {
	n, err := strconv.Atoi(count)
	c.Count = uint(n)
	return err
}

func (c *CreatesRelationComponent) ComponentType() uint64 {
	return ecs.CreatesRelationComponentType
}

func (c *CreatesRelationComponent) IsRelationalComponent() bool {
	return true
}

func (c *CreatesRelationComponent) GetEntity() ecs.Entity {
	return c.Entity
}
