package converter

import (
	"tiketo/dto"
	"tiketo/entity"
	"tiketo/util"
	"time"
)

func TicketEntityToDto(ticket *entity.Ticket) *dto.Ticket {
	if ticket == nil {
		return nil
	}

	t := &dto.Ticket{
		ID:          ticket.ID,
		Name:        ticket.Name,
		Description: ticket.Description,
		Price:       ticket.Price,
		Image:       util.PathTicketImageDir + ticket.Image,
		Quantity:    ticket.Quantity,
		CreatedAt:   ticket.CreatedAt.Format(time.DateTime),
		User:        UserEntityToNameOnlyDto(ticket.User),
	}

	return t
}

func TicketEntitiesToDto(tickets []entity.Ticket) []dto.Ticket {
	if len(tickets) == 0 {
		return nil
	}

	lt := len(tickets)

	t := make([]dto.Ticket, len(tickets))

	for i := range lt {
		convert := TicketEntityToDto(&tickets[i])
		t[i] = *convert
	}

	return t
}
