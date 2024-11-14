package ecs

import (
	"errors"
	"github.com/google/uuid"
	"log"
)

type World interface {
	AddEntity(e Entity) error
	AddEntities(e []Entity) error
	GetCharacterEntities() []Entity
	GetPlayerCharacterEntities() []Entity
	GetMapEntities() []Entity
	GetItemEntities() []Entity
	GetFactionEntities() []Entity
	GetInventoryEntities() []Entity
	GetEntityByUuid(uuid uuid.UUID) (Entity, bool)
	GetMapEntityByUuid(uuid uuid.UUID) (Entity, bool)
	GetItemEntityByUuid(uuid uuid.UUID) (Entity, bool)
	GetCharacterEntityByUuid(uuid uuid.UUID) (Entity, bool)
}

type BaseWorld struct {
	systems  []System
	entities []Entity

	UuidToEntity          map[uuid.UUID]Entity
	UuidToItemEntity      map[uuid.UUID]Entity
	UuidToCharacterEntity map[uuid.UUID]Entity
	UuidToMapEntity       map[uuid.UUID]Entity
	UuidToFactionEntity   map[uuid.UUID]Entity
	UuidToInventoryEntity map[uuid.UUID]Entity
}

func NewBaseWorld() BaseWorld {
	return BaseWorld{
		UuidToEntity:          make(map[uuid.UUID]Entity),
		UuidToItemEntity:      make(map[uuid.UUID]Entity),
		UuidToCharacterEntity: make(map[uuid.UUID]Entity),
		UuidToMapEntity:       make(map[uuid.UUID]Entity),
		UuidToFactionEntity:   make(map[uuid.UUID]Entity),
		UuidToInventoryEntity: make(map[uuid.UUID]Entity),
	}
}

func (w *BaseWorld) AddEntity(e Entity) error {
	// Check UUID uniqueness
	_, match := w.UuidToEntity[e.GetId()]
	if match {
		return errors.New("entity with UUID already exists")
	}

	w.UuidToEntity[e.GetId()] = e

	// @todo Add mutual exclusive component check? YES, SHOULD BE CHECKED
	// @todo Add filter check? YES, DO

	if e.HasComponentType(CharacterComponentType) {
		w.UuidToCharacterEntity[e.GetId()] = e
	} else if e.HasComponentType(MapComponentType) {
		w.UuidToMapEntity[e.GetId()] = e
	} else if e.HasComponentType(FactionComponentType) {
		w.UuidToFactionEntity[e.GetId()] = e
	} else if e.HasComponentType(SlotsComponentType) {
		w.UuidToInventoryEntity[e.GetId()] = e
	} else if e.HasComponentType(ItemComponentType) {
		w.UuidToItemEntity[e.GetId()] = e
	}

	w.entities = append(w.entities, e)
	return nil
}

func (w *BaseWorld) AddEntities(e []Entity) error {
	for _, entity := range e {
		if err := w.AddEntity(entity); err != nil {
			log.Fatal(err)
			return err
		}
	}
	return nil
}

func (w *BaseWorld) GetCharacterEntities() []Entity {
	return w.getEntityValuesOfMap(w.UuidToCharacterEntity)
}

func (w *BaseWorld) GetPlayerCharacterEntities() []Entity {
	return w.getFilteredEntityValuesOfMap(w.UuidToCharacterEntity,
		func(entity Entity) bool { return entity.HasComponentType(PlayerComponentType) })
}

func (w *BaseWorld) GetMapEntities() []Entity {
	return w.getEntityValuesOfMap(w.UuidToMapEntity)
}

func (w *BaseWorld) GetItemEntities() []Entity {
	return w.getEntityValuesOfMap(w.UuidToItemEntity)
}

func (w *BaseWorld) GetFactionEntities() []Entity {
	return w.getEntityValuesOfMap(w.UuidToFactionEntity)
}

func (w *BaseWorld) GetInventoryEntities() []Entity {
	return w.getEntityValuesOfMap(w.UuidToInventoryEntity)
}

func (w *BaseWorld) GetEntityByUuid(uuid uuid.UUID) (Entity, bool) {
	entity, ok := w.UuidToEntity[uuid]
	return entity, ok
}

func (w *BaseWorld) GetMapEntityByUuid(uuid uuid.UUID) (Entity, bool) {
	entity, ok := w.UuidToMapEntity[uuid]
	return entity, ok
}

func (w *BaseWorld) GetItemEntityByUuid(uuid uuid.UUID) (Entity, bool) {
	entity, ok := w.UuidToItemEntity[uuid]
	return entity, ok
}

func (w *BaseWorld) GetCharacterEntityByUuid(uuid uuid.UUID) (Entity, bool) {
	entity, ok := w.UuidToCharacterEntity[uuid]
	return entity, ok
}

func (w *BaseWorld) getEntityValuesOfMap(dict map[uuid.UUID]Entity) []Entity {
	values := make([]Entity, len(dict))
	index := 0
	for _, value := range dict {
		values[index] = value
		index++
	}
	return values
}

func (w *BaseWorld) getFilteredEntityValuesOfMap(dict map[uuid.UUID]Entity, filter func(entity Entity) bool) []Entity {
	values := make([]Entity, 0)
	for _, value := range dict {
		if filter(value) {
			values = append(values, value)
		}
	}
	return values
}
