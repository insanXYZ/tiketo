package service

import (
	"context"
	"errors"
	"path/filepath"
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

func (t *TicketService) HandleGetUserTickets(ctx context.Context, claims jwt.MapClaims) ([]entity.Ticket, error) {
	var tickets []entity.Ticket

	idUser := claims["sub"].(string)

	err := t.ticketRepository.FindUserTickets(ctx, t.db, idUser, &tickets)
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

	err = sage.Validate(req.ImageFile)
	if err != nil {
		return err
	}

	file, err := req.ImageFile.Open()
	if err != nil {
		return err
	}

	defer file.Close()

	id := uuid.NewString()
	filename := util.GenerateFilenameTicket(id, filepath.Ext(req.ImageFile.Filename))

	err = util.SaveTicketImage(file, filename)
	if err != nil {
		return err
	}

	ticket := &entity.Ticket{
		ID:          id,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Quantity:    req.Quantity,
		UserID:      claims["sub"].(string),
		Image:       filename,
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

	err = t.ticketRepository.Take(ctx, t.db, ticket)
	if err != nil {
		return err
	}

	go util.DeleteTicketImage(ticket.Image)

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

		err = util.SaveTicketImage(f, util.GenerateFilenameTicket(ticket.ID, filepath.Ext(req.ImageFile.Filename)))
		if err != nil {
			return err
		}

		ticket.Image = req.ImageFile.Filename
	}

	if req.Name != nil {
		ticket.Name = *req.Name
	}

	if req.Description != nil {
		ticket.Description = *req.Description
	}

	if req.Price != nil {
		ticket.Price = *req.Price
	}

	if req.Quantity != nil {
		ticket.Quantity = *req.Quantity
	}

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

func (t *TicketService) HandleGetTickets(ctx context.Context, req *dto.GetTIckets) ([]entity.Ticket, error) {
	var tickets []entity.Ticket

	err := t.ticketRepository.FindPagingWithJoinUser(ctx, t.db, &tickets, req.Page)
	return tickets, err
}
