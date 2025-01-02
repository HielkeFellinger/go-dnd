package session

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/hielkefellinger/go-dnd/app/game_engine"
	"html"
	"log"
)

type campaignClient struct {
	Id   string
	Lead bool
	Conn *websocket.Conn
	Pool *baseCampaignPool
}

func (c *campaignClient) GetId() string {
	return c.Id
}

func (c *campaignClient) IsLead() bool {
	return c.Lead
}

func (c *campaignClient) Read() {
	defer func() {
		// Check if pool has been initialised
		if c.Pool != nil {
			c.Pool.Unregister <- c
		}
		c.Pool = nil
		err := c.Conn.Close()
		if err != nil {
			log.Println(err)
			return
		}
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
		var rawMessage game_engine.EventMessage
		err = json.Unmarshal(rawEvent, &rawMessage)
		if err != nil {
			log.Printf("Message failed with error: %+v\n", err.Error())
			return
		}

		// Make "safe"
		parsedMessage := game_engine.NewEventMessage()
		parsedMessage.Body = html.EscapeString(rawMessage.Body)
		parsedMessage.Source = c.Id
		parsedMessage.Type = rawMessage.Type
		parsedMessage.Destinations = rawMessage.Destinations

		// Update with user credentials and send to channel
		if c.Pool != nil {
			c.Pool.Receive <- parsedMessage
		}
	}
}
