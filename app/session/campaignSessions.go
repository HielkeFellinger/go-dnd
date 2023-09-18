package session

import (
	"github.com/hielkefellinger/go-dnd/app/models"
	"log"
)

var runningCampaignSessionsContainer = initCampaignSessionsContainer()

type campaignSessionsContainer struct {
	Register   chan *campaignPool
	Unregister chan *campaignPool
	Pools      map[*campaignPool]bool
}

func initCampaignSessionsContainer() *campaignSessionsContainer {
	c := &campaignSessionsContainer{
		Register:   make(chan *campaignPool),
		Unregister: make(chan *campaignPool),
		Pools:      make(map[*campaignPool]bool),
	}
	go c.Run()

	return c
}

func (c *campaignSessionsContainer) Run() {
	for {
		select {
		case pool := <-c.Register:
			// Skip if already added
			match := false
			for campaignPool := range c.Pools {
				if campaignPool.Id == pool.Id {
					match = true
					break
				}
			}
			if !match {
				c.Pools[pool] = true
				log.Printf("Adding new Campaign Pool `%d` total count : %d", pool.Id, len(c.Pools))
			}
			break
		case pool := <-c.Unregister:
			delete(c.Pools, pool)
			log.Printf("Removing Campaign Pool `%d` total count : %d", pool.Id, len(c.Pools))
			break
		}
	}
}

func (c *campaignSessionsContainer) initAndRegisterCampaignPool(campaign models.Campaign) {
	pool := initCampaignPool(campaign.ID, campaign.Lead.Name)
	go pool.Run()
	runningCampaignSessionsContainer.Register <- pool
}

func (c *campaignSessionsContainer) addClientToCampaignPool(id uint, client *campaignClient) bool {
	pool := c.getCampaignPoolById(id)
	if pool != nil {
		client.Pool = pool
		pool.Register <- client
	}
	return false
}

func (c *campaignSessionsContainer) getCampaignPoolById(id uint) *campaignPool {
	for pool := range c.Pools {
		if pool != nil && pool.Id == id {
			return pool
		}
	}
	return nil
}

func (c *campaignSessionsContainer) IsCampaignRunning(id uint) bool {
	return c.getCampaignPoolById(id) != nil
}

// Public functions

func IsCampaignRunning(id uint) bool {
	return runningCampaignSessionsContainer.IsCampaignRunning(id)
}
