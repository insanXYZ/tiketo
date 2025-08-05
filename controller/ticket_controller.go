package controller

import (
	"tiketo/dto"
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

func (t *TicketController) CreateTicket(c echo.Context) error {
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
