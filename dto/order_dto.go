package dto

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
