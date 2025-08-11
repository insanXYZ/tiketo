package entity

import (
	"database/sql"
	"time"
)

type Status string

const (
	Unpaid Status = "unpaid"
	Paid   Status = "paid"
	Cancel Status = "cancel"
)

type Order struct {
	ID          string       `gorm:"column:id;primaryKey"`
	Status      Status       `gorm:"column:status"`
	UserID      string       `gorm:"column:user_id"`
	Total       int          `gorm:"column:total"`
	CreatedAt   time.Time    `gorm:"column:created_at;autoCreateTime"`
	PaidAt      sql.NullTime `gorm:"column:paid_at"`
	OrderDetail *OrderDetail `gorm:"foreignKey:order_id;references:id"`
	User        *User        `gorm:"foreignKey:user_id;references:id"`
}
