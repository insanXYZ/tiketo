package controller

import (
	"net/http"
	"tiketo/dto"
	"tiketo/dto/converter"
	"tiketo/dto/message"
	"tiketo/middleware"
	"tiketo/service"
	"tiketo/util/httpresponse"
	"time"

	"github.com/golang-jwt/jwt/v5"
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

func (u *UserController) RegisterRoutes(e *echo.Group) {
	e.POST("/login", u.Login)
	e.POST("/register", u.Register)
	e.GET("/refresh", u.Refresh, middleware.HasRefToken)
	e.GET("/me", u.GetCurrentUser, middleware.HasAccToken)
}

func (u *UserController) Login(c echo.Context) error {
	req := new(dto.Login)

	err := c.Bind(req)
	if err != nil {
		return httpresponse.Error(c, message.ErrBind, err)
	}

	accToken, refToken, err := u.userService.HandleLogin(c.Request().Context(), req)
	if err != nil {
		return httpresponse.Error(c, message.ErrLogin, err)
	}

	cookie := &http.Cookie{
		Value:    refToken,
		Path:     "/api/refresh",
		Secure:   true,
		Name:     "refresh-token",
		HttpOnly: true,
		Expires:  time.Now().Add(24 * time.Hour),
	}

	c.SetCookie(cookie)

	return httpresponse.Success(c, "success login", echo.Map{
		"access-token": accToken,
	})

}

func (u *UserController) Register(c echo.Context) error {
	req := new(dto.Register)

	err := c.Bind(req)
	if err != nil {
		return httpresponse.Error(c, message.ErrBind, err)
	}

	err = u.userService.HandleRegister(c.Request().Context(), req)
	if err != nil {
		return httpresponse.Error(c, message.ErrRegister, err)
	}

	return httpresponse.Success(c, message.SuccessRegister, nil)
}

func (u *UserController) Refresh(c echo.Context) error {
	claims := c.Get("user").(jwt.MapClaims)

	accToken, err := u.userService.HandleRefresh(c.Request().Context(), claims)
	if err != nil {
		return httpresponse.Error(c, message.ErrRefresh, err)
	}

	return httpresponse.Success(c, message.SuccessRefresh, echo.Map{
		"access-token": accToken,
	})
}

func (u *UserController) GetCurrentUser(c echo.Context) error {
	claims := c.Get("user").(jwt.MapClaims)

	user, err := u.userService.HandleGetCurrentUser(c.Request().Context(), claims)
	if err != nil {
		return httpresponse.Error(c, message.ErrGetCurrentUser, err)
	}

	return httpresponse.Success(c, message.SuccessGetCurrentUser, converter.UserEntityToDto(user))
}
