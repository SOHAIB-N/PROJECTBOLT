package models

import (
	"gorm.io/gorm"
)

type MenuItem struct {
	gorm.Model
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Category    string  `json:"category"`
	ImageURL    string  `json:"image_url"`
	Available   bool    `json:"available" gorm:"default:true"`
}