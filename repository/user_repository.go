package repository

import (
	"tiketo/entity"
)

type UserRepositoryInterface interface {
	RepositoryInterface[*entity.User]
}

type UserRepository struct {
	Repository[*entity.User]
}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}
