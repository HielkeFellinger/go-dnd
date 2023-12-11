package ecs

type World interface {
}

type BaseWorld struct {
	Systems  []System
	Entities []Entity
}
