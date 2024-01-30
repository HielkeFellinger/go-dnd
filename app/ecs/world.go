package ecs

import (
	"github.com/google/uuid"
)

type World interface {
	AddEntity(e Entity)
	AddEntities(e []Entity)
	GetCharacterEntities() []Entity
	GetMapEntities() []Entity
	GetEntityByUuid(uuid uuid.UUID) Entity
	GetMapEntityByUuid(uuid uuid.UUID) Entity
	GetCharacterEntityByUuid(uuid uuid.UUID) Entity
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

func (w *BaseWorld) GetEntityByUuid(uuid uuid.UUID) Entity {
	return w.UuidToEntity[uuid]
}

func (w *BaseWorld) GetMapEntityByUuid(uuid uuid.UUID) Entity {
	return w.UuidToMapEntity[uuid]
}

func (w *BaseWorld) GetCharacterEntityByUuid(uuid uuid.UUID) Entity {
	return w.UuidToCharacterEntity[uuid]
}

func (w *BaseWorld) getEntityValuesOfMap(dict map[uuid.UUID]Entity) []Entity {
	var values []Entity
	for _, v := range dict {
		values = append(values, v)
	}

	return values
}
