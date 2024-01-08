package models

type CampaignScreenContent struct {
	Tabs    []CampaignTabItem
	Content []CampaignContentItem
}

func NewCampaignScreenContent() CampaignScreenContent {
	return CampaignScreenContent{
		Tabs:    make([]CampaignTabItem, 0),
		Content: make([]CampaignContentItem, 0),
	}
}

type CampaignContentItem struct {
	Id   string
	Html string
}

type CampaignTabItem struct {
	Id   string
	Html string
}
