package repository

import (
	"tiketo/entity"
)

type TicketRepository struct {
	Repository[*entity.Ticket]
}

func NewTicketRepository() *TicketRepository {
	return &TicketRepository{}
}
