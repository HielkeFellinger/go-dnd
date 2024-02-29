package models

type CampaignMap struct {
	Id          string
	Name        string
	Description string
	X           uint
	Y           uint
	Active      bool
	Image       CampaignImage
}

type CampaignImage struct {
	Name string `json:"Name"`
	Url  string `json:"Url"`
}

type CampaignMapCellContent struct {
	Id      string
	Visible bool
	Health  uint
	Image   CampaignImage
}
