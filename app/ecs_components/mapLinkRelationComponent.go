package ecs_components

import (
	"github.com/google/uuid"
	"github.com/hielkefellinger/go-dnd/app/ecs"
	"strconv"
)

type MapLinkRelationComponent struct {
	ecs.BaseComponent
	SourcePosition      *PositionComponent
	DestinationPosition *PositionComponent
	Entity              ecs.Entity
}

func NewMapLinkRelationComponent() ecs.Component {
	return &MapLinkRelationComponent{
		BaseComponent: ecs.BaseComponent{Id: uuid.New()},
	}
}

func (c *MapLinkRelationComponent) LoadFromRawComponentRelation(raw ecs.RawComponent, entity ecs.Entity) error {
	loadedValues := 1
	c.Entity = entity

	// Load the position (empty or missing X or Y is on the standby "bench")
	if sourcePos, err := c.loadPositionFromXAndY(raw, "source_x", "source_y"); err != nil {
		loadedValues += 2
		c.SourcePosition = sourcePos
	} else {
		return err
	}

	if destinationPos, err := c.loadPositionFromXAndY(raw, "destination_x", "destination_y"); err != nil {
		loadedValues += 2
		c.DestinationPosition = destinationPos
	} else {
		return err
	}

	return c.CheckValuesParsedFromRaw(loadedValues, raw)
}

func (c *MapLinkRelationComponent) loadPositionFromXAndY(raw ecs.RawComponent, x string, y string) (*PositionComponent, error) {
	var position = NewPositionComponent().(*PositionComponent)
	if value, ok := raw.Params[x]; ok {
		if err := position.XFromString(value); err != nil {
			return nil, err
		}
	}
	if value, ok := raw.Params[y]; ok {
		if err := position.YFromString(value); err != nil {
			return nil, err
		}
	}

	return position, nil
}

func (c *MapLinkRelationComponent) ParseToRawComponent() (ecs.RawComponent, error) {
	rawComponent := ecs.RawComponent{
		ComponentType: ecs.TypeNameToNthBit[c.ComponentType()].Name,
		Params: map[string]string{
			"entity":        c.Entity.GetId().String(),
			"source_x":      strconv.Itoa(int(c.SourcePosition.X)),
			"source_y":      strconv.Itoa(int(c.SourcePosition.Y)),
			"destination_x": strconv.Itoa(int(c.DestinationPosition.X)),
			"destination_y": strconv.Itoa(int(c.DestinationPosition.Y)),
		},
	}
	return rawComponent, nil
}

func (c *MapLinkRelationComponent) ComponentType() uint64 {
	return ecs.MapLinkRelationComponentType
}

func (c *MapLinkRelationComponent) IsRelationalComponent() bool {
	return true
}

func (c *MapLinkRelationComponent) AllowMultipleOfType() bool {
	return true
}

func (c *MapLinkRelationComponent) GetEntity() ecs.Entity {
	return c.Entity
}
