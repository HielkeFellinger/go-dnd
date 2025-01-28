package models

type CampaignInventory struct {
	Id                string
	Description       string
	Name              string
	Slots             string
	Size              uint
	ReadOnly          bool
	ShowDetailButtons bool
	Items             []CampaignInventoryItem
	LinkedInventories []CampaignLinkedInventory
	Characters        CampaignCharacters
}

type CampaignLinkedInventory struct {
	Id          string
	Description string
	Name        string
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
	Image       CampaignImage
}

type CampaignInventoryItemRange struct {
	Min string
	Max string
}
