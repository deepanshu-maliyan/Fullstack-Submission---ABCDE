package models

import (
	"time"
)

type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Username  string    `json:"username" gorm:"unique;not null"`
	Password  string    `json:"password" gorm:"not null"`
	Token     string    `json:"token"`
	CartID    uint      `json:"cart_id"`
	CreatedAt time.Time `json:"created_at"`
	Cart      Cart      `json:"cart" gorm:"foreignKey:UserID"`
	Orders    []Order   `json:"orders" gorm:"foreignKey:UserID"`
}

type Item struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"not null"`
	Status    string    `json:"status" gorm:"default:active"`
	Image     string    `json:"image"`
	CreatedAt time.Time `json:"created_at"`
}

type Cart struct {
	ID        uint       `json:"id" gorm:"primaryKey"`
	UserID    uint       `json:"user_id" gorm:"not null"`
	Name      string     `json:"name"`
	Status    string     `json:"status" gorm:"default:active"`
	CreatedAt time.Time  `json:"created_at"`
	CartItems []CartItem `json:"cart_items" gorm:"foreignKey:CartID"`
}

type CartItem struct {
	CartID uint `json:"cart_id" gorm:"primaryKey"`
	ItemID uint `json:"item_id" gorm:"primaryKey"`
	Cart   Cart `json:"cart" gorm:"foreignKey:CartID"`
	Item   Item `json:"item" gorm:"foreignKey:ItemID"`
}

type Order struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	CartID    uint      `json:"cart_id" gorm:"not null"`
	UserID    uint      `json:"user_id" gorm:"not null"`
	CreatedAt time.Time `json:"created_at"`
	Cart      Cart      `json:"cart" gorm:"foreignKey:CartID"`
}
