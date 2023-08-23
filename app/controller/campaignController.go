package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hielkefellinger/go-dnd/app/helpers"
	"github.com/hielkefellinger/go-dnd/app/initializers"
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

	var body struct {
		Title         string `form:"title"`
		Description   string `form:"description"`
		Password      string `form:"password"`
		PasswordCheck string `form:"passwordCheck"`
	}

	if c.Bind(&body) != nil {
		templateMap[errMessage], templateMap[errTitle] = "Failed to read request", "Error"
		c.HTML(http.StatusBadRequest, template, templateMap)
		return
	}

	// Check password
	if body.PasswordCheck != body.Password {
		templateMap[errMessage], templateMap[errTitle] = "Passwords do not match", "Error"
		c.HTML(http.StatusBadRequest, template, templateMap)
		return
	}

	hashByteArray, err := helpers.HashPassword(body.Password)
	if err != nil {
		templateMap[errMessage], templateMap[errTitle] = "Password could not be hashed", "Error"
		c.HTML(http.StatusBadRequest, template, templateMap)
		return
	}

	campaign := models.Campaign{
		Private:     false,
		Title:       body.Title,
		Description: body.Description,
		Password:    string(hashByteArray),
		LeadId:      int(rawUser.(models.User).ID),
	}
	result := initializers.DB.Create(&campaign)
	if result.Error != nil {
		templateMap[errMessage], templateMap[errTitle] = "User could not created", "Error"
		c.HTML(http.StatusBadRequest, template, templateMap)
		return
	}

	// Redirect
	c.Redirect(http.StatusCreated, fmt.Sprintf("/campaign/session/%d", campaign.ID))
}
