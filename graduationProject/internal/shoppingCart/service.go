package shoppingCart

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/gamze.sakallioglu/learningGo/bitirme-projesi-gamzesakallioglu/internal/api"
)

type shoppingCartService struct {
	repo Repository
}

type Service interface {
	AddItemToCart(ctx context.Context, shoppingCartItem *api.ShoppingCart, mux *sync.Mutex) error
	UpdateCartItem(ctx context.Context, shoppingCartItem *api.ShoppingCart, mux *sync.Mutex) error
	GetShoppingCart(ctx context.Context, productID string, customerID string) (*api.ShoppingCart, error)
	IncreaseQuantity(ctx context.Context, productID string, customerID string, quantitiy int) error
	ListProductsInCart(ctx context.Context, customerID string) ([]*api.ShoppingCart, error)
	DeleteItemFromCart(ctx context.Context, productID string, customerID string) error
}

func NewShoppingCartService(repo Repository) Service {
	if repo == nil {
		return nil
	}

	return &shoppingCartService{repo: repo}
}

func (p *shoppingCartService) UpdateCartItem(ctx context.Context, shoppingCartItem *api.ShoppingCart, mux *sync.Mutex) error {

	mux.Lock()
	defer mux.Unlock()

	productExist, _ := p.repo.GetProductById(ctx, shoppingCartItem.Product.ID)

	if productExist == nil {
		return errors.New("product does not exist")
	}

	if productExist.StockNumber < int(shoppingCartItem.Quantity) {
		return fmt.Errorf("what customer asked is more than actual stock. stock: %v", shoppingCartItem.Product.StockNumber)
	}

	shoppingCartExits, _ := p.GetShoppingCart(ctx, shoppingCartItem.Product.ID, shoppingCartItem.Customer.ID)
	if shoppingCartExits != nil {
		p.repo.UpdateQuantity(ctx, shoppingCartItem.Product.ID, shoppingCartItem.Customer.ID, 0, int(shoppingCartItem.Quantity))
		return nil
	}

	return errors.New("the item is not in the cart")
}

func (p *shoppingCartService) AddItemToCart(ctx context.Context, shoppingCartItem *api.ShoppingCart, mux *sync.Mutex) error {

	mux.Lock()
	defer mux.Unlock()
	productExist, _ := p.repo.GetProductById(ctx, shoppingCartItem.Product.ID)

	if productExist == nil {
		return errors.New("product does not exist")
	}

	if productExist.StockNumber < int(shoppingCartItem.Quantity) {
		return fmt.Errorf("what customer asked is more than actual stock. stock: %v", shoppingCartItem.Product.StockNumber)
	}

	shoppingCartExits, _ := p.GetShoppingCart(ctx, shoppingCartItem.Product.ID, shoppingCartItem.Customer.ID)
	if shoppingCartExits != nil {
		p.IncreaseQuantity(ctx, shoppingCartItem.Product.ID, shoppingCartItem.Customer.ID, int(shoppingCartItem.Quantity))
		return nil
	}

	p.repo.AddItemToCart(ctx, shoppingCartItem.Product.ID, shoppingCartItem.Customer.ID, int(shoppingCartItem.Quantity))
	return nil
}

func (p *shoppingCartService) GetShoppingCart(ctx context.Context, productID string, customerID string) (*api.ShoppingCart, error) {

	if len(productID) <= 0 {
		return nil, errors.New("product ID should be given")
	}

	if len(customerID) <= 0 {
		return nil, errors.New("customer ID should be given")
	}

	shoppingCartExist, _ := p.repo.GetShoppingCart(ctx, productID, customerID)
	if shoppingCartExist == nil {
		return nil, fmt.Errorf("shopping cart does not exist")
	}

	shoppingCartExistApi := ShoppingCartToResponse(shoppingCartExist)

	return shoppingCartExistApi, nil
}

func (p *shoppingCartService) IncreaseQuantity(ctx context.Context, productID string, customerID string, quantity int) error {

	shoppingCart, _ := p.GetShoppingCart(ctx, productID, customerID)
	newQuantity := int(shoppingCart.Quantity) + int(quantity)

	p.repo.UpdateQuantity(ctx, productID, customerID, quantity, newQuantity)
	return nil
}

func (p *shoppingCartService) ListProductsInCart(ctx context.Context, customerID string) ([]*api.ShoppingCart, error) {

	shoppingCartList, err := p.repo.GetAllProductsInCart(ctx, customerID)
	if err != nil {
		return nil, err
	}

	shoppingCartListApi := ShoppingCartsToResponses(shoppingCartList)
	return shoppingCartListApi, nil
}

func (p *shoppingCartService) DeleteItemFromCart(ctx context.Context, productID string, customerID string) error {

	err := p.repo.DeleteItemFromCart(ctx, productID, customerID)
	if err != nil {
		return err
	}

	return nil

}
