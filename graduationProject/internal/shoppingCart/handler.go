package shoppingCart

import (
	"net/http"
	"sync"

	"github.com/gamze.sakallioglu/learningGo/bitirme-projesi-gamzesakallioglu/httpErrors"
	"github.com/gamze.sakallioglu/learningGo/bitirme-projesi-gamzesakallioglu/internal/api"
	mwCustomer "github.com/gamze.sakallioglu/learningGo/bitirme-projesi-gamzesakallioglu/internal/middleware/customer"
	"github.com/gamze.sakallioglu/learningGo/bitirme-projesi-gamzesakallioglu/pkg/config"
	"github.com/gamze.sakallioglu/learningGo/bitirme-projesi-gamzesakallioglu/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-openapi/strfmt"
)

type shoppingCartHandler struct {
	cfg     *config.Config
	service Service
}

func NewShoppingCartHandler(r *gin.RouterGroup, cfg *config.Config, service Service) {
	shoppingCartHandler := shoppingCartHandler{cfg: cfg, service: service}

	r.POST("/carts", mwCustomer.AuthMiddleware(cfg.JWTConfig.SecretKey), shoppingCartHandler.addProductToCart)
	r.PUT("/carts", mwCustomer.AuthMiddleware(cfg.JWTConfig.SecretKey), shoppingCartHandler.updateCartItem)
	r.GET("/carts", mwCustomer.AuthMiddleware(cfg.JWTConfig.SecretKey), shoppingCartHandler.listProductsInCart)
	r.DELETE("/carts/item/:id", mwCustomer.AuthMiddleware(cfg.JWTConfig.SecretKey), shoppingCartHandler.deleteItemFromCart)

}

func (p *shoppingCartHandler) updateCartItem(c *gin.Context) {

	user, err := utils.GetUserFromCtx(c)
	if err != nil {
		c.JSON(httpErrors.ErrorResponse(httpErrors.NewRestError(http.StatusBadRequest, err.Error(), nil)))
		return
	}
	if user == nil {
		c.JSON(httpErrors.ErrorResponse(httpErrors.NewRestError(http.StatusBadRequest, "user cannot found", nil)))
	}

	var shoppingCartItem api.ShoppingCart
	if err := c.Bind(&shoppingCartItem); err != nil {
		c.JSON(httpErrors.ErrorResponse(httpErrors.NewRestError(http.StatusBadRequest, "something's wrong! check your request body", nil)))
		return
	}

	format := strfmt.Default
	err = shoppingCartItem.Validate(format)
	if err != nil {
		c.JSON(httpErrors.ErrorResponse(httpErrors.NewRestError(http.StatusBadRequest, err.Error(), nil)))
		return
	}

	customer := api.Customer{ID: user.UserId}
	shoppingCartItem.Customer = &customer

	mux := &sync.Mutex{}
	err = p.service.UpdateCartItem(c.Request.Context(), &shoppingCartItem, mux)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, "Product has been updated")

}

func (p *shoppingCartHandler) deleteItemFromCart(c *gin.Context) {

	user, err := utils.GetUserFromCtx(c)
	if err != nil {
		c.JSON(httpErrors.ErrorResponse(httpErrors.NewRestError(http.StatusBadRequest, err.Error(), nil)))
		return
	}
	if user == nil {
		c.JSON(httpErrors.ErrorResponse(httpErrors.NewRestError(http.StatusBadRequest, "user cannot found", nil)))
	}

	productID := c.Param("id")
	if len(productID) <= 0 {
		c.JSON(httpErrors.ErrorResponse(httpErrors.NewRestError(http.StatusBadRequest, "Failed to delete product from cart. ID cannot be empty", nil)))
		return
	}

	err = p.service.DeleteItemFromCart(c.Request.Context(), productID, user.UserId)
	if err != nil {
		c.JSON(httpErrors.ErrorResponse(httpErrors.NewRestError(http.StatusBadRequest, err.Error(), nil)))
	}

	c.JSON(http.StatusOK, "product has been deleted from cart")
}

func (p *shoppingCartHandler) addProductToCart(c *gin.Context) {

	var shoppingCartItem api.ShoppingCart
	if err := c.Bind(&shoppingCartItem); err != nil {
		c.JSON(httpErrors.ErrorResponse(httpErrors.NewRestError(http.StatusBadRequest, "something's wrong! check your request body", nil)))
		return
	}

	format := strfmt.Default
	err := shoppingCartItem.Validate(format)
	if err != nil {
		c.JSON(httpErrors.ErrorResponse(httpErrors.NewRestError(http.StatusBadRequest, err.Error(), nil)))
		return
	}

	user, err := utils.GetUserFromCtx(c)
	if err != nil {
		c.JSON(httpErrors.ErrorResponse(httpErrors.NewRestError(http.StatusBadRequest, err.Error(), nil)))
		return
	}
	if user == nil {
		c.JSON(httpErrors.ErrorResponse(httpErrors.NewRestError(http.StatusBadRequest, "user cannot found", nil)))
	}

	// Customer will be added to shoppingCart object
	customer := api.Customer{ID: user.UserId}
	shoppingCartItem.Customer = &customer

	mux := &sync.Mutex{}
	err = p.service.AddItemToCart(c.Request.Context(), &shoppingCartItem, mux)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, "Product has been added to your shopping cart")

}

func (p *shoppingCartHandler) listProductsInCart(c *gin.Context) {

	user, err := utils.GetUserFromCtx(c)
	if err != nil {
		c.JSON(httpErrors.ErrorResponse(httpErrors.NewRestError(http.StatusBadRequest, err.Error(), nil)))
		return
	}
	if user == nil {
		c.JSON(httpErrors.ErrorResponse(httpErrors.NewRestError(http.StatusBadRequest, "user cannot found", nil)))
	}

	shoppingCartList, err := p.service.ListProductsInCart(c.Request.Context(), user.UserId)
	if err != nil {
		c.JSON(httpErrors.ErrorResponse(httpErrors.NewRestError(http.StatusBadRequest, err.Error(), nil)))
	}

	c.JSON(http.StatusOK, shoppingCartList)
}
