package session

import (
	"fmt"
	"log"
)

type campaignPool struct {
	Id         uint
	LeadId     string
	Register   chan *campaignClient
	Unregister chan *campaignClient
	Clients    map[*campaignClient]bool
	Transmit   chan message
}

func initCampaignPool(id uint, leadId string) *campaignPool {
	return &campaignPool{
		Id:         id,
		LeadId:     leadId,
		Register:   make(chan *campaignClient),
		Unregister: make(chan *campaignClient),
		Clients:    make(map[*campaignClient]bool),
		Transmit:   make(chan message),
	}
}

func (pool *campaignPool) Run() {
	for {
		select {
		case client := <-pool.Register:
			pool.Clients[client] = true
			log.Printf("Size of Connection Pool `%d`: %d", pool.Id, len(pool.Clients))

			pool.SendMessage(message{
				Source: "Server", Type: MSG_TYPE_USER_JOIN,
				Body: fmt.Sprintf("User '%s' Joins the game", client.Id),
			})

			break
		case client := <-pool.Unregister:
			// Test if user is Lead, if so close the pool
			if client.Lead {
				// Close game
				pool.SendMessage(message{Source: "Server", Type: MSG_TYPE_GAME_CLOSE, Body: "Closing Game!"})
				return
			}

			delete(pool.Clients, client)
			log.Printf("Size of Connection Pool `%d`: %d", pool.Id, len(pool.Clients))

			break
		case message := <-pool.Transmit:

			// Check Commands like whisper or dice

			// Pass-trough

			// Just broadcast for now
			pool.SendMessage(message)

			break
		}
	}
}

func (pool *campaignPool) SendMessage(message message) {
	for client := range pool.Clients {
		// Skip message on clients who are not recipient
		if message.Destinations != nil && len(message.Destinations) > 1 && !contains(message.Destinations, client.Id) {
			continue
		}

		// Send JSON to clients
		err := client.Conn.WriteJSON(message)
		if err != nil {
			// Log failure
		}
	}
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}
