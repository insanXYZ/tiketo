package main

import (
	"os"
	"tiketo/db"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err.Error())
	}

	err = os.MkdirAll("assets/image/ticket", os.ModePerm)
	if err != nil {
		panic(err.Error())
	}

	// init db
	pg, err := db.NewPostgresClient()
	if err != nil {
		panic(err.Error())
	}

	redisClient, err := db.NewRedisClient()
	if err != nil {
		panic(err.Error())
	}

	e := echo.New()

	e.Logger.Fatal(e.Start(os.Getenv("APP_PORT")))

}
