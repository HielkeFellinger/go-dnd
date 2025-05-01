package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hielkefellinger/go-dnd/app/models"
	"github.com/hielkefellinger/go-dnd/app/session"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/html"
	"net/http"
	"slices"
)

func CampaignSelectPage(c *gin.Context) {
	user := getUserFromContextOrRedirectToLoginPage(c)
	if user.ID == 0 {
		return
	}

	campaignSelectPage(c, http.StatusOK, gin.H{}, user)
}

func campaignSelectPage(c *gin.Context, code int, templateMap gin.H, user models.User) {
	templateMap["title"] = "GO-DND Campaign Select"
	templateMap["user"] = user

	// Retrieve campaigns
	var service = models.CampaignService{}
	userCampaigns, err := service.RetrieveCampaignsLinkedToUser(user)
	if err != nil {
		templateMap[errMessage], templateMap[errTitle] = "Could not retrieve campaigns", "Error"
		c.HTML(http.StatusUnauthorized, "campaignSelect.html", templateMap)
		return
	}

	// Test if active sessions linked to user are active
	userCampaignIds := make([]uint, len(userCampaigns))
	for key, campaign := range userCampaigns {
		userCampaignIds[key] = campaign.ID
		userCampaigns[key].Active = session.IsCampaignRunning(campaign.ID)
		userCampaigns[key].UserIsLead = campaign.LeadID == user.ID
	}
	templateMap["userCampaigns"] = userCampaigns

	// Get active non-user aligned campaigns
	allActiveNonUserCampaignIds := session.GetRunningCampaignIds(userCampaignIds)
	otherActiveCampaigns, err := service.RetrieveCampaignsByIds(allActiveNonUserCampaignIds)
	for key := range otherActiveCampaigns {
		otherActiveCampaigns[key].Active = true
	}
	templateMap["otherCampaigns"] = otherActiveCampaigns

	c.HTML(code, "campaignSelect.html", templateMap)
}

func CampaignNewPage(c *gin.Context) {
	templateMap := gin.H{}
	templateMap["title"] = "GO-DND Create a Campaign"

	user := getUserFromContextOrRedirectToLoginPage(c)
	if user.ID == 0 {
		return
	}
	templateMap["user"] = user

	c.HTML(http.StatusOK, "campaignAdd.html", templateMap)
}

func CampaignNew(c *gin.Context) {
	templateMap := gin.H{}
	templateMap["title"] = "GO-DND Create Campaign"
	const template = "campaignAdd.html"

	// Block non post content
	if c.Request.Method != http.MethodPost {
		templateMap[errMessage], templateMap[errTitle] = "Failed to read request", "Error"
		c.HTML(http.StatusBadRequest, template, templateMap)
		return
	}

	user := getUserFromContextOrRedirectToLoginPage(c)
	if user.ID == 0 {
		return
	}
	templateMap["user"] = user

	// Parse body to model
	var campaign models.Campaign
	if c.Bind(&campaign) != nil {
		templateMap[errMessage], templateMap[errTitle] = "Failed to read request", "Error"
		c.HTML(http.StatusBadRequest, template, templateMap)
		return
	}
	campaign.Title = html.EscapeString(campaign.Title)
	campaign.Description = html.EscapeString(campaign.Description)

	// Attempt to insert campaign
	var service = models.CampaignService{}
	campaign.Lead = user
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

	// Retrieve user (set by middleware)
	user := getUserFromContextOrRedirectToLoginPage(c)
	if user.ID == 0 {
		return
	}
	templateMap["user"] = user

	var campaign models.Campaign
	id := c.Params.ByName("id")
	models.DB.Preload("Users").Preload("Lead").First(&campaign, id)
	if campaign.ID == 0 {
		templateMap[errMessage], templateMap[errTitle] = "Failed find campaign", "Error"
		campaignSelectPage(c, http.StatusNotFound, templateMap, user)
		return
	}
	if user.ID == campaign.LeadID {
		campaign.UserIsLead = true
	}
	templateMap["campaign"] = campaign

	// Check if user is linked to campaign, if not redirect
	if !campaign.UserIsLead && !slices.Contains(campaign.Users, user) {
		c.Redirect(http.StatusFound, fmt.Sprintf("/campaign/login/%d", campaign.ID))
		return
	}

	if !campaign.UserIsLead && !session.IsCampaignRunning(campaign.ID) {
		redirectTemplateMap := gin.H{}
		redirectTemplateMap[errMessage], redirectTemplateMap[errTitle] = "Selected campaign is currently inactive!", "Error"
		campaignSelectPage(c, http.StatusBadRequest, redirectTemplateMap, user)
		return
	}

	templateMap["blockMenu"] = true
	templateMap["title"] = fmt.Sprintf("GO-DND Campaign %s", campaign.Title)
	c.HTML(http.StatusOK, "campaign.html", templateMap)
}

func CampaignSessionAuthorize(c *gin.Context) {
	var body struct {
		Key string `form:"key"`
	}

	templateMap := gin.H{}

	// Retrieve campaign & user (set by middleware)
	user := getUserFromContextOrRedirectToLoginPage(c)
	if user.ID == 0 {
		return
	}
	templateMap["user"] = user

	var campaign models.Campaign
	id := c.Params.ByName("id")
	models.DB.Preload("Users").Preload("Lead").First(&campaign, id)
	if campaign.ID == 0 {
		templateMap[errMessage], templateMap[errTitle] = "Failed to find campaign", "Error"
		campaignSelectPage(c, http.StatusNotFound, templateMap, user)
		return
	}
	if user.ID == campaign.LeadID {
		campaign.UserIsLead = true
	}
	templateMap["campaign"] = campaign

	// Check if campaign is currently running
	if !session.IsCampaignRunning(campaign.ID) {
		redirectTemplateMap := gin.H{}
		redirectTemplateMap[errMessage], redirectTemplateMap[errTitle] = "Selected campaign is currently inactive!", "Error"
		campaignSelectPage(c, http.StatusBadRequest, redirectTemplateMap, user)
		return
	}

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
		redirectTemplateMap := gin.H{}
		redirectTemplateMap[errMessage], redirectTemplateMap[errTitle] = "Authorized, but campaign is currently inactive!", "Error"
		campaignSelectPage(c, http.StatusBadRequest, redirectTemplateMap, user)
		return
	}

	// Redirect  (After auth successful authorisation)
	c.Redirect(http.StatusFound, fmt.Sprintf("/campaign/session/%d", campaign.ID))
}

func getUserFromContextOrRedirectToLoginPage(c *gin.Context) models.User {
	rawUser, exists := c.Get("user")
	if !exists {
		c.Redirect(http.StatusFound, "/u/login")
		c.AbortWithStatus(http.StatusUnauthorized)
		return models.User{}
	}
	return rawUser.(models.User)
}
