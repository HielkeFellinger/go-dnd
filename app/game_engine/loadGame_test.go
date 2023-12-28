package game_engine

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLoadingGame(t *testing.T) {
	// Arrange

	// Act
	world := loadGame(SpaceGameTest)

	// Assert
	assert.NotNil(t, world)
}
