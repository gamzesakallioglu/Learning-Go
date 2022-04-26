package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProductCategory struct {
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
	Name        *string        `gorm:"unique"`
	Description string
	IsParent    bool
	ParentID    string
}

type ProductCategories []ProductCategory
