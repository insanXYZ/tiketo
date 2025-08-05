package dto

import "mime/multipart"

type CreateTicket struct {
	Name        string                `json:"name,omitempty"`
	Description string                `json:"description,omitempty"`
	Price       int                   `json:"price,omitempty"`
	Amount      int                   `json:"amount,omitempty"`
	ImageFile   *multipart.FileHeader `json:"-"`
}
