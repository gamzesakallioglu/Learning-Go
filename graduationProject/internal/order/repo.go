package order

import (
	"context"
	"time"

	"github.com/gamze.sakallioglu/learningGo/bitirme-projesi-gamzesakallioglu/internal/models"
	"github.com/gamze.sakallioglu/learningGo/bitirme-projesi-gamzesakallioglu/internal/product"
	"github.com/gamze.sakallioglu/learningGo/bitirme-projesi-gamzesakallioglu/internal/shoppingCart"
	"github.com/gamze.sakallioglu/learningGo/bitirme-projesi-gamzesakallioglu/utils"
	"gorm.io/gorm"
)

type Repository interface {
	GetShoppingCart(ctx context.Context, customerID string) (*models.ShoppingCarts, error)
	GetProductByID(ctx context.Context, productID *string) (*models.Product, error)
	UpdateProductStockNumber(ctx context.Context, productID *string, newStockNumber int) error
	CompleteOrder(ctx context.Context, orderHeader *models.OrderHeader, orderDetails []*models.OrderDetail) (string, error)
	GetAllOrders(ctx context.Context, customerID string) ([]*models.OrderHeader, error)
	GetOrderByID(ctx context.Context, orderID string, customerID string) (*models.OrderHeader, error)
	CancelOrder(ctx context.Context, orderID string, customerID string) error
	Migration()
}

type orderRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &orderRepository{db: db}
}

// Migrations keeps db schema up to date
func (p *orderRepository) Migration() {
	p.db.AutoMigrate(&models.OrderHeader{})
	p.db.AutoMigrate(&models.OrderDetail{})
}

func (p *orderRepository) GetOrderByID(ctx context.Context, orderID string, customerID string) (*models.OrderHeader, error) {

	// added 'and customer_id' to query. Customers should see only their order headers

	var orderHeader *models.OrderHeader
	if err := p.db.WithContext(ctx).Preload("OrderDetails").Preload("OrderDetails.Product").Where("id = ? AND customer_id = ?", orderID, customerID).First(&orderHeader).Error; err != nil {
		return nil, err
	}
	return orderHeader, nil
}

func (p *orderRepository) CancelOrder(ctx context.Context, orderID string, customerID string) error {

	orderHeader, err := p.GetOrderByID(ctx, orderID, customerID)
	if err != nil {
		return err
	}

	if err := p.db.Model(&models.OrderHeader{}).Where("id = ?", orderHeader.ID).Update("cancel_date", time.Now()).Error; err != nil {
		return err
	}

	for _, orderDetail := range orderHeader.OrderDetails {
		productID := orderDetail.Product.ID.String()
		pr, _ := p.GetProductByID(ctx, &productID)
		newStockNumber := pr.StockNumber + orderDetail.Quantity
		err = p.UpdateProductStockNumber(ctx, &productID, newStockNumber)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *orderRepository) GetAllOrders(ctx context.Context, customerID string) ([]*models.OrderHeader, error) {

	var orderHeaders []*models.OrderHeader
	if err := p.db.Preload("OrderDetails").Preload("OrderDetails.Product").Preload("OrderDetails.Product.Category").Find(&orderHeaders, "customer_id = ?", customerID).Error; err != nil {
		return nil, err
	}

	return orderHeaders, nil
}

func (p *orderRepository) GetShoppingCart(ctx context.Context, customerID string) (*models.ShoppingCarts, error) {

	shoppingCartList, err := shoppingCart.NewRepository(p.db).GetAllProductsInCart(ctx, customerID)
	if err != nil {
		return nil, err
	}
	return shoppingCartList, nil
}

func (p *orderRepository) GetProductByID(ctx context.Context, productID *string) (*models.Product, error) {

	productRepo := product.NewRepository(p.db)
	product, err := productRepo.GetProductByID(ctx, productID)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (p *orderRepository) UpdateProductStockNumber(ctx context.Context, productID *string, newStockNumber int) error {

	productRepo := product.NewRepository(p.db)
	_, err := productRepo.GetProductByID(ctx, productID)
	if err != nil {
		return err
	}

	productRepo.UpdateStockNumber(ctx, productID, newStockNumber)

	return nil
}

func (p *orderRepository) CompleteOrder(ctx context.Context, orderHeader *models.OrderHeader, orderDetails []*models.OrderDetail) (string, error) {

	productRepo := product.NewRepository(p.db)
	shoppingCartRepo := shoppingCart.NewRepository(p.db)

	if err := p.db.Create(&orderHeader).Error; err != nil {
		return "", err
	}

	orderHeaderID := orderHeader.ID

	for _, orderDetail := range orderDetails {

		orderDetail.OrderHeader.ID = orderHeaderID
		if err := p.db.Create(&orderDetail).Error; err != nil {
			return "", err
		}

		product, _ := productRepo.GetProductByID(ctx, &orderDetail.ProductID)
		newStockNumber := product.StockNumber - orderDetail.Quantity
		productRepo.UpdateStockNumber(ctx, &orderDetail.ProductID, newStockNumber)

	}

	if err := shoppingCartRepo.DeleteAllItemsInTheCart(ctx, *orderHeader.CustomerID); err != nil {
		return "", err
	}

	orderHeaderIDString := orderHeaderID.String()
	orderNumber := utils.GetMD5Hash(&orderHeaderIDString)

	return orderNumber, nil

}
