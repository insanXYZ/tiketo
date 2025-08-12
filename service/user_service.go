package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"tiketo/dto"
	"tiketo/entity"
	"tiketo/repository"
	"tiketo/util"
	"tiketo/util/logger"
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
	ExpAccToken = 15 * time.Minute
	ExpRefToken = 24 * time.Hour
)

func (u *UserService) HandleLogin(ctx context.Context, req *dto.Login) (accToken, refToken string, err error) {
	logger.EnteringMethod("UserService.HandleLogin")
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
		"exp":  time.Now().Add(ExpAccToken).Unix(),
		"sub":  user.ID,
		"name": user.Name,
	}

	accToken, err = util.GenerateAccToken(claims)

	if err != nil {
		return
	}

	claims["exp"] = time.Now().Add(ExpRefToken).Unix()

	refToken, err = util.GenerateRefToken(claims)

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
		return err
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

	c := util.BuildClaims(
		claims["name"].(string),
		claims["sub"].(string),
		time.Now().Add(ExpAccToken).Unix(),
	)

	return util.GenerateAccToken(c)
}

func (u *UserService) HandleGetCurrentUser(ctx context.Context, claims jwt.MapClaims) (*entity.User, error) {

	key := fmt.Sprintf("user-%s", claims["sub"].(string))

	val, err := u.redis.Get(ctx, key).Bytes()
	if err == nil {
		logger.Info(nil, "Retrieve user from redis")

		dst := new(entity.User)
		if err := json.Unmarshal(val, dst); err == nil {
			return dst, nil
		}
	}

	user := &entity.User{
		ID: claims["sub"].(string),
	}

	err = u.userRepository.Take(ctx, u.db, user)
	if err != nil {
		return nil, err
	}

	expSetUser := time.Duration(5 * time.Minute)

	logger.Info(nil, "Trying set user to redis")

	err = u.redis.Set(ctx, key, user, expSetUser).Err()
	if err != nil {
		logger.Warn(nil, "Err redis set on UserService.HandleGetCurrentUser :", err.Error())
	}

	return user, nil
}
