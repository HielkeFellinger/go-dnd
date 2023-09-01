package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"time"
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
	for i := 0; i < 7; i++ {
		ws.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Session '%s' Message '%d'", id, i)))
		time.Sleep(time.Second)
	}
}
