package ecs

import "github.com/google/uuid"

const (
	PositionComponentType uint64 = 1 << 0
	AreaComponentType     uint64 = 1 << 1
	RangeComponentType    uint64 = 1 << 2
	DamageComponentType   uint64 = 1 << 3
	RestoreComponentType  uint64 = 1 << 4
	ItemComponentType     uint64 = 1 << 5
	AmountComponentType   uint64 = 1 << 6
	WeightComponentType   uint64 = 1 << 7
	SlotsComponentType    uint64 = 1 << 8
	LevelComponentType    uint64 = 1 << 9
	TypeComponentType     uint64 = 1 << 10
	ValutaComponentType   uint64 = 1 << 11
	ResourceComponentType uint64 = 1 << 12
)

type Component interface {
	ComponentType() uint64
}

type BaseComponent struct {
	Id uuid.UUID
}

type TransportComponent struct {
	BaseComponent
	Name        string
	Description string
}

type TurnDistanceComponent struct {
	BaseComponent
	Distance uint
}

type VisibilityComponent struct {
	BaseComponent
	Hidden bool
}

// 		Stats? Player (Enemy/ Faction), Health.. Exp.. Lvl

// Relation Component

type HasRelationComponent struct {
	BaseComponent
	Count  uint
	Entity BaseEntity
}

type RequirementComponent struct {
	BaseComponent
	Count  uint
	Entity BaseEntity
}

type FilterMode uint64

const (
	AllowFilter FilterMode = 1 << 0
	BlockFilter FilterMode = 1 << 1
)

type FilterComponent struct {
	BaseComponent
	Mode  FilterMode
	Value Component
}