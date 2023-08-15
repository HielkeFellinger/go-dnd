package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/hielkefellinger/go-dnd/app/models"
	"net/http"
)

func CampaignSelectPage(c *gin.Context) {
	templateMap := gin.H{}
	templateMap["title"] = "GO-DND Campaign Select"

	rawUser, exists := c.Get("user")
	if exists {
		templateMap["user"] = rawUser.(models.User)
	}

	c.HTML(
		http.StatusOK,
		"campaignSelect.html",
		templateMap,
	)
}
