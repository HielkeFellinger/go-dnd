package session

import (
	"fmt"
	"github.com/gorilla/websocket"
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
		_, p, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		log.Println(string(p))

		message := message{Source: c.Id, Type: MSG_TYPE_CHAT_BROADCAST, Body: string(p)}
		c.Pool.Transmit <- message
		fmt.Printf("Message Received: %+v\n", message)
	}
}
