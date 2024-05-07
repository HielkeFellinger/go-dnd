package session

import (
	"fmt"
	"github.com/hielkefellinger/go-dnd/app/game_engine"
	"log"
	"slices"
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
	log.Printf("Pool transmission Message: '%+v' Source: '%s' & Destination: '%v'", message.Id, message.Source, message.Destinations)
	pool.transmitMessage(message)
}

func (pool *baseCampaignPool) GetAllClientIds(filterOut ...string) []string {
	userIds := make([]string, 0)
	for client := range pool.Clients {

		// Filter out if applicable
		if len(filterOut) > 0 && slices.Contains(filterOut, client.Id) {
			continue
		}
		userIds = append(userIds, client.Id)
	}
	return userIds
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

			var transmitMessage = game_engine.NewEventMessage()
			transmitMessage.Type = game_engine.TypeUserJoin
			transmitMessage.Source = game_engine.ServerUser
			transmitMessage.Body = fmt.Sprintf("User '%s' Joins the content", client.Id)
			pool.transmitMessage(transmitMessage)

			pool.updateCharacterRibbon(client.Id)
			break
		case client := <-pool.Unregister:
			// Test if user is Lead, if so close the pool
			if client.Lead {
				// Send closing EventMessage
				var transmitMessage = game_engine.NewEventMessage()
				transmitMessage.Type = game_engine.TypeGameClose
				transmitMessage.Source = game_engine.ServerUser
				transmitMessage.Body = "Closing Game!"
				pool.transmitMessage(transmitMessage)

				// Close content; and remove from session container @todo Save state?
				for client := range pool.Clients {
					delete(pool.Clients, client)
				}
				runningCampaignSessionsContainer.Unregister <- pool
				return
			} else {
				pool.Clients[client] = false
				delete(pool.Clients, client)
			}

			log.Printf("Size of Connection Pool `%d`: %d", pool.Id, len(pool.Clients))
			pool.updateCharacterRibbon(client.Id)

			break
		case eventMessage := <-pool.Receive:
			log.Printf("Received Channel Message ID: '%s'", eventMessage.Id)

			err := pool.Engine.GetEventMessageHandler().HandleEventMessage(eventMessage, pool)
			if err != nil {
				log.Printf("Message failed with error: %+v\n", err.Error())
			}
			break
		}
	}
}

// Update Character Ribbon of other players
func (pool *baseCampaignPool) updateCharacterRibbon(clientId string) {
	var reloadCharacterMessage = game_engine.NewEventMessage()
	reloadCharacterMessage.Type = game_engine.TypeLoadCharacters
	reloadCharacterMessage.Source = "server"
	reloadCharacterMessage.Destinations = pool.GetAllClientIds(clientId)
	err := pool.Engine.GetEventMessageHandler().HandleEventMessage(reloadCharacterMessage, pool)
	if err != nil {
		log.Printf("Message failed with error: %+v\n", err.Error())
	}
}

func (pool *baseCampaignPool) transmitMessage(message game_engine.EventMessage) {
	for client := range pool.Clients {
		// Skip EventMessage on clients who are not recipient; empty Destinations is message to all
		if message.Destinations != nil && len(message.Destinations) > 0 && !slices.Contains(message.Destinations, client.Id) {
			continue
		}

		// Send JSON to clients
		err := client.Conn.WriteJSON(message)
		if err != nil {
			log.Printf("Error while sending Message ID: '%s'. Error: '%s'", message.Id, err.Error())
		}
	}
	log.Printf("Message Transmitted. ID: '%s' and Type: '%v'", message.Id, message.Type)
}
