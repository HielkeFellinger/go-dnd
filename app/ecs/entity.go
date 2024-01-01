package ecs

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
)

type Entity interface {
	GetId() uuid.UUID
	GetName() string
	GetVersion() uint
	HasComponentType(ct uint64) bool
	AddComponent(c Component) error
	LoadFromRawEntity(raw RawEntity) error
	hasCircularRef(uuid uuid.UUID) bool
}

type BaseEntity struct {
	Id                             uuid.UUID
	Name                           string
	Version                        uint // Placeholder; take changes into account
	Description                    string
	ComponentFlag                  uint64
	Components                     []Component
	RefEntities                    []Entity
	uuidToComponent                map[uuid.UUID]Component
	uuidToRelEntity                map[uuid.UUID]Entity
	componentTypeToComponentArrMap map[uint64][]Component
}

func NewEntity() BaseEntity {
	return BaseEntity{
		Id:                             uuid.New(),
		uuidToComponent:                make(map[uuid.UUID]Component),
		uuidToRelEntity:                make(map[uuid.UUID]Entity),
		componentTypeToComponentArrMap: make(map[uint64][]Component),
	}
}

func (e *BaseEntity) GetId() uuid.UUID {
	return e.Id
}

func (e *BaseEntity) GetVersion() uint {
	return e.Version
}

func (e *BaseEntity) GetName() string {
	return e.Name
}

func (e *BaseEntity) LoadFromRawEntity(raw RawEntity) error {
	e.WithName(raw.Name).WithDescription(raw.Description)
	return nil
}

func (e *BaseEntity) HasComponentType(ct uint64) bool {
	_, ok := e.componentTypeToComponentArrMap[ct]
	return ok
}

func (e *BaseEntity) AddComponent(c Component) error {
	// Check uniqueness
	if !e.isComponentUuidUnique(c) {
		return errors.New(fmt.Sprintf("Skip adding component '%s' due to duplicate "+
			"component UUID in Entity: '%s'", c.GetId().String(), e.Id.String()))
	}
	if !e.isComponentTypeAllowedToBeAdded(c) {
		return errors.New(fmt.Sprintf("Skip adding component '%d' due to duplicate "+
			"component type in Entity: '%s'", c.ComponentType(), e.Id.String()))
	}

	// Check Nesting
	err := e.checkIfRelationalComponentIsAllowedToBeAdded(c)
	if err != nil {
		return err
	}

	e.Components = append(e.Components, c)
	e.uuidToComponent[c.GetId()] = c
	if arr, ok := e.componentTypeToComponentArrMap[c.ComponentType()]; ok {
		arr = append(arr, c)
	} else {
		e.componentTypeToComponentArrMap[c.ComponentType()] = []Component{c}
	}

	return nil
}

// Recursive Check if an Uuid is present is children @todo Optimize; is this REALLY needed?
func (e *BaseEntity) hasCircularRef(uuid uuid.UUID) bool {
	for _, relEntity := range e.uuidToRelEntity {
		if relEntity.GetId() == uuid {
			return true
		}
		if relEntity.hasCircularRef(uuid) {
			return true
		}
	}
	return false
}

func (e *BaseEntity) WithName(name string) *BaseEntity {
	e.Name = name
	return e
}

func (e *BaseEntity) WithDescription(description string) *BaseEntity {
	e.Description = description
	return e
}

func (e *BaseEntity) checkIfRelationalComponentIsAllowedToBeAdded(c Component) error {
	if c.IsRelationalComponent() {
		if relEntity, ok := c.(RelationalComponent); ok {
			childEntity := relEntity.GetEntity()
			if childEntity != nil {
				if childEntity.GetId() == e.GetId() {
					return errors.New(fmt.Sprintf("Skip adding component '%s' due to containing "+
						"its parent Entity: '%s'", c.GetId().String(), e.Id.String()))
				}
				if relEntity.GetEntity().hasCircularRef(e.Id) {
					return errors.New(fmt.Sprintf("Skip adding component '%s' due to containing "+
						"its parent Entity: '%s' in one of its RelationComponents", c.GetId().String(), e.Id.String()))
				}
			} else {
				return errors.New(fmt.Sprintf("Skip adding component '%s' due to not containing an Entity.",
					c.GetId().String()))
			}

			// Add to map
			e.uuidToRelEntity[childEntity.GetId()] = childEntity
		} else {
			return errors.New(fmt.Sprintf("Skip adding component '%s' due not implementing RelationalComponent interface.",
				c.GetId().String()))
		}
	}
	return nil
}

func (e *BaseEntity) isComponentTypeAllowedToBeAdded(c Component) bool {
	if !c.IsRelationalComponent() {
		for _, component := range e.Components {
			if component.ComponentType() == c.ComponentType() {
				return false
			}
		}
	}
	return true
}

func (e *BaseEntity) isComponentUuidUnique(c Component) bool {
	_, match := e.uuidToComponent[c.GetId()]
	return !match
}
