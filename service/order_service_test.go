package service_test

import (
	"testing"
	"tiketo/db"
	"tiketo/repository"
	"tiketo/service"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestOrderService(t *testing.T) {
	mockOrderRepository := repository.NewMockOrderRepository(t)
	mockTicketRepository := repository.NewMockTicketRepository(t)
	mockRedis := db.NewMockRedisClient()

	t.Run("Success get history orders", func(t *testing.T) {
		mockOrderRepository.On("FindAllOrderHistoryUser", mock.Anything, mock.Anything, mock.AnythingOfType("*[]entity.Order"), mock.Anything).Return(nil).Once()

		orderService := service.NewOrderService(mockOrderRepository, nil, nil, mockTicketRepository, mockRedis, nil)
		_, err := orderService.HandleGetHistoryOrders(t.Context(), jwt.MapClaims{
			"sub": "user-1",
		})

		assert.NoError(t, err)
	})

}
