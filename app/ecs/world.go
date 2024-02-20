package ecs

import (
	"github.com/google/uuid"
)

type World interface {
	AddEntity(e Entity)
	AddEntities(e []Entity)
	GetCharacterEntities() []Entity
	GetMapEntities() []Entity
	GetEntityByUuid(uuid uuid.UUID) (Entity, bool)
	GetMapEntityByUuid(uuid uuid.UUID) (Entity, bool)
	GetCharacterEntityByUuid(uuid uuid.UUID) (Entity, bool)
}

type BaseWorld struct {
	systems  []System
	entities []Entity

	UuidToEntity          map[uuid.UUID]Entity
	UuidToCharacterEntity map[uuid.UUID]Entity
	UuidToMapEntity       map[uuid.UUID]Entity
}

func NewBaseWorld() BaseWorld {
	return BaseWorld{
		UuidToEntity:          make(map[uuid.UUID]Entity),
		UuidToCharacterEntity: make(map[uuid.UUID]Entity),
		UuidToMapEntity:       make(map[uuid.UUID]Entity),
	}
}

func (w *BaseWorld) AddEntity(e Entity) {
	w.UuidToEntity[e.GetId()] = e

	// @todo Add mutual exclusive component check?
	// @todo Add filter check?

	if e.HasComponentType(CharacterComponentType) {
		w.UuidToCharacterEntity[e.GetId()] = e
	}
	if e.HasComponentType(MapComponentType) {
		w.UuidToMapEntity[e.GetId()] = e
	}

	w.entities = append(w.entities, e)
}

func (w *BaseWorld) AddEntities(e []Entity) {
	for _, entity := range e {
		w.AddEntity(entity)
	}
}

func (w *BaseWorld) GetCharacterEntities() []Entity {
	return w.getEntityValuesOfMap(w.UuidToCharacterEntity)
}

func (w *BaseWorld) GetMapEntities() []Entity {
	return w.getEntityValuesOfMap(w.UuidToMapEntity)
}

func (w *BaseWorld) GetEntityByUuid(uuid uuid.UUID) (Entity, bool) {
	entity, ok := w.UuidToEntity[uuid]
	return entity, ok
}

func (w *BaseWorld) GetMapEntityByUuid(uuid uuid.UUID) (Entity, bool) {
	entity, ok := w.UuidToMapEntity[uuid]
	return entity, ok
}

func (w *BaseWorld) GetCharacterEntityByUuid(uuid uuid.UUID) (Entity, bool) {
	entity, ok := w.UuidToCharacterEntity[uuid]
	return entity, ok
}

func (w *BaseWorld) getEntityValuesOfMap(dict map[uuid.UUID]Entity) []Entity {
	var values []Entity
	for _, v := range dict {
		values = append(values, v)
	}

	return values
}
