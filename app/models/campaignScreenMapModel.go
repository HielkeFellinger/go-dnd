package models

type CampaignScreenMapItems struct {
	MapId    string
	Elements map[string]CampaignScreenMapItemElement
}

type CampaignScreenMapItemElement struct {
	Id          string   `json:"Id"`
	Type        string   `json:"Type"`
	EntityName  string   `json:"EntityName"`
	EntityId    string   `json:"EntityId"`
	Hidden      bool     `json:"Hidden"`
	Controllers []string `json:"Controllers"`
	MapId       string   `json:"MapId"`
	Html        string   `json:"Html"`
	Position    CampaignScreenMapPosition
	Image       CampaignImage
	Health      CampaignScreenMapItemHealth
}

func (mi *CampaignScreenMapItemElement) HasHealth() bool {
	return mi.Health.Total != 0
}

type CampaignScreenMapPosition struct {
	X string `json:"X"`
	Y string `json:"Y"`
}

type CampaignScreenMapItemHealth struct {
	Total   uint
	Current int
}
