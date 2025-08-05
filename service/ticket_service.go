package service

import (
	"context"
	"tiketo/dto"
	"tiketo/entity"
	"tiketo/repository"
	"tiketo/util"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/insanXYZ/sage"
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

func (t *TicketService) HandleCreateTicket(ctx context.Context, claims jwt.MapClaims, req *dto.CreateTicket) error {
	err := util.ValidateStruct(req)
	if err != nil {
		return err
	}

	err = sage.Validate(req.ImageFile)
	if err != nil {
		return err
	}

	file, err := req.ImageFile.Open()
	if err != nil {
		return err
	}

	defer file.Close()

	err = util.SaveTicketImage(file, req.ImageFile.Filename)
	if err != nil {
		return err
	}

	ticket := &entity.Ticket{
		ID:          uuid.NewString(),
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Amount:      req.Amount,
		UserID:      claims["sub"].(string),
		Image:       req.ImageFile.Filename,
	}

	return t.ticketRepository.Create(ctx, t.db, ticket)
}
