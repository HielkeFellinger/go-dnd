package models

type CampaignMap struct {
	Id          string
	Name        string
	Description string
	X           uint
	Y           uint
	Image       CampaignMapImage
}

type CampaignMapImage struct {
	Name string
	Url  string
}
