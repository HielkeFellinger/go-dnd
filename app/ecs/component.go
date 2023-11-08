package ecs

import "github.com/google/uuid"

const (
	PositionComponentType uint64 = 1 << 0
	AreaComponentType     uint64 = 1 << 1
	RangeComponentType    uint64 = 1 << 2
	DamageComponentType   uint64 = 1 << 3
	RestoreComponentType  uint64 = 1 << 4
	ItemComponentType     uint64 = 1 << 5
	ValueComponentType    uint64 = 1 << 6
	WeightComponentType   uint64 = 1 << 7
)

type Component interface {
	ComponentType() uint64
}

type BaseComponent struct {
	Id uuid.UUID
}

// Relation Component

type HasRelationComponent struct {
	BaseComponent
	Count  int
	Entity BaseEntity
}

type StatusComponent struct {
	BaseComponent
	AllowNegative bool
	CurrentValue  int
	MaxValue      int
	MinValue      int
}
