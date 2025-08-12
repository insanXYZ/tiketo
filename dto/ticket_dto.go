package dto

import (
	"mime/multipart"
)

type Ticket struct {
	ID          string `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Price       int    `json:"price"`
	Image       string `json:"image"`
	Quantity    int    `json:"quantity,omitempty"`
	CreatedAt   string `json:"created_at,omitempty"`
	User        *User  `json:"created_by,omitempty"`
}

type CreateTicket struct {
	Name        string                `json:"name" form:"name" validate:"required,min=3"`
	Description string                `json:"description" form:"description" validate:"required,min=3"`
	Price       int                   `json:"price" form:"price" validate:"required"`
	Quantity    int                   `json:"quantity" form:"quantity" validate:"required"`
	ImageFile   *multipart.FileHeader `validate:"isImage"`
}

type DeleteTicket struct {
	Id string `param:"id" validate:"required"`
}

type GetTicket struct {
	Id string `param:"id" validate:"required"`
}

type UpdateTicket struct {
	ID          string                `param:"id"`
	Name        *string               `json:"name" validate:"omitempty,min=3,max=100"`
	Description *string               `json:"description" validate:"omitempty,min=3,max=255"`
	Price       *int                  `json:"price"`
	Quantity    *int                  `json:"quantity"`
	ImageFile   *multipart.FileHeader `json:"-" validate:"omitempty,isImage"`
}
