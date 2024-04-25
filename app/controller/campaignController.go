package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hielkefellinger/go-dnd/app/models"
	"github.com/hielkefellinger/go-dnd/app/session"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"slices"
)

const FailedAuthMessage = "Failed authenticate user"
const CampaignSelectHtmlFile = "campaignSelect.html"

func CampaignSelectPage(c *gin.Context) {
	templateMap := gin.H{}
	templateMap["title"] = "GO-DND Campaign Select"

	rawUser, exists := c.Get("user")
	if !exists {
		templateMap[errMessage], templateMap[errTitle] = FailedAuthMessage, "Error"
		c.HTML(http.StatusUnauthorized, CampaignSelectHtmlFile, templateMap)
	}
	templateMap["user"] = rawUser.(models.User)

	// Retrieve campaigns
	var service = models.CampaignService{}
	userCampaigns, err := service.RetrieveCampaignsLinkedToUser(rawUser.(models.User))
	if err != nil {
		templateMap[errMessage], templateMap[errTitle] = err.Error(), "Error"
		c.HTML(http.StatusUnauthorized, CampaignSelectHtmlFile, templateMap)
	}

	// Test if active sessions linked to user are active
	userCampaignIds := make([]uint, len(userCampaigns))
	for key, campaign := range userCampaigns {
		userCampaignIds[key] = campaign.ID
		userCampaigns[key].Active = session.IsCampaignRunning(campaign.ID)
		userCampaigns[key].UserIsLead = campaign.LeadID == rawUser.(models.User).ID
	}
	templateMap["userCampaigns"] = userCampaigns

	// Get active non-user aligned campaigns
	allActiveNonUserCampaignIds := session.GetRunningCampaignIds(userCampaignIds)
	otherActiveCampaigns, err := service.RetrieveCampaignsByIds(allActiveNonUserCampaignIds)
	for key := range otherActiveCampaigns {
		otherActiveCampaigns[key].Active = true
	}
	templateMap["otherCampaigns"] = otherActiveCampaigns

	c.HTML(http.StatusOK, CampaignSelectHtmlFile, templateMap)
}

func CampaignNewPage(c *gin.Context) {
	templateMap := gin.H{}
	templateMap["title"] = "GO-DND Create a Campaign"

	rawUser, exists := c.Get("user")
	if !exists {
		templateMap[errMessage], templateMap[errTitle] = FailedAuthMessage, "Error"
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
		templateMap[errMessage], templateMap[errTitle] = FailedAuthMessage, "Error"
		c.HTML(http.StatusUnauthorized, CampaignSelectHtmlFile, templateMap)
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
		templateMap[errMessage], templateMap[errTitle] = FailedAuthMessage, "Error"
		c.HTML(http.StatusUnauthorized, CampaignSelectHtmlFile, templateMap)
		return
	}
	user := rawUser.(models.User)
	templateMap["user"] = user

	rawCampaign, exists := c.Get("campaign")
	if !exists {
		templateMap[errMessage], templateMap[errTitle] = "Failed find campaign", "Error"
		c.HTML(http.StatusNotFound, CampaignSelectHtmlFile, templateMap)
		return
	}
	campaign := rawCampaign.(models.Campaign)
	if user.ID == campaign.LeadID {
		campaign.UserIsLead = true
	}
	templateMap["campaign"] = campaign

	// Check if user is linked to campaign
	if !campaign.UserIsLead && !slices.Contains(campaign.Users, user) {
		templateMap["ID"] = campaign.ID
		templateMap["campaignTitle"] = campaign.Title
		templateMap["title"] = fmt.Sprintf("GO-DND Campaign %s", campaign.Title)
		c.HTML(http.StatusUnauthorized, "campaignLogin.html", templateMap)
		return
	}

	templateMap["blockMenu"] = true
	templateMap["title"] = fmt.Sprintf("GO-DND Campaign %s", rawCampaign.(models.Campaign).Title)

	if !campaign.UserIsLead && !session.IsCampaignRunning(campaign.ID) {
		c.Redirect(http.StatusFound, "/campaign/select")
		return
	}

	templateMap["blockMenu"] = true
	templateMap["title"] = fmt.Sprintf("GO-DND Campaign %s", rawCampaign.(models.Campaign).Title)
	c.HTML(http.StatusOK, "campaign.html", templateMap)
}

func CampaignSessionAuthorize(c *gin.Context) {
	var body struct {
		Key string `form:"key"`
	}

	templateMap := gin.H{}

	// Retrieve campaign & user (set by middleware)
	rawUser, exists := c.Get("user")
	if !exists {
		templateMap[errMessage], templateMap[errTitle] = FailedAuthMessage, "Error"
		c.HTML(http.StatusUnauthorized, CampaignSelectHtmlFile, templateMap)
		return
	}
	user := rawUser.(models.User)
	templateMap["user"] = user

	rawCampaign, exists := c.Get("campaign")
	if !exists {
		templateMap[errMessage], templateMap[errTitle] = "Failed find campaign", "Error"
		c.HTML(http.StatusNotFound, CampaignSelectHtmlFile, templateMap)
		return
	}
	campaign := rawCampaign.(models.Campaign)
	if user.ID == campaign.LeadID {
		campaign.UserIsLead = true
	}
	templateMap["campaign"] = campaign

	// Only check if user is not already linked
	failure := false
	if !campaign.UserIsLead && !slices.Contains(campaign.Users, user) {
		// Check auth.
		if c.Bind(&body) != nil {
			templateMap[errMessage], templateMap[errTitle] = "Failed to read request", "Error"
			failure = true
		}

		if body.Key != "" {
			if errBcrypt := bcrypt.CompareHashAndPassword([]byte(campaign.Password), []byte(body.Key)); errBcrypt != nil {
				templateMap[errMessage], templateMap[errTitle] = "Invalid Campaign Key!", "Error"
				failure = true
			}
		} else {
			failure = true
		}

		// Add user on missing failure
		if !failure {
			var service = models.CampaignService{}
			if err := service.AddUserToCampaign(user, campaign); err != nil {
				templateMap[errMessage], templateMap[errTitle] = err.Error(), "Error"
				failure = true
			}
		}

		if failure {
			templateMap["ID"] = campaign.ID
			templateMap["campaignTitle"] = campaign.Title
			templateMap["title"] = fmt.Sprintf("GO-DND Campaign %s", campaign.Title)
			c.HTML(http.StatusBadRequest, "campaignLogin.html", templateMap)
			return
		}
	}

	if !campaign.UserIsLead && !session.IsCampaignRunning(campaign.ID) {
		c.Redirect(http.StatusFound, "/campaign/select")
		return
	}

	templateMap["blockMenu"] = true
	templateMap["title"] = fmt.Sprintf("GO-DND Campaign %s", rawCampaign.(models.Campaign).Title)
	c.HTML(http.StatusOK, "campaign.html", templateMap)
}
