package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func CampaignSelectPage(c *gin.Context) {
	c.HTML(
		http.StatusOK,
		"campaignSelect.html",
		gin.H{"title": "GO-DND Campaign Select"},
	)
}
