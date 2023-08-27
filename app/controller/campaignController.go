package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hielkefellinger/go-dnd/app/models"
	"net/http"
)

func CampaignSelectPage(c *gin.Context) {
	templateMap := gin.H{}
	templateMap["title"] = "GO-DND Campaign Select"

	rawUser, exists := c.Get("user")
	if !exists {
		templateMap[errMessage], templateMap[errTitle] = "Failed authenticate user", "Error"
		c.HTML(http.StatusUnauthorized, "campaignSelect.html", templateMap)
	}
	templateMap["user"] = rawUser.(models.User)

	// Retrieve campaigns

	c.HTML(http.StatusOK, "campaignSelect.html", templateMap)
}

func CampaignNewPage(c *gin.Context) {
	templateMap := gin.H{}
	templateMap["title"] = "GO-DND Create a Campaign"

	rawUser, exists := c.Get("user")
	if !exists {
		templateMap[errMessage], templateMap[errTitle] = "Failed authenticate user", "Error"
		c.HTML(http.StatusUnauthorized, "campaign.html", templateMap)
	}
	templateMap["user"] = rawUser.(models.User)

	c.HTML(http.StatusOK, "campaign.html", templateMap)
}

func CampaignNew(c *gin.Context) {
	templateMap := gin.H{}
	templateMap["title"] = "GO-DND Create Campaign"
	const template = "crudCampaign.html"

	// Block non post content
	if c.Request.Method != http.MethodPost {
		templateMap[errMessage], templateMap[errTitle] = "Failed to read request", "Error"
		c.HTML(http.StatusBadRequest, template, templateMap)
		return
	}

	rawUser, exists := c.Get("user")
	if !exists {
		templateMap[errMessage], templateMap[errTitle] = "Failed authenticate user", "Error"
		c.HTML(http.StatusUnauthorized, "campaign.html", templateMap)
	}
	templateMap["user"] = rawUser.(models.User)

	// Parse body to model
	var campaign models.Campaign
	if c.Bind(&campaign) != nil {
		templateMap[errMessage], templateMap[errTitle] = "Failed to read request", "Error"
		c.HTML(http.StatusBadRequest, template, templateMap)
		return
	}

	// Attempt to insert campaign
	var service = models.CampaignService{}
	campaign.LeadId = int(rawUser.(models.User).ID)
	campaign.Private = false
	if err := service.InsertCampaign(&campaign); err != nil {
		templateMap[errMessage], templateMap[errTitle] = err.Error(), "Error"
		c.HTML(http.StatusBadRequest, template, templateMap)
		return
	}

	// Redirect  (After creating a successful campaign)
	c.Redirect(http.StatusCreated, fmt.Sprintf("/campaign/session/%d", campaign.ID))
}
