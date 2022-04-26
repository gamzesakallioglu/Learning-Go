package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type OrderDetail struct {
	ID            uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	OrderHeaderID string
	OrderHeader   OrderHeader
	ProductID     string
	Product       Product
	Quantity      int
	PricePerItem  decimal.Decimal
}
