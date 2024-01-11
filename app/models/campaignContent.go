package models

type CampaignMap struct {
	Id          string
	Name        string
	Description string
	X           uint
	Y           uint
	Enabled     bool
	Image       CampaignImage
}

type CampaignImage struct {
	Name string
	Url  string
}

type CampaignMapCellContent struct {
	Id      string
	Visible bool
	Health  uint
	Image   CampaignImage
}
