package order

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/gamze.sakallioglu/learningGo/bitirme-projesi-gamzesakallioglu/internal/api"
	"github.com/gamze.sakallioglu/learningGo/bitirme-projesi-gamzesakallioglu/internal/models"
)

type orderService struct {
	repo Repository
}

type Service interface {
	CompleteOrder(ctx context.Context, orderInfo *api.CompleteOrder, customerID string, mux *sync.Mutex) (string, error)
	GetProductByID(ctx context.Context, ID *string) (*models.Product, error)
	ListOrders(ctx context.Context, customerID string) ([]*api.Order, error)
	GetOrderByID(ctx context.Context, orderID string, customerID string) (*api.Order, error)
	CancelOrder(ctx context.Context, orderID string, customerID string) error
}

func NewOrderService(repo Repository) Service {
	if repo == nil {
		return nil
	}

	return &orderService{repo: repo}
}

func (p *orderService) CancelOrder(ctx context.Context, orderID string, customerID string) error {

	ordr, err := p.GetOrderByID(ctx, orderID, customerID)
	if err != nil {
		return err
	}

	if ordr.CancelDate[:4] != "0001" {
		return errors.New("this order already canceled")
	}

	orderDate, _ := time.Parse("2006-01-02 15:04:05", ordr.OrderDate)
	orderDatePlus14Days := orderDate.Add(14 * 24 * time.Hour)
	if orderDatePlus14Days.Format("2006-01-02 15:04:05") < time.Now().Format("2006-01-02 15:04:05") {
		return errors.New("orders cannot be canceled after 14 days of order date")
	}

	err = p.repo.CancelOrder(ctx, orderID, customerID)
	if err != nil {
		return err
	}

	return nil

}

func (p *orderService) GetOrderByID(ctx context.Context, orderID string, customerID string) (*api.Order, error) {

	orderHeader, err := p.repo.GetOrderByID(ctx, orderID, customerID)
	if err != nil {
		return nil, err
	}

	orderApi := OrderToResponse(orderHeader)
	return orderApi, err
}

func (p *orderService) ListOrders(ctx context.Context, customerID string) ([]*api.Order, error) {

	orderHeaders, err := p.repo.GetAllOrders(ctx, customerID)
	if err != nil {
		return nil, err
	}

	orderHeadersApi := OrdersToResponse(orderHeaders)
	return orderHeadersApi, nil
}

func (p *orderService) CompleteOrder(ctx context.Context, orderInfo *api.CompleteOrder, customerID string, mux *sync.Mutex) (string, error) {

	mux.Lock()
	defer mux.Unlock()

	shoppingCarts, err := p.repo.GetShoppingCart(ctx, customerID)
	if err != nil {
		return "", err
	}

	for _, pr := range *shoppingCarts {
		productExist, _ := p.GetProductByID(ctx, &pr.ProductID)

		if productExist == nil {
			return "", fmt.Errorf(*productExist.Name, " product does not exist")
		}

		if productExist.StockNumber < int(pr.Quantity) {
			return "", fmt.Errorf(*productExist.Name, " what customer asked is more than actual stock. stock: %v", productExist.StockNumber)
		}
	}

	orderHeader, orderDetails := ShoppingCartToOrderHeaderAndDetail(orderInfo, customerID, shoppingCarts)

	orderNumber, err := p.repo.CompleteOrder(ctx, orderHeader, orderDetails)
	if err != nil {
		return "", err
	}

	return orderNumber, nil
}

func (p *orderService) GetProductByID(ctx context.Context, ID *string) (*models.Product, error) {
	if len(*ID) < 1 {
		return nil, errors.New("id cannot be empty")
	}

	product, err := p.repo.GetProductByID(ctx, ID)
	if err != nil {
		return nil, err
	}

	return product, nil
}
