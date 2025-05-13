package models

import (
	"errors"
	"github.com/hielkefellinger/go-dnd/app/helpers"
	"gorm.io/gorm"
)

type UserType string

type User struct {
	gorm.Model
	Name          string `gorm:"unique" form:"username"`
	Password      string `form:"password"`
	PasswordCheck string `gorm:"-" form:"passwordCheck"`
	Type          UserType
}

type UserService struct{}

func (service UserService) InsertUser(user *User) error {
	// Check password match
	if user.PasswordCheck != user.Password {
		return errors.New("passwords do not match")
	}

	// Hash & update password
	hashByteArray, err := helpers.HashPassword(user.Password)
	if err != nil {
		return errors.New("password could not be hashed")
	}
	user.Password = string(hashByteArray)

	// Attempt to save
	if result := DB.Create(&user); result.Error != nil {
		return errors.New("user could not created")
	}

	return nil
}
