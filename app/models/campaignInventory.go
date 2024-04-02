package models

type CampaignInventory struct {
	Id    string
	Items []CampaignInventoryItem
}

type CampaignInventoryItem struct {
	Id          string
	Count       uint
	Name        string
	Description string
	Damage      string
	Restore     string
	Range       CampaignInventoryItemRange
}

type CampaignInventoryItemRange struct {
	Min string
	Max string
}
