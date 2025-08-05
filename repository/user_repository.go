package repository

import (
	"tiketo/entity"
)

type UserRepository struct {
	Repository[*entity.User]
}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}
