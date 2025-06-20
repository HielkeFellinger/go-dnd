package ecs_components

import (
	"github.com/google/uuid"
	"github.com/hielkefellinger/go-dnd/app/ecs"
	"strconv"
)

type HealthComponent struct {
	ecs.BaseComponent
	Damage    uint `yaml:"damage"`
	Temporary uint `yaml:"temporary"`
	Maximum   uint `yaml:"maximum"`
}

func NewHealthComponent() ecs.Component {
	return &HealthComponent{
		BaseComponent: ecs.BaseComponent{Id: uuid.New()},
	}
}

func (c *HealthComponent) LoadFromRawComponent(raw ecs.RawComponent) error {
	loadedValues := 0
	if value, ok := raw.Params["damage"]; ok {
		if err := c.DamageFromString(value); err != nil {
			return err
		}
		loadedValues++
	}
	if value, ok := raw.Params["temporary"]; ok {
		if err := c.TemporaryFromString(value); err != nil {
			return err
		}
		loadedValues++
	}
	if value, ok := raw.Params["maximum"]; ok {
		if err := c.MaximumFromString(value); err != nil {
			return err
		}
		loadedValues++
	}

	return c.CheckValuesParsedFromRaw(loadedValues, raw)
}

func (c *HealthComponent) ParseToRawComponent() (ecs.RawComponent, error) {
	rawComponent := ecs.RawComponent{
		ComponentType: ecs.TypeNameToNthBit[c.ComponentType()].Name,
		Params: map[string]string{
			"damage":    strconv.Itoa(int(c.Damage)),
			"temporary": strconv.Itoa(int(c.Temporary)),
			"maximum":   strconv.Itoa(int(c.Maximum)),
		},
	}
	return rawComponent, nil
}

func (c *HealthComponent) DamageFromString(damage string) error {
	n, err := strconv.Atoi(damage)
	c.Damage = uint(n)
	return err
}

func (c *HealthComponent) TemporaryFromString(temporary string) error {
	n, err := strconv.Atoi(temporary)
	c.Temporary = uint(n)
	return err
}

func (c *HealthComponent) MaximumFromString(maximum string) error {
	n, err := strconv.Atoi(maximum)
	c.Maximum = uint(n)
	return err
}

func (c *HealthComponent) ComponentType() uint64 {
	return ecs.HealthComponentType
}

func (c *HealthComponent) IsLessThanValue(value int) bool {
	return int(c.Maximum+c.Temporary)-int(c.Damage) < value
}

func (c *HealthComponent) IsMoreThanValue(value int) bool {
	return int(c.Maximum+c.Temporary)-int(c.Damage) > value
}
