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
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct {
	userRepository *repository.UserRepository
	db             *gorm.DB
	redis          *redis.Client
}

func NewUserService(userRepository *repository.UserRepository, db *gorm.DB, redis *redis.Client) *UserService {
	return &UserService{
		userRepository: userRepository,
		db:             db,
		redis:          redis,
	}
}

const (
	Exp_Acc_Token = 15 * time.Minute
	Exp_Ref_Token = 24 * time.Hour
)

func (u *UserService) HandleLogin(ctx context.Context, req *dto.Login) (accToken, refToken string, err error) {

	err = util.Validator.Struct(req)
	if err != nil {
		return
	}

	user := &entity.User{
		Email: req.Email,
	}

	err = u.userRepository.Take(ctx, u.db, user)
	if err != nil {
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return
	}

	claims := jwt.MapClaims{
		"exp":  time.Now().Add(Exp_Acc_Token).Unix(),
		"sub":  user.ID,
		"name": user.Name,
	}

	accToken, err = util.GenerateJWT(claims)

	if err != nil {
		return
	}

	claims["exp"] = time.Now().Add(Exp_Ref_Token).Unix()

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

	err = u.userRepository.Take(ctx, u.db, user)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err // error email was used
	}

	b, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.ID = uuid.NewString()
	user.Name = req.Name
	user.Password = string(b)

	return u.userRepository.Create(ctx, u.db, user)

}

func (u *UserService) HandleRefresh(ctx context.Context, claims jwt.MapClaims) (string, error) {
	return util.GenerateJWT(
		util.BuildClaims(
			claims["name"].(string),
			claims["sub"].(string),
			time.Now().Add(Exp_Acc_Token).Unix(),
		),
	)
}
