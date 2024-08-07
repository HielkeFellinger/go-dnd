package game_engine

import (
	"bytes"
	"fmt"
	"log"
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
		err := e.handlePersistDataEvents(message, pool)
		if err != nil {
			return err
		}
	}

	if message.Type == TypeLoadFullGame || (message.Type >= TypeLoadCharacters && message.Type <= TypeLoadCharactersDetails) {
		err := e.handleLoadCharacterEvents(message, pool)
		if err != nil {
			return err
		}
	}

	if message.Type == TypeLoadFullGame || (message.Type >= TypeLoadMap && message.Type <= TypeUpsertMap) {
		err := e.handleMapLoadEvents(message, pool)
		if err != nil {
			return err
		}
	}

	if message.Type >= TypeUpdateCharacterHealth && message.Type <= TypeUpdateCharacterUsers {
		err := e.handleUpdateCharacterEvents(message, pool)
		if err != nil {
			return err
		}
	}

	if message.Type >= TypeUpdateMapEntity && message.Type <= TypeChangeMapBackgroundImage {
		err := e.handleMapUpdateEvents(message, pool)
		if err != nil {
			return err
		}
	}

	if message.Type >= TypeManageMaps && message.Type <= TypeManageCampaign {
		err := e.handleManagementEvents(message, pool)
		if err != nil {
			return err
		}
	}

	if message.Type >= TypeLoadUpsertItem && message.Type <= TypeUpsertItem {
		err := e.handleItemEvents(message, pool)
		if err != nil {
			return err
		}
	}

	if message.Type >= TypeChatBroadcast && message.Type <= TypeChatWhisper {
		err := e.handleChatEventMessage(message, pool)
		if err != nil {
			return err
		}
	}

	return nil
}

func (e *baseEventMessageHandler) handleLoadHtmlBodyMultipleTemplateFiles(fileNames []string, templateName string, data map[string]any) string {
	files := make([]string, 0)
	for _, fileName := range fileNames {
		files = append(files, fmt.Sprintf("web/templates/%s", fileName))
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
	tmpl := template.Must(template.ParseFiles(fmt.Sprintf("web/templates/%s", fileName)))
	err := tmpl.ExecuteTemplate(&buf, templateName, data)
	if err != nil {
		log.Printf("Error parsing %s `%s`", fileName, err.Error())
	}
	return string(buf.Bytes())
}
