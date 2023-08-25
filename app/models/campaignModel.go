package models

import (
	"errors"
	"github.com/hielkefellinger/go-dnd/app/helpers"
	"gorm.io/gorm"
)

type Campaign struct {
	gorm.Model
	Private       bool
	Title         string `form:"title"`
	Description   string `form:"description"`
	Password      string `form:"password"`
	PasswordCheck string `gorm:"-" form:"passwordCheck"`
	LeadId        int
	Lead          User        `gorm:"foreignKey:LeadId"`
	Users         []User      `gorm:"many2many:campaign_users;"`
	Characters    []Character `gorm:"many2many:campaign_characters;"`
}

type CampaignService struct{}

func (service *CampaignService) InsertCampaign(campaign *Campaign) error {
	// Check password match
	if campaign.PasswordCheck != campaign.Password {
		return errors.New("passwords do not match")
	}

	// Hash & update password
	hashByteArray, err := helpers.HashPassword(campaign.Password)
	if err != nil {
		return errors.New("password could not be hashed")
	}
	campaign.Password = string(hashByteArray)

	// Attempt to save
	result := DB.Create(&campaign)
	if result.Error != nil {
		return errors.New("user could not created")
	}

	return nil
}
