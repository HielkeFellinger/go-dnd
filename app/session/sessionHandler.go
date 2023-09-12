package session

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/hielkefellinger/go-dnd/app/models"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// SEC FAIL/DANGER THIS DOES BYPASS ORIGIN CHECK!!
	CheckOrigin: func(r *http.Request) bool { return true },
}

func ServeSessionWS(c *gin.Context) {
	jsonReturn := gin.H{}

	// Retrieve campaign & user (set by middleware)
	rawUser, exists := c.Get("user")
	if !exists {
		jsonReturn[errMessage], jsonReturn[errTitle] = "Failed authenticate user", "Error"
		c.JSON(http.StatusUnauthorized, jsonReturn)
	}
	user := rawUser.(models.User)

	rawCampaign, exists := c.Get("campaign")
	if !exists {
		jsonReturn[errMessage], jsonReturn[errTitle] = "Failed find campaign", "Error"
		c.JSON(http.StatusNotFound, jsonReturn)
	}
	campaign := rawCampaign.(models.Campaign)

	// Check Pools if active or lead
	campaignRunning := runningCampaignSessionsContainer.IsCampaignRunning(campaign.ID)

	if !campaignRunning && campaign.LeadID != user.ID {
		// No lead and campaign not running
		jsonReturn[errMessage], jsonReturn[errTitle] = "Campaign not running!", "Error"
		c.JSON(http.StatusNotFound, jsonReturn)
	}
	if !campaignRunning && campaign.LeadID == user.ID {
		// Start Campaign! (user is lead of campaign and can init campaign session)
		log.Printf("Starting campaign '%d'", campaign.ID)
		runningCampaignSessionsContainer.initAndRegisterCampaignPool(campaign)
	}

	// Upgrade Connection
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		jsonReturn[errMessage], jsonReturn[errTitle] = err.Error(), "Error"
		c.JSON(http.StatusServiceUnavailable, jsonReturn)
	}

	// Create Client and link to pool
	client := &campaignClient{
		Id:   user.Name,
		Lead: campaign.LeadID == user.ID,
		Conn: ws,
	}
	runningCampaignSessionsContainer.addClientToCampaignPool(campaign.ID, client)
	client.Read()
}
