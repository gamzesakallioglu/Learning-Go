package product

import (
	"context"

	"github.com/gamze.sakallioglu/learningGo/bitirme-projesi-gamzesakallioglu/internal/models"
	"gorm.io/gorm"
)

type Repository interface {
	GetProductByName(ctx context.Context, name *string) (*models.Product, error)
	GetProductByID(ctx context.Context, id *string) (*models.Product, error)
	CreateProduct(ctx context.Context, product *models.Product) error
	SearchProducts(ctx context.Context, key *string, pageSize int, page int, sorting string) (*models.Products, int64)
	UpdateStockNumber(ctx context.Context, id *string, stockNumber int) error
	GetAllProducts(ctx context.Context, pagesize int, page int, sorting string) (*models.Products, int64)
	DeleteProductByID(ctx context.Context, id *string) error
	UpdateProductByID(ctx context.Context, id *string, productUpdated *models.Product) error
	Migration()
}

type productRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &productRepository{db: db}
}

// Migrations keeps db schema up to date
func (p *productRepository) Migration() {
	p.db.AutoMigrate(&models.Product{})
}

func (p *productRepository) SearchProducts(ctx context.Context, key *string, pageSize int, page int, sorting string) (*models.Products, int64) {

	if page == 0 {
		page = 1
	}
	offset := (page - 1) * pageSize
	var products *models.Products
	var totalcount int64
	keyToWrite := "%" + *key + "%"
	p.db.Where("name LIKE ?", keyToWrite).Or("stock_code LIKE ?", keyToWrite).Offset(offset).Limit(pageSize).Order(sorting).Preload("Category").Find(&products).Count(&totalcount)

	return products, totalcount
}

func (p *productRepository) UpdateProductByID(ctx context.Context, id *string, productUpdated *models.Product) error {

	if err := p.db.Where("id = ?", id).Updates(productUpdated).Error; err != nil {
		return err
	}
	return nil
}

func (p *productRepository) DeleteProductByID(ctx context.Context, id *string) error {

	if err := p.db.Where("id = ?", id).Delete(&models.Product{}).Error; err != nil {
		return err
	}
	return nil
}

func (p *productRepository) GetAllProducts(ctx context.Context, pagesize int, page int, sorting string) (*models.Products, int64) {

	if page <= 0 {
		page = 1
	}

	offset := (page - 1) * pagesize
	var products *models.Products
	var totalcount int64
	p.db.Offset(offset).Limit(pagesize).Order(sorting).Preload("Category").Find(&products).Count(&totalcount)

	return products, totalcount
}

func (p *productRepository) GetProductByName(ctx context.Context, name *string) (*models.Product, error) {
	var product *models.Product
	if err := p.db.WithContext(ctx).Where("name = ?", *name).First(&product).Error; err != nil {
		return nil, err
	}
	return product, nil
}

func (p *productRepository) GetProductByID(ctx context.Context, id *string) (*models.Product, error) {
	var product *models.Product
	if err := p.db.WithContext(ctx).Where("id = ?", *id).First(&product).Error; err != nil {
		return nil, err
	}
	return product, nil
}

func (p *productRepository) CreateProduct(ctx context.Context, product *models.Product) error {
	// if there is a product has the same name - already checked in the service. Only insert
	result := p.db.Select("Name", "Price", "StockNumber", "StockCode", "Description", "CategoryID").Create(&product)
	if result != nil {
		return result.Error
	}

	return nil

}

func (p *productRepository) UpdateStockNumber(ctx context.Context, id *string, stockNumber int) error {

	if err := p.db.Model(&models.Product{}).Where("id = ?", *id).Update("stock_number", stockNumber).Error; err != nil {
		return err
	}

	return nil
}
