package ecs_components

import (
	"github.com/google/uuid"
	"github.com/hielkefellinger/go-dnd/app/ecs"
	"strconv"
)

type AmountComponent struct {
	ecs.BaseComponent
	Amount int `yaml:"amount"`
}

func NewAmountComponent() ecs.Component {
	return &AmountComponent{
		BaseComponent: ecs.BaseComponent{Id: uuid.New()},
	}
}

func (c *AmountComponent) LoadFromRawComponent(raw ecs.RawComponent) error {
	loadedValues := 0
	if value, ok := raw.Params["amount"]; ok {
		if err := c.AmountFromString(value); err != nil {
			return err
		}
		loadedValues++
	}

	return c.CheckValuesParsedFromRaw(loadedValues, raw)
}

func (c *AmountComponent) ParseToRawComponent() (ecs.RawComponent, error) {
	rawComponent := ecs.RawComponent{
		ComponentType: ecs.TypeNameToNthBit[c.ComponentType()].Name,
		Params: map[string]string{
			"amount": strconv.Itoa(c.Amount),
		},
	}
	return rawComponent, nil
}

func (c *AmountComponent) AmountFromString(amount string) error {
	n, err := strconv.Atoi(amount)
	c.Amount = n
	return err
}

func (c *AmountComponent) ComponentType() uint64 {
	return ecs.AmountComponentType
}

func (c *AmountComponent) IsLessThanValue(value int) bool {
	return c.Amount < value
}

func (c *AmountComponent) IsMoreThanValue(value int) bool {
	return c.Amount > value
}
