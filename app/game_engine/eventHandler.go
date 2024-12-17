package game_engine

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"text/template"
)

type EventMessageHandler interface {
	HandleEventMessage(message EventMessage, pool CampaignPool) error
}

type baseEventMessageHandler struct {
}

func (e *baseEventMessageHandler) HandleEventMessage(message EventMessage, pool CampaignPool) error {
	log.Printf("Message Handler Parsing ID: '%+v' of Type: '%d' \n", message.Id, message.Type)

	if message.Type == TypeGameSave {
		log.Printf("- Save Game Type: '%d' Message: '%s'", message.Type, message.Id)
		return e.handlePersistDataEvents(message, pool)
	}

	if message.Type == TypeLoadFullGame {
		log.Printf("- Load Full Game Type: '%d' Message: '%s'", message.Type, message.Id)
		if err := e.loadCharacters(message, pool); err != nil {
			return err
		}
		if err := e.typeLoadMap(message, pool); err != nil {
			return err
		}
		if err := e.typeLoadMapEntities(message, pool); err != nil {
			return err
		}
		return nil
	}

	// Player Character Load Events
	if message.Type >= TypeLoadCharacters && message.Type <= TypeLoadCharactersDetails {
		log.Printf("- Char. Load Event Type: '%d' Message: '%s'", message.Type, message.Id)

		if message.Type == TypeLoadCharacters {
			return e.loadCharacters(message, pool)
		} else if message.Type == TypeLoadCharactersDetails {
			return e.loadCharactersDetails(message, pool)
		}
		return errors.New(fmt.Sprintf("message of type '%d' is not recognised by 'handleLoadCharacterEvents()'", message.Type))
	}

	// Player Character Update Event(s)
	if message.Type >= TypeUpdateCharacterHealth && message.Type <= TypeUpdateCharacterUsers {
		log.Printf("- Char. Update Event Type: '%d' Message: '%s'", message.Type, message.Id)

		if message.Type == TypeUpdateCharacterHealth {
			return e.typeUpdateCharacterHealth(message, pool)
		} else if message.Type == TypeUpdateCharacterUsers {
			return e.typeUpdateCharacterUsers(message, pool)
		}
		return errors.New(fmt.Sprintf("message of type '%d' is not recognised as a 'Player Character Update Event'", message.Type))
	}

	// Player Item Details Load Event(s)
	if message.Type == TypeLoadItemDetails {
		log.Printf("- Item. Event Type: '%d' Message: '%s'", message.Type, message.Id)
		return e.typeLoadItemDetails(message, pool)
	}

	// Player Map Load Event(s)
	if message.Type >= TypeLoadMap && message.Type <= TypeLoadMapEntity {
		log.Printf("- Map. Event Type: '%d' Message: '%s'", message.Type, message.Id)

		if message.Type == TypeLoadMap {
			return e.typeLoadMap(message, pool)
		} else if message.Type == TypeLoadMapEntities {
			return e.typeLoadMapEntities(message, pool)
		} else if message.Type == TypeLoadMapEntity {
			return e.typeLoadMapEntity(message, pool)
		}
		return errors.New(fmt.Sprintf("message of type '%d' is not recognised as a 'Player Map Load Event'", message.Type))
	}

	// Player Map Update/Interaction Event(s)
	if message.Type >= TypeUpdateMapEntity && message.Type <= TypeChangeMapBackgroundImage {
		log.Printf("- Map Update Event Type: '%d' Message: '%s'", message.Type, message.Id)

		if message.Type == TypeUpdateMapEntity {
			return e.typeUpdateMapEntity(message, pool)
		} else if message.Type == TypeUpdateMapVisibility {
			return e.typeUpdateMapVisibility(message, pool)
		} else if message.Type == TypeAddMapItem {
			return e.typeAddMapItem(message, pool)
		} else if message.Type == TypeRemoveMapItem {
			return e.typeRemoveMapItem(message, pool)
		} else if message.Type == TypeSignalMapItem {
			return e.typeSignalMapItem(message, pool)
		} else if message.Type == TypeChangeMapBackgroundImage {
			return e.typeChangeMapBackgroundImage(message, pool)
		}
		return errors.New(fmt.Sprintf("message of type '%d' is not recognised as a 'Player Map Update/Interaction Event'", message.Type))
	}

	// Management Overview Event(s)
	if message.Type >= TypeManagementOverviewStart && message.Type <= TypeManagementOverviewEnd {
		log.Printf("- Game Management Events Type: '%d' Message: '%s'", message.Type, message.Id)

		if message.Type == TypeManageMaps {
			return e.typeManageMaps(message, pool)
		} else if message.Type == TypeManageCharacters {
			return e.typeManageCharacters(message, pool)
		} else if message.Type == TypeManageInventory {
			return e.typeManageInventory(message, pool)
		} else if message.Type == TypeManageItems {
			return e.typeManageItems(message, pool)
		} else if message.Type == TypeManageCampaign {
			return e.typeManageCampaign(message, pool)
		}
		return errors.New(fmt.Sprintf("message of type '%d' is not recognised as a 'Management Overview Event'", message.Type))
	}
	// Management CRUD Event(s)
	if message.Type > TypeManagementCrudStart && message.Type < TypeManagementCrudEnd {
		log.Printf("- Game Management CRUD Events Type: '%d' Message: '%s'", message.Type, message.Id)

		if message.Type == TypeLoadUpsertMap { // Maps
			return e.typeLoadUpsertMap(message, pool)
		} else if message.Type == TypeUpsertMap {
			return e.typeUpsertMap(message, pool)
		} else if message.Type == TypeLoadUpsertItem { // Items
			return e.typeLoadUpsertItem(message, pool)
		} else if message.Type == TypeUpsertItem {
			return e.typeUpsertItem(message, pool)
		} else if message.Type == TypeLoadUpsertCharacter { // Characters
			return e.typeLoadUpsertCharacter(message, pool)
		} else if message.Type == TypeUpsertCharacter {
			return e.typeUpsertCharacter(message, pool)
		} else if message.Type == TypeLoadUpsertInventory { // Inventories
			return e.typeLoadUpsertInventory(message, pool)
		} else if message.Type == TypeUpsertInventory {
			return e.typeUpsertInventory(message, pool)
		} else if message.Type == TypeAddItemToInventory {
			return e.typeAddItemToInventory(message, pool)
		} else if message.Type == TypeRemoveItemFromInventory {
			return e.typeRemoveItemFromInventory(message, pool)
		}
		return errors.New(fmt.Sprintf("message of type '%d' is not recognised as a 'Management CRUD Event'", message.Type))
	}

	// Chat
	if message.Type >= TypeChatBroadcast && message.Type <= TypeChatWhisper {
		return e.handleChatEventMessage(message, pool)
	}

	return errors.New(fmt.Sprintf("message of type '%d' is not recognised by server", message.Type))
}

func (e *baseEventMessageHandler) handleLoadHtmlBodyMultipleTemplateFiles(fileNames []string, templateName string, data map[string]any) string {
	files := make([]string, 0)
	for _, fileName := range fileNames {
		files = append(files, fmt.Sprintf(os.Getenv("TEMPLATE_DIR")+"%s", fileName))
	}

	var buf bytes.Buffer
	tmpl := template.Must(template.ParseFiles(files...))
	err := tmpl.ExecuteTemplate(&buf, templateName, data)
	if err != nil {
		log.Printf("Error parsing %v `%s`", fileNames, err.Error())
	}
	return string(buf.Bytes())
}

func (e *baseEventMessageHandler) handleLoadHtmlBody(fileName string, templateName string, data map[string]any) string {
	var buf bytes.Buffer
	tmpl := template.Must(template.ParseFiles(fmt.Sprintf(os.Getenv("TEMPLATE_DIR")+"%s", fileName)))
	err := tmpl.ExecuteTemplate(&buf, templateName, data)
	if err != nil {
		log.Printf("Error parsing %s `%s`", fileName, err.Error())
	}
	return string(buf.Bytes())
}
