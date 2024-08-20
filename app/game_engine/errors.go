package game_engine

import (
	"encoding/json"
	"errors"
	"github.com/hielkefellinger/go-dnd/app/models"
)

func SendManagementError(title string, message string, pool CampaignPool) error {
	body, err := json.Marshal(models.ManagementError{
		Title:   title,
		Message: message,
	})
	if err != nil {
		return err
	}

	errorMessage := NewEventMessage()
	errorMessage.Source = pool.GetLeadId()
	errorMessage.Type = TypeManagementError
	errorMessage.Body = string(body)
	errorMessage.Destinations = append(errorMessage.Destinations, pool.GetLeadId())

	pool.TransmitEventMessage(errorMessage)
	return errors.New(message)
}
