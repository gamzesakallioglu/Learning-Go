package productCategory

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/gamze.sakallioglu/learningGo/bitirme-projesi-gamzesakallioglu/internal/api"
	"github.com/gamze.sakallioglu/learningGo/bitirme-projesi-gamzesakallioglu/internal/models"
	"github.com/gamze.sakallioglu/learningGo/bitirme-projesi-gamzesakallioglu/pkg/pagination"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
)

var catName1 = "Style & Fashion"
var catName2 = "Electronics"
var catName3 = "Furniture"
var catName4 = "Health & Wellness"
var catName5 = "Pet Supplies"
var catName6 = "Pet Clothing"
var catName7 = "Vitamins"

var id1, _ = uuid.Parse("60b98def-cd93-4abc-aa1a-c31dc8b18172")
var id2, _ = uuid.Parse("2c5c2c18-de40-459a-bd19-787f46140e43")
var id3, _ = uuid.Parse("6baf7c15-25c3-4bce-af73-0d4cf5ca8545")
var id4, _ = uuid.Parse("d6ee7eb1-b113-469f-b831-047fd28ddccb")
var id5, _ = uuid.Parse("52804d99-a33c-4d8e-9603-869bc88901a3")
var id6, _ = uuid.Parse("123e4567-e89b-12d3-a456-426614174000")
var id7, _ = uuid.Parse("4be2f5e7-833b-4f46-8ce9-02a3f39bc74b")

var givenCategory1 = models.ProductCategory{

	ID:        id1,
	Name:      &catName1,
	CreatedAt: time.Now(),
}
var givenCategory2 = models.ProductCategory{
	ID:        id2,
	Name:      &catName2,
	CreatedAt: time.Now(),
}
var givenCategory3 = models.ProductCategory{
	ID:        id3,
	Name:      &catName3,
	CreatedAt: time.Now(),
}
var givenCategory4 = models.ProductCategory{
	ID:        id4,
	Name:      &catName4,
	CreatedAt: time.Now(),
}
var givenCategory5 = models.ProductCategory{
	ID:        id5,
	Name:      &catName5,
	CreatedAt: time.Now(),
}
var givenCategory6 = models.ProductCategory{
	ID:        id6,
	Name:      &catName6,
	CreatedAt: time.Now(),
}
var givenCategory7 = models.ProductCategory{
	ID:        id7,
	Name:      &catName7,
	CreatedAt: time.Now(),
}

var givenCategories = models.ProductCategories{
	givenCategory1,
	givenCategory2,
	givenCategory3,
	givenCategory4,
	givenCategory5,
	givenCategory6,
	givenCategory7,
}
var givenCatSlice = givenCategories[2:4]

var catToresponse = ProductCategoriesToResponses(&givenCategories)
var catToresponseWithParam = ProductCategoriesToResponses(&givenCatSlice)

var pgForNoParam = pagination.Pagination{
	Page:       1,
	PageSize:   10,
	Sorting:    "ID desc",
	TotalRows:  7,
	TotalPages: 1,
	Rows:       catToresponse,
}
var pgWithParams = pagination.Pagination{
	Page:       2,
	PageSize:   2,
	Sorting:    "ID desc",
	TotalRows:  7,
	TotalPages: 4,
	Rows:       catToresponseWithParam,
}

func TestService(usecase *testing.T) {
	usecase.Run("Get All", func(t *testing.T) {

		mockRepo := &mockRepository{items: givenCategories}
		s := NewProductCategoryService(mockRepo)
		t.Run("GetAll Method Test", func(t *testing.T) {
			tests := []struct {
				name       string
				page       int
				pageSize   int
				wantReturn *pagination.Pagination
				wantErr    bool
			}{
				{name: "GetAllWithoutParams", page: 0, pageSize: 0, wantReturn: &pgForNoParam, wantErr: false},
				{name: "GetAllWithPageAndPageSize", page: 2, pageSize: 2, wantReturn: &pgWithParams, wantErr: false},
			}

			ctx := context.Background()
			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {

					gotCats, err := s.GetAllProductCategories(ctx, tt.pageSize, tt.page)
					if (err != nil) != tt.wantErr {
						t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
						return
					}
					if diff := cmp.Diff(gotCats, tt.wantReturn); diff != "" {
						t.Errorf("Get() mismatch (-want +got):\n%s", diff)
					}
				})
			}
		})
	})

	usecase.Run("GetByName", func(t *testing.T) {

		emptyArg := ""
		notEmptyArg := "Electronics"
		notEmptyArgFalse := "Some false category"
		givenCategory2Api := ProductCategoryToResponse(&givenCategory2)
		mockRepo := &mockRepository{items: givenCategories}
		s := NewProductCategoryService(mockRepo)
		t.Run("GetAll Method Test", func(t *testing.T) {
			tests := []struct {
				name       string
				arg        *string
				wantReturn *api.ProductCategory
				wantErr    bool
			}{
				{name: "NoCategoryName", arg: &emptyArg, wantReturn: nil, wantErr: true},
				{name: "GetWithCategoryName", arg: &notEmptyArg, wantReturn: givenCategory2Api, wantErr: false},
				{name: "GetWithFalseCategoryName", arg: &notEmptyArgFalse, wantReturn: nil, wantErr: true},
			}

			ctx := context.Background()
			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {

					gotCats, err := s.GetProductCategoryByName(ctx, tt.arg)
					if (err != nil) != tt.wantErr {
						t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
						return
					}
					if diff := cmp.Diff(gotCats, tt.wantReturn); diff != "" {
						t.Errorf("Get() mismatch (-want +got):\n%s", diff)
					}
				})
			}
		})
	})

	usecase.Run("GetByID", func(t *testing.T) {

		emptyArg := ""
		notEmptyArg := "2c5c2c18-de40-459a-bd19-787f46140e43"
		notEmptyArgFalse := "e475593a-a365-4070-9d6b-01ff5edaf823"
		givenCategory2Api := ProductCategoryToResponse(&givenCategory2)
		mockRepo := &mockRepository{items: givenCategories}
		s := NewProductCategoryService(mockRepo)
		t.Run("GetAll Method Test", func(t *testing.T) {
			tests := []struct {
				name       string
				arg        *string
				wantReturn *api.ProductCategory
				wantErr    bool
			}{
				{name: "NoCategoryID", arg: &emptyArg, wantReturn: nil, wantErr: true},
				{name: "GetWithCategoryID", arg: &notEmptyArg, wantReturn: givenCategory2Api, wantErr: false},
				{name: "GetWithFalseCategoryID", arg: &notEmptyArgFalse, wantReturn: nil, wantErr: true},
			}

			ctx := context.Background()
			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {

					gotCats, err := s.GetProductCategoryByID(ctx, tt.arg)
					if (err != nil) != tt.wantErr {
						t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
						return
					}
					if diff := cmp.Diff(gotCats, tt.wantReturn); diff != "" {
						t.Errorf("Get() mismatch (-want +got):\n%s", diff)
					}
				})
			}
		})
	})

	usecase.Run("DeleteByID", func(t *testing.T) {

		emptyArg := ""
		notEmptyArg := "2c5c2c18-de40-459a-bd19-787f46140e43"
		notEmptyArgFalse := "e475593a-a365-4070-9d6b-01ff5edaf823"
		mockRepo := &mockRepository{items: givenCategories}
		s := NewProductCategoryService(mockRepo)
		t.Run("GetAll Method Test", func(t *testing.T) {
			tests := []struct {
				name       string
				arg        *string
				wantReturn *api.ProductCategory
				wantErr    bool
			}{
				{name: "NoCategoryID", arg: &emptyArg, wantReturn: nil, wantErr: true},
				{name: "DeleteWithCategoryID", arg: &notEmptyArg, wantReturn: nil, wantErr: false},
				{name: "DeleteWithFalseCategoryID", arg: &notEmptyArgFalse, wantReturn: nil, wantErr: true},
			}

			ctx := context.Background()
			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {

					err := s.DeleteProductCategoryByID(ctx, tt.arg)
					if (err != nil) != tt.wantErr {
						t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
						return
					}
				})
			}
		})
	})

}

type mockRepository struct {
	items models.ProductCategories
}

func (m *mockRepository) CreateProductCategoryForBulk(ctx context.Context, category *models.ProductCategory) error {

	return nil
}

func (m *mockRepository) GetProductCategoryByName(ctx context.Context, name *string) (*models.ProductCategory, error) {

	if len(*name) <= 0 {
		return nil, errors.New("name cannot be empty")
	}

	for _, category := range m.items {
		if *category.Name == *name {
			return &category, nil
		}
	}
	return nil, errors.New("category not found")
}

func (m *mockRepository) GetAllProductCategories(ctx context.Context, pagesize int, page int, sorting string) (*models.ProductCategories, int64) {

	if pagesize <= 0 {
		pagesize = 10
	}

	if page <= 0 {
		page = 1
	}

	offset := (page - 1) * pagesize
	until := offset + pagesize
	if until > len(m.items) {
		until = len(m.items)
	}

	categories := m.items[offset:until]

	return &categories, int64(len(m.items))
}

func (m *mockRepository) DeleteProductCategoryByID(ctx context.Context, id *string) error {

	var indexSearch int
	if len(*id) <= 0 {
		return errors.New("id cannot be empty")
	}

	for i, category := range m.items {
		if category.ID.String() == *id {
			indexSearch = i
		}
	}

	if indexSearch == 0 {
		return errors.New("record not found")
	}

	m.items = append(m.items[:indexSearch], m.items[indexSearch+1:len(m.items)]...)

	return nil
}
func (m *mockRepository) GetProductCategoryByID(ctx context.Context, id *string) (*models.ProductCategory, error) {

	if len(*id) <= 0 {
		return nil, errors.New("id cannot be empty")
	}

	for _, category := range m.items {

		if category.ID.String() == *id {
			return &category, nil
		}
	}
	return nil, errors.New("category not found")
}

func (m *mockRepository) Migration() {

}
