package ecs_builders

import (
	"github.com/hielkefellinger/go-dnd/app/ecs"
	"github.com/hielkefellinger/go-dnd/app/ecs_components"
)

func BuildBasicBlockerEntity() (ecs.Entity, error) {
	entity := ecs.NewEntity()

	component := ecs_components.NewBlockerComponent()
	if errAddComponent := entity.AddComponent(component); errAddComponent != nil {
		return nil, errAddComponent
	}

	// Hide by default
	hidden := ecs_components.NewVisibilityComponent().(*ecs_components.VisibilityComponent)
	hidden.Hidden = true
	if errAddComponent := entity.AddComponent(hidden); errAddComponent != nil {
		return nil, errAddComponent
	}

	return &entity, nil
}

func AddBasicBlockerToMap(world ecs.World, mapEntity ecs.Entity, posX uint, posY uint) (string, error) {
	basicEntity, errBasicBlocker := BuildBasicBlockerEntity()
	if errBasicBlocker != nil {
		return "", errBasicBlocker
	}

	// Add Entity
	if errAddEntity := world.AddEntity(basicEntity); errAddEntity != nil {
		return "", errAddEntity
	}

	mapItemRelation := ecs_components.NewMapItemRelationComponent().(*ecs_components.MapItemRelationComponent)
	var position = ecs_components.NewPositionComponent().(*ecs_components.PositionComponent)
	position.X = posX
	position.Y = posY
	mapItemRelation.Position = position

	if blockerMatch, ok := world.GetEntityByUuid(basicEntity.GetId()); ok {
		mapItemRelation.Entity = blockerMatch
	}

	if mapMatch, ok := world.GetEntityByUuid(mapEntity.GetId()); ok {
		if errAddRelComponent := mapMatch.AddComponent(mapItemRelation); errAddRelComponent != nil {
			return "", errAddRelComponent
		}
	}

	return mapItemRelation.GetId().String(), nil
}
