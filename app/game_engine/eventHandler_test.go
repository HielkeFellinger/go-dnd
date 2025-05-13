package game_engine

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// Minimal Regression ""Test""
func TestHandleEventMessageLoadFullGame(t *testing.T) {
	// Arrange
	t.Setenv("TEMPLATE_DIR", "../../web/templates/")
	testPool := initTestCampaignPool()
	baseEventMessageHndlr := baseEventMessageHandler{}
	messageStartGame := NewEventMessage()
	messageStartGame.Type = TypeLoadFullGame
	messageStartGame.Source = playerOneId
	messageStartGame.ReloadDateTime()

	// Act
	result := baseEventMessageHndlr.HandleEventMessage(messageStartGame, testPool)

	// Assert
	assert.Nil(t, result) // No Error
	assert.Equal(t, countEventMessageIf(testPool.Messages, func(m *EventMessage) bool {
		return m.Type == TypeLoadCharacters // (501)
	}), 1)
	assert.Equal(t, countEventMessageIf(testPool.Messages, func(m *EventMessage) bool {
		return m.Type == TypeLoadMap // (531)
	}), 1)
	assert.Equal(t, countEventMessageIf(testPool.Messages, func(m *EventMessage) bool {
		return m.Type == TypeLoadMapEntities // (532)
	}), 1)
}

func TestHandleEventMessageEmpty(t *testing.T) {
	// Arrange
	t.Setenv("TEMPLATE_DIR", "../../web/templates/")
	testPool := initTestCampaignPool()
	baseEventMessageHndlr := baseEventMessageHandler{}
	messageStartGame := NewEventMessage()

	// Act
	result := baseEventMessageHndlr.HandleEventMessage(messageStartGame, testPool)

	// Assert
	assert.Error(t, result)
}

func countEventMessageIf(list []EventMessage, fn func(m *EventMessage) bool) int {
	count := 0
	for _, m := range list {
		if fn(&m) {
			count++
		}
	}
	return count
}
