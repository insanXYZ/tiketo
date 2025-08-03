package repository

import (
	"tiketo/entity"

	"gorm.io/gorm"
)

type UserRepository struct {
	Repository[*entity.User]
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		Repository: Repository[*entity.User]{
			db: db,
		},
	}
}
