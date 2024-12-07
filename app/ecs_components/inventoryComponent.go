package ecs_components

import (
	"github.com/google/uuid"
	"github.com/hielkefellinger/go-dnd/app/ecs"
	"strconv"
)

type InventoryComponent struct {
	ecs.BaseComponent
	Slots uint `yaml:"slots"`
}

func NewInventoryComponent() ecs.Component {
	return &InventoryComponent{
		BaseComponent: ecs.BaseComponent{Id: uuid.New()},
	}
}

func (c *InventoryComponent) LoadFromRawComponent(raw ecs.RawComponent) error {
	loadedValues := 0
	if value, ok := raw.Params["slots"]; ok {
		if err := c.SlotsFromString(value); err != nil {
			return err
		}
		loadedValues++
	}

	return c.CheckValuesParsedFromRaw(loadedValues, raw)
}

func (c *InventoryComponent) ParseToRawComponent() (ecs.RawComponent, error) {
	rawComponent := ecs.RawComponent{
		ComponentType: ecs.TypeNameToNthBit[c.ComponentType()].Name,
		Params: map[string]string{
			"slots": strconv.Itoa(int(c.Slots)),
		},
	}
	return rawComponent, nil
}

func (c *InventoryComponent) SlotsFromString(slots string) error {
	n, err := strconv.Atoi(slots)
	c.Slots = uint(n)
	return err
}

func (c *InventoryComponent) ComponentType() uint64 {
	return ecs.InventoryComponentType
}
