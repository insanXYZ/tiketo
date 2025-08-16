package service_test

import (
	"bytes"
	"context"
	"io"
	"mime/multipart"
	"testing"
	"tiketo/db"
	"tiketo/dto"
	"tiketo/repository"
	"tiketo/service"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func NewMockPNGFileHeader() (*multipart.FileHeader, error) {
	pngHeader := []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}
	filename := "default_mock_image.png"
	fileField := "image"

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile(fileField, filename)
	if err != nil {
		return nil, err
	}

	_, err = io.Copy(part, bytes.NewReader(pngHeader))
	if err != nil {
		return nil, err
	}

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	reqBody := bytes.NewReader(body.Bytes())
	req, err := multipart.NewReader(reqBody, writer.Boundary()).ReadForm(10 * 1024 * 1024)
	if err != nil {
		return nil, err
	}

	files := req.File[fileField]
	if len(files) == 0 {
		return nil, multipart.ErrMessageTooLarge
	}

	files[0].Header.Set("Content-Type", "image/png")

	return files[0], nil
}

func TestTicketService(t *testing.T) {
	mockTicketRepository := repository.NewMockTicketRepository(t)
	mockRedis := db.NewMockRedisClient()

	userClaims := jwt.MapClaims{
		"sub": "user-1",
	}

	t.Run("Success create ticket", func(t *testing.T) {
		file, _ := NewMockPNGFileHeader()

		mockTicketRepository.On("Create", mock.Anything, mock.Anything, mock.AnythingOfType("*entity.Ticket")).Return(nil).Once()
		reqTicket := &dto.CreateTicket{
			Name:        "ticket-1",
			Description: "description ticket-1",
			Price:       10000,
			Quantity:    10,
			ImageFile:   file,
		}

		ticketService := service.NewTicketService(mockTicketRepository, nil, mockRedis)
		err := ticketService.HandleCreateTicket(context.Background(), userClaims, reqTicket)
		assert.NoError(t, err)
	})

	t.Run("Failed create ticket", func(t *testing.T) {
		file, _ := NewMockPNGFileHeader()

		mockTicketRepository.On("Create", mock.Anything, mock.Anything, mock.AnythingOfType("*entity.Ticket")).Return(nil)
		reqTicket := &dto.CreateTicket{
			Name:        "ticket-1",
			Description: "description ticket-1",
			Price:       0,
			Quantity:    10,
			ImageFile:   file,
		}

		ticketService := service.NewTicketService(mockTicketRepository, nil, mockRedis)
		err := ticketService.HandleCreateTicket(context.Background(), userClaims, reqTicket)
		assert.Error(t, err)
	})

	t.Run("Success delete ticket", func(t *testing.T) {
		mockTicketRepository.On("Take", mock.Anything, mock.Anything, mock.AnythingOfType("*entity.Ticket")).Return(nil).Once()
		mockTicketRepository.On("Delete", mock.Anything, mock.Anything, mock.AnythingOfType("*entity.Ticket")).Return(nil).Once()

		ticketService := service.NewTicketService(mockTicketRepository, nil, mockRedis)
		err := ticketService.HandleDelete(t.Context(), jwt.MapClaims{
			"sub": "user-1",
		}, &dto.DeleteTicket{
			Id: uuid.NewString(),
		})

		assert.NoError(t, err)
	})

	t.Run("Failed delete ticket", func(t *testing.T) {
		mockTicketRepository.On("Take", mock.Anything, mock.Anything, mock.AnythingOfType("*entity.Ticket")).Return(gorm.ErrRecordNotFound).Once()
		mockTicketRepository.On("Delete", mock.Anything, mock.Anything, mock.AnythingOfType("*entity.Ticket")).Return(nil)

		ticketService := service.NewTicketService(mockTicketRepository, nil, mockRedis)
		err := ticketService.HandleDelete(t.Context(), jwt.MapClaims{
			"sub": "user-1",
		}, &dto.DeleteTicket{
			Id: uuid.NewString(),
		})

		assert.Error(t, err)
	})
}
