package ecs_components

import (
	"errors"
	"github.com/google/uuid"
	"github.com/hielkefellinger/go-dnd/app/ecs"
)

type FilterRelationComponent struct {
	ecs.BaseComponent
	Mode   ecs.FilterMode
	Entity ecs.Entity
}

func NewFilterRelationComponent() ecs.Component {
	return &FilterRelationComponent{
		BaseComponent: ecs.BaseComponent{Id: uuid.New()},
		Mode:          ecs.UnknownFilterMode,
	}
}

func (c *FilterRelationComponent) LoadFromRawComponentRelation(raw ecs.RawComponent, entity ecs.Entity) error {
	loadedValues := 0
	if value, ok := raw.Params["mode"]; ok {
		if err := c.ModeFromString(value); err != nil {
			return err
		}
		loadedValues++
	}
	c.Entity = entity
	loadedValues++

	return c.CheckValuesParsedFromRaw(loadedValues, raw)
}

func (c *FilterRelationComponent) ModeFromString(mode string) error {
	c.Mode = MapStringToFilterMode(mode)
	if c.Mode == ecs.UnknownFilterMode {
		return errors.New("no or an Unknown FilterMode has been supplied")
	}
	return nil
}

func (c *FilterRelationComponent) ComponentType() uint64 {
	return ecs.FilterRelationComponentType
}
