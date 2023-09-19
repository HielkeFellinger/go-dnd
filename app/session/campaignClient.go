package session

import (
	"github.com/gorilla/websocket"
	"html"
	"log"
)

type campaignClient struct {
	Id   string
	Lead bool
	Conn *websocket.Conn
	Pool *campaignPool
}

func (c *campaignClient) Read() {
	defer func() {
		c.Pool.Unregister <- c
		c.Pool = nil
		c.Conn.Close()
	}()

	for {
		// user input unsafe!
		_, p, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		pSafe := html.EscapeString(string(p))

		log.Println(pSafe)

		// Do some other actions based on client input

		// Check the rights handle lead actions
		if c.Lead {

		}

		message := message{Source: c.Id, Type: TypeChatBroadcast, Body: pSafe}
		c.Pool.Transmit <- message
		log.Printf("Message Received: %+v\n", message)
	}
}
