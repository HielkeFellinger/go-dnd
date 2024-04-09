package models

type CampaignScreenMapItems struct {
	MapId    string
	Elements map[string]CampaignScreenMapItemElement
}

type SetActivity struct {
	Id     string `json:"Id"`
	Active bool   `json:"Active"`
}

type AddMapItem struct {
	EntityId string `json:"EntityId"`
	MapId    string `json:"MapId"`
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
