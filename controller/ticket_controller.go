package controller

import (
	"tiketo/dto"
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

func (t *TicketController) RegisterRoutes(c *echo.Echo) {
	c.GET("/tickets/:id", t.Get)
	c.GET("/tickets", t.GetAll)

	hasAcc := c.Group("/me", middleware.HasAccToken)
	hasAcc.GET("/tickets", t.GetUserTickets)
	hasAcc.POST("/tickets", t.Create)
	hasAcc.DELETE("/tickets/:id", t.Delete)
	hasAcc.PUT("/tickets/:id", t.Update)
}

func (t *TicketController) GetUserTickets(c echo.Context) error {
	claims := c.Get("user").(jwt.MapClaims)

	tickets, err := t.ticketService.HandleGetUserTickets(c.Request().Context(), claims)
	if err != nil {
		return httpresponse.Error(c, "failed get user tickets", err)
	}

	return httpresponse.Success(c, "success get user tickets", tickets)
}

func (t *TicketController) Get(c echo.Context) error {
	req := new(dto.GetTicket)
	err := c.Bind(req)
	if err != nil {
		return httpresponse.Error(c, "error bind", err)
	}

	ticket, err := t.ticketService.HandleGet(c.Request().Context(), req)
	if err != nil {
		return httpresponse.Error(c, "failed get ticket", err)
	}
	return httpresponse.Success(c, "success get ticket", ticket)
}

func (t *TicketController) GetAll(c echo.Context) error {
	tickets, err := t.ticketService.HandleGetAll(c.Request().Context())
	if err != nil {
		return httpresponse.Error(c, "failed get all tickets", err)
	}

	return httpresponse.Success(c, "success get all tickets", tickets)
}

func (t *TicketController) Create(c echo.Context) error {
	claims := c.Get("user").(jwt.MapClaims)
	req := new(dto.CreateTicket)

	err := c.Bind(req)
	if err != nil {
		return httpresponse.Error(c, "error bind", err)
	}

	fileHeader, err := c.FormFile("image")
	if err != nil {
		return httpresponse.Error(c, "error get file", err)
	}

	req.ImageFile = fileHeader

	err = t.ticketService.HandleCreateTicket(c.Request().Context(), claims, req)
	if err != nil {
		return httpresponse.Error(c, "failed create ticket", err)
	}

	return httpresponse.Success(c, "success create ticket", nil)
}

func (t *TicketController) Delete(c echo.Context) error {
	claims := c.Get("user").(jwt.MapClaims)
	req := new(dto.DeleteTicket)

	err := c.Bind(req)
	if err != nil {
		return httpresponse.Error(c, "error binding", err)
	}

	err = t.ticketService.HandleDelete(c.Request().Context(), claims, req)
	if err != nil {
		return httpresponse.Error(c, "failed delete ticket", err)
	}

	return httpresponse.Success(c, "success delete ticket", nil)
}

func (t *TicketController) Update(c echo.Context) error {
	claims := c.Get("user").(jwt.MapClaims)
	req := new(dto.UpdateTicket)

	err := c.Bind(req)
	if err != nil {
		return httpresponse.Error(c, "error bind", err)
	}

	err = t.ticketService.HandleUpdate(c.Request().Context(), claims, req)
	if err != nil {
		return httpresponse.Error(c, "failed update ticket", err)
	}

	return httpresponse.Success(c, "success update ticket", nil)
}
