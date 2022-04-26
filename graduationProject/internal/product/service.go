package product

import (
	"context"
	"encoding/hex"
	"errors"

	"github.com/gamze.sakallioglu/learningGo/bitirme-projesi-gamzesakallioglu/internal/api"
	"github.com/gamze.sakallioglu/learningGo/bitirme-projesi-gamzesakallioglu/internal/models"
	"github.com/gamze.sakallioglu/learningGo/bitirme-projesi-gamzesakallioglu/pkg/pagination"
	"github.com/shopspring/decimal"
)

type productService struct {
	repo Repository
}

type Service interface {
	CreateProduct(ctx context.Context, product *api.Product) error
	GetProductByName(ctx context.Context, name *string) (*models.Product, error)
	GetProductByID(ctx context.Context, id *string) (*models.Product, error)
	SearchProducts(ctx context.Context, key *string, pageSize int, page int) (*pagination.Pagination, error)
	GetAllProducts(ctx context.Context, pageSize int, page int) (*pagination.Pagination, error)
	DeleteProductByID(ctx context.Context, id *string) error
	UpdateProductByID(ctx context.Context, id *string, productUpdated *api.Product) error
}

func NewProductService(repo Repository) Service {
	if repo == nil {
		return nil
	}

	return &productService{repo: repo}
}

func (p *productService) SearchProducts(ctx context.Context, key *string, pageSize int, page int) (*pagination.Pagination, error) {

	if page <= 0 {
		page = 1
	}
	sorting := "ID desc"
	products, count := p.repo.SearchProducts(ctx, key, pageSize, page, sorting)
	productsApi := ProductsToResponses(products)

	pagination := pagination.NewPagination(page, pageSize, count, sorting)
	pagination.Rows = productsApi

	return pagination, nil

}

func (p *productService) UpdateProductByID(ctx context.Context, id *string, productUpdated *api.Product) error {

	cat, err := p.GetProductByID(ctx, id)
	if err != nil {
		return err
	}
	if cat == nil {
		return errors.New("product doesn't exist")
	}

	productUpdatedApi := ResponseToProduct(productUpdated)
	err = p.repo.UpdateProductByID(ctx, id, productUpdatedApi)
	if err != nil {
		return err
	}

	return nil
}

func (p *productService) DeleteProductByID(ctx context.Context, id *string) error {

	cat, err := p.GetProductByID(ctx, id)
	if err != nil {
		return err
	}
	if cat == nil {
		return errors.New("product doesn't exist")
	}

	err = p.repo.DeleteProductByID(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (p *productService) GetAllProducts(ctx context.Context, page int, pagesize int) (*pagination.Pagination, error) {

	if page <= 0 {
		page = 1
	}
	sorting := "ID desc"
	products, count := p.repo.GetAllProducts(ctx, pagesize, page, sorting)
	productsApi := ProductsToResponses(products)

	pagination := pagination.NewPagination(page, pagesize, count, sorting)
	pagination.Rows = productsApi

	return pagination, nil
}

func (p *productService) CreateProduct(ctx context.Context, product *api.Product) error {

	price := decimal.NewFromFloat(product.Price)
	var categoryID [16]byte
	categoryIDbArr, _ := hex.DecodeString(product.Category.ID) // category id to []byte. Then [16]byte
	copy(categoryID[:], categoryIDbArr)

	productToGo := models.Product{Name: product.Name, Description: product.Description, Price: price, StockNumber: int(product.StockNumber), StockCode: product.StockCode, CategoryID: product.Category.ID}
	err := p.repo.CreateProduct(ctx, &productToGo)
	if err != nil {
		return err
	}

	return nil
}

func (p *productService) GetProductByName(ctx context.Context, name *string) (*models.Product, error) {
	if len(*name) < 1 {
		return nil, errors.New("name cannot be empty")
	}

	product, err := p.repo.GetProductByName(ctx, name)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (p *productService) GetProductByID(ctx context.Context, id *string) (*models.Product, error) {
	if len(*id) < 1 {
		return nil, errors.New("id cannot be empty")
	}

	product, err := p.repo.GetProductByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return product, nil
}
