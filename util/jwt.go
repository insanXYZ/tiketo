package util

import (
	"os"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWT(claims jwt.MapClaims) (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func BuildClaims(name, sub string, exp int64) jwt.MapClaims {
	return jwt.MapClaims{
		"name": name,
		"sub":  sub,
		"exp":  exp,
	}
}
