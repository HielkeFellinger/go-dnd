package ecs

import (
	"errors"
	"github.com/google/uuid"
	"log"
	"slices"
)

type World interface {
	AddEntity(e Entity) error
	RemoveEntity(e Entity) error
	AddEntities(e []Entity) error
	GetCharacterEntities() []Entity
	GetPlayerCharacterEntities() []Entity
	GetMapEntities() []Entity
	GetItemEntities() []Entity
	GetFactionEntities() []Entity
	GetInventoryEntities() []Entity
	GetMapContentEntities() []Entity
	GetOtherEntities() []Entity
	GetEntityByUuid(uuid uuid.UUID) (Entity, bool)
	GetMapEntityByUuid(uuid uuid.UUID) (Entity, bool)
	GetItemEntityByUuid(uuid uuid.UUID) (Entity, bool)
	GetOtherEntityByUuid(uuid uuid.UUID) (Entity, bool)
	GetMapContentEntityByUuid(uuid uuid.UUID) (Entity, bool)
	GetInventoryEntityByUuid(uuid uuid.UUID) (Entity, bool)
	GetCharacterEntityByUuid(uuid uuid.UUID) (Entity, bool)
	DoOtherEntitiesHaveARelationToSpecificEntity(matchEntity Entity, ignoredUidsFilter []uuid.UUID) bool
	GetAllEntitiesWithRelationToSpecificEntity(matchEntity Entity, ignoredUidsFilter []uuid.UUID) []Entity
}

type BaseWorld struct {
	systems  []System
	entities []Entity

	UuidToEntity           map[uuid.UUID]Entity
	UuidToItemEntity       map[uuid.UUID]Entity
	UuidToCharacterEntity  map[uuid.UUID]Entity
	UuidToMapEntity        map[uuid.UUID]Entity
	UuidToFactionEntity    map[uuid.UUID]Entity
	UuidToInventoryEntity  map[uuid.UUID]Entity
	UuidToMapContentEntity map[uuid.UUID]Entity
	UuidToOtherEntity      map[uuid.UUID]Entity
}

func NewBaseWorld() BaseWorld {
	return BaseWorld{
		UuidToEntity:           make(map[uuid.UUID]Entity),
		UuidToItemEntity:       make(map[uuid.UUID]Entity),
		UuidToCharacterEntity:  make(map[uuid.UUID]Entity),
		UuidToMapEntity:        make(map[uuid.UUID]Entity),
		UuidToFactionEntity:    make(map[uuid.UUID]Entity),
		UuidToInventoryEntity:  make(map[uuid.UUID]Entity),
		UuidToMapContentEntity: make(map[uuid.UUID]Entity),
		UuidToOtherEntity:      make(map[uuid.UUID]Entity),
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
	} else if e.HasComponentType(InventoryComponentType) {
		w.UuidToInventoryEntity[e.GetId()] = e
	} else if e.HasComponentType(ItemComponentType) {
		w.UuidToItemEntity[e.GetId()] = e
	} else if e.HasComponentType(MapContentComponentType) || e.HasComponentType(BlockerComponentType) {
		w.UuidToMapContentEntity[e.GetId()] = e
	} else {
		// Add Leftover
		w.UuidToOtherEntity[e.GetId()] = e
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

func (w *BaseWorld) RemoveEntity(e Entity) error {
	rawItem, match := w.UuidToEntity[e.GetId()]
	if !match {
		return errors.New("entity with UUID already removed")
	}

	// Cleanup ref's
	delete(w.UuidToEntity, e.GetId())
	delete(w.UuidToItemEntity, e.GetId())
	delete(w.UuidToCharacterEntity, e.GetId())
	delete(w.UuidToMapEntity, e.GetId())
	delete(w.UuidToFactionEntity, e.GetId())
	delete(w.UuidToInventoryEntity, e.GetId())
	delete(w.UuidToMapContentEntity, e.GetId())
	delete(w.UuidToOtherEntity, e.GetId())

	// Remove from entity
	w.entities = slices.DeleteFunc(w.entities, func(e Entity) bool { return e.GetId() == rawItem.GetId() })

	return nil
}

func (w *BaseWorld) DoOtherEntitiesHaveARelationToSpecificEntity(matchEntity Entity, ignoredUidsFilter []uuid.UUID) bool {
	for _, entity := range w.entities {
		// Skip refs to self and possible filtered Entities
		if entity.GetId() == matchEntity.GetId() || slices.Contains(ignoredUidsFilter, entity.GetId()) {
			continue
		}
		if entity.HasRelationWithEntityByUuid(matchEntity.GetId()) {
			return true
		}
	}
	return false
}

func (w *BaseWorld) GetAllEntitiesWithRelationToSpecificEntity(matchEntity Entity, ignoredUidsFilter []uuid.UUID) []Entity {
	entitiesWithRelation := make([]Entity, 0)

	for _, entity := range w.entities {
		// Skip refs to self and possible filtered Entities
		if entity.GetId() == matchEntity.GetId() || slices.Contains(ignoredUidsFilter, entity.GetId()) {
			continue
		}
		if entity.HasRelationWithEntityByUuid(matchEntity.GetId()) {
			entitiesWithRelation = append(entitiesWithRelation, entity)
		}
	}

	return entitiesWithRelation
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

func (w *BaseWorld) GetMapContentEntities() []Entity {
	return w.getEntityValuesOfMap(w.UuidToMapContentEntity)
}

func (w *BaseWorld) GetOtherEntities() []Entity {
	return w.getEntityValuesOfMap(w.UuidToOtherEntity)
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

func (w *BaseWorld) GetInventoryEntityByUuid(uuid uuid.UUID) (Entity, bool) {
	entity, ok := w.UuidToInventoryEntity[uuid]
	return entity, ok
}

func (w *BaseWorld) GetCharacterEntityByUuid(uuid uuid.UUID) (Entity, bool) {
	entity, ok := w.UuidToCharacterEntity[uuid]
	return entity, ok
}

func (w *BaseWorld) GetOtherEntityByUuid(uuid uuid.UUID) (Entity, bool) {
	entity, ok := w.UuidToOtherEntity[uuid]
	return entity, ok
}

func (w *BaseWorld) GetMapContentEntityByUuid(uuid uuid.UUID) (Entity, bool) {
	entity, ok := w.UuidToMapContentEntity[uuid]
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
