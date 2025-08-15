package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"tiketo/db"
	"tiketo/dto"
	"tiketo/entity"
	"tiketo/repository"
	"tiketo/util"
	"tiketo/util/logger"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/insanXYZ/sage"
	"gorm.io/gorm"
)

type TicketServiceInterface interface {
	HandleGetUserTickets(context.Context, jwt.MapClaims) ([]entity.Ticket, error)
	HandleCreateTicket(context.Context, jwt.MapClaims, *dto.CreateTicket) error
	HandleDelete(context.Context, jwt.MapClaims, *dto.DeleteTicket) error
	HandleUpdate(context.Context, jwt.MapClaims, *dto.UpdateTicket) error
	HandleGetTicket(context.Context, *dto.GetTicket) (*entity.Ticket, error)
	HandleGetTickets(context.Context, *dto.GetTIckets) ([]entity.Ticket, error)
}

type TicketService struct {
	ticketRepository *repository.TicketRepository
	db               *gorm.DB
	redis            db.RedisInterface
}

func NewTicketService(repository *repository.TicketRepository, db *gorm.DB, redis db.RedisInterface) *TicketService {
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

	err = t.ticketRepository.Delete(ctx, t.db, ticket)
	if err != nil {
		return err
	}

	key := fmt.Sprintf("ticket-%s", ticket.ID)
	t.redis.Del(ctx, key)

	return nil
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

	err = t.ticketRepository.Save(ctx, t.db, ticket)
	if err != nil {
		return err
	}

	key := fmt.Sprintf("ticket-%s", req.ID)
	t.redis.Set(ctx, key, ticket, db.ExpRedis)

	return nil
}

func (t *TicketService) HandleGetTicket(ctx context.Context, req *dto.GetTicket) (*entity.Ticket, error) {
	err := util.ValidateStruct(req)
	if err != nil {
		return nil, err
	}

	key := fmt.Sprintf("ticket=%s", req.Id)

	b, err := t.redis.Get(ctx, key)
	if err == nil {
		ticket := new(entity.Ticket)
		if err := json.Unmarshal(b, ticket); err == nil {
			return ticket, nil
		}
	}

	ticket := &entity.Ticket{
		ID: req.Id,
	}

	err = t.ticketRepository.TakeWithUser(ctx, t.db, ticket)
	if err != nil {
		return nil, err
	}

	expSetTicket := time.Duration(5 * time.Minute)

	err = t.redis.Set(ctx, key, ticket, expSetTicket)
	if err != nil {
		logger.Warn(nil, "Err set redis on TicketService.HandleGetTicket")
	}

	return ticket, nil
}

func (t *TicketService) HandleGetTickets(ctx context.Context, req *dto.GetTIckets) ([]entity.Ticket, error) {
	var tickets []entity.Ticket

	err := t.ticketRepository.FindPagingWithJoinUser(ctx, t.db, &tickets, req.Page)
	return tickets, err
}
