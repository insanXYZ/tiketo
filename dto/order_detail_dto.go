package dto

type OrderDetail struct {
	Quantity uint    `json:"quantity,omitempty"`
	Ticket   *Ticket `json:"ticket,omitempty"`
}
