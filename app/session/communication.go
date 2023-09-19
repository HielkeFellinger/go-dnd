package session

type event int

const (
	TypeGameClose event = 0
	TypeGameStart event = 1

	TypeUserJoin event = 400

	TypeLoadCharacters  event = 500
	TypeAddCharacter    event = 501
	TypeRemoveCharacter event = 502

	TypeLoadMap event = 600

	TypeChatBroadcast event = 800
	TypeChatServerMsg event = 801
	TypeChatCommand   event = 802
	TypeChatWhisper   event = 802
)

type message struct {
	Source       string   `json:"source"`
	Destinations []string `json:"-"`
	Type         event    `json:"type"`
	Body         string   `json:"body"`
	DateTime     string   `json:"dateTime"`
}
