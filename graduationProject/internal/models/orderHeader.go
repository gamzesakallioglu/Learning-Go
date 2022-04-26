package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type OrderHeader struct {
	ID             uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	OrderDetails   []OrderDetail
	OrderNumber    string
	CustomerID     *string
	Customer       Customer
	OrderDate      time.Time
	ShippingDate   time.Time
	ReceiveDate    time.Time
	CancelDate     time.Time
	PaymentDueDate time.Time
	PaymentDate    time.Time
	OrderStatus    string
	PaymentStatus  string
	PaymentType    string
	PhoneNumber    *string
	Address        *string
	OrderTotal     decimal.Decimal
}
