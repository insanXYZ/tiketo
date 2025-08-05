package dto

import "mime/multipart"

type CreateTicket struct {
	Name        string                `json:"name"`
	Description string                `json:"description"`
	Price       int                   `json:"price"`
	Amount      int                   `json:"amount"`
	ImageFile   *multipart.FileHeader `json:"-"`
}

type DeleteTicket struct {
	Id string `param:"id"`
}

type GetTicket struct {
	Id string `param:"id"`
}

type UpdateTicket struct {
	ID          string                `param:"id"`
	Name        string                `json:"name"`
	Description string                `json:"description"`
	Price       int                   `json:"price"`
	Amount      int                   `json:"amount"`
	ImageFile   *multipart.FileHeader `json:"-"`
}
