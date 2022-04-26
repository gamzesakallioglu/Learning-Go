package product

import (
	"github.com/gamze.sakallioglu/learningGo/bitirme-projesi-gamzesakallioglu/internal/api"
	"github.com/gamze.sakallioglu/learningGo/bitirme-projesi-gamzesakallioglu/internal/models"
	"github.com/gamze.sakallioglu/learningGo/bitirme-projesi-gamzesakallioglu/internal/productCategory"
	"github.com/shopspring/decimal"
)

func ProductsToResponses(pcs *models.Products) []*api.Product {
	products := make([]*api.Product, 0)

	for _, pc := range *pcs {
		products = append(products, ProductToResponse(&pc))
	}

	return products
}

func ProductToResponse(p *models.Product) *api.Product {

	price, _ := p.Price.Float64()

	return &api.Product{
		ID:          p.ID.String(),
		Name:        p.Name,
		Description: p.Description,
		Price:       price,
		StockCode:   p.StockCode,
		StockNumber: int64(p.StockNumber),
		Category:    productCategory.ProductCategoryToResponse(&p.Category),
	}
}

func ResponseToProduct(p *api.Product) *models.Product {

	price := decimal.NewFromFloat(p.Price)

	return &models.Product{Name: p.Name, Description: p.Description, Price: price, StockNumber: int(p.StockNumber), StockCode: p.StockCode, CategoryID: p.Category.ID}

}
