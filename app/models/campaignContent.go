package models

type CampaignMap struct {
	Id          string
	Name        string
	Description string
	X           uint
	Y           uint
	Active      bool
	ActiveImage CampaignImage
	Images      []CampaignImage
}

type CampaignImage struct {
	Name   string `json:"Name"`
	Url    string `json:"Url"`
	Id     string
	Active bool
}

type CampaignMapCellContent struct {
	Id      string
	Visible bool
	Health  uint
	Image   CampaignImage
}
