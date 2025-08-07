package service

import (
	"context"
	"errors"
	"tiketo/dto"
	"tiketo/entity"
	"tiketo/repository"
	"tiketo/util"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
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

func (t *TicketService) HandleGetUserTickets(ctx context.Context, claims jwt.MapClaims) ([]entity.Ticket, error) {
	tickets := make([]entity.Ticket, 0, 10)

	idUser := claims["sub"].(string)

	err := t.ticketRepository.FindUserTickets(ctx, t.db, idUser, tickets)
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	return tickets, err
}

func (t *TicketService) HandleCreateTicket(ctx context.Context, claims jwt.MapClaims, req *dto.CreateTicket) error {
	err := util.ValidateStruct(req)
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
		Quantity:    req.Quantity,
		UserID:      claims["sub"].(string),
		Image:       req.ImageFile.Filename,
	}

	return t.ticketRepository.Create(ctx, t.db, ticket)
}

func (t *TicketService) HandleDelete(ctx context.Context, claims jwt.MapClaims, req *dto.DeleteTicket) error {
	err := util.ValidateStruct(req)
	if err != nil {
		return err
	}

	ticket := &entity.Ticket{
		ID:     req.Id,
		UserID: claims["sub"].(string),
	}

	return t.ticketRepository.Delete(ctx, t.db, ticket)

}

func (t *TicketService) HandleUpdate(ctx context.Context, claims jwt.MapClaims, req *dto.UpdateTicket) error {
	err := util.ValidateStruct(req)
	if err != nil {
		return err
	}

	ticket := &entity.Ticket{
		ID:     req.ID,
		UserID: claims["sub"].(string),
	}

	err = t.ticketRepository.Take(ctx, t.db, ticket)
	if err != nil {
		return err
	}

	if req.ImageFile != nil {
		f, err := req.ImageFile.Open()
		if err != nil {
			return err
		}

		err = util.SaveTicketImage(f, req.ImageFile.Filename)
		if err != nil {
			return err
		}

		ticket.Image = req.ImageFile.Filename
	}

	ticket.Name = req.Name
	ticket.Description = req.Description
	ticket.Price = req.Price
	ticket.Quantity = req.Quantity

	return t.ticketRepository.Save(ctx, t.db, ticket)
}

func (t *TicketService) HandleGet(ctx context.Context, req *dto.GetTicket) (*entity.Ticket, error) {
	err := util.ValidateStruct(req)
	if err != nil {
		return nil, err
	}

	ticket := &entity.Ticket{
		ID: req.Id,
	}

	err = t.ticketRepository.TakeWithUser(ctx, t.db, ticket)
	return ticket, err
}

func (t *TicketService) HandleGetAll(ctx context.Context) ([]entity.Ticket, error) {
	var tickets []entity.Ticket

	err := t.ticketRepository.FindWithUser(ctx, t.db, &tickets)
	return tickets, err
}
