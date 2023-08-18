package models

import "gorm.io/gorm"

type UserType string

const (
	ADMIN   UserType = "USER_ADMIN"
	REGULAR UserType = "USER_ADMIN"
)

type User struct {
	gorm.Model
	Name     string `gorm:"unique"`
	Password string
	Type     UserType
}
