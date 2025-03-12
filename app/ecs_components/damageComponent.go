package ecs_components

import (
	"github.com/google/uuid"
	"github.com/hielkefellinger/go-dnd/app/ecs"
	"golang.org/x/net/html"
)

type DamageComponent struct {
	ecs.BaseComponent
	Amount string `yaml:"amount"`
}

func NewDamageComponent() ecs.Component {
	return &DamageComponent{
		BaseComponent: ecs.BaseComponent{Id: uuid.New()},
	}
}

func (c *DamageComponent) LoadFromRawComponent(raw ecs.RawComponent) error {
	loadedValues := 0
	if value, ok := raw.Params["amount"]; ok {
		c.Amount = value
		loadedValues++
	}

	return c.CheckValuesParsedFromRaw(loadedValues, raw)
}

func (c *DamageComponent) ParseToRawComponent() (ecs.RawComponent, error) {
	rawComponent := ecs.RawComponent{
		ComponentType: ecs.TypeNameToNthBit[c.ComponentType()].Name,
		Params: map[string]string{
			"amount": html.EscapeString(html.UnescapeString(c.Amount)),
		},
	}
	return rawComponent, nil
}

func (c *DamageComponent) ComponentType() uint64 {
	return ecs.DamageComponentType
}
