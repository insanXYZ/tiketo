package converter

import (
	"tiketo/dto"
	"tiketo/entity"
)

func TicketEntityToDto(ticket *entity.Ticket) *dto.Ticket {
	return &dto.Ticket{
		ID:          ticket.ID,
		Name:        ticket.Name,
		Description: ticket.Description,
		Price:       ticket.Price,
		Image:       "",
	}
}
