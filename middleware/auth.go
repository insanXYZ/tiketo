package middleware

import (
	"os"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

var HasAccToken, HasRefToken echo.MiddlewareFunc

func SetMiddleware() {
	HasAccToken = buildJwtMiddleware(os.Getenv("ACC_JWT_SECRET"), withTokenLookup("header:Authorization"))
	HasRefToken = buildJwtMiddleware(os.Getenv("REF_JWT_SECRET"), withTokenLookup("header:cookie:refresh-token"))
}

type Option func(*echojwt.Config)

func withTokenLookup(token string) Option {
	return func(e *echojwt.Config) {
		e.TokenLookup = token
	}
}

func buildJwtMiddleware(key string, options ...Option) echo.MiddlewareFunc {
	cfg := echojwt.Config{
		SigningKey: []byte(key),
		SuccessHandler: func(c echo.Context) {
			c.Set("user", c.Get("user").(*jwt.Token).Claims.(jwt.MapClaims))
		},
	}

	for _, v := range options {
		v(&cfg)
	}

	return echojwt.WithConfig(cfg)
}
