package session

import (
	"bytes"
	"fmt"
	"github.com/hielkefellinger/go-dnd/app/models"
	"html/template"
	"log"
	"strings"
	"time"
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
				Source: "Server", Type: TypeUserJoin,
				Body: fmt.Sprintf("User '%s' Joins the game", client.Id),
			})

			break
		case client := <-pool.Unregister:
			// Test if user is Lead, if so close the pool
			if client.Lead {
				// Send closing message
				pool.SendMessage(message{Source: "Server", Type: TypeGameClose, Body: "Closing Game!"})

				// Close game; and remove from session container @todo Save state?
				for client := range pool.Clients {
					delete(pool.Clients, client)
				}
				runningCampaignSessionsContainer.Unregister <- pool
				return
			}

			delete(pool.Clients, client)
			log.Printf("Size of Connection Pool `%d`: %d", pool.Id, len(pool.Clients))

			break
		case message := <-pool.Transmit:
			if message.DateTime == "" {
				now := time.Now()
				message.DateTime = fmt.Sprintf("%d-%d-%d %d:%d:%d", now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second())
			}

			// Check Commands like whisper or dice

			// HTML over Websocket Test!
			if strings.Contains(message.Body, "char") {
				message.Type = TypeLoadCharacters
				chars := []models.Character{
					{Name: "Kaas - 1"}, {Name: "Kaas - 2"},
				}
				data := make(map[string]any)
				data["chars"] = chars

				var buf bytes.Buffer
				tmpl := template.Must(template.ParseFiles("web/templates/test.html"))
				err := tmpl.ExecuteTemplate(&buf, "chars", data)
				if err != nil {
					log.Printf("Error parsing test.html `%s`", err.Error())
				}
				message.Body = string(buf.Bytes())
			}

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
