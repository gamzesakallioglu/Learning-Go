package shoppingCart

import (
	"context"

	"github.com/gamze.sakallioglu/learningGo/bitirme-projesi-gamzesakallioglu/internal/models"
	"github.com/gamze.sakallioglu/learningGo/bitirme-projesi-gamzesakallioglu/internal/product"
	"gorm.io/gorm"
)

type Repository interface {
	GetShoppingCart(ctx context.Context, productID string, customerID string) (*models.ShoppingCart, error)
	UpdateQuantity(ctx context.Context, productID string, customerID string, increaseAmount int, newQuantity int) error
	//UpdateItemQuantityInCart(ctx context.Context, productID string, customerID string, newQuantity int) error
	AddItemToCart(ctx context.Context, productID string, customerID string, quantity int) error
	GetProductById(ctx context.Context, productID string) (*models.Product, error)
	GetAllProductsInCart(ctx context.Context, customerID string) (*models.ShoppingCarts, error)
	DeleteItemFromCart(ctx context.Context, productID string, customerID string) error
	DeleteAllItemsInTheCart(ctx context.Context, customerID string) error
	Migration()
}

type shoppingCartRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &shoppingCartRepository{db: db}
}

// Migrations keeps db schema up to date
func (p *shoppingCartRepository) Migration() {
	p.db.AutoMigrate(&models.ShoppingCart{})
}

func (p *shoppingCartRepository) GetAllProductsInCart(ctx context.Context, customerID string) (*models.ShoppingCarts, error) {
	var shoppingCartList *models.ShoppingCarts
	if err := p.db.Preload("Customer").Preload("Product").Preload("Product.Category").Find(&shoppingCartList, "customer_id = ?", customerID).Error; err != nil {
		return nil, err
	}

	return shoppingCartList, nil
}

func (p *shoppingCartRepository) GetShoppingCart(ctx context.Context, productID string, customerID string) (*models.ShoppingCart, error) {
	var shoppingCart *models.ShoppingCart
	if err := p.db.WithContext(ctx).Where("product_id = ? AND customer_id = ?", productID, customerID).First(&shoppingCart).Error; err != nil {
		return nil, err
	}
	return shoppingCart, nil
}

func (p *shoppingCartRepository) GetProductById(ctx context.Context, productID string) (*models.Product, error) {

	productRepo := product.NewRepository(p.db)
	product, err := productRepo.GetProductByID(ctx, &productID)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (p *shoppingCartRepository) UpdateQuantity(ctx context.Context, productID string, customerID string, increaseAmount int, newQuantity int) error {

	err := p.db.Model(&models.ShoppingCart{}).Where("customer_id = ? AND product_id = ?", customerID, productID).Update("quantity", newQuantity).Error
	if err != nil {
		return err
	}

	return nil
}

func (p *shoppingCartRepository) AddItemToCart(ctx context.Context, productID string, customerID string, quantity int) error {

	var shoppingCart models.ShoppingCart

	err := p.db.Where(models.ShoppingCart{ProductID: productID, CustomerID: customerID}).
		Attrs(models.ShoppingCart{ProductID: productID, CustomerID: customerID, Quantity: quantity}).
		FirstOrCreate(&shoppingCart).Error

	if err != nil {
		return err
	}

	return nil
}

func (p *shoppingCartRepository) DeleteItemFromCart(ctx context.Context, productID string, customerID string) error {

	if err := p.db.Unscoped().Where("customer_id = ? AND product_id = ?", customerID, productID).Delete(&models.ShoppingCart{}).Error; err != nil {
		return err
	}

	return nil

}

func (p *shoppingCartRepository) DeleteAllItemsInTheCart(ctx context.Context, customerID string) error {

	if err := p.db.Unscoped().Where("customer_id = ?", customerID).Delete(&models.ShoppingCart{}).Error; err != nil {
		return err
	}

	return nil
}
