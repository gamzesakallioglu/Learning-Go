package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Customer struct {
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
	Name        *string
	Email       *string `gorm:"unique"`
	Password    *string
	Address     string
	Phone       string
	PassiveDate time.Time
}
