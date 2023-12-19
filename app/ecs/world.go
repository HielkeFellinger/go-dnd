package ecs

import "github.com/google/uuid"

type World interface {
}

type BaseWorld struct {
	systems  []System
	entities []Entity

	UuidToEntity map[uuid.UUID]Entity
}

func (w *BaseWorld) AddEntity(e Entity) {

	// Add more checks and mapping + functions

	w.entities = append(w.entities, e)
}

func (w *BaseWorld) AddEntities(e []Entity) {

	// Add more checks and mapping + functions
	w.entities = append(w.entities, e...)
}
