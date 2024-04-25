package models

import (
	"errors"
	"github.com/hielkefellinger/go-dnd/app/helpers"
	"gorm.io/gorm"
	"log"
)

type Campaign struct {
	gorm.Model
	Private       bool
	Active        bool   `gorm:"-"`
	UserIsLead    bool   `gorm:"-"`
	Title         string `form:"title"`
	GameFile      string
	Description   string `form:"description"`
	Password      string `form:"password"`
	PasswordCheck string `gorm:"-" form:"passwordCheck"`
	LeadID        uint
	Lead          User        `gorm:"foreignKey:LeadID"`
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
	if result := DB.Create(&campaign); result.Error != nil {
		return errors.New("campaign could not created")
	}

	return nil
}

func (service *CampaignService) AddUserToCampaign(user User, campaign Campaign) error {
	if err := DB.Model(&campaign).Association("Users").Append(&user); err != nil {
		return err
	}

	return nil
}

func (service *CampaignService) RetrieveCampaignsLinkedToUser(user User) ([]Campaign, error) {
	var campaigns []Campaign
	var campaignIds []uint

	DB.Distinct("campaign_id").
		Table("campaign_users").
		Where("user_id = ?", user.ID).
		Find(&campaignIds)

	log.Println(campaignIds)

	result := DB.Where(&Campaign{LeadID: user.ID}).
		Or("id IN ?", campaignIds).
		Preload("Users").
		Find(&campaigns)

	return campaigns, result.Error
}

func (service *CampaignService) RetrieveCampaignsByIds(ids []uint) ([]Campaign, error) {
	// Insure a return
	if len(ids) == 0 {
		return make([]Campaign, 0), nil
	}

	var campaigns []Campaign
	result := DB.Where("id IN ?", ids).
		Find(&campaigns)

	return campaigns, result.Error
}
