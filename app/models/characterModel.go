package models

import "gorm.io/gorm"

type Character struct {
	gorm.Model
	Name string
}
