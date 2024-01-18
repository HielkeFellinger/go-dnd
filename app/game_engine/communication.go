package game_engine

import (
	"github.com/google/uuid"
	"time"
)

type EventType int

const (
	TypeGameClose EventType = 0
	TypeGameStart EventType = 1

	TypeUserJoin EventType = 400

	TypeLoadFullGame EventType = 500

	TypeLoadCharacters  EventType = 501
	TypeAddCharacter    EventType = 502
	TypeRemoveCharacter EventType = 503

	TypeLoadMap         EventType = 531
	TypeLoadMapEntities EventType = 532
	TypeLoadMapEntity   EventType = 533
	TypeAddMap          EventType = 534
	TypeRemoveMap       EventType = 535

	// Set attention to a specific tab

	TypeChatBroadcast EventType = 800
	TypeChatServerMsg EventType = 801
	TypeChatCommand   EventType = 802
	TypeChatWhisper   EventType = 802
)

type EventMessage struct {
	Id           uuid.UUID `json:"-"`
	Source       string    `json:"source"`
	Destinations []string  `json:"-"`
	Type         EventType `json:"type"`
	Body         string    `json:"body"`
	DateTime     string    `json:"dateTime"`
}

func NewEventMessage() EventMessage {
	m := EventMessage{Id: uuid.New()}
	m.ReloadDateTime()
	return m
}

func (m *EventMessage) ReloadDateTime() {
	now := time.Now()
	m.DateTime = now.Format("2006-01-02 15:04:05")
}
