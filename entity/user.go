package entity

import (
	"time"
)

type User struct {
	ID        string    `gorm:"column:id;primaryKey"`
	Name      string    `gorm:"column:name"`
	Email     string    `gorm:"column:email"`
	Password  string    `gorm:"column:password"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"columd:updated_at;autoCreateTime;autoUpdateTime"`
	Tickets   []Ticket  `gorm:"foreignKey:user_id;references:id"`
}
