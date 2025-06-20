package ecs

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
)

type Component interface {
	ComponentType() uint64
	AllowMultipleOfType() bool
	LoadFromRawComponent(raw RawComponent) error
	ParseToRawComponent() (RawComponent, error)
	IsRelationalComponent() bool
	IsLessThanValue(value int) bool
	IsMoreThanValue(value int) bool
	GetId() uuid.UUID
}

type RelationalComponent interface {
	Component
	LoadFromRawComponentRelation(raw RawComponent, entity Entity) error
	GetEntity() Entity
}

type BaseComponent struct {
	Version uint      // Placeholder; take changes into account
	Id      uuid.UUID `yaml:"id"`
}

func (c *BaseComponent) GetId() uuid.UUID {
	return c.Id
}

func (c *BaseComponent) IsRelationalComponent() bool {
	return false
}

func (c *BaseComponent) CheckValuesParsedFromRaw(loadedValues int, raw RawComponent) error {
	if loadedValues != len(raw.Params) {
		return errors.New(fmt.Sprintf("Mismatch on Type: (%v) '%v'. Loaded '%d' items total '%d'. Values: '%v'",
			c.ComponentType(), raw.ComponentType, loadedValues, len(raw.Params), raw.Params))
	}
	return nil
}

func (c *BaseComponent) LoadFromRawComponent(raw RawComponent) error {
	return errors.New(fmt.Sprintf("loadFromRawComponent(raw RawComponent) not implemented. Raw type: '%s'", raw.ComponentType))
}

func (c *BaseComponent) ParseToRawComponent() (RawComponent, error) {
	return RawComponent{}, errors.New(fmt.Sprintf("ParseToRawComponent() not implemented. on type: '%d'", c.ComponentType()))
}

func (c *BaseComponent) AllowMultipleOfType() bool {
	return false
}

func (c *BaseComponent) ComponentType() uint64 {
	return UnknownComponentType
}

func (c *BaseComponent) IsLessThanValue(value int) bool {
	return false
}

func (c *BaseComponent) IsMoreThanValue(value int) bool {
	return false
}

type FilterMode uint64

const (
	UnknownFilterMode       FilterMode = 1 << 0
	AllowFilterMode         FilterMode = 1 << 1
	BlockFilterMode         FilterMode = 1 << 2
	LessThanFilterMode      FilterMode = 1 << 3
	MoreThanFilterMode      FilterMode = 1 << 4
	ShouldHaveFilterMode    FilterMode = 1 << 5
	ShouldNotHaveFilterMode FilterMode = 1 << 6
)
