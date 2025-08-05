package converter

import (
	"tiketo/dto"
	"tiketo/entity"
)

func UserEntityToDto(entity *entity.User) *dto.User {
	return &dto.User{
		Name:  entity.Name,
		Email: entity.Email,
	}
}
