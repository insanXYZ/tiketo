package dto

type CreateOrder struct {
	TicketID string `json:"ticket_id" validate:"required"`
	Quantity int    `json:"quantity" validate:"required,min=1"`
}
