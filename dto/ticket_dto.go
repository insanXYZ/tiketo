package dto

import (
	"mime/multipart"
)

type Ticket struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       int    `json:"price"`
	Image       string `json:"image"`
	Amount      int    `json:"amount"`
	CreatedAt   string `json:"created_at"`
	User        *User  `json:"created_by"`
}

type CreateTicket struct {
	Name        string                `json:"name" validate:"required,min=3"`
	Description string                `json:"description" validate:"required,min=3"`
	Price       int                   `json:"price" validate:"required"`
	Amount      int                   `json:"amount" validate:"required"`
	ImageFile   *multipart.FileHeader `json:"-" validate:"isImage"`
}

type DeleteTicket struct {
	Id string `param:"id" validate:"required"`
}

type GetTicket struct {
	Id string `param:"id" validate:"required"`
}

type UpdateTicket struct {
	ID          string                `param:"id"`
	Name        string                `json:"name" validate:"omitempty,min=3,max=100"`
	Description string                `json:"description" validate:"omitempty,min=3,max=255"`
	Price       int                   `json:"price"`
	Amount      int                   `json:"amount"`
	ImageFile   *multipart.FileHeader `json:"-" validate:"omitempty,isImage"`
}
