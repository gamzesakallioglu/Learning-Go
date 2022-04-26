package shoppingCart

import (
	"github.com/gamze.sakallioglu/learningGo/bitirme-projesi-gamzesakallioglu/internal/api"
	"github.com/gamze.sakallioglu/learningGo/bitirme-projesi-gamzesakallioglu/internal/customer"
	"github.com/gamze.sakallioglu/learningGo/bitirme-projesi-gamzesakallioglu/internal/models"
	"github.com/gamze.sakallioglu/learningGo/bitirme-projesi-gamzesakallioglu/internal/product"
)

func ShoppingCartsToResponses(pcs *models.ShoppingCarts) []*api.ShoppingCart {
	productCategories := make([]*api.ShoppingCart, 0)

	for _, pc := range *pcs {
		productCategories = append(productCategories, ShoppingCartToResponse(&pc))
	}

	return productCategories
}

func ShoppingCartToResponse(p *models.ShoppingCart) *api.ShoppingCart {

	return &api.ShoppingCart{
		ID:       p.ID.String(),
		Product:  product.ProductToResponse(&p.Product),
		Customer: customer.CustomerToResponse(&p.Customer),
		Quantity: int64(p.Quantity),
	}
}
