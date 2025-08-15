package service_test

import (
	"context"
	"testing"
	"tiketo/db"
	"tiketo/repository"
	"tiketo/service"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUserService(t *testing.T) {
	mockUserRepository := repository.NewMockUserRepository(t)
	mockRedis := db.NewMockRedisClient()

	t.Run("Success get current user", func(t *testing.T) {
		mockUserRepository.On("Take", mock.Anything, mock.Anything, mock.AnythingOfType("*entity.User")).Return(nil)

		userService := service.NewUserService(mockUserRepository, nil, mockRedis)
		_, err := userService.HandleGetCurrentUser(context.Background(), jwt.MapClaims{
			"sub": "user-1",
		})
		assert.NoError(t, err)
	})

}
