package ecs

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"slices"
)

type Entity interface {
	GetId() uuid.UUID
	GetName() string
	SetName(name string)
	GetVersion() uint
	GetDescription() string
	SetDescription(description string)
	HasComponentType(ct uint64) bool
	HasComponentByUuid(uuid uuid.UUID) bool
	GetComponentByUuid(uuid uuid.UUID) (Component, bool)
	RemoveComponentByUuid(uuid uuid.UUID) bool
	AddComponent(c Component) error
	LoadFromRawEntity(raw RawEntity) error
	GetAllComponents() []Component
	GetAllComponentsOfType(ct uint64) []Component
	hasCircularRef(uuid uuid.UUID) bool
	HasRelationWithEntityByUuid(uuid uuid.UUID) bool
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

func (e *BaseEntity) SetName(name string) {
	e.Name = name
}

func (e *BaseEntity) GetDescription() string {
	return e.Description
}

func (e *BaseEntity) SetDescription(description string) {
	e.Description = description
}

func (e *BaseEntity) LoadFromRawEntity(raw RawEntity) error {
	e.WithName(raw.Name).WithDescription(raw.Description)
	return nil
}

func (e *BaseEntity) HasComponentType(ct uint64) bool {
	_, ok := e.componentTypeToComponentArrMap[ct]
	return ok
}

func (e *BaseEntity) HasComponentByUuid(uuid uuid.UUID) bool {
	_, ok := e.uuidToComponent[uuid]
	return ok
}

func (e *BaseEntity) GetComponentByUuid(uuid uuid.UUID) (Component, bool) {
	component, ok := e.uuidToComponent[uuid]
	return component, ok
}

func (e *BaseEntity) RemoveComponentByUuid(uuid uuid.UUID) bool {
	// Get the component if it exists
	if component, ok := e.GetComponentByUuid(uuid); ok {
		// Clean up maps/dictionaries
		if items := e.GetAllComponentsOfType(component.ComponentType()); len(items) <= 1 {
			delete(e.componentTypeToComponentArrMap, component.ComponentType())
		} else {
			array := e.componentTypeToComponentArrMap[component.ComponentType()]
			if index := slices.Index(array, component); index > -1 {
				e.componentTypeToComponentArrMap[component.ComponentType()] = slices.Delete(array, index, index+1)
			}
		}
		delete(e.uuidToComponent, uuid)
		if component.IsRelationalComponent() {
			relComp := component.(RelationalComponent)
			delete(e.uuidToRelEntity, relComp.GetEntity().GetId())
			if index := slices.Index(e.RefEntities, relComp.GetEntity()); index > -1 {
				e.RefEntities = slices.Delete(e.RefEntities, index, index+1)
			}
		}
		// Remove Item
		if index := slices.Index(e.Components, component); index > -1 {
			e.Components = slices.Delete(e.Components, index, index+1)
		}
		return true
	}

	return false
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
		e.componentTypeToComponentArrMap[c.ComponentType()] = append(arr, c)
	} else {
		e.componentTypeToComponentArrMap[c.ComponentType()] = []Component{c}
	}

	return nil
}

func (e *BaseEntity) HasRelationWithEntityByUuid(uuid uuid.UUID) bool {
	_, ok := e.uuidToRelEntity[uuid]
	return ok
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

func (e *BaseEntity) GetAllComponentsOfType(ct uint64) []Component {
	if e.HasComponentType(ct) {
		components := e.componentTypeToComponentArrMap[ct]
		if components != nil && len(components) > 0 {
			return components
		}
	}

	return make([]Component, 0)
}

func (e *BaseEntity) GetAllComponents() []Component {
	return e.Components
}

func (e *BaseEntity) isComponentTypeAllowedToBeAdded(c Component) bool {
	if !c.AllowMultipleOfType() {
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
