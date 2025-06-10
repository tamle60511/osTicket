package api

import (
	"ecommerce/config"
	"ecommerce/internal/api/rest"
	"ecommerce/internal/api/rest/handlers"
	"ecommerce/internal/domain"
	"ecommerce/internal/helper"
	"log"

	"github.com/gofiber/fiber/v3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func StartServer(c config.Appconfig) {
	app := fiber.New()
	db, err := gorm.Open(postgres.Open(c.DSN), &gorm.Config{})
	if err != nil {
		log.Println("failed to connect database")
	}
	auth := helper.SetupAuth(c.AppSecret)
	db.AutoMigrate(&domain.User{})
	rh := &rest.RestHandler{
		App:  app,
		DB:   db,
		Auth: auth,
	}
	SetupRoutes(rh)
	app.Listen(c.ServerPort)
}

func SetupRoutes(rh *rest.RestHandler) {
	handlers.SetupUserRoutes(rh)
}
