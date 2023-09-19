package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hielkefellinger/go-dnd/app/models"
	"github.com/hielkefellinger/go-dnd/app/session"
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

	// Test if active
	for key, campaign := range userCampaigns {
		userCampaigns[key].Active = session.IsCampaignRunning(campaign.ID)
		userCampaigns[key].UserIsLead = campaign.LeadID == rawUser.(models.User).ID
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
	campaign.Lead = rawUser.(models.User)
	campaign.Private = false
	if err := service.InsertCampaign(&campaign); err != nil {
		templateMap[errMessage], templateMap[errTitle] = err.Error(), "Error"
		c.HTML(http.StatusBadRequest, template, templateMap)
		return
	}

	// Redirect  (After creating a successful campaign)
	c.Redirect(http.StatusFound, fmt.Sprintf("/campaign/session/%d", campaign.ID))
}

func CampaignSessionPage(c *gin.Context) {
	templateMap := gin.H{}

	// Retrieve campaign & user (set by middleware)
	rawUser, exists := c.Get("user")
	if !exists {
		templateMap[errMessage], templateMap[errTitle] = "Failed authenticate user", "Error"
		c.HTML(http.StatusUnauthorized, "campaignSelect.html", templateMap)
	}
	user := rawUser.(models.User)

	rawCampaign, exists := c.Get("campaign")
	if !exists {
		templateMap[errMessage], templateMap[errTitle] = "Failed find campaign", "Error"
		c.HTML(http.StatusNotFound, "campaignSelect.html", templateMap)
	}
	campaign := rawCampaign.(models.Campaign)

	if user.ID == campaign.LeadID {
		campaign.UserIsLead = true
	}

	templateMap["user"] = user
	templateMap["campaign"] = campaign
	templateMap["title"] = fmt.Sprintf("GO-DND Campaign %s", rawCampaign.(models.Campaign).Title)

	// Check if campaign is active (A (ws) Pool must be created by the Lead)
	// - If user != lead && user !in campaign users; fail
	// If user == lead: 	Load extra Campaign data -> Entity component system?? (ECS)

	c.HTML(http.StatusOK, "campaign.html", templateMap)
}
