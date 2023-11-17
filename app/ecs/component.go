package ecs

import "github.com/google/uuid"

const (
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
)

type Component interface {
	ComponentType() uint64
}

type BaseComponent struct {
	Id uuid.UUID
}

// 		Stats? Player (Enemy/ Faction), Exp..

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
