package dto

import (
	"database/sql"
	"time"
)

type Order struct {
	ID          string       `json:"id,omitempty"`
	Status      string       `json:"status,omitempty"`
	Total       int          `json:"total,omitempty"`
	CreatedAt   time.Time    `json:"created_at"`
	PaidAt      sql.NullTime `json:"paid_at"`
	OrderDetail *OrderDetail `json:"order_detail,omitempty"`
}

type CreateOrder struct {
	TicketID string `json:"ticket_id" validate:"required"`
	Quantity int    `json:"quantity" validate:"required,min=1"`
}

type GetOrder struct {
	TicketID string `param:"id"`
}

type AfterPayment struct {
	TransactionStatus string `json:"transaction_status" validate:"required"`
	SignatureKey      string `json:"signature_key" validate:"required"`
	OrderId           string `json:"order_id" validate:"required"`
	StatusCode        string `json:"status_code" validate:"required"`
	GrossAmount       string `json:"gross_amount" validate:"required"`
}
