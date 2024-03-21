package models

type CampaignCharacter struct {
	Id          string
	Name        string
	Description string
	Image       CampaignImage
	Health      CampaignCharacterHealth
	Inventories []CampaignInventory
}

type CampaignCharacterHealth struct {
	Damage             string
	TemporaryHitPoints string
	MaximumHitPoints   string
}
