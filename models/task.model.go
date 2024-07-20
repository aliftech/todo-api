package models

import "gorm.io/gorm"

type Tasks struct {
	gorm.Model
	Userid      uint
	Title       string
	Description string
	Due         string
	Status      string
}
