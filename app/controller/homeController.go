package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/hielkefellinger/go-dnd/app/models"
	"net/http"
)

func HomePage(c *gin.Context) {
	templateMap := gin.H{}
	templateMap["title"] = "GO-DND Home Page"

	rawUser, exists := c.Get("user")
	if exists {
		templateMap["user"] = rawUser.(models.User)
	}

	c.HTML(
		http.StatusOK,
		"index.html",
		templateMap,
	)
}
