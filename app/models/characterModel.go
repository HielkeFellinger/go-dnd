package models

import "gorm.io/gorm"

type Character struct {
	gorm.Model
	Name    string
	Id      string        `gorm:"-"`
	Image   CampaignImage `gorm:"-"`
	Visible bool          `gorm:"-"`
	Online  bool          `gorm:"-"`
}

type Characters []Character

func (c Characters) Len() int {
	return len(c)
}

func (c Characters) Less(i, j int) bool {
	return c[i].Name < c[j].Name // Compare names for ordering
}

func (c Characters) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}
