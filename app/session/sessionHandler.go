package session

import (
	"fmt"
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

	// Upgrade Connection
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		jsonReturn[errMessage], jsonReturn[errTitle] = err.Error(), "Error"
		c.JSON(http.StatusServiceUnavailable, jsonReturn)
	}

	defer func(ws *websocket.Conn) {
		err := ws.Close()
		if err != nil {
			jsonReturn[errMessage], jsonReturn[errTitle] = err.Error(), "Error"
			c.JSON(http.StatusInternalServerError, jsonReturn)
		}
	}(ws)

	writeErr := ws.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Welcome '%s' to campaign: '%s'", user.Name, campaign.Title)))
	if writeErr != nil {
		jsonReturn[errMessage], jsonReturn[errTitle] = writeErr.Error(), "Error"
		c.JSON(http.StatusServiceUnavailable, jsonReturn)
	}

	for {
		// Read message
		_, p, err := ws.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		log.Println(string(p))

		// Reflect message
		if err := ws.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Reflect message: '%s'", p))); err != nil {
			log.Println(err)
			return
		}
	}
}
