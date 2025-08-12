package controller

import (
	"tiketo/dto"
	"tiketo/dto/converter"
	"tiketo/dto/message"
	"tiketo/middleware"
	"tiketo/service"
	"tiketo/util/httpresponse"
	"tiketo/util/logger"

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
	e.POST("/orders/callback", o.AfterPayment)
	e.POST("/orders", o.CreateOrder, middleware.HasAccToken)
	e.GET("/me/orders", o.GetHistoryOrders, middleware.HasAccToken)
	e.GET("/me/orders/:id", o.GetHistoryOrder, middleware.HasAccToken)
}

func (o *OrderController) AfterPayment(c echo.Context) error {
	logger.EnteringMethod("OrderController.AfterPayment")

	req := new(dto.AfterPayment)
	err := c.Bind(req)
	if err != nil {
		logger.Warn(nil, "Err bind :", err.Error())
		return httpresponse.Error(c, message.ErrBind, err)
	}

	err = o.orderService.HandleAfterPayment(c.Request().Context(), req)
	if err != nil {
		logger.Warn(nil, "Err HandleAfterPayment :", err.Error())
		return httpresponse.Error(c, message.ErrPaymentProcess, err)
	}

	return httpresponse.Success(c, message.SuccessPaymentProcess, nil)
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

	return httpresponse.Success(c, message.SuccessGetHistoryOrder, converter.OrderEntityToDto(order))
}

func (o *OrderController) GetHistoryOrders(c echo.Context) error {
	claims := c.Get("user").(jwt.MapClaims)

	orders, err := o.orderService.HandleGetHistoryOrders(c.Request().Context(), claims)
	if err != nil {
		return httpresponse.Error(c, message.ErrGetHistoryOrders, err)
	}

	return httpresponse.Success(c, message.SuccessGetHistoryOrders, converter.OrderEntitiesToDto(orders))
}

func (o *OrderController) CreateOrder(c echo.Context) error {
	claims := c.Get("user").(jwt.MapClaims)
	req := new(dto.CreateOrder)

	err := c.Bind(req)
	if err != nil {
		return httpresponse.Error(c, message.ErrBind, err)
	}

	snapResponse, err := o.orderService.HandleCreate(c.Request().Context(), claims, req)
	if err != nil {
		return httpresponse.Error(c, message.ErrCreateOrder, err)
	}

	return httpresponse.Success(c, message.SuccessCreateTicket, jwt.MapClaims{
		"token":        snapResponse.Token,
		"redirect_url": snapResponse.RedirectURL,
	})
}
