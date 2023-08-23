package models

import "gorm.io/gorm"

type Campaign struct {
	gorm.Model
	Private     bool
	Title       string
	Description string
	Password    string
	LeadId      int
	Lead        User        `gorm:"foreignKey:LeadId"`
	Users       []User      `gorm:"many2many:campaign_users;"`
	Characters  []Character `gorm:"many2many:campaign_characters;"`
}
