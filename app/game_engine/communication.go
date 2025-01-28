package game_engine

import (
	"encoding/json"
	"github.com/google/uuid"
	"time"
)

const ServerUser string = "server"

type EventType int

const (
	TypeGameClose EventType = 0
	TypeGameStart EventType = 1
	TypeGameSave  EventType = 2

	TypeUserJoin  EventType = 400
	TypeUserLeave EventType = 401

	TypeLoadFullGame EventType = 500

	TypeLoadCharacters        EventType = 501
	TypeLoadCharactersDetails EventType = 504

	TypeUpdateCharacterHealth EventType = 511
	TypeUpdateCharacterUsers  EventType = 512

	TypeLoadItemDetails EventType = 521

	TypeLoadMap         EventType = 531
	TypeLoadMapEntities EventType = 532
	TypeLoadMapEntity   EventType = 533

	TypeUpdateMapEntity          EventType = 543
	TypeUpdateMapVisibility      EventType = 544
	TypeAddMapItem               EventType = 545
	TypeRemoveMapItem            EventType = 546
	TypeSignalMapItem            EventType = 547
	TypeChangeMapBackgroundImage EventType = 548

	// Management Events
	TypeManagementOverviewStart EventType = 550
	TypeManageMaps              EventType = 551
	TypeManageCharacters        EventType = 552
	TypeManageInventory         EventType = 553
	TypeManageItems             EventType = 554
	TypeManageCampaign          EventType = 555
	TypeManageImages            EventType = 556
	TypeManagementOverviewEnd   EventType = 560

	TypeManagementCrudStart EventType = 5500

	TypeLoadUpsertMap EventType = 5511
	TypeRemoveMap     EventType = 5512
	TypeUpsertMap     EventType = 5513

	TypeLoadUpsertItem EventType = 5541
	TypeUpsertItem     EventType = 5542
	TypeRemoveItem     EventType = 5543

	TypeLoadUpsertCharacter EventType = 5521
	TypeUpsertCharacter     EventType = 5522
	TypeRemoveCharacter     EventType = 5523
	TypeCloneCharacter      EventType = 5524

	TypeLoadUpsertInventory             EventType = 5531
	TypeUpsertInventory                 EventType = 5532
	TypeRemoveInventory                 EventType = 5533
	TypeCloneInventory                  EventType = 5534
	TypeAddItemToInventory              EventType = 5535
	TypeRemoveItemFromInventory         EventType = 5536
	TypeUpdateItemCountInventory        EventType = 5537
	TypeMoveItemCountBetweenInventories EventType = 5538

	TypeManagementCrudEnd EventType = 5600

	TypeManagementError EventType = 5999

	// Set attention to a specific tab

	TypeChatBroadcast   EventType = 800
	TypeChatServerMsg   EventType = 801
	TypeChatCommandRoll EventType = 802
	TypeChatWhisper     EventType = 803
)

type EventMessage struct {
	Id           uuid.UUID `json:"-"`
	Source       string    `json:"source"`
	Destinations []string  `json:"-"`
	Type         EventType `json:"type"`
	Body         string    `json:"body"`
	DateTime     string    `json:"dateTime"`
}

type EventMessageIdBody struct {
	Id   string `json:"Id"`
	Html string `json:"Html"`
}

func (midBody *EventMessageIdBody) ToBodyString() string {
	rawJsonBytes, err := json.Marshal(midBody)
	if err == nil {
		return string(rawJsonBytes)
	}

	return ""
}

func NewEventMessage() EventMessage {
	m := EventMessage{Id: uuid.New()}
	m.ReloadDateTime()
	m.Destinations = make([]string, 0)
	return m
}

func (m *EventMessage) ReloadDateTime() {
	now := time.Now()
	m.DateTime = now.Format("2006-01-02 15:04:05")
}
