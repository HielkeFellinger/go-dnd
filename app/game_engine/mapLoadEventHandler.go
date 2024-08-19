package game_engine

import (
	"encoding/json"
	"errors"
	"github.com/hielkefellinger/go-dnd/app/ecs"
	"github.com/hielkefellinger/go-dnd/app/ecs_components"
	"github.com/hielkefellinger/go-dnd/app/ecs_model_translation"
	"github.com/hielkefellinger/go-dnd/app/helpers"
	"github.com/hielkefellinger/go-dnd/app/models"
	"golang.org/x/net/html"
	"log"
	"slices"
	"strconv"
)

func (e *baseEventMessageHandler) handleMapLoadEvents(message EventMessage, pool CampaignPool) error {
	log.Printf("- Map. Event Type: '%d' Message: '%s'", message.Type, message.Id)
	if message.Type == TypeLoadMap || message.Type == TypeLoadFullGame {
		e.typeLoadMap(message, pool)
	}
	if message.Type == TypeLoadMapEntities || message.Type == TypeLoadFullGame {
		e.typeLoadMapEntities(message, pool)
	}
	if message.Type == TypeLoadUpsertMap {
		return e.typeLoadUpsertMap(message, pool)
	}
	if message.Type == TypeUpsertMap {
		return e.typeUpsertMap(message, pool)
	}
	if message.Type == TypeLoadMapEntity {
		return e.typeLoadMapEntity(message, pool)
	}

	return nil
}

func (e *baseEventMessageHandler) typeUpsertMap(message EventMessage, pool CampaignPool) error {
	// Undo escaping
	clearedBody := html.UnescapeString(message.Body)

	// Check if user is lead
	if message.Source != pool.GetLeadId() {
		return errors.New("modifying items is not allowed as non-lead")
	}

	var mapUpdateRequest MapUpsertRequest
	err := json.Unmarshal([]byte(clearedBody), &mapUpdateRequest)
	if err != nil {
		return err
	}

	// Check if it needs to be updated; or inserted
	var newEntry = mapUpdateRequest.Id == ""
	var mapEntity ecs.BaseEntity
	if !newEntry {
		mapUuid, err := parseStingToUuid(mapUpdateRequest.Id)
		if err != nil {
			return err
		}
		mapEntityFromWorld, match := pool.GetEngine().GetWorld().GetMapEntityByUuid(mapUuid)
		if !match || mapEntityFromWorld == nil {
			return e.sendManagementError("Error", "failure of loading MAP by UUID", pool)
		}
		mapEntity = *mapEntityFromWorld.(*ecs.BaseEntity)
	} else {
		mapEntity = ecs.NewEntity()
		if addErr := mapEntity.AddComponent(ecs_components.NewMapComponent()); addErr != nil {
			return e.sendManagementError("Error", addErr.Error(), pool)
		}
		if addErr := mapEntity.AddComponent(ecs_components.NewAreaComponent()); addErr != nil {
			return e.sendManagementError("Error", addErr.Error(), pool)
		}
	}

	// Parse Area
	mapAreaComponent := ecs_components.NewAreaComponent().(*ecs_components.AreaComponent)
	if mapUpdateRequest.X != "" {
		if areaError := mapAreaComponent.WidthFromString(mapUpdateRequest.X); areaError != nil {
			return e.sendManagementError("Error", areaError.Error(), pool)
		}
		if mapAreaComponent.Width < 2 {
			return e.sendManagementError("Error", "Invalid width (X). Should be > 2", pool)
		}
	}
	if mapUpdateRequest.Y != "" {
		if areaError := mapAreaComponent.LengthFromString(mapUpdateRequest.Y); areaError != nil {
			return e.sendManagementError("Error", areaError.Error(), pool)
		}
		if mapAreaComponent.Length < 2 {
			return e.sendManagementError("Error", "Invalid length (Y). Should be > 2", pool)
		}
	}
	if newEntry && (mapAreaComponent.Length < 2 || mapAreaComponent.Width < 2) {
		return e.sendManagementError("Error", "Invalid length and or width set. Should be > 2", pool)
	}

	// Block editing of active maps.
	mapComponents := mapEntity.GetAllComponentsOfType(ecs.MapComponentType)
	if len(mapComponents) > 1 {
		mapComponent := mapComponents[0].(*ecs_components.MapComponent)
		if mapComponent.Active {
			return e.sendManagementError("Error", "Can not modify active maps", pool)
		}
	}

	// Update if add image
	if mapUpdateRequest.Image != (helpers.FileUpload{}) && mapUpdateRequest.ImageName != "" {
		// Handle file-upload
		link, imageErr := helpers.SaveImageToCampaign(mapUpdateRequest.Image, pool.GetId(), mapUpdateRequest.ImageName)
		if imageErr != nil {
			return e.sendManagementError("Error", imageErr.Error(), pool)
		}
		imageComponent := ecs_components.NewImageComponent().(*ecs_components.ImageComponent)
		imageComponent.Name = mapUpdateRequest.ImageName
		imageComponent.Active = newEntry
		imageComponent.Url = link
		imageComponent.Version = 1
		if addErr := mapEntity.AddComponent(imageComponent); addErr != nil {
			return e.sendManagementError("Error", addErr.Error(), pool)
		}
	} else if newEntry {
		return e.sendManagementError("Error", "Missing image on new Map", pool)
	}

	// Update the needed removal of images
	if !newEntry && len(mapUpdateRequest.RemoveImages) > 0 {
		var imageComponents = mapEntity.GetAllComponentsOfType(ecs.ImageComponentType)
		if len(imageComponents) <= len(mapUpdateRequest.RemoveImages) {
			return e.sendManagementError("Error", "Can not delete all images; needs at least one", pool)
		}

		for index := range imageComponents {
			imageComponent := imageComponents[index].(*ecs_components.ImageComponent)
			if slices.Contains(mapUpdateRequest.RemoveImages, imageComponent.GetId().String()) {
				mapEntity.RemoveComponentByUuid(imageComponent.GetId())
			}
		}

		// Reload minus removed; set at least one of the leftover backgrounds to active
		imageComponents = mapEntity.GetAllComponentsOfType(ecs.ImageComponentType)
		var firstImageComponent *ecs_components.ImageComponent
		var hasActive = false
		for index := range imageComponents {
			imageComponent := imageComponents[index].(*ecs_components.ImageComponent)
			if index == 0 {
				firstImageComponent = imageComponent
			}
			hasActive = hasActive || imageComponent.Active
		}
		if !hasActive {
			firstImageComponent.Active = true
		}
	}

	// Update Basis Entity Properties
	if mapUpdateRequest.Name != "" && mapUpdateRequest.Name != mapEntity.Name {
		log.Printf("Update Name")
		mapEntity.Name = mapUpdateRequest.Name
	} else if newEntry {
		return e.sendManagementError("Error", "Map should have a name", pool)
	}
	if mapUpdateRequest.Description != "" && mapUpdateRequest.Description != mapEntity.Description {
		log.Printf("Update Description")
		mapEntity.Description = mapUpdateRequest.Description
	}

	// Update Area
	if mapAreaComponent.Length >= 2 || mapAreaComponent.Width >= 2 {
		areaComponent := mapEntity.GetAllComponentsOfType(ecs.AreaComponentType)[0].(*ecs_components.AreaComponent)
		if mapAreaComponent.Length != areaComponent.Length && mapAreaComponent.Length >= 2 {
			areaComponent.Length = mapAreaComponent.Length
		}
		if mapAreaComponent.Width != areaComponent.Width && mapAreaComponent.Width >= 2 {
			areaComponent.Width = mapAreaComponent.Width
		}
	}

	// Add new Map Entity
	if newEntry {
		if addErr := pool.GetEngine().GetWorld().AddEntity(&mapEntity); addErr != nil {
			return e.sendManagementError("Error", addErr.Error(), pool)
		}
	}

	loadUpsertMapMessage := NewEventMessage()
	loadUpsertMapMessage.Source = pool.GetLeadId()
	loadUpsertMapMessage.Body = mapEntity.GetId().String()
	if typeLoadUpsertMapErr := e.typeLoadUpsertMap(loadUpsertMapMessage, pool); typeLoadUpsertMapErr != nil {
		return e.sendManagementError("Error", typeLoadUpsertMapErr.Error(), pool)
	}
	if newEntry {
		if typeManageMapsErr := e.typeManageMaps(loadUpsertMapMessage, pool); typeManageMapsErr != nil {
			return e.sendManagementError("Error", typeManageMapsErr.Error(), pool)
		}
	}
	return nil
}

func (e *baseEventMessageHandler) typeLoadUpsertMap(message EventMessage, pool CampaignPool) error {
	// Undo escaping
	clearedBody := html.UnescapeString(message.Body)

	// Check if user is lead
	if message.Source != pool.GetLeadId() {
		return errors.New("modifying maps is not allowed as non-lead")
	}

	data := make(map[string]any)

	// Check if there is an existing map with the supplied uuid
	uuidItemFilter, err := parseStingToUuid(clearedBody)
	if err == nil {
		mapCandidate, match := pool.GetEngine().GetWorld().GetEntityByUuid(uuidItemFilter)
		if match && mapCandidate.HasComponentType(ecs.MapComponentType) {
			data["Map"] = ecs_model_translation.MapEntityToCampaignMapModel(mapCandidate)
		}
	}

	rawJsonBytes, err := json.Marshal(
		e.handleLoadHtmlBodyMultipleTemplateFiles([]string{"campaignUpsertMap.html", "diceSpinnerSvg.html"},
			"campaignUpsertMap", data))
	if err != nil {
		return err
	}

	loadItemMessage := NewEventMessage()
	loadItemMessage.Source = message.Source
	loadItemMessage.Type = TypeLoadUpsertMap
	loadItemMessage.Body = string(rawJsonBytes)
	loadItemMessage.Destinations = append(loadItemMessage.Destinations, pool.GetLeadId())
	pool.TransmitEventMessage(loadItemMessage)

	return nil
}

func (e *baseEventMessageHandler) typeLoadMapEntity(message EventMessage, pool CampaignPool) error {
	// No filter in body equals no map entity to load
	if len(message.Body) == 0 {
		return nil
	}

	// Validate Filter
	uuidMapItemFilter, err := parseStingToUuid(message.Body)
	if err != nil {
		return err
	}

	var transmitMessage = NewEventMessage()
	transmitMessage.Type = TypeLoadMapEntity
	transmitMessage.Source = message.Source

	isLead := message.Source == pool.GetLeadId()

	mapEntities := pool.GetEngine().GetWorld().GetMapEntities()
	for _, mapEntity := range mapEntities {
		// Only get the map with the relevant entity
		if !mapEntity.HasComponentByUuid(uuidMapItemFilter) {
			continue
		}

		componentMap := ecs_model_translation.MapEntityToCampaignMapModel(mapEntity)
		mapItem, ok := mapEntity.GetComponentByUuid(uuidMapItemFilter)
		if !ok || mapItem == nil {
			return errors.New("failure loading MapItem")
		}

		// Translate Entity to controlling
		var mapItemModel = ecs_model_translation.MapItemEntityToCampaignMapItemElement(mapItem, componentMap.Id)

		// Lead Message
		mapItemModel = e.buildMapItem(mapItemModel, true, true)
		transmitMessage.Body = string(e.parseObjectToJson(mapItemModel))
		transmitMessage.Destinations = make([]string, 1)
		transmitMessage.Destinations[0] = pool.GetLeadId()
		pool.TransmitEventMessage(transmitMessage)

		// Get list of controlling players
		controllingPlayers := make([]string, 0)
		if componentMap.Active {
			controllingPlayers = append(controllingPlayers, mapItemModel.Controllers...)
		}

		// Controlling users Message
		if len(controllingPlayers) > 0 {
			mapItemModel = e.buildMapItem(mapItemModel, false, true)
			transmitMessage.Body = string(e.parseObjectToJson(mapItemModel))
			transmitMessage.Destinations = controllingPlayers
			pool.TransmitEventMessage(transmitMessage)
		}

		// - @todo team visibility
		// Check hidden non owner
		if mapItemModel.Hidden && !isLead && !slices.Contains(mapItemModel.Controllers, message.Source) {
			continue
		}

		managingPlayersFilter := append(controllingPlayers, pool.GetLeadId())
		nonControllingPlayers := pool.GetAllClientIds(managingPlayersFilter...)
		if componentMap.Active && len(nonControllingPlayers) > 0 {
			mapItemModel = e.buildMapItem(mapItemModel, false, false)
			transmitMessage.Body = string(e.parseObjectToJson(mapItemModel))
			transmitMessage.Destinations = nonControllingPlayers

			log.Printf("    Non-controlling players: %v", transmitMessage.Destinations)
			pool.TransmitEventMessage(transmitMessage)
		}
	}
	return nil
}

func (e *baseEventMessageHandler) typeLoadMapEntities(message EventMessage, pool CampaignPool) {
	var transmitMessage = NewEventMessage()
	transmitMessage.Type = TypeLoadMapEntities
	transmitMessage.Source = message.Source
	transmitMessage.Destinations = append(transmitMessage.Destinations, message.Source)

	// Check if is GM:
	isLead := message.Source == pool.GetLeadId()

	mapEntities := pool.GetEngine().GetWorld().GetMapEntities()
	for _, mapEntity := range mapEntities {
		// Only show filtered form body
		if len(message.Body) > 0 && mapEntity.GetId().String() != message.Body {
			continue
		}

		componentMap := ecs_model_translation.MapEntityToCampaignMapModel(mapEntity)

		// Only show enabled maps for player
		if !componentMap.Active && !isLead {
			continue
		}

		// load models
		mapItems := mapEntity.GetAllComponentsOfType(ecs.MapItemRelationComponentType)
		log.Printf("The mapItems: %v", mapItems)
		mapItemsModel := models.CampaignScreenMapItems{
			MapId:    componentMap.Id,
			Elements: make(map[string]models.CampaignScreenMapItemElement, len(mapItems)),
		}

		// Translate all items
		for _, mapItem := range mapItems {
			var mapItemModel = ecs_model_translation.MapItemEntityToCampaignMapItemElement(mapItem, mapItemsModel.MapId)

			// - @todo team visibility
			// Check hidden non owner
			if mapItemModel.Hidden && !isLead && !slices.Contains(mapItemModel.Controllers, message.Source) {
				continue
			}

			mapItemModel = e.buildMapItem(mapItemModel, isLead, isLead || slices.Contains(mapItemModel.Controllers, message.Source))
			mapItemsModel.Elements[mapItemModel.Id] = mapItemModel
		}

		rawJsonBytes, err := json.Marshal(mapItemsModel)
		if err != nil {
			log.Printf("Error parsing Loading Map Item content `%s`", err.Error())
		}

		transmitMessage.Body = string(rawJsonBytes)
		pool.TransmitEventMessage(transmitMessage)
	}
}

func (e *baseEventMessageHandler) typeLoadMap(message EventMessage, pool CampaignPool) {
	var transmitMessage = NewEventMessage()
	transmitMessage.Type = TypeLoadMap
	transmitMessage.Source = message.Source
	transmitMessage.Destinations = append(transmitMessage.Destinations, message.Source)

	// Check if is GM:
	isLead := message.Source == pool.GetLeadId()

	// @todo Load Focus Map Related Details
	// - Gray out non-present players;

	var campaignScreenContent = models.NewCampaignScreenContent()
	mapEntities := pool.GetEngine().GetWorld().GetMapEntities()
	charEntities := pool.GetEngine().GetWorld().GetCharacterEntities()
	for _, mapEntity := range mapEntities {
		// Only handle if id is requested from body
		if len(message.Body) > 0 && mapEntity.GetId().String() != message.Body {
			continue
		}

		// Translate
		componentMap := ecs_model_translation.MapEntityToCampaignMapModel(mapEntity)

		// Only show enabled maps for player (if not directly requested by lead)
		if !componentMap.Active && (!isLead || (isLead && len(message.Body) == 0)) {
			continue
		}

		var data = e.buildMapData(componentMap, isLead, charEntities)
		var content = models.CampaignContentItem{}
		var tab = models.CampaignTabItem{}
		tab.Id = componentMap.Id
		content.Id = componentMap.Id
		tab.Html = e.handleLoadHtmlBody("campaignSelector.html", "campaignSelector", data)
		content.Html = e.handleLoadHtmlBody("campaignContentMap.html", "campaignContentMap", data)

		campaignScreenContent.Tabs = append(campaignScreenContent.Tabs, tab)
		campaignScreenContent.Content = append(campaignScreenContent.Content, content)
	}

	rawJsonBytes, err := json.Marshal(campaignScreenContent)
	if err != nil {
		log.Printf("Error parsing Loading Map content `%s`", err.Error())
	}

	transmitMessage.Body = string(rawJsonBytes)
	pool.TransmitEventMessage(transmitMessage)
}

func (e *baseEventMessageHandler) parseObjectToJson(object any) []byte {
	rawJsonBytes, err := json.Marshal(object)
	if err != nil {
		log.Printf("Error parsing content to Json `%s`", err.Error())
	}
	return rawJsonBytes
}

func (e *baseEventMessageHandler) buildMapItem(mapItemModel models.CampaignScreenMapItemElement, isLead bool, hasControl bool) models.CampaignScreenMapItemElement {
	data := make(map[string]any)
	data["id"] = mapItemModel.Id
	data["mapId"] = mapItemModel.MapId
	data["hidden"] = mapItemModel.Hidden
	data["lead"] = isLead
	data["hasControl"] = hasControl
	data["entityName"] = mapItemModel.EntityName
	data["backgroundImage"] = mapItemModel.Image.Url

	if mapItemModel.HasHealth() {
		data["healthBar"] = true
		data["healthTotal"] = mapItemModel.Health.Total
		data["healthCurrent"] = mapItemModel.Health.Current

		// @todo health colour changing progressbar
	}

	mapItemModel.Html = e.handleLoadHtmlBody("campaignContentMapCell.html", "campaignContentMapCell", data)
	return mapItemModel
}

func (e *baseEventMessageHandler) buildMapData(model models.CampaignMap, isLead bool, characters []ecs.Entity) map[string]any {
	data := make(map[string]any)
	data["type"] = "Map"
	data["id"] = model.Id
	data["name"] = model.Name
	data["lead"] = isLead
	data["active"] = model.Active
	data["backgroundImage"] = model.ActiveImage
	data["altImage"] = model.Images
	data["characters"] = characters

	xVal := make([]string, model.X)
	yVal := make([]string, model.Y)
	for i := range xVal {
		xVal[i] = strconv.Itoa(i)
	}
	for i := range yVal {
		yVal[i] = strconv.Itoa(i)
	}
	data["x"] = xVal
	data["y"] = yVal
	return data
}

type MapUpsertRequest struct {
	Id           string             `json:"Id"`
	Name         string             `json:"Name"`
	Description  string             `json:"Description"`
	X            string             `json:"X"`
	Y            string             `json:"Y"`
	ImageName    string             `json:"ImageName"`
	RemoveImages []string           `json:"RemoveImages"`
	Image        helpers.FileUpload `json:"Image"`
}
