package ecs_components

import (
	"github.com/google/uuid"
	"github.com/hielkefellinger/go-dnd/app/ecs"
)

type LevelComponent struct {
	ecs.BaseComponent
	Level uint `yaml:"level"`
}

func NewLevelComponent() LevelComponent {
	return LevelComponent{
		BaseComponent: ecs.BaseComponent{Id: uuid.New()},
	}
}

func (c *LevelComponent) WithLevel(level uint) *LevelComponent {
	c.Level = level
	return c
}

func (c *LevelComponent) ComponentType() uint64 {
	return ecs.LevelComponentType
}
