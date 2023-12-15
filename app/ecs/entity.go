package ecs

import (
	"github.com/google/uuid"
)

type Entity interface {
	AddComponent(c Component)
	LoadFromRawEntity(raw RawEntity) error
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

func (e *BaseEntity) LoadFromRawEntity(raw RawEntity) error {
	e.WithName(raw.Name).WithDescription(raw.Description)
	return nil
}

func (e *BaseEntity) WithName(name string) *BaseEntity {
	e.Name = name
	return e
}

func (e *BaseEntity) WithDescription(description string) *BaseEntity {
	e.Description = description
	return e
}

func (e *BaseEntity) AddComponent(c Component) {
	e.Components = append(e.Components, c)
}
