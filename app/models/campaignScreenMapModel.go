package models

type CampaignScreenMapItems struct {
	MapId    string
	Elements map[string]CampaignScreenMapItemElement
}

type CampaignScreenMapItemElement struct {
	Id          string   `json:"Id"`
	EntityName  string   `json:"EntityName"`
	EntityId    string   `json:"EntityId"`
	Hidden      bool     `json:"Hidden"`
	Controllers []string `json:"Controllers"`
	MapId       string   `json:"MapId"`
	Html        string   `json:"Html"`
	Position    CampaignScreenMapPosition
	Image       CampaignImage
}

type CampaignScreenMapPosition struct {
	X string `json:"X"`
	Y string `json:"Y"`
}
