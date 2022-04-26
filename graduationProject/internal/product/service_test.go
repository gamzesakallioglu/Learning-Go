package product

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/gamze.sakallioglu/learningGo/bitirme-projesi-gamzesakallioglu/internal/models"
	"github.com/gamze.sakallioglu/learningGo/bitirme-projesi-gamzesakallioglu/pkg/pagination"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
)

var productName1 = "Dell PC"
var productName2 = "Gimcat Cat Food"
var productName3 = "Vitamin D"
var productName4 = "Table for 4"

var id1, _ = uuid.Parse("60b98def-cd93-4abc-aa1a-c31dc8b18172")
var id2, _ = uuid.Parse("2c5c2c18-de40-459a-bd19-787f46140e43")
var id3, _ = uuid.Parse("6baf7c15-25c3-4bce-af73-0d4cf5ca8545")
var id4, _ = uuid.Parse("d6ee7eb1-b113-469f-b831-047fd28ddccb")

var givenProduct1 = models.Product{

	ID:          id1,
	Name:        &productName1,
	StockNumber: 500,
	CreatedAt:   time.Now(),
}
var givenProduct2 = models.Product{
	ID:          id2,
	Name:        &productName2,
	StockNumber: 800,
	CreatedAt:   time.Now(),
}
var givenProduct3 = models.Product{
	ID:          id3,
	Name:        &productName3,
	StockNumber: 300,
	CreatedAt:   time.Now(),
}
var givenProduct4 = models.Product{
	ID:          id4,
	Name:        &productName4,
	StockNumber: 100,
	CreatedAt:   time.Now(),
}

var givenProducts = models.Products{
	givenProduct1,
	givenProduct2,
	givenProduct3,
	givenProduct4,
}
var givenProductSlice = givenProducts[2:4]
var givenProductSliceForSearch = givenProducts[1:2]

var productToresponse = ProductsToResponses(&givenProducts)
var productToresponseWithParam = ProductsToResponses(&givenProductSlice)
var productToresponseForSearch = ProductsToResponses(&givenProductSliceForSearch)

var pgForNoParam = pagination.Pagination{
	Page:       1,
	PageSize:   10,
	Sorting:    "ID desc",
	TotalRows:  4,
	TotalPages: 1,
	Rows:       productToresponse,
}
var pgWithParams = pagination.Pagination{
	Page:       2,
	PageSize:   2,
	Sorting:    "ID desc",
	TotalRows:  4,
	TotalPages: 2,
	Rows:       productToresponseWithParam,
}
var pgWithParamsForSearch = pagination.Pagination{
	Page:       1,
	PageSize:   1,
	Sorting:    "ID desc",
	TotalRows:  1,
	TotalPages: 1,
	Rows:       productToresponseForSearch,
}

func TestService(usecase *testing.T) {
	usecase.Run("Get All", func(t *testing.T) {

		mockRepo := &mockRepository{items: givenProducts}
		s := NewProductService(mockRepo)
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

					gotCats, err := s.GetAllProducts(ctx, tt.pageSize, tt.page)
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
		notEmptyArg := "Gimcat Cat Food"
		notEmptyArgFalse := "Some false product"
		mockRepo := &mockRepository{items: givenProducts}
		s := NewProductService(mockRepo)
		t.Run("GetAll Method Test", func(t *testing.T) {
			tests := []struct {
				name       string
				arg        *string
				wantReturn *models.Product
				wantErr    bool
			}{
				{name: "NoProductName", arg: &emptyArg, wantReturn: nil, wantErr: true},
				{name: "GetWithProductName", arg: &notEmptyArg, wantReturn: &givenProduct2, wantErr: false},
				{name: "GetWithFalseProductName", arg: &notEmptyArgFalse, wantReturn: nil, wantErr: true},
			}

			ctx := context.Background()
			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {

					gotCats, err := s.GetProductByName(ctx, tt.arg)
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
		mockRepo := &mockRepository{items: givenProducts}
		s := NewProductService(mockRepo)
		t.Run("Get Method Test", func(t *testing.T) {
			tests := []struct {
				name       string
				arg        *string
				wantReturn *models.Product
				wantErr    bool
			}{
				{name: "NoProductID", arg: &emptyArg, wantReturn: nil, wantErr: true},
				{name: "GetWithProductID", arg: &notEmptyArg, wantReturn: &givenProduct2, wantErr: false},
				{name: "GetWithFalseProductID", arg: &notEmptyArgFalse, wantReturn: nil, wantErr: true},
			}

			ctx := context.Background()
			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {

					gotCats, err := s.GetProductByID(ctx, tt.arg)
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

	usecase.Run("SearchByName", func(t *testing.T) {

		notEmptyArg := "GiMCat"
		mockRepo := &mockRepository{items: givenProducts}
		s := NewProductService(mockRepo)
		t.Run("Search Method Test", func(t *testing.T) {
			tests := []struct {
				name       string
				arg        *string
				wantReturn *pagination.Pagination
				wantErr    bool
			}{
				{name: "GetWithProductName", arg: &notEmptyArg, wantReturn: &pgWithParamsForSearch, wantErr: false},
			}

			ctx := context.Background()
			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {

					gotCats, err := s.SearchProducts(ctx, tt.arg, 1, 1)
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
		mockRepo := &mockRepository{items: givenProducts}
		s := NewProductService(mockRepo)
		t.Run("Delete Method Test", func(t *testing.T) {
			tests := []struct {
				name       string
				arg        *string
				wantReturn *models.Product
				wantErr    bool
			}{
				{name: "NoCategoryID", arg: &emptyArg, wantReturn: nil, wantErr: true},
				{name: "DeleteWithCategoryID", arg: &notEmptyArg, wantReturn: nil, wantErr: false},
				{name: "DeleteWithFalseCategoryID", arg: &notEmptyArgFalse, wantReturn: nil, wantErr: true},
			}

			ctx := context.Background()
			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {

					err := s.DeleteProductByID(ctx, tt.arg)
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
	items models.Products
}

func (m *mockRepository) CreateProduct(ctx context.Context, product *models.Product) error {

	return nil
}

func (m *mockRepository) GetProductByName(ctx context.Context, name *string) (*models.Product, error) {

	if len(*name) <= 0 {
		return nil, errors.New("name cannot be empty")
	}

	for _, category := range m.items {
		if *category.Name == *name {
			return &category, nil
		}
	}
	return nil, errors.New("product not found")
}

func (m *mockRepository) UpdateStockNumber(ctx context.Context, id *string, stockNumber int) error {
	if len(*id) <= 0 {
		return nil
	}

	for i, product := range m.items {

		if product.ID.String() == *id {
			product.StockNumber = stockNumber
			m.items[i] = product
			return nil
		}
	}
	return errors.New("product not found")
}

func (m *mockRepository) GetAllProducts(ctx context.Context, pagesize int, page int, sorting string) (*models.Products, int64) {

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

func (m *mockRepository) DeleteProductByID(ctx context.Context, id *string) error {

	var indexSearch int
	if len(*id) <= 0 {
		return errors.New("id cannot be empty")
	}

	for i, product := range m.items {
		if product.ID.String() == *id {
			indexSearch = i
		}
	}

	if indexSearch == 0 {
		return errors.New("record not found")
	}

	m.items = append(m.items[:indexSearch], m.items[indexSearch+1:len(m.items)]...)

	return nil
}
func (m *mockRepository) GetProductByID(ctx context.Context, id *string) (*models.Product, error) {

	if len(*id) <= 0 {
		return nil, errors.New("id cannot be empty")
	}

	for _, product := range m.items {

		if product.ID.String() == *id {
			return &product, nil
		}
	}
	return nil, errors.New("product not found")
}

func (m *mockRepository) SearchProducts(ctx context.Context, key *string, pageSize int, page int, sorting string) (*models.Products, int64) {

	if pageSize <= 0 {
		pageSize = 10
	}

	if page <= 0 {
		page = 1
	}

	offset := (page - 1) * pageSize
	until := offset + pageSize
	var productsToShow models.Products
	for _, product := range m.items {
		if strings.Contains(strings.ToLower(*product.Name), strings.ToLower(*key)) {
			productsToShow = append(productsToShow, product)
		}
	}
	if until > len(productsToShow) {
		until = len(productsToShow)
	}

	categories := productsToShow[offset:until]

	return &categories, int64(len(productsToShow))
}

func (m *mockRepository) UpdateProductByID(ctx context.Context, id *string, productUpdated *models.Product) error {

	if len(*id) <= 0 {
		return nil
	}

	for i, product := range m.items {

		if product.ID.String() == *id {
			product.Name = productUpdated.Name
			product.StockNumber = productUpdated.StockNumber
			m.items[i] = product
			return nil
		}
	}
	return errors.New("product not found")

}

func (m *mockRepository) Migration() {

}
