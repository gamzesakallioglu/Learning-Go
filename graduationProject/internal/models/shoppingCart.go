package models

import (
	"time"

	"github.com/google/uuid"
)

type ShoppingCart struct {
	ID         uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	CustomerID string
	Customer   Customer
	ProductID  string
	Product    Product
	Quantity   int
}

type ShoppingCarts []ShoppingCart
