package ecs

import "github.com/google/uuid"

type World interface {
	AddEntity(e Entity)
	AddEntities(e []Entity)
}

type BaseWorld struct {
	systems  []System
	entities []Entity

	UuidToEntity          map[uuid.UUID]Entity
	UuidToCharacterEntity map[uuid.UUID]Entity
	UuidToMapEntity       map[uuid.UUID]Entity
}

func (w *BaseWorld) AddEntity(e Entity) {
	w.UuidToEntity[e.GetId()] = e

	w.entities = append(w.entities, e)
}

func (w *BaseWorld) AddEntities(e []Entity) {
	for _, entity := range e {
		w.AddEntity(entity)
	}
}
