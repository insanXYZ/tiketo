package repository

import (
	"context"
	"tiketo/entity"
	"tiketo/util/logger"

	"gorm.io/gorm"
)

type TicketRepositoryInterface interface {
	RepositoryInterface[*entity.Ticket]
	TakeWithUser(context.Context, *gorm.DB, *entity.Ticket) error
	FindPagingWithJoinUser(context.Context, *gorm.DB, *[]entity.Ticket, int) error
	FindUserTickets(context.Context, *gorm.DB, string, *[]entity.Ticket) error
}

type TicketRepository struct {
	Repository[*entity.Ticket]
}

func NewTicketRepository() *TicketRepository {
	return &TicketRepository{}
}

func (t *TicketRepository) TakeWithUser(ctx context.Context, db *gorm.DB, dst *entity.Ticket) error {
	err := db.WithContext(ctx).Joins("User").Take(dst).Error
	if err != nil {
		logger.WarnMethod("TicketRepository.TakeWithUser", err)
	}
	return err
}

func (t *TicketRepository) FindPagingWithJoinUser(ctx context.Context, db *gorm.DB, dst *[]entity.Ticket, page int) error {
	numberPerPage := 10

	err := db.WithContext(ctx).Joins("User").Limit(10).Offset(page * numberPerPage).Find(dst).Error
	if err != nil {
		logger.WarnMethod("TicketRepository.FindPagingWithJoinUser", err)
	}
	return err
}

func (t *TicketRepository) FindUserTickets(ctx context.Context, db *gorm.DB, id string, dst *[]entity.Ticket) error {
	err := db.WithContext(ctx).Where("user_id = ?", id).Find(dst).Error
	if err != nil {
		logger.WarnMethod("TicketRepository.FindUserTickets", err)
	}
	return err
}
