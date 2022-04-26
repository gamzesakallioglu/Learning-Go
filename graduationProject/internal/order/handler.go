package order

import (
	"fmt"
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

type orderHandler struct {
	cfg     *config.Config
	service Service
}

func NewOrderHandler(r *gin.RouterGroup, cfg *config.Config, service Service) {
	orderHandler := orderHandler{cfg: cfg, service: service}

	r.POST("/orders", mwCustomer.AuthMiddleware(cfg.JWTConfig.SecretKey), orderHandler.completeOrder)
	r.GET("/orders/cancel/:orderID", mwCustomer.AuthMiddleware(cfg.JWTConfig.SecretKey), orderHandler.cancelOrder)
	r.GET("/orders", mwCustomer.AuthMiddleware(cfg.JWTConfig.SecretKey), orderHandler.listOrders)

}

func (p orderHandler) cancelOrder(c *gin.Context) {
	user, err := utils.GetUserFromCtx(c)
	if err != nil {
		c.JSON(httpErrors.ErrorResponse(httpErrors.NewRestError(http.StatusBadRequest, err.Error(), nil)))
		return
	}
	if user == nil {
		c.JSON(httpErrors.ErrorResponse(httpErrors.NewRestError(http.StatusBadRequest, "user cannot found", nil)))
	}

	orderID := c.Param("orderID")
	if len(orderID) <= 0 {
		c.JSON(httpErrors.ErrorResponse(httpErrors.NewRestError(http.StatusBadRequest, "Failed to cancel order. ID cannot be empty", nil)))
		return
	}

	err = p.service.CancelOrder(c.Request.Context(), orderID, user.UserId)
	if err != nil {
		c.JSON(httpErrors.ErrorResponse(httpErrors.NewRestError(http.StatusBadRequest, err.Error(), nil)))
		return
	}

	c.JSON(http.StatusOK, "order has been canceled")

}
func (p orderHandler) completeOrder(c *gin.Context) {

	var orderInfo api.CompleteOrder
	if err := c.Bind(&orderInfo); err != nil {
		c.JSON(httpErrors.ErrorResponse(httpErrors.NewRestError(http.StatusBadRequest, "something's wrong! check your request body", nil)))
		return
	}

	format := strfmt.Default
	err := orderInfo.Validate(format)
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

	mux := &sync.Mutex{}
	orderNumber, err := p.service.CompleteOrder(c.Request.Context(), &orderInfo, user.UserId, mux)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	message := fmt.Sprintf("Order has been completed. Order number: %s", orderNumber)
	c.JSON(http.StatusOK, message)
}

func (p orderHandler) listOrders(c *gin.Context) {

	user, err := utils.GetUserFromCtx(c)
	if err != nil {
		c.JSON(httpErrors.ErrorResponse(httpErrors.NewRestError(http.StatusBadRequest, err.Error(), nil)))
		return
	}
	if user == nil {
		c.JSON(httpErrors.ErrorResponse(httpErrors.NewRestError(http.StatusBadRequest, "user cannot found", nil)))
	}

	shoppingCartList, err := p.service.ListOrders(c.Request.Context(), user.UserId)
	if err != nil {
		c.JSON(httpErrors.ErrorResponse(httpErrors.NewRestError(http.StatusBadRequest, err.Error(), nil)))
	}

	c.JSON(http.StatusOK, shoppingCartList)

}
