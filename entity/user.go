package entity

import (
	"encoding/json"
	"time"
)

type User struct {
	ID        string    `gorm:"column:id;primaryKey" json:"id,omitempty"`
	Name      string    `gorm:"column:name" json:"name,omitempty"`
	Email     string    `gorm:"column:email" json:"email,omitempty"`
	Password  string    `gorm:"column:password" json:"password,omitempty"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at,omitempty"`
	UpdatedAt time.Time `gorm:"columd:updated_at;autoCreateTime;autoUpdateTime" json:"updated_at,omitempty"`
	Tickets   []Ticket  `gorm:"foreignKey:user_id;references:id" json:"tickets,omitempty"`
}

func (u *User) MarshalBinary() (data []byte, err error) {
	return json.Marshal(u)
}

func (u *User) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, u)
}
