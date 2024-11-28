package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name     string `json:"name"`
	Email    string `json:"email" gorm:"unique"`
	Password string `json:"password"`
	Phone    string `json:"phone"`
	Orders   []Order
}

type Admin struct {
	gorm.Model
	Email           string `gorm:"unique"`
	Password        string
	LastLoginIP     string
	TwoFactorSecret string
	LastLoginAt     time.Time
}