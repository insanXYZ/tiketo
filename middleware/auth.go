package middleware

import (
	"net/http"
	"os"
	"tiketo/dto/message"
	"tiketo/util/httpresponse"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

var HasAccToken, HasRefToken echo.MiddlewareFunc

func InitMiddleware() {
	HasAccToken = echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(os.Getenv("ACC_JWT_SECRET")),
		SuccessHandler: func(c echo.Context) {
			c.Set("user", c.Get("user").(*jwt.Token).Claims.(jwt.MapClaims))
		},
		TokenLookup: "header:Authorization",
	})

	HasRefToken = func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			reftoken, err := c.Cookie("refresh-token")
			if err != nil {
				return httpresponse.Error(c, message.ErrMalformedJwt, nil, http.StatusUnauthorized)
			}

			token, err := jwt.Parse(reftoken.Value, func(t *jwt.Token) (any, error) {
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, echojwt.ErrJWTInvalid
				}

				return []byte(os.Getenv("REF_JWT_SECRET")), nil
			})

			if err != nil {
				return httpresponse.Error(c, message.ErrExpiredJwt, nil, http.StatusUnauthorized)
			}

			claims, ok := token.Claims.(jwt.MapClaims)

			if !ok {
				return httpresponse.Error(c, message.ErrExpiredJwt, nil, http.StatusUnauthorized)
			}

			c.Set("user", claims)

			return next(c)
		}
	}
}
