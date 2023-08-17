package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/hielkefellinger/go-dnd/app/models"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func createGinTestContext(w *httptest.ResponseRecorder) *gin.Context {
	gin.SetMode(gin.TestMode)

	ctx, engine := gin.CreateTestContext(w)
	engine.LoadHTMLGlob("../../web/templates/*")
	engine.Static("/assets", "../../web/assets")

	ctx.Request = &http.Request{Header: make(http.Header)}

	return ctx
}

func TestHomePageWithoutUser(t *testing.T) {
	// Arrange
	w := httptest.NewRecorder()
	ctx := createGinTestContext(w)

	// Act
	HomePage(ctx)

	// Assert
	assert.EqualValues(t, http.StatusOK, w.Code)
	bodyString := w.Body.String()
	assert.Equal(t, strings.Contains(bodyString, "<h1> Test the world </h1>"), true)
	assert.Equal(t, strings.Contains(bodyString, "/u/login"), true)
	assert.Equal(t, strings.Contains(bodyString, "/u/register"), true)
}

func TestHomePageWithUser(t *testing.T) {
	// Arrange
	w := httptest.NewRecorder()
	ctx := createGinTestContext(w)

	// - Add a user
	ctx.Set("user", models.User{Name: "TEST_USER"})

	// Act
	HomePage(ctx)

	// Assert
	assert.EqualValues(t, http.StatusOK, w.Code)
	bodyString := w.Body.String()
	assert.Equal(t, strings.Contains(bodyString, "<h1> Test the world </h1>"), true)
	assert.Equal(t, strings.Contains(bodyString, "/u/logout"), true)
	assert.Equal(t, strings.Contains(bodyString, "/campaign/select"), true)
}
