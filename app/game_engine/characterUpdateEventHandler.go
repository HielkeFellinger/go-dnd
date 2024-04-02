package game_engine

import (
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/hielkefellinger/go-dnd/app/ecs"
	"github.com/hielkefellinger/go-dnd/app/ecs_components"
	"github.com/hielkefellinger/go-dnd/app/models"
	"golang.org/x/net/html"
	"log"
	"strconv"
)

func (e *baseEventMessageHandler) handleUpdateCharacterEvents(message EventMessage, pool CampaignPool) error {
	log.Printf("- Char. Update Event Type: '%d' Message: '%s'", message.Type, message.Id)
	if message.Type == TypeUpdateCharacterHealth {
		return e.updateCharacterHealth(message, pool)
	}

	return nil
}

func (e *baseEventMessageHandler) updateCharacterHealth(message EventMessage, pool CampaignPool) error {
	// Undo escaping
	clearedBody := html.UnescapeString(message.Body)

	// Attempt to parse the campaign screen map messageMapItem
	var characterHealth models.CampaignCharacterHealth
	err := json.Unmarshal([]byte(clearedBody), &characterHealth)
	if err != nil {
		return err
	}

	// Validate UUID Filter form message
	var uuidCharFilter uuid.UUID
	if savedUuid, err := uuid.Parse(characterHealth.Id); err == nil {
		uuidCharFilter = savedUuid
	} else {
		return err
	}

	// Test if Character exists
	var charEntity ecs.Entity
	charEntity, ok := pool.GetEngine().GetWorld().GetCharacterEntityByUuid(uuidCharFilter)
	if !ok || charEntity == nil {
		return errors.New("filter UUID has no match")
	}

	// Parse health values
	damage, errD := strconv.Atoi(characterHealth.Damage)
	if errD != nil {
		damage = 0
	}
	temp, errT := strconv.Atoi(characterHealth.TemporaryHitPoints)
	if errT != nil {
		temp = 0
	}
	maxHP, errM := strconv.Atoi(characterHealth.MaximumHitPoints)
	if errM != nil {
		maxHP = 0
	}

	// @todo Check if source is allowed to update

	// Change or add the HealthComponent
	healthComponents := charEntity.GetAllComponentsOfType(ecs.HealthComponentType)
	if len(healthComponents) >= 1 {
		healthComponent := healthComponents[0].(*ecs_components.HealthComponent)
		healthComponent.Damage = uint(damage)
		healthComponent.Temporary = uint(temp)
		healthComponent.Maximum = uint(maxHP)
	} else {
		healthComponent := ecs_components.NewHealthComponent().(*ecs_components.HealthComponent)
		healthComponent.Damage = uint(damage)
		healthComponent.Temporary = uint(temp)
		healthComponent.Maximum = uint(maxHP)
		errAdd := charEntity.AddComponent(healthComponent)
		if errAdd != nil {
			return errAdd
		}
	}

	// Trigger Visual Updates
	var reloadCharDetailMessage = NewEventMessage()
	reloadCharDetailMessage.Source = ServerUser
	reloadCharDetailMessage.Body = characterHealth.Id
	loadCharErr := e.loadCharactersDetails(reloadCharDetailMessage, pool)
	if loadCharErr != nil {
		return loadCharErr
	}

	// Update possible Map Entities
	mapEntities := pool.GetEngine().GetWorld().GetMapEntities()
	for _, mapEntity := range mapEntities {

		// Only get the map with the relevant relation to entity
		if !mapEntity.HasRelationWithEntityByUuid(uuidCharFilter) {
			continue
		}

		for _, mapItem := range mapEntity.GetAllComponentsOfType(ecs.MapItemRelationComponentType) {
			mapItemRelComponent := mapItem.(*ecs_components.MapItemRelationComponent)

			if mapItemRelComponent.Entity.GetId() == uuidCharFilter {
				var reloadMapItemMessage = NewEventMessage()
				reloadMapItemMessage.Source = ServerUser
				reloadMapItemMessage.Body = mapItemRelComponent.Id.String()
				reloadMapItemErr := e.typeLoadMapEntity(reloadMapItemMessage, pool)
				if reloadMapItemErr != nil {
					return reloadMapItemErr
				}
			}
		}
	}

	return nil
}
