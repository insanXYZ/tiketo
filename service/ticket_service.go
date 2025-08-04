package service

import (
	"tiketo/repository"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type TicketService struct {
	ticketRepository *repository.TicketRepository
	db               *gorm.DB
	redis            *redis.Client
}

func NewTicketService(repository *repository.TicketRepository, db *gorm.DB, redis *redis.Client) *TicketService {
	return &TicketService{
		ticketRepository: repository,
		db:               db,
		redis:            redis,
	}
}
