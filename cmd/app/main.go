package main

import (
	"errors"
	"os"
	"tiketo/controller"
	"tiketo/db"
	"tiketo/middleware"
	"tiketo/repository"
	"tiketo/service"
	"tiketo/util/logger"

	"github.com/golang-migrate/migrate/v4"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err.Error())
	}

	logger.InitLogger()
	middleware.InitMiddleware()

	err = os.MkdirAll("assets/image/ticket", os.ModePerm)
	if err != nil {
		logger.Fatal(nil, "Failed create directory assets/image/ticket", err.Error())
	}

	// init db
	gorm, err := db.NewGorm()
	if err != nil {
		logger.Fatal(nil, "Failed connect database", err.Error())
	}

	err = db.Migrate()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		logger.Fatal(nil, "Failed migrate table database", err.Error())
	}

	redisClient, err := db.NewRedisClient()
	if err != nil {
		logger.Fatal(nil, "Failed connect redis", err.Error())
	}

	logger.Info(nil, "Load all repository")
	userRepository := repository.NewUserRepository()
	ticketRepository := repository.NewTicketRepository()
	orderRepository := repository.NewOrderRepository()
	orderDetailRepository := repository.NewOrderDetailRepository()

	logger.Info(nil, "Load all service")
	userService := service.NewUserService(userRepository, gorm, redisClient)
	ticketService := service.NewTicketService(ticketRepository, gorm, redisClient)
	orderService := service.NewOrderService(orderRepository, orderDetailRepository, userRepository, ticketRepository, redisClient, gorm)

	logger.Info(nil, "Load all controller")
	userController := controller.NewUserController(userService)
	ticketController := controller.NewTicketController(ticketService)
	orderController := controller.NewOrderController(orderService)

	logger.Info(nil, "Init app")
	e := echo.New()
	e.Use(middleware.LoggingRequest)

	logger.Info(nil, "Register all route")
	api := e.Group("/api")
	userController.RegisterRoutes(api)
	ticketController.RegisterRoutes(api)
	orderController.RegisterRoutes(api)

	logger.Info(nil, "App started successfully")
	e.Logger.Fatal(e.Start(os.Getenv("APP_PORT")))
}
