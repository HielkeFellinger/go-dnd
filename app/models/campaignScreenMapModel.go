package models

type CampaignScreenMapItems struct {
	MapId    string
	Elements map[string]CampaignScreenMapItemElement
}

type CampaignScreenMapItemElement struct {
	Id          string
	EntityName  string
	EntityId    string
	Controllers []string
	MapId       string
	Html        string
	Position    CampaignScreenMapPosition
	Image       CampaignImage
}

type CampaignScreenMapPosition struct {
	X uint
	Y uint
}
