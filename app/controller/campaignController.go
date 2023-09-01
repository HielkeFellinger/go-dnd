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
	var service = models.CampaignService{}
	userCampaigns, err := service.RetrieveCampaignsLinkedToUser(rawUser.(models.User))
	if err != nil {
		templateMap[errMessage], templateMap[errTitle] = err.Error(), "Error"
		c.HTML(http.StatusUnauthorized, "campaignSelect.html", templateMap)
	}
	templateMap["userCampaigns"] = userCampaigns

	c.HTML(http.StatusOK, "campaignSelect.html", templateMap)
}

func CampaignNewPage(c *gin.Context) {
	templateMap := gin.H{}
	templateMap["title"] = "GO-DND Create a Campaign"

	rawUser, exists := c.Get("user")
	if !exists {
		templateMap[errMessage], templateMap[errTitle] = "Failed authenticate user", "Error"
		c.HTML(http.StatusUnauthorized, "index.html", templateMap)
	}
	templateMap["user"] = rawUser.(models.User)

	c.HTML(http.StatusOK, "campaignCrud.html", templateMap)
}

func CampaignNew(c *gin.Context) {
	templateMap := gin.H{}
	templateMap["title"] = "GO-DND Create Campaign"
	const template = "campaignCrud.html"

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

func CampaignSessionPage(c *gin.Context) {
	templateMap := gin.H{}

	// Retrieve campaign ID
	id := c.Params.ByName("id")

	rawUser, exists := c.Get("user")
	if !exists {
		templateMap[errMessage], templateMap[errTitle] = "Failed authenticate user", "Error"
		c.HTML(http.StatusUnauthorized, "campaignSelect.html", templateMap)
	}
	templateMap["user"] = rawUser.(models.User)
	templateMap["campaignId"] = id
	templateMap["title"] = fmt.Sprintf("GO-DND Campaign %s", id)

	// Load Campaign

	c.HTML(http.StatusOK, "campaign.html", templateMap)
}
