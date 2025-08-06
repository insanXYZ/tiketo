package dto

import (
	"mime/multipart"
	"time"
)

type Ticket struct {
	ID          string    `gorm:"column:id"`
	Name        string    `gorm:"column:name"`
	Description string    `gorm:"column:description"`
	Price       int       `gorm:"column:price"`
	Image       string    `gorm:"column:image"`
	Amount      int       `gorm:"column:amount"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime"`
	User        *User     `gorm:"foreignKey:user_id;references:id"`
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
