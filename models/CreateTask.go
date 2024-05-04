package models

import "gorm.io/gorm"

type Tasks struct {
	gorm.Model
	Title       string
	Description string
	Status      string
}
