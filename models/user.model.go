package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email    string `gorm:"unique"`
	Fullname string
	Photo    string
	Password string
}
