package models

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
}

func GetNewCampaignCharacter() CampaignCharacter {
	return CampaignCharacter{
		Inventories: make([]CampaignInventory, 0),
		Controllers: make([]string, 0),
	}
}

type CampaignCharacterHealth struct {
	Id                 string `json:"id"`
	Damage             string `json:"damage"`
	TemporaryHitPoints string `json:"temp"`
	MaximumHitPoints   string `json:"max"`
}
