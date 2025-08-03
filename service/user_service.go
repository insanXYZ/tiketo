package service

import (
	"context"
	"errors"
	"tiketo/dto"
	"tiketo/entity"
	"tiketo/repository"
	"tiketo/util"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
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

func (u *UserService) HandleLogin(ctx context.Context, req *dto.Login) (accToken, refToken string, err error) {

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

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
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

func (u *UserService) HandleRegister(ctx context.Context, req *dto.Register) error {
	err := util.Validator.Struct(req)
	if err != nil {
		return err
	}

	user := &entity.User{
		Email: req.Email,
	}

	err = u.userRepository.Take(ctx, user)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err // error email was used
	}

	b, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Name = req.Name
	user.Password = string(b)

	return u.userRepository.Create(ctx, user)

}
