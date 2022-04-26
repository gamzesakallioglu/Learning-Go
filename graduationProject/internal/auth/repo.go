package auth

import (
	"context"

	"github.com/gamze.sakallioglu/learningGo/bitirme-projesi-gamzesakallioglu/internal/models"
	"gorm.io/gorm"
)

// Repository encapsulates the logic to access basket from the data source.
type Repository interface {
	GetCustomer(ctx context.Context, email *string, password *string) *models.Customer
	GetUser(ctx context.Context, email *string, password *string) *models.User
	CheckCustomerExists(ctx context.Context, email *string) *models.Customer
	CheckUserExists(ctx context.Context, email *string) *models.User
	CreateCustomer(ctx context.Context, customer *models.Customer) error
	CreateUser(ctx context.Context, user *models.User) error
}

type authRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) Repository {
	if db == nil {
		return nil
	}

	return &authRepository{db: db}
}

func (a authRepository) GetCustomer(ctx context.Context, email *string, password *string) *models.Customer {
	var customer *models.Customer
	if err := a.db.WithContext(ctx).Where("email = ? AND password = ?", *email, *password).First(&customer).Error; err != nil {
		return nil
	}
	return customer
}

func (a authRepository) GetUser(ctx context.Context, email *string, password *string) *models.User {
	var user *models.User
	if err := a.db.WithContext(ctx).Where("email = ? AND password = ?", *email, *password).First(&user).Error; err != nil {
		return nil
	}
	return user
}

func (a authRepository) CheckCustomerExists(ctx context.Context, email *string) *models.Customer {
	var customer *models.Customer
	if err := a.db.WithContext(ctx).Where("email = ?", *email).First(&customer).Error; err != nil {
		return nil
	}
	return customer
}

func (a authRepository) CheckUserExists(ctx context.Context, email *string) *models.User {
	var user *models.User
	if err := a.db.WithContext(ctx).Where("email = ?", *email).First(&user).Error; err != nil {
		return nil
	}
	return user
}

func (a authRepository) CreateCustomer(ctx context.Context, customer *models.Customer) error {
	// if there is a customer has the same email - already checked in the service. Only insert
	result := a.db.Select("Id", "Name", "Email", "Password", "Address", "Phone").Create(&customer)
	if result != nil {
		return result.Error
	}

	return nil
}

func (a authRepository) CreateUser(ctx context.Context, user *models.User) error {
	// if there is a customer has the same email - already checked in the service. Only insert
	result := a.db.Select("Id", "Name", "Email", "Password", "Phone", "user_role").Create(&user)
	if result != nil {
		return result.Error
	}

	return nil
}
