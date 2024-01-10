package models

type CampaignMap struct {
	Id    string
	X     string
	Y     string
	Image CampaignMapImage
}

type CampaignMapImage struct {
	Name string
	Url  string
}
