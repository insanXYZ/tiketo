package service

import (
	"context"
	"tiketo/dto"
	"tiketo/entity"
	"tiketo/repository"
	"tiketo/util"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
)

type UserService struct {
	userRepository *repository.UserRepository
	redis          *redis.Client
}

func NewUserService(userRepository *repository.UserRepository) *UserService {
	return &UserService{
		userRepository: userRepository,
	}
}

func (u *UserService) Login(ctx context.Context, req *dto.Login) (accToken, refToken string, err error) {

	err = util.Validator.Struct(req)
	if err != nil {
		return
	}

	user := &entity.User{
		Email: req.Email,
	}

	err = u.userRepository.Take(ctx, user)
	if err != nil {
		return
	}

	claims := jwt.MapClaims{
		"exp":  time.Now().Add(15 * time.Minute).Unix(),
		"sub":  user.ID,
		"name": user.Name,
	}

	accToken, err = util.GenerateJWT(claims)

	if err != nil {
		return
	}

	claims["exp"] = time.Now().Add(24 * time.Hour).Unix()

	refToken, err = util.GenerateJWT(claims)

	return

}
