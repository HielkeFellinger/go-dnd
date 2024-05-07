package ecs

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
)

const (
	UnknownComponentType      uint64 = 0
	PositionComponentType     uint64 = 1 << 0
	AreaComponentType         uint64 = 1 << 1
	RangeComponentType        uint64 = 1 << 2
	DamageComponentType       uint64 = 1 << 3
	RestoreComponentType      uint64 = 1 << 4
	ItemComponentType         uint64 = 1 << 5
	AmountComponentType       uint64 = 1 << 6
	WeightComponentType       uint64 = 1 << 7
	SlotsComponentType        uint64 = 1 << 8
	LevelComponentType        uint64 = 1 << 9
	TypeComponentType         uint64 = 1 << 10
	ValutaComponentType       uint64 = 1 << 11
	ResourceComponentType     uint64 = 1 << 12
	TransportComponentType    uint64 = 1 << 13
	TurnDistanceComponentType uint64 = 1 << 14
	VisibilityComponentType   uint64 = 1 << 15
	HealthComponentType       uint64 = 1 << 16
	StatComponentType         uint64 = 1 << 17
	FactionComponentType      uint64 = 1 << 18
	CharacterComponentType    uint64 = 1 << 19
	MapComponentType          uint64 = 1 << 20
	ImageComponentType        uint64 = 1 << 21
	PlayerComponentType       uint64 = 1 << 22

	/* Relational ComponentTypes */

	ControlsRelationComponentType uint64 = 1 << 40
	HasRelationComponentType      uint64 = 1 << 41
	RequiresRelationComponentType uint64 = 1 << 42
	CreatesRelationComponentType  uint64 = 1 << 43
	FilterRelationComponentType   uint64 = 1 << 44
	MapItemRelationComponentType  uint64 = 1 << 45
)

type Component interface {
	ComponentType() uint64
	AllowMultipleOfType() bool
	LoadFromRawComponent(raw RawComponent) error
	IsRelationalComponent() bool
	GetId() uuid.UUID
}

type RelationalComponent interface {
	Component
	LoadFromRawComponentRelation(raw RawComponent, entity Entity) error
	GetEntity() Entity
}

type BaseComponent struct {
	Version uint // Placeholder; take changes into account
	Id      uuid.UUID
}

func (c *BaseComponent) GetId() uuid.UUID {
	return c.Id
}

func (c *BaseComponent) IsRelationalComponent() bool {
	return false
}

func (c *BaseComponent) CheckValuesParsedFromRaw(loadedValues int, raw RawComponent) error {
	if loadedValues != len(raw.Params) {
		return errors.New(fmt.Sprintf("Mismatch on Type: (%v)'%v'. Loaded '%d' items total '%d'. Values: '%v'",
			c.ComponentType(), raw.ComponentType, loadedValues, len(raw.Params), raw.Params))
	}
	return nil
}

func (c *BaseComponent) LoadFromRawComponent(raw RawComponent) error {
	return errors.New(fmt.Sprintf("loadFromRawComponent(raw RawComponent) not implemented. Raw type: '%s'", raw.ComponentType))
}

func (c *BaseComponent) AllowMultipleOfType() bool {
	return false
}

func (c *BaseComponent) ComponentType() uint64 {
	return UnknownComponentType
}

type FilterMode uint64

const (
	UnknownFilterMode  FilterMode = 1 << 0
	AllowFilterMode    FilterMode = 1 << 1
	BlockFilterMode    FilterMode = 1 << 2
	LessThanFilterMode FilterMode = 1 << 3
	MoreThanFilterMode FilterMode = 1 << 4
)
