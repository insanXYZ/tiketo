package service_test

import (
	"context"
	"testing"
	"tiketo/db"
	"tiketo/dto"
	"tiketo/entity"
	"tiketo/repository"
	"tiketo/service"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func TestUserService(t *testing.T) {
	mockUserRepository := repository.NewMockUserRepository(t)
	mockRedis := db.NewMockRedisClient()

	t.Run("Success get current user", func(t *testing.T) {
		mockUserRepository.On("Take", mock.Anything, mock.Anything, mock.AnythingOfType("*entity.User")).Return(nil).Once()

		userService := service.NewUserService(mockUserRepository, nil, mockRedis)
		_, err := userService.HandleGetCurrentUser(context.Background(), jwt.MapClaims{
			"sub": "user-1",
		})
		assert.NoError(t, err)
	})

	t.Run("Failed get current user", func(t *testing.T) {
		mockUserRepository.On("Take", mock.Anything, mock.Anything, mock.AnythingOfType("*entity.User")).Return(gorm.ErrRecordNotFound).Once()

		userService := service.NewUserService(mockUserRepository, nil, mockRedis)
		_, err := userService.HandleGetCurrentUser(context.Background(), jwt.MapClaims{
			"sub": "user-2",
		})
		assert.Error(t, err)
	})

	t.Run("Success login", func(t *testing.T) {
		reqPw := "12345678"
		bcryptPw, _ := bcrypt.GenerateFromPassword([]byte(reqPw), bcrypt.DefaultCost)

		mockUserRepository.On("Take", mock.Anything, mock.Anything, mock.AnythingOfType("*entity.User")).Run(func(args mock.Arguments) {
			user := args.Get(2).(*entity.User)
			user.Password = string(bcryptPw)
		}).Return(nil).Once()

		userService := service.NewUserService(mockUserRepository, nil, mockRedis)
		_, _, err := userService.HandleLogin(context.Background(), &dto.Login{
			Email:    "johndoe@example.com",
			Password: reqPw,
		})

		assert.NoError(t, err)
	})

	t.Run("Failed login", func(t *testing.T) {
		validPw := "12345678"
		bcryptPw, _ := bcrypt.GenerateFromPassword([]byte(validPw), bcrypt.DefaultCost)

		mockUserRepository.On("Take", mock.Anything, mock.Anything, mock.AnythingOfType("*entity.User")).Run(func(args mock.Arguments) {
			user := args.Get(2).(*entity.User)
			user.Password = string(bcryptPw)
		}).Return(nil).Once()

		userService := service.NewUserService(mockUserRepository, nil, mockRedis)
		_, _, err := userService.HandleLogin(context.Background(), &dto.Login{
			Email:    "johndoe@example.com",
			Password: "87654321",
		})

		assert.Error(t, err)
	})

	t.Run("Success register", func(t *testing.T) {
		mockUserRepository.On("Take", mock.Anything, mock.Anything, mock.AnythingOfType("*entity.User")).Return(gorm.ErrRecordNotFound).Once()
		mockUserRepository.On("Create", mock.Anything, mock.Anything, mock.AnythingOfType("*entity.User")).Return(nil).Once()

		userService := service.NewUserService(mockUserRepository, nil, mockRedis)
		err := userService.HandleRegister(context.Background(), &dto.Register{
			Name:     "john",
			Email:    "johndoe@example.com",
			Password: "12345678",
		})

		assert.NoError(t, err)
	})

	t.Run("Failed register", func(t *testing.T) {
		mockUserRepository.On("Take", mock.Anything, mock.Anything, mock.AnythingOfType("*entity.User")).Return(nil).Once()
		mockUserRepository.On("Create", mock.Anything, mock.Anything, mock.AnythingOfType("*entity.User")).Return(nil)

		userService := service.NewUserService(mockUserRepository, nil, mockRedis)
		err := userService.HandleRegister(context.Background(), &dto.Register{
			Name:     "john",
			Email:    "johndoe@example.com",
			Password: "12345678",
		})

		assert.Error(t, err)
	})
}
