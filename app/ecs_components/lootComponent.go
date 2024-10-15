package ecs_components

import (
	"github.com/google/uuid"
	"github.com/hielkefellinger/go-dnd/app/ecs"
	"strconv"
)

type LootComponent struct {
	ecs.BaseComponent
	PlayerVisible bool `yaml:"playerVisible"`
}

func NewLootComponent() ecs.Component {
	return &LootComponent{
		BaseComponent: ecs.BaseComponent{Id: uuid.New()},
	}
}

func (c *LootComponent) LoadFromRawComponent(raw ecs.RawComponent) error {
	loadedValues := 0
	if value, ok := raw.Params["playerVisible"]; ok {
		if err := c.PlayerVisibleFromString(value); err != nil {
			return err
		}
		loadedValues++
	}

	return c.CheckValuesParsedFromRaw(loadedValues, raw)
}

func (c *LootComponent) ParseToRawComponent() (ecs.RawComponent, error) {
	rawComponent := ecs.RawComponent{
		ComponentType: ecs.TypeNameToNthBit[c.ComponentType()].Name,
		Params: map[string]string{
			"playerVisible": strconv.FormatBool(c.PlayerVisible),
		},
	}
	return rawComponent, nil
}

func (c *LootComponent) PlayerVisibleFromString(bool string) error {
	b, err := strconv.ParseBool(bool)
	c.PlayerVisible = b
	return err
}

func (c *LootComponent) ComponentType() uint64 {
	return ecs.LootComponentType
}
