package models

type CampaignInventory struct {
	Id         string
	Name       string
	Size       uint
	Items      []CampaignInventoryItem
	Characters CampaignCharacters
}

type CampaignInventoryItem struct {
	Id          string
	Count       uint
	Name        string
	Description string
	Damage      string
	Restore     string
	Range       CampaignInventoryItemRange
	Weight      string
}

type CampaignInventoryItemRange struct {
	Min string
	Max string
}
