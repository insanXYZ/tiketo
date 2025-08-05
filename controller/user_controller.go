package controller

import (
	"net/http"
	"tiketo/dto"
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

func (u *UserController) RegisterRoutes(e *echo.Echo) {
	e.POST("/login", u.Login)
	e.POST("/register", u.Register)
	e.GET("/refresh", u.Refresh)
}

func (u *UserController) Login(c echo.Context) error {
	req := new(dto.Login)

	err := c.Bind(req)
	if err != nil {
		return httpresponse.Error(c, "bind err", err)
	}

	accToken, refToken, err := u.userService.HandleLogin(c.Request().Context(), req)

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

func (u *UserController) Register(c echo.Context) error {
	req := new(dto.Register)

	err := c.Bind(req)
	if err != nil {
		return err
	}

	err = u.userService.HandleRegister(c.Request().Context(), req)
	if err != nil {
		return httpresponse.Error(c, "failed register", err)
	}

	return httpresponse.Success(c, "success register", nil)
}

func (u *UserController) Refresh(c echo.Context) error {
	claims := c.Get("user").(jwt.MapClaims)

	accToken, err := u.userService.HandleRefresh(c.Request().Context(), claims)
	if err != nil {
		return httpresponse.Error(c, "failed refresh", err)
	}

	return httpresponse.Success(c, "success refresh", echo.Map{
		"access_token": accToken,
	})
}
