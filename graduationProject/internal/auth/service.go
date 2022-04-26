package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/gamze.sakallioglu/learningGo/bitirme-projesi-gamzesakallioglu/internal/api"
	"github.com/gamze.sakallioglu/learningGo/bitirme-projesi-gamzesakallioglu/internal/models"
	"github.com/gamze.sakallioglu/learningGo/bitirme-projesi-gamzesakallioglu/utils"
)

type authService struct {
	repo Repository
}

type Service interface {
	GetCustomer(ctx context.Context, email *string, password *string) (*models.Customer, error)
	GetUser(ctx context.Context, email *string, password *string) (*models.User, error)
	CheckCustomerExists(ctx context.Context, email *string) (*models.Customer, error)
	CheckUserExists(ctx context.Context, email *string) (*models.User, error)
	CreateCustomer(ctx context.Context, customer *api.CustomerSignUp) (*models.Customer, error)
	CreateUser(ctx context.Context, customer *api.UserSignUp) (*models.User, error)
}

func NewAuthService(repo Repository) Service {
	if repo == nil {
		return nil
	}

	return &authService{repo: repo}
}

func (a authService) GetCustomer(ctx context.Context, email *string, password *string) (*models.Customer, error) {
	passwordEncoded := utils.GetMD5Hash(password)

	if len(*email) == 0 {
		return nil, fmt.Errorf("email cannot be empty")
	}

	customer := a.repo.GetCustomer(ctx, email, &passwordEncoded)
	return customer, nil
}

func (a authService) GetUser(ctx context.Context, email *string, password *string) (*models.User, error) {
	passwordEncoded := utils.GetMD5Hash(password)

	if len(*email) == 0 {
		return nil, fmt.Errorf("email cannot be empty")
	}

	user := a.repo.GetUser(ctx, email, &passwordEncoded)
	return user, nil
}

func (a authService) CheckCustomerExists(ctx context.Context, email *string) (*models.Customer, error) {
	if len(*email) == 0 {
		return nil, fmt.Errorf("email cannot be empty")
	}

	customer := a.repo.CheckCustomerExists(ctx, email)
	return customer, nil
}

func (a authService) CheckUserExists(ctx context.Context, email *string) (*models.User, error) {
	if len(*email) == 0 {
		return nil, fmt.Errorf("email cannot be empty")
	}

	user := a.repo.CheckUserExists(ctx, email)
	return user, nil
}

func (a authService) CreateCustomer(ctx context.Context, customer *api.CustomerSignUp) (*models.Customer, error) {

	passwordEncoded := utils.GetMD5Hash(customer.Password)
	customerToGo := models.Customer{
		Name:     customer.Name,
		Email:    customer.Email,
		Address:  customer.Address,
		Password: &passwordEncoded,
		Phone:    customer.Phone,
	}
	customerExisted, err := a.CheckCustomerExists(ctx, customer.Email)

	if err != nil {
		return nil, err
	}

	if customerExisted != nil {
		return nil, errors.New("this email already taken by another user")
	}

	err = a.repo.CreateCustomer(ctx, &customerToGo)
	if err != nil {
		return nil, err
	}

	return &customerToGo, nil

}

func (a authService) CreateUser(ctx context.Context, user *api.UserSignUp) (*models.User, error) {

	passwordEncoded := utils.GetMD5Hash(user.Password)
	userToGo := models.User{
		Name:     user.Name,
		Email:    user.Email,
		Password: &passwordEncoded,
		Phone:    user.Phone,
		UserRole: user.Role,
	}
	userExisted, err := a.CheckUserExists(ctx, user.Email)

	if err != nil {
		return nil, err
	}

	if userExisted != nil {
		return nil, errors.New("this email already taken by another user")
	}

	err = a.repo.CreateUser(ctx, &userToGo)
	if err != nil {
		return nil, err
	}

	return &userToGo, nil

}
