package session

const (
	MSG_TYPE_GAME_CLOSE int = 0
	MSG_TYPE_GAME_START int = 1

	MSG_TYPE_USER_JOIN int = 400

	MSG_TYPE_CHAT_BROADCAST int = 800
)

type message struct {
	Source       string   `json:"source"`
	Destinations []string `json:"-"`
	Type         int      `json:"type"`
	Body         string   `json:"body"`
}
