package game_engine

import (
	"log"
	"slices"
)

var leadId = "lead"
var playerOneId = "playerOne"
var playerTwoId = "playerTwo"
var playerThreeId = "playerThree"

func initTestCampaignPool() *testCampaignPool {
	lead := testCampaignClient{
		Id:   "lead",
		Lead: true,
	}
	playerOne := testCampaignClient{
		Id:   "playerOne",
		Lead: false,
	}
	playerTwo := testCampaignClient{
		Id:   "playerTwo",
		Lead: false,
	}
	playerThree := testCampaignClient{
		Id:   "playerThree",
		Lead: false,
	}

	testPool := testCampaignPool{
		Id:       0,
		LeadId:   "lead",
		Clients:  make(map[*testCampaignClient]bool),
		Messages: make([]EventMessage, 0),
	}
	testPool.Engine = initTestGameEngine()
	testPool.Clients[&lead] = true
	testPool.Clients[&playerOne] = true
	testPool.Clients[&playerTwo] = true
	testPool.Clients[&playerThree] = false

	return &testPool
}

func initTestGameEngine() Engine {
	var baseEngineInst = baseEngine{}

	baseEngineInst.World = loadGame(SpaceGameTest)
	baseEngineInst.EventHandler = &baseEventMessageHandler{}

	return &baseEngineInst
}

type testCampaignPool struct {
	Id       uint
	LeadId   string
	Clients  map[*testCampaignClient]bool
	Engine   Engine
	Messages []EventMessage
}

type testCampaignClient struct {
	Id   string
	Lead bool
}

func (c *testCampaignClient) GetId() string {
	return c.Id
}

func (c *testCampaignClient) IsLead() bool {
	return c.Lead
}

func (pool *testCampaignPool) GetId() uint {
	return pool.Id
}

func (pool *testCampaignPool) GetLeadId() string {
	return pool.LeadId
}

func (pool *testCampaignPool) GetEngine() Engine {
	return pool.Engine
}

func (pool *testCampaignPool) TransmitEventMessage(message EventMessage) {
	log.Printf("TEST Pool transmission Message: '%+v' Source: '%s' & Destination: '%v'",
		message.Id, message.Source, message.Destinations)
	pool.Messages = append(pool.Messages, message)
}

func (pool *testCampaignPool) GetAllClientIds(filterOut ...string) []string {
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
