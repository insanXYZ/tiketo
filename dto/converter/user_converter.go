package converter

import (
	"tiketo/dto"
	"tiketo/entity"
)

func UserEntityToDto(user *entity.User) *dto.User {
	if user == nil {
		return nil
	}

	u := &dto.User{
		Name:  user.Name,
		Email: user.Email,
	}

	return u
}

func UserEntityToNameOnlyDto(user *entity.User) *dto.User {
	if user == nil {
		return nil
	}

	return &dto.User{
		Name: user.Name,
	}
}
