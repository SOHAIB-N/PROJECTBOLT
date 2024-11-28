package models

import (
	"time"

	"gorm.io/gorm"
)

type Order struct {
	gorm.Model
	UserID       uint
	User         User
	OrderNumber  string `gorm:"unique"`
	Items        []OrderItem
	TotalAmount  float64
	Status       string // pending, confirmed, preparing, ready, delivered
	PaymentID    string
	PaymentStatus string
	IsDelivery   bool
	RoomNumber   string
	DeliveryTime *time.Time
}

type OrderItem struct {
	gorm.Model
	OrderID    uint
	MenuItemID uint
	MenuItem   MenuItem
	Quantity   int
	Price      float64
	TotalPrice float64
}

type Payment struct {
	gorm.Model
	OrderID       uint
	Order         Order
	Amount        float64
	TransactionID string
	PaymentMethod string
	Status        string
	PaymentData   string // Store additional payment gateway data
}