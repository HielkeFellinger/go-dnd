package models

import "github.com/hielkefellinger/go-dnd/app/ecs"

type CampaignCharacter struct {
	Id          string
	Name        string
	Description string
	Level       string
	Hidden      bool
	Image       CampaignImage
	Health      CampaignCharacterHealth
	Inventories []CampaignInventory
	Controllers []string
	Images      []CampaignImage
}

type CampaignDropdownCharacter struct {
	Id       string
	Name     string
	Selected bool
	Source   ecs.Entity
}

type CampaignCharacters []CampaignCharacter

func (c CampaignCharacters) Len() int {
	return len(c)
}

func (c CampaignCharacters) Less(i, j int) bool {
	return c[i].Name < c[j].Name // Compare names for ordering
}

func (c CampaignCharacters) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

func GetNewCampaignCharacter() CampaignCharacter {
	return CampaignCharacter{
		Inventories: make([]CampaignInventory, 0),
		Controllers: make([]string, 0),
		Images:      make([]CampaignImage, 0),
	}
}

type CampaignCharacterHealth struct {
	Id                 string `json:"id"`
	Damage             string `json:"damage"`
	TemporaryHitPoints string `json:"temp"`
	MaximumHitPoints   string `json:"max"`
}
