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
	"slices"
	"strconv"
)

func (e *baseEventMessageHandler) handleUpdateCharacterEvents(message EventMessage, pool CampaignPool) error {
	log.Printf("- Char. Update Event Type: '%d' Message: '%s'", message.Type, message.Id)
	if message.Type == TypeUpdateCharacterHealth {
		return e.typeUpdateCharacterHealth(message, pool)
	} else if message.Type == TypeUpdateCharacterUsers {
		return e.typeUpdateCharacterUsers(message, pool)
	}

	return nil
}

func (e *baseEventMessageHandler) typeUpdateCharacterUsers(message EventMessage, pool CampaignPool) error {
	// Undo escaping
	clearedBody := html.UnescapeString(message.Body)

	// Parse body
	var characterToPlayerLink characterToPlayer
	err := json.Unmarshal([]byte(clearedBody), &characterToPlayerLink)
	if err != nil {
		return err
	}

	// Validate UUID Filter form message
	uuidCharFilter, err := parseStingToUuid(characterToPlayerLink.Id)
	if err != nil {
		return err
	}

	// Test if Character exists
	var charEntity ecs.Entity
	charEntity, ok := pool.GetEngine().GetWorld().GetCharacterEntityByUuid(uuidCharFilter)
	if !ok || charEntity == nil {
		return errors.New("filter UUID has no match")
	}

	// Check if player is (already) linked to campaign
	playerHasControl := false
	var playerComponent *ecs_components.PlayerComponent
	playerComponents := charEntity.GetAllComponentsOfType(ecs.PlayerComponentType)
	for _, component := range playerComponents {
		tmpPlayerComponent := component.(*ecs_components.PlayerComponent)
		if tmpPlayerComponent.Name == characterToPlayerLink.PlayerName {
			// Remove possible duplicates
			if playerHasControl {
				charEntity.RemoveComponentByUuid(tmpPlayerComponent.Id)
			}
			playerHasControl = true
			playerComponent = tmpPlayerComponent
		}
	}

	// Clean-up or add player
	updatedChar := false
	if characterToPlayerLink.Connect {
		// Player is already connected
		if playerHasControl {
			return nil
		}

		// Add the player
		newPlayerComponent := ecs_components.NewPlayerComponent().(*ecs_components.PlayerComponent)
		newPlayerComponent.Name = characterToPlayerLink.PlayerName
		if failedToAdd := charEntity.AddComponent(newPlayerComponent); failedToAdd != nil {
			return failedToAdd
		}
		updatedChar = true
	} else if playerHasControl && playerComponent != nil {

		// Always ensure a character remains a player character
		if len(playerComponents) == 1 {
			if failedToAdd := charEntity.AddComponent(ecs_components.NewPlayerComponent()); failedToAdd != nil {
				return failedToAdd
			}
		}
		// Remove (if it has a match)
		charEntity.RemoveComponentByUuid(playerComponent.Id)
		updatedChar = true
	}

	// >> Ensure it remains a playable character

	// Send updates
	if updatedChar {
		// Update ribbon
		var reloadCharRibbon = NewEventMessage()
		reloadCharRibbon.Source = ServerUser
		reloadCharRibbon.Type = TypeLoadCharacters
		e.loadCharacters(reloadCharRibbon, pool)

		// Update possible maps with this char
		if err := e.updateAllPossibleMapsOfChar(pool, uuidCharFilter); err != nil {
			return err
		}

		// Remove Char details on players who should not see it anymore
		if !characterToPlayerLink.Connect {
			var closeCharDetailMessage = NewEventMessage()
			closeCharDetailMessage.Source = characterToPlayerLink.PlayerName
			closeCharDetailMessage.Type = TypeLoadCharactersDetails
			closeCharDetailMessage.Body = charEntity.GetId().String()
			if err := e.loadCharactersDetails(closeCharDetailMessage, pool); err != nil {
				return err
			}
		}
	}

	return nil
}

func (e *baseEventMessageHandler) typeUpdateCharacterHealth(message EventMessage, pool CampaignPool) error {
	// Undo escaping
	clearedBody := html.UnescapeString(message.Body)

	// Attempt to parse the campaign screen map messageMapItem
	var characterHealth models.CampaignCharacterHealth
	err := json.Unmarshal([]byte(clearedBody), &characterHealth)
	if err != nil {
		return err
	}

	// Validate UUID Filter form message
	uuidCharFilter, err := parseStingToUuid(characterHealth.Id)
	if err != nil {
		return err
	}

	// Test if Character exists
	var charEntity ecs.Entity
	charEntity, ok := pool.GetEngine().GetWorld().GetCharacterEntityByUuid(uuidCharFilter)
	if !ok || charEntity == nil {
		return errors.New("filter UUID has no match")
	}

	// Check if source is allowed to update
	controllingPlayers := make([]string, 0)
	controllingPlayers = append(controllingPlayers, pool.GetLeadId())
	playerComponents := charEntity.GetAllComponentsOfType(ecs.PlayerComponentType)
	for _, component := range playerComponents {
		tmpPlayerComponent := component.(*ecs_components.PlayerComponent)
		if !slices.Contains(controllingPlayers, tmpPlayerComponent.Name) {
			controllingPlayers = append(controllingPlayers, tmpPlayerComponent.Name)
		}
	}
	if !slices.Contains(controllingPlayers, message.Source) {
		return errors.New("player is not allowed to update character health")
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

	// @todo check if player died?

	// Trigger Visual Updates
	var reloadCharDetailMessage = NewEventMessage()
	reloadCharDetailMessage.Source = ServerUser
	reloadCharDetailMessage.Body = characterHealth.Id
	loadCharErr := e.loadCharactersDetails(reloadCharDetailMessage, pool)
	if loadCharErr != nil {
		return loadCharErr
	}

	// Update all maps with char
	if err := e.updateAllPossibleMapsOfChar(pool, uuidCharFilter); err != nil {
		return err
	}

	return nil
}

func (e *baseEventMessageHandler) updateAllPossibleMapsOfChar(pool CampaignPool, uuidCharFilter uuid.UUID) error {
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

type characterToPlayer struct {
	Id         string `json:"id"`
	Connect    bool   `json:"connect"`
	PlayerName string `json:"playerName"`
}
