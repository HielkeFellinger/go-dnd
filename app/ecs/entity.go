package ecs

import "github.com/google/uuid"

type Entity interface {
	AddComponent(c Component)
}

type BaseEntity struct {
	Id          uuid.UUID
	Name        string
	Description string
	Components  []Component
}

func NewEntity() BaseEntity {
	return BaseEntity{
		Id: uuid.New(),
	}
}

func (e *BaseEntity) widthName(name string) *BaseEntity {
	e.Name = name
	return e
}

func (e *BaseEntity) widthDescription(description string) *BaseEntity {
	e.Description = description
	return e
}

func (e *BaseEntity) AddComponent(c Component) {
	e.Components = append(e.Components, c)
}
