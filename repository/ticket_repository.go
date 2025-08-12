package repository

import (
	"context"
	"tiketo/entity"

	"gorm.io/gorm"
)

type TicketRepository struct {
	Repository[*entity.Ticket]
}

func NewTicketRepository() *TicketRepository {
	return &TicketRepository{}
}

func (t *TicketRepository) TakeWithUser(ctx context.Context, db *gorm.DB, dst *entity.Ticket) error {
	return db.WithContext(ctx).Joins("User").Take(dst).Error
}

func (t *TicketRepository) FindWithUser(ctx context.Context, db *gorm.DB, dst *[]entity.Ticket) error {
	return db.WithContext(ctx).Joins("User").Find(dst).Error
}

func (t *TicketRepository) FindUserTickets(ctx context.Context, db *gorm.DB, id string, dst *[]entity.Ticket) error {
	return db.WithContext(ctx).Where("user_id = ?", id).Find(dst).Error
}
