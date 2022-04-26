package customer

import (
	"github.com/gamze.sakallioglu/learningGo/bitirme-projesi-gamzesakallioglu/internal/models"
	"gorm.io/gorm"
)

type CustomerRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *CustomerRepository {
	return &CustomerRepository{db: db}
}

// Migrations keeps db schema up to date
func (b *CustomerRepository) Migration() {
	b.db.AutoMigrate(&models.Customer{})
}
