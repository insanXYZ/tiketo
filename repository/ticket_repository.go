package repository

import (
	"context"
	"tiketo/dst"
	"tiketo/entity"

	"gorm.io/gorm"
)

type TicketRepository struct {
	Repository[*dst.Ticket]
}

func NewTicketRepository() *TicketRepository {
	return &TicketRepository{}
}

func (t *TicketRepository) TakeWithUser(ctx context.Context, db *gorm.DB, dst *entity.Ticket) error {
	return db.WithContext(ctx).Joins("User").Take(dst).Error
}

func (t *TicketRepository) FindWithUser(ctx context.Context, db *gorm.DB, dst []entity.Ticket) error {
	return db.WithContext(ctx).Joins("User").Find(dst).Error
}

func (t *TicketRepository) FindUserTickets(ctx context.Context, db *gorm.DB, model *entity.User, dst []entity.Ticket) error {
	return db.WithContext(ctx).Model(model).Find(dst).Error
}
