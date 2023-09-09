package session

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
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

	// Retrieve campaign ID
	id := c.Params.ByName("id")

	// Upgrade Connection
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		jsonReturn[errMessage], jsonReturn[errTitle] = err.Error(), "Error"
		c.JSON(http.StatusUnauthorized, jsonReturn)
	}
	defer ws.Close()

	ws.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Welcome to campaign: '%s'", id)))

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
