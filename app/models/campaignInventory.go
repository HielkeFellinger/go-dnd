package models

func GetCampaignInventory() CampaignInventory {
	return CampaignInventory{
		Size:              0,
		Slots:             "0",
		ShowDetailButtons: true,
		Items:             make([]CampaignInventoryItem, 0),
		LinkedInventories: make([]CampaignLinkedInventory, 0),
	}
}

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

func GetCampaignInventoryItem() CampaignInventoryItem {
	return CampaignInventoryItem{
		Images: make([]CampaignImage, 0),
	}
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
	Images      []CampaignImage
}

type CampaignInventoryItemRange struct {
	Min string
	Max string
}
