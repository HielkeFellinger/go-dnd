package game_engine

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLoadingGame(t *testing.T) {
	// Arrange

	// Act

	world := loadGame()

	// Assert
	assert.NotNil(t, world)
}
