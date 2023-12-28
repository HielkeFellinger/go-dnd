package session

import (
	"encoding/json"
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
		_, rawEvent, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		log.Println(string(rawEvent))

		// Check the rights handle lead actions
		var parsedMessage eventMessage
		err = json.Unmarshal(rawEvent, &parsedMessage)
		if err != nil {
			log.Printf("Message failed with error: %+v\n", err.Error())
			return
		}

		// Make "safe"
		parsedMessage.Body = html.EscapeString(parsedMessage.Body)

		// Update with user credentials and send to channel
		parsedMessage.Source = c.Id
		log.Printf("Message Received: %+v\n", parsedMessage)
		c.Pool.Transmit <- parsedMessage
	}
}
