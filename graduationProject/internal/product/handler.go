package product

import (
	"net/http"

	"github.com/gamze.sakallioglu/learningGo/bitirme-projesi-gamzesakallioglu/httpErrors"
	"github.com/gamze.sakallioglu/learningGo/bitirme-projesi-gamzesakallioglu/internal/api"
	mwUserAdmin "github.com/gamze.sakallioglu/learningGo/bitirme-projesi-gamzesakallioglu/internal/middleware/userAdmin"
	"github.com/gamze.sakallioglu/learningGo/bitirme-projesi-gamzesakallioglu/pkg/config"
	"github.com/gamze.sakallioglu/learningGo/bitirme-projesi-gamzesakallioglu/pkg/pagination"
	"github.com/gin-gonic/gin"
	"github.com/go-openapi/strfmt"
)

type productHandler struct {
	cfg     *config.Config
	service Service
}

func NewProductHandler(r *gin.RouterGroup, cfg *config.Config, service Service) {
	productHandler := productHandler{cfg: cfg, service: service}

	r.POST("/products", mwUserAdmin.AuthMiddleware(cfg.JWTConfig.SecretKey), productHandler.createProduct)
	r.PUT("/products/:id", mwUserAdmin.AuthMiddleware(cfg.JWTConfig.SecretKey), productHandler.updateProductByID)
	r.GET("/products", productHandler.getAllProducts)
	r.GET("/products/:key", productHandler.searchProducts)
	r.DELETE("products/:id", mwUserAdmin.AuthMiddleware(cfg.JWTConfig.SecretKey), productHandler.deleteProductByID)

}

func (p *productHandler) updateProductByID(c *gin.Context) {

	id := c.Param("id")
	if len(id) <= 0 {
		c.JSON(httpErrors.ErrorResponse(httpErrors.NewRestError(http.StatusBadRequest, "id cannot be null", nil)))
		return
	}

	var productUpdated api.Product
	if err := c.Bind(&productUpdated); err != nil {
		c.JSON(httpErrors.ErrorResponse(httpErrors.NewRestError(http.StatusBadRequest, "something's wrong! check your request body", nil)))
		return
	}

	format := strfmt.Default
	err := productUpdated.Validate(format)
	if err != nil {
		c.JSON(httpErrors.ErrorResponse(httpErrors.NewRestError(http.StatusBadRequest, err.Error(), nil)))
		return
	}

	err = p.service.UpdateProductByID(c.Request.Context(), &id, &productUpdated)
	if err != nil {
		c.JSON(httpErrors.ErrorResponse(httpErrors.NewRestError(http.StatusBadRequest, err.Error(), nil)))
		return
	}
	c.JSON(http.StatusOK, "product has been updated")
}

func (p *productHandler) searchProducts(c *gin.Context) {

	page, pageSize := pagination.GetPaginationParametersFromRequest(c)
	key := c.Param("key")

	if len(key) <= 0 {
		c.JSON(httpErrors.ErrorResponse(httpErrors.NewRestError(http.StatusBadRequest, "key cannot be null", nil)))
		return
	}

	pg, err := p.service.SearchProducts(c.Request.Context(), &key, pageSize, page)
	if err != nil {
		c.JSON(httpErrors.ErrorResponse(httpErrors.NewRestError(http.StatusBadRequest, err.Error(), nil)))
		return
	}
	c.JSON(http.StatusOK, pg)
}
func (p *productHandler) deleteProductByID(c *gin.Context) {

	id := c.Param("id")
	if len(id) <= 0 {
		c.JSON(httpErrors.ErrorResponse(httpErrors.NewRestError(http.StatusBadRequest, "id cannot be null", nil)))
		return
	}

	err := p.service.DeleteProductByID(c.Request.Context(), &id)
	if err != nil {
		c.JSON(httpErrors.ErrorResponse(httpErrors.NewRestError(http.StatusBadRequest, err.Error(), nil)))
		return
	}
	c.JSON(http.StatusOK, "product has been deleted")
}

func (p *productHandler) getAllProducts(c *gin.Context) {

	page, pageSize := pagination.GetPaginationParametersFromRequest(c)

	pagination, err := p.service.GetAllProducts(c.Request.Context(), page, pageSize)
	if err != nil {
		c.JSON(httpErrors.ErrorResponse(httpErrors.NewRestError(http.StatusBadRequest, err.Error(), nil)))
		return
	}

	c.JSON(http.StatusOK, pagination)

}

func (p *productHandler) createProduct(c *gin.Context) {
	var product api.Product
	if err := c.Bind(&product); err != nil {
		c.JSON(httpErrors.ErrorResponse(httpErrors.NewRestError(http.StatusBadRequest, "something's wrong! check your request body", nil)))
		return
	}

	format := strfmt.Default
	err := product.Validate(format)
	if err != nil {
		c.JSON(httpErrors.ErrorResponse(httpErrors.NewRestError(http.StatusBadRequest, err.Error(), nil)))
		return
	}

	err = p.service.CreateProduct(c.Request.Context(), &product)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, "Product has been created")
}
