package main

import (
	"os"
	"tiketo/controller"
	"tiketo/db"
	"tiketo/middleware"
	"tiketo/repository"
	"tiketo/service"
	"tiketo/util/logger"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err.Error())
	}

	middleware.InitMiddleware()
	logger.InitLogger()

	err = db.Migrate()
	if err != nil {
		panic(err.Error())
	}

	err = os.MkdirAll("assets/image/ticket", os.ModePerm)
	if err != nil {
		panic(err.Error())
	}

	// init db
	gorm, err := db.NewGorm()
	if err != nil {
		panic(err.Error())
	}

	redisClient, err := db.NewRedisClient()
	if err != nil {
		panic(err.Error())
	}

	userRepository := repository.NewUserRepository()
	ticketRepository := repository.NewTicketRepository()

	userService := service.NewUserService(userRepository, gorm, redisClient)
	ticketService := service.NewTicketService(ticketRepository, gorm, redisClient)

	userController := controller.NewUserController(userService)
	ticketController := controller.NewTicketController(ticketService)

	e := echo.New()
	api := e.Group("/api")
	userController.RegisterRoutes(api)
	ticketController.RegisterRoutes(api)

	e.Logger.Fatal(e.Start(os.Getenv("APP_PORT")))
}
