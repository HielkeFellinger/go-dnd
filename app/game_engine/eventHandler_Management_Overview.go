package game_engine

import (
	"encoding/json"
	"errors"
	"github.com/hielkefellinger/go-dnd/app/ecs"
	"github.com/hielkefellinger/go-dnd/app/ecs_components"
	"github.com/hielkefellinger/go-dnd/app/ecs_model_translation"
	"github.com/hielkefellinger/go-dnd/app/models"
	"slices"
	"sort"
)

func (e *baseEventMessageHandler) typeManageItems(message EventMessage, pool CampaignPool) error {
	if message.Source != pool.GetLeadId() {
		return errors.New("managing game Item elements is not allowed as non-lead")
	}

	// Get all Items and parse them to CampaignInventoryItem's
	data := make(map[string]any)
	parsedItems := make([]*models.CampaignInventoryItem, 0)
	allItemEntities := pool.GetEngine().GetWorld().GetItemEntities()
	for _, itemEntity := range allItemEntities {
		parsedItems = append(parsedItems, ecs_model_translation.ItemEntityToCampaignInventoryItem(itemEntity, 0))
	}
	sort.Slice(parsedItems, func(i, j int) bool {
		return parsedItems[i].Name < parsedItems[j].Name
	})
	data["Items"] = parsedItems

	rawJsonBytes, err := json.Marshal(
		e.handleLoadHtmlBody("manageItems.html", "manageItems", data))
	if err != nil {
		return err
	}

	// Send
	manageItems := NewEventMessage()
	manageItems.Source = message.Source
	manageItems.Type = TypeManageItems
	manageItems.Body = string(rawJsonBytes)
	manageItems.Destinations = append(manageItems.Destinations, pool.GetLeadId())
	pool.TransmitEventMessage(manageItems)

	return nil
}

func (e *baseEventMessageHandler) typeManageInventory(message EventMessage, pool CampaignPool) error {
	if message.Source != pool.GetLeadId() {
		return errors.New("managing game Inventory elements is not allowed as non-lead")
	}

	data := make(map[string]any)

	inventories := pool.GetEngine().GetWorld().GetInventoryEntities()
	characters := pool.GetEngine().GetWorld().GetCharacterEntities()
	parsedInventories := make([]models.CampaignInventory, 0)
	pcInventories := make([]models.CampaignInventory, 0)
	for _, inventoryEntity := range inventories {
		inventoryModel := ecs_model_translation.InventoryEntityToCampaignInventoryModel(inventoryEntity)
		isOwnedByPc := false

		// Link to characters
		for _, characterEntity := range characters {
			if characterEntity.HasRelationWithEntityByUuid(inventoryEntity.GetId()) {
				isOwnedByPc = isOwnedByPc || characterEntity.HasComponentType(ecs.PlayerComponentType)
				inventoryModel.Characters = append(inventoryModel.Characters,
					ecs_model_translation.CharacterEntityToCampaignCharacterModel(characterEntity))
			}
		}
		sort.Sort(inventoryModel.Characters)
		if isOwnedByPc {
			pcInventories = append(pcInventories, inventoryModel)
		} else {
			parsedInventories = append(parsedInventories, inventoryModel)
		}
	}
	sort.Slice(parsedInventories, func(i, j int) bool {
		return parsedInventories[i].Name < parsedInventories[j].Name
	})
	sort.Slice(pcInventories, func(i, j int) bool {
		return pcInventories[i].Name < pcInventories[j].Name
	})
	data["Inventories"] = parsedInventories
	data["PcInventories"] = pcInventories

	rawJsonBytes, err := json.Marshal(
		e.handleLoadHtmlBodyMultipleTemplateFiles([]string{"manageInventories.html",
			"manageInventorySelectionBox.html", "diceSpinnerSvg.html"}, "manageInventories", data))
	if err != nil {
		return err
	}

	manageInventories := NewEventMessage()
	manageInventories.Source = message.Source
	manageInventories.Type = TypeManageInventory
	manageInventories.Body = string(rawJsonBytes)
	manageInventories.Destinations = append(manageInventories.Destinations, pool.GetLeadId())
	pool.TransmitEventMessage(manageInventories)

	return nil
}

func (e *baseEventMessageHandler) typeManageCampaign(message EventMessage, pool CampaignPool) error {
	if message.Source != pool.GetLeadId() {
		return errors.New("managing the campaign is not allowed as non-lead")
	}

	data := make(map[string]any)

	// CRUD the basic campaign details
	service := models.CampaignService{}
	campaign, err := service.RetrieveCampaignsById(pool.GetId())
	if err != nil {
		return err
	}
	data["campaign"] = campaign

	// Campaign users
	campaignUsers := make([]string, 0)
	for _, user := range campaign.Users {
		if user.Name != pool.GetLeadId() {
			campaignUsers = append(campaignUsers, user.Name)
		}
	}
	sort.Strings(campaignUsers)
	data["campaignUsers"] = campaignUsers

	// Get all the player chars and see who controls them
	characters := pool.GetEngine().GetWorld().GetPlayerCharacterEntities()
	charControllers := make([]charUserController, 0)
	for _, character := range characters {
		charUserControl := newCharUserController()
		charUserControl.Id = character.GetId().String()
		charUserControl.Name = character.GetName()

		// Get list of controlling users
		listOfControllingUserNames := make([]string, 0)
		playerComponents := character.GetAllComponentsOfType(ecs.PlayerComponentType)
		for index := 0; index < len(playerComponents); index++ {
			playerComponent := playerComponents[index].(*ecs_components.PlayerComponent)
			listOfControllingUserNames = append(listOfControllingUserNames, playerComponent.Name)
		}

		for _, user := range campaign.Users {
			if user.Name != pool.GetLeadId() {
				charUserControl.ControllingPlayers[user.Name] = slices.Contains(listOfControllingUserNames, user.Name)
			}
		}

		var image *ecs_components.ImageComponent
		var imageDetails = character.GetAllComponentsOfType(ecs.ImageComponentType)
		if imageDetails != nil && len(imageDetails) > 0 {
			image = imageDetails[0].(*ecs_components.ImageComponent)
		} else {
			// Set default
			image = ecs_components.NewMissingImageComponent()
		}
		charUserControl.Image = models.CampaignImage{
			Name: image.Name,
			Url:  image.Url,
		}

		charControllers = append(charControllers, charUserControl)
	}
	data["charToPlayers"] = charControllers

	rawJsonBytes, err := json.Marshal(
		e.handleLoadHtmlBodyMultipleTemplateFiles([]string{"manageCampaign.html", "diceSpinnerSvg.html"},
			"manageCampaign", data))
	if err != nil {
		return err
	}

	manageCampaign := NewEventMessage()
	manageCampaign.Source = message.Source
	manageCampaign.Type = TypeManageCampaign
	manageCampaign.Body = string(rawJsonBytes)
	manageCampaign.Destinations = append(manageCampaign.Destinations, pool.GetLeadId())
	pool.TransmitEventMessage(manageCampaign)

	return nil
}

func (e *baseEventMessageHandler) typeManageCharacters(message EventMessage, pool CampaignPool) error {
	if message.Source != pool.GetLeadId() {
		return errors.New("managing game Characters elements is not allowed as non-lead")
	}

	charEntities := pool.GetEngine().GetWorld().GetCharacterEntities()
	allPlayers := pool.GetAllClientIds()

	var playerChars models.Characters
	var nonPlayerChars models.Characters
	for _, charEntity := range charEntities {

		charModel := models.Character{
			Id:          charEntity.GetId().String(),
			Name:        charEntity.GetName(),
			Description: charEntity.GetDescription(),
		}

		// Check if (one of n of controlling) player(s) is online
		isPC := false
		for _, rawPlayerComponent := range charEntity.GetAllComponentsOfType(ecs.PlayerComponentType) {
			isPC = true
			playerComponent := rawPlayerComponent.(*ecs_components.PlayerComponent)
			charModel.Online = charModel.Online || slices.Contains(allPlayers, playerComponent.Name)
		}

		var charDetails = charEntity.GetAllComponentsOfType(ecs.CharacterComponentType)
		if charDetails != nil && len(charDetails) > 0 {
			characterComponent := charDetails[0].(*ecs_components.CharacterComponent)
			charModel.Name = characterComponent.Name
			charModel.Description = characterComponent.Description
		}

		var image *ecs_components.ImageComponent
		var imageDetails = charEntity.GetAllComponentsOfType(ecs.ImageComponentType)
		if imageDetails != nil && len(imageDetails) > 0 {
			image = imageDetails[0].(*ecs_components.ImageComponent)
		} else {
			// Set default
			image = ecs_components.NewMissingImageComponent()
		}

		charModel.Image = models.CampaignImage{
			Name: image.Name,
			Url:  image.Url,
		}

		if isPC {
			playerChars = append(playerChars, charModel)
		} else {
			nonPlayerChars = append(nonPlayerChars, charModel)
		}
	}

	// Sort the list to always show the same order
	sort.Sort(nonPlayerChars)
	sort.Sort(playerChars)

	data := make(map[string]any)
	data["pc_chars"] = playerChars
	data["npc_chars"] = nonPlayerChars

	rawJsonBytes, err := json.Marshal(
		e.handleLoadHtmlBodyMultipleTemplateFiles(
			[]string{"manageCharacters.html", "manageCharacterSelectionBox.html"},
			"manageCharacters", data))
	if err != nil {
		return err
	}

	manageChars := NewEventMessage()
	manageChars.Source = message.Source
	manageChars.Type = TypeManageCharacters
	manageChars.Body = string(rawJsonBytes)
	manageChars.Destinations = append(manageChars.Destinations, pool.GetLeadId())
	pool.TransmitEventMessage(manageChars)

	return nil
}

func (e *baseEventMessageHandler) typeManageMaps(message EventMessage, pool CampaignPool) error {
	if message.Source != pool.GetLeadId() {
		return errors.New("managing game Map elements is not allowed as non-lead")
	}

	// Update possible Map Entities
	mapEntries := make([]models.CampaignMap, 0)
	mapEntities := pool.GetEngine().GetWorld().GetMapEntities()
	sort.Slice(mapEntities, func(i, j int) bool {
		return mapEntities[i].GetName() < mapEntities[j].GetName()
	})
	for _, mapEntity := range mapEntities {
		mapEntries = append(mapEntries, ecs_model_translation.MapEntityToCampaignMapModel(mapEntity))
	}

	data := make(map[string]any)
	data["Maps"] = mapEntries
	rawJsonBytes, err := json.Marshal(
		e.handleLoadHtmlBodyMultipleTemplateFiles(
			[]string{"manageMaps.html", "manageMapSelectionBox.html"},
			"manageMaps", data))
	if err != nil {
		return err
	}

	manageMaps := NewEventMessage()
	manageMaps.Source = message.Source
	manageMaps.Type = TypeManageMaps
	manageMaps.Body = string(rawJsonBytes)
	manageMaps.Destinations = append(manageMaps.Destinations, pool.GetLeadId())
	pool.TransmitEventMessage(manageMaps)

	return nil
}

type charUserController struct {
	Id                 string
	Name               string
	Image              models.CampaignImage
	ControllingPlayers map[string]bool
}

func newCharUserController() charUserController {
	return charUserController{
		ControllingPlayers: make(map[string]bool),
	}
}
