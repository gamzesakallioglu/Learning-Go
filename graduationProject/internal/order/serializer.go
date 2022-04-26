package order

import (
	"time"

	"github.com/gamze.sakallioglu/learningGo/bitirme-projesi-gamzesakallioglu/internal/api"
	"github.com/gamze.sakallioglu/learningGo/bitirme-projesi-gamzesakallioglu/internal/models"
	"github.com/gamze.sakallioglu/learningGo/bitirme-projesi-gamzesakallioglu/internal/product"
	"github.com/shopspring/decimal"
)

func OrderToOrderHeaderAndDetail(order *api.Order) (*models.OrderHeader, []*models.OrderDetail) {

	var orderDetails []*models.OrderDetail
	var total = int64(0)
	totalPrice := decimal.NewFromInt(total)

	for _, product := range order.Products {

		pricePerItem := decimal.NewFromFloat(product.PricePerItem)
		totalForItem := product.PricePerItem * float64(product.Quantity)
		totalPrice.Add(decimal.NewFromFloat(totalForItem))

		orderDetails = append(orderDetails, &models.OrderDetail{
			ProductID:    product.Product.ID,
			Quantity:     int(product.Quantity),
			PricePerItem: pricePerItem,
		})
	}

	orderHeader := &models.OrderHeader{
		PhoneNumber:    order.PhoneNumber,
		Address:        order.Address,
		CustomerID:     &order.Customer.ID,
		OrderStatus:    "awaiting shipment",
		PaymentStatus:  "completed",
		PaymentDueDate: time.Now().Add(24 * time.Hour),
		OrderDate:      time.Now(),
		PaymentDate:    time.Now(),
		OrderTotal:     totalPrice,
	}

	return orderHeader, orderDetails

}

func ShoppingCartToOrderHeaderAndDetail(orderInfo *api.CompleteOrder, customerID string, shoppingCarts *models.ShoppingCarts) (*models.OrderHeader, []*models.OrderDetail) {

	var orderDetails []*models.OrderDetail
	var total = int64(0)
	totalPrice := decimal.NewFromInt(total)

	for _, shoppingCart := range *shoppingCarts {

		product := shoppingCart.Product

		pricePerItem := product.Price
		totalForItem := product.Price.InexactFloat64() * float64(shoppingCart.Quantity)
		totalPrice = totalPrice.Add(decimal.NewFromFloat(totalForItem))

		orderDetails = append(orderDetails, &models.OrderDetail{
			ProductID:    product.ID.String(),
			Quantity:     int(shoppingCart.Quantity),
			PricePerItem: pricePerItem,
		})
	}

	orderHeader := &models.OrderHeader{
		PhoneNumber:    orderInfo.PhoneNumber,
		Address:        orderInfo.Address,
		CustomerID:     &customerID,
		OrderStatus:    "awaiting shipment",
		PaymentStatus:  "completed",
		PaymentType:    *&orderInfo.PaymentType,
		PaymentDueDate: time.Now().Add(24 * time.Hour),
		OrderDate:      time.Now(),
		PaymentDate:    time.Now(),
		OrderTotal:     totalPrice,
	}

	return orderHeader, orderDetails

}

func OrderToResponse(orderHeader *models.OrderHeader) *api.Order {

	var products []*api.OrderProduct
	for _, orderDetail := range orderHeader.OrderDetails {
		products = append(products, orderDetailToOrderProduct(&orderDetail))
	}
	return &api.Order{
		ID:             orderHeader.ID.String(),
		Address:        orderHeader.Address,
		PhoneNumber:    orderHeader.PhoneNumber,
		OrderNumber:    orderHeader.OrderNumber,
		OrderStatus:    orderHeader.OrderStatus,
		OrderTotal:     orderHeader.OrderTotal.InexactFloat64(),
		PaymentDueDate: orderHeader.PaymentDueDate.Format("2006-01-02 15:04:05"),
		OrderDate:      orderHeader.OrderDate.Format("2006-01-02 15:04:05"),
		PaymentDate:    orderHeader.PaymentDate.Format("2006-01-02 15:04:05"),
		ReceiveDate:    orderHeader.ReceiveDate.Format("2006-01-02 15:04:05"),
		ShippingDate:   orderHeader.ShippingDate.Format("2006-01-02 15:04:05"),
		CancelDate:     orderHeader.CancelDate.String(),
		PaymentStatus:  orderHeader.PaymentStatus,
		PaymentType:    orderHeader.PaymentType,
		Products:       products,
	}
}

func OrdersToResponse(orderHeaders []*models.OrderHeader) []*api.Order {

	var ordersApi []*api.Order

	for _, orderHeader := range orderHeaders {
		ordersApi = append(ordersApi, OrderToResponse(orderHeader))
	}

	return ordersApi

}

func orderDetailToOrderProduct(orderDetail *models.OrderDetail) *api.OrderProduct {

	return &api.OrderProduct{
		Product:      product.ProductToResponse(&orderDetail.Product),
		Quantity:     int64(orderDetail.Quantity),
		PricePerItem: orderDetail.PricePerItem.InexactFloat64(),
	}
}
