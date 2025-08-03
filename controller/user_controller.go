package controller

import (
	"net/http"
	"tiketo/dto"
	"tiketo/service"
	"tiketo/util/httpresponse"
	"time"

	"github.com/labstack/echo/v4"
)

type UserController struct {
	userService *service.UserService
}

func NewUserController(userService *service.UserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

func (u *UserController) RegisterRoutes(e *echo.Echo) {
	e.POST("/login", u.handleLogin)
	e.POST("/register", u.handleRegister)
}

func (u *UserController) handleLogin(c echo.Context) error {
	var req *dto.Login

	err := c.Bind(req)
	if err != nil {
		return httpresponse.Error(c, "bind err", err)
	}

	accToken, refToken, err := u.userService.Login(c.Request().Context(), req)

	cookie := &http.Cookie{
		Value:    "Bearer " + refToken,
		Path:     "/api/refresh",
		Secure:   true,
		Name:     "Refresh-Token",
		HttpOnly: true,
		Expires:  time.Now().Add(24 * time.Hour),
	}

	c.SetCookie(cookie)

	return httpresponse.Success(c, "success login", echo.Map{
		"access_token": "Bearer " + accToken,
	})

}

func (u *UserController) handleRegister(c echo.Context) error {
	return nil
}
