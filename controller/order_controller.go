package controller

import (
	"tiketo/dto"
	"tiketo/dto/message"
	"tiketo/middleware"
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

func (o *OrderController) RegisterRoutes(e *echo.Group) {
	e.POST("/orders", o.CreateOrder, middleware.HasAccToken)
	e.GET("/me/orders", o.GetHistoryOrders, middleware.HasAccToken)
	e.GET("/me/orders/:id", o.GetHistoryOrder, middleware.HasAccToken)
}

func (o *OrderController) GetHistoryOrder(c echo.Context) error {
	claims := c.Get("user").(jwt.MapClaims)
	req := new(dto.GetOrder)

	err := c.Bind(req)
	if err != nil {
		return httpresponse.Error(c, message.ErrBind, err)
	}

	order, err := o.orderService.HandleGetHistoryOrder(c.Request().Context(), claims, req)
	if err != nil {
		return httpresponse.Error(c, message.ErrGetHistoryOrder, err)
	}

	return httpresponse.Success(c, message.SuccessGetHistoryOrder, order)
}

func (o *OrderController) GetHistoryOrders(c echo.Context) error {
	claims := c.Get("user").(jwt.MapClaims)

	orders, err := o.orderService.HandleGetHistoryOrders(c.Request().Context(), claims)
	if err != nil {
		return httpresponse.Error(c, message.ErrGetHistoryOrders, err)
	}

	return httpresponse.Success(c, message.SuccessGetHistoryOrders, orders)
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
