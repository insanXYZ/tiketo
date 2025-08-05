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

func (t *TicketRepository) TakeWithUser(ctx context.Context, db *gorm.DB, entity *entity.Ticket) error {
	return db.WithContext(ctx).Joins("User").Take(entity).Error
}

func (t *TicketRepository) FindWithUser(ctx context.Context, db *gorm.DB, entity []entity.Ticket) error {
	return db.WithContext(ctx).Joins("User").Find(entity).Error
}
