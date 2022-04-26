package productCategory

import (
	"context"
	"errors"

	"github.com/gamze.sakallioglu/learningGo/bitirme-projesi-gamzesakallioglu/internal/api"
	"github.com/gamze.sakallioglu/learningGo/bitirme-projesi-gamzesakallioglu/internal/models"
	"github.com/gamze.sakallioglu/learningGo/bitirme-projesi-gamzesakallioglu/pkg/pagination"
)

type productCategoryService struct {
	repo Repository
}

type Service interface {
	CreateProductCategoryForBulk(ctx context.Context, category *api.ProductCategory) error
	GetProductCategoryByName(ctx context.Context, name *string) (*api.ProductCategory, error)
	GetAllProductCategories(ctx context.Context, page int, pageSize int) (*pagination.Pagination, error)
	DeleteProductCategoryByID(ctx context.Context, id *string) error
	GetProductCategoryByID(ctx context.Context, id *string) (*api.ProductCategory, error)
}

func NewProductCategoryService(repo Repository) Service {
	if repo == nil {
		return nil
	}

	return &productCategoryService{repo: repo}
}

func (p *productCategoryService) DeleteProductCategoryByID(ctx context.Context, id *string) error {

	cat, err := p.GetProductCategoryByID(ctx, id)
	if err != nil {
		return err
	}
	if cat == nil {
		return errors.New("category doesn't exist")
	}

	err = p.repo.DeleteProductCategoryByID(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (p *productCategoryService) GetProductCategoryByID(ctx context.Context, id *string) (*api.ProductCategory, error) {

	cat, err := p.repo.GetProductCategoryByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if cat == nil {
		return nil, errors.New("category does not exist")
	}

	catApi := ProductCategoryToResponse(cat)
	return catApi, nil
}

func (p *productCategoryService) CreateProductCategoryForBulk(ctx context.Context, category *api.ProductCategory) error {

	categoryToGo := models.ProductCategory{Name: category.Name, Description: category.Description, IsParent: category.IsParent, ParentID: category.ParentID}
	err := p.repo.CreateProductCategoryForBulk(ctx, &categoryToGo)
	if err != nil {
		return err
	}

	return nil
}

func (p *productCategoryService) GetProductCategoryByName(ctx context.Context, name *string) (*api.ProductCategory, error) {
	if len(*name) < 1 {
		return nil, errors.New("name cannot be empty")
	}

	category, err := p.repo.GetProductCategoryByName(ctx, name)
	if err != nil {
		return nil, err
	}

	categoryApi := ProductCategoryToResponse(category)

	return categoryApi, nil
}

func (p *productCategoryService) GetAllProductCategories(ctx context.Context, page int, pagesize int) (*pagination.Pagination, error) {

	sorting := "ID desc"
	if page <= 0 {
		page = 1
	}
	productCategories, count := p.repo.GetAllProductCategories(ctx, pagesize, page, sorting)
	productCategoriesApi := ProductCategoriesToResponses(productCategories)

	pagination := pagination.NewPagination(page, pagesize, count, sorting)
	pagination.Rows = productCategoriesApi

	return pagination, nil
}
