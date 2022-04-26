package user

import (
	"github.com/gamze.sakallioglu/learningGo/bitirme-projesi-gamzesakallioglu/internal/models"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Migrations keeps db schema up to date
func (b *UserRepository) Migration() {
	b.db.AutoMigrate(&models.User{})
}
