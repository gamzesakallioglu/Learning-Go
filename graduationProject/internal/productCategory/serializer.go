package productCategory

import (
	"github.com/gamze.sakallioglu/learningGo/bitirme-projesi-gamzesakallioglu/internal/api"
	"github.com/gamze.sakallioglu/learningGo/bitirme-projesi-gamzesakallioglu/internal/models"
)

func ProductCategoriesToResponses(pcs *models.ProductCategories) []*api.ProductCategory {
	productCategories := make([]*api.ProductCategory, 0)

	for _, pc := range *pcs {
		productCategories = append(productCategories, ProductCategoryToResponse(&pc))
	}

	return productCategories
}

func ProductCategoryToResponse(p *models.ProductCategory) *api.ProductCategory {

	return &api.ProductCategory{
		ID:          p.ID.String(),
		Name:        p.Name,
		Description: p.Description,
		IsParent:    p.IsParent,
		ParentID:    p.ParentID,
	}
}
