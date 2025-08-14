package entity

import (
	"encoding/json"
	"time"
)

type Ticket struct {
	ID          string    `gorm:"column:id"`
	Name        string    `gorm:"column:name"`
	Description string    `gorm:"column:description"`
	UserID      string    `gorm:"column:user_id"`
	Price       int       `gorm:"column:price"`
	Image       string    `gorm:"column:image"`
	Quantity    int       `gorm:"column:quantity"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   time.Time `gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`
	User        *User     `gorm:"foreignKey:user_id;references:id"`
}

func (t *Ticket) MarshalBinary() (data []byte, err error) {
	return json.Marshal(t)
}

func (t *Ticket) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, t)
}
