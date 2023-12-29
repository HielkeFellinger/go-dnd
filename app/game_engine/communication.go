package game_engine

type EventType int

const (
	TypeGameClose EventType = 0
	TypeGameStart EventType = 1

	TypeUserJoin EventType = 400

	TypeLoadGame        EventType = 500
	TypeLoadCharacters  EventType = 501
	TypeAddCharacter    EventType = 502
	TypeRemoveCharacter EventType = 503

	TypeLoadMap EventType = 600

	TypeChatBroadcast EventType = 800
	TypeChatServerMsg EventType = 801
	TypeChatCommand   EventType = 802
	TypeChatWhisper   EventType = 802
)

type EventMessage struct {
	Source       string    `json:"source"`
	Destinations []string  `json:"-"`
	Type         EventType `json:"type"`
	Body         string    `json:"body"`
	DateTime     string    `json:"dateTime"`
}
