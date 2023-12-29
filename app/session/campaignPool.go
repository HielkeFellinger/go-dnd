package session

import (
	"fmt"
	"github.com/hielkefellinger/go-dnd/app/game_engine"
	"log"
	"time"
)

type baseCampaignPool struct {
	Id         uint
	LeadId     string
	Register   chan *campaignClient
	Unregister chan *campaignClient
	Clients    map[*campaignClient]bool
	Receive    chan game_engine.EventMessage
	Engine     game_engine.Engine
}

func (pool *baseCampaignPool) GetId() uint {
	return pool.Id
}

func (pool *baseCampaignPool) GetLeadId() string {
	return pool.LeadId
}

func (pool *baseCampaignPool) GetEngine() game_engine.Engine {
	return pool.Engine
}

func (pool *baseCampaignPool) TransmitEventMessage(message game_engine.EventMessage) {
	log.Printf("Pool internal: %+v\n", message)
	pool.transmitMessage(message)
}

func initCampaignPool(id uint, leadId string) *baseCampaignPool {
	return &baseCampaignPool{
		Id:         id,
		LeadId:     leadId,
		Register:   make(chan *campaignClient),
		Unregister: make(chan *campaignClient),
		Clients:    make(map[*campaignClient]bool),
		Receive:    make(chan game_engine.EventMessage),
	}
}

func (pool *baseCampaignPool) Run() {
	for {
		select {
		case client := <-pool.Register:
			pool.Clients[client] = true
			log.Printf("Size of Connection Pool `%d`: %d", pool.Id, len(pool.Clients))

			pool.transmitMessage(game_engine.EventMessage{
				Source: "Server", Type: game_engine.TypeUserJoin,
				Body: fmt.Sprintf("User '%s' Joins the content", client.Id),
			})

			break
		case client := <-pool.Unregister:
			// Test if user is Lead, if so close the pool
			if client.Lead {
				// Send closing EventMessage
				pool.transmitMessage(game_engine.EventMessage{Source: "Server", Type: game_engine.TypeGameClose, Body: "Closing Game!"})

				// Close content; and remove from session container @todo Save state?
				for client := range pool.Clients {
					delete(pool.Clients, client)
				}
				runningCampaignSessionsContainer.Unregister <- pool
				return
			}

			delete(pool.Clients, client)
			log.Printf("Size of Connection Pool `%d`: %d", pool.Id, len(pool.Clients))

			break
		case eventMessage := <-pool.Receive:
			log.Printf("Message Received Channel: %+v\n", eventMessage)
			if eventMessage.DateTime == "" {
				now := time.Now()
				eventMessage.DateTime = fmt.Sprintf("%d-%d-%d %d:%d:%d", now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second())
			}

			err := pool.Engine.GetEventMessageHandler().HandleEventMessage(eventMessage, pool)
			if err != nil {
				// @todo Handle error
			}
			break
		}
	}
}

func (pool *baseCampaignPool) transmitMessage(message game_engine.EventMessage) {
	log.Printf("Message Transmit Channel: %+v\n", message)
	if message.DateTime == "" {
		now := time.Now()
		message.DateTime = fmt.Sprintf("%d-%d-%d %d:%d:%d", now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second())
	}

	for client := range pool.Clients {
		// Skip EventMessage on clients who are not recipient
		if message.Destinations != nil && len(message.Destinations) > 1 && !contains(message.Destinations, client.Id) {
			continue
		}

		// Send JSON to clients
		err := client.Conn.WriteJSON(message)
		if err != nil {
			// @todo Log failure
		}
	}
	log.Printf("Message(s) Transmitted: %+v\n", message)
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}
