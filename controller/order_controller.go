package controller

import (
	"tiketo/dto"
	"tiketo/dto/message"
	"tiketo/service"
	"tiketo/util/httpresponse"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type OrderController struct {
	orderService *service.OrderService
}

func NewOrderController(orderService *service.OrderService) *OrderController {
	return &OrderController{
		orderService: orderService,
	}
}

func (o *OrderController) CreateOrder(c echo.Context) error {
	claims := c.Get("user").(jwt.MapClaims)
	req := new(dto.CreateOrder)

	err := c.Bind(req)
	if err != nil {
		return httpresponse.Error(c, message.ErrBind, err)
	}

	m, err := o.orderService.HandleCreate(c.Request().Context(), claims, req)
	if err != nil {
		return httpresponse.Error(c, message.ErrCreateOrder, err)
	}

	return httpresponse.Success(c, message.SuccessCreateTicket, m)
}
