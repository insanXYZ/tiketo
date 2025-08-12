package controller

import (
	"tiketo/dto"
	"tiketo/dto/converter"
	"tiketo/dto/message"
	"tiketo/middleware"
	"tiketo/service"
	"tiketo/util/httpresponse"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type TicketController struct {
	ticketService *service.TicketService
}

func NewTicketController(ticketService *service.TicketService) *TicketController {
	return &TicketController{
		ticketService: ticketService,
	}
}

func (t *TicketController) RegisterRoutes(c *echo.Group) {
	c.GET("/tickets/:id", t.GetTicket)
	c.GET("/tickets", t.GetTickets)

	hasAcc := c.Group("/me", middleware.HasAccToken)
	hasAcc.GET("/tickets", t.GetUserTickets)
	hasAcc.POST("/tickets", t.CreateTicket)
	hasAcc.DELETE("/tickets/:id", t.DeleteTicket)
	hasAcc.PUT("/tickets/:id", t.UpdateTicket)
}

func (t *TicketController) GetUserTickets(c echo.Context) error {
	claims := c.Get("user").(jwt.MapClaims)

	tickets, err := t.ticketService.HandleGetUserTickets(c.Request().Context(), claims)
	if err != nil {
		return httpresponse.Error(c, message.ErrGetUserTickets, err)
	}

	return httpresponse.Success(c, message.SuccessGetUserTickets, converter.TicketEntitiesToDto(tickets))
}

func (t *TicketController) GetTicket(c echo.Context) error {
	req := new(dto.GetTicket)
	err := c.Bind(req)
	if err != nil {
		return httpresponse.Error(c, message.ErrBind, err)
	}

	ticket, err := t.ticketService.HandleGet(c.Request().Context(), req)
	if err != nil {
		return httpresponse.Error(c, message.ErrGetTicket, err)
	}
	return httpresponse.Success(c, message.SuccessGetTicket, converter.TicketEntityToDto(ticket))
}

func (t *TicketController) GetTickets(c echo.Context) error {
	tickets, err := t.ticketService.HandleGetTickets(c.Request().Context())
	if err != nil {
		return httpresponse.Error(c, message.ErrGetTickets, err)
	}

	return httpresponse.Success(c, message.SuccessGetTickets, converter.TicketEntitiesToDto(tickets))
}

func (t *TicketController) CreateTicket(c echo.Context) error {
	claims := c.Get("user").(jwt.MapClaims)
	req := new(dto.CreateTicket)

	err := c.Bind(req)
	if err != nil {
		return httpresponse.Error(c, message.ErrBind, err)
	}

	fileHeader, err := c.FormFile("image")
	if err != nil {
		return httpresponse.Error(c, message.ErrGetFormFile, err)
	}

	req.ImageFile = fileHeader

	err = t.ticketService.HandleCreateTicket(c.Request().Context(), claims, req)
	if err != nil {
		return httpresponse.Error(c, message.ErrCreateTicket, err)
	}

	return httpresponse.Success(c, message.SuccessCreateTicket, nil)
}

func (t *TicketController) DeleteTicket(c echo.Context) error {
	claims := c.Get("user").(jwt.MapClaims)
	req := new(dto.DeleteTicket)

	err := c.Bind(req)
	if err != nil {
		return httpresponse.Error(c, message.ErrBind, err)
	}

	err = t.ticketService.HandleDelete(c.Request().Context(), claims, req)
	if err != nil {
		return httpresponse.Error(c, message.ErrDeleteTicket, err)
	}

	return httpresponse.Success(c, message.SuccessDeleteTicket, nil)
}

func (t *TicketController) UpdateTicket(c echo.Context) error {
	claims := c.Get("user").(jwt.MapClaims)
	req := new(dto.UpdateTicket)

	err := c.Bind(req)
	if err != nil {
		return httpresponse.Error(c, message.ErrBind, err)
	}

	err = t.ticketService.HandleUpdate(c.Request().Context(), claims, req)
	if err != nil {
		return httpresponse.Error(c, message.ErrUpdateTicket, err)
	}

	return httpresponse.Success(c, message.SuccessUpdateTicket, nil)
}
